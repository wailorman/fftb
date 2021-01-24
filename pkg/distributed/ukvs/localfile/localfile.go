package localfile

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailorman/fftb/pkg/distributed/ukvs"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

var autoFlushFrequency = time.Duration(5 * time.Second)
var autoPruneFrequency = time.Duration(30 * time.Second)

var errExpired = errors.New("Container expired")
var errEmptyValue = errors.New("Container has empty value")

// Client _
type Client struct {
	storagePath          string
	registry             *Registry
	lock                 sync.Mutex
	ctx                  context.Context
	changedFromLastFlush bool
	closed               chan struct{}
}

// Container _
type Container struct {
	Val       []byte     `json:"val"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Registry _
type Registry struct {
	Version string                `json:"version"`
	Data    map[string]*Container `json:"data"`
}

// NewClient _
func NewClient(ctx context.Context, storagePath string) (*Client, error) {
	c := &Client{
		storagePath: storagePath,
		registry: &Registry{
			Version: "1",
			Data:    make(map[string]*Container),
		},
		lock:   sync.Mutex{},
		ctx:    ctx,
		closed: make(chan struct{}),
	}

	if err := c.Init(); err != nil {
		return nil, errors.Wrap(err, "Initialization error")
	}

	go func() {
		for range time.Tick(autoFlushFrequency) {
			select {
			case <-c.ctx.Done():
				c.Flush()
				return
			default:
				c.Flush()
			}
		}
	}()

	go func() {
		for range time.Tick(autoPruneFrequency) {
			select {
			case <-c.ctx.Done():
				c.Prune()
				return
			default:
				c.Prune()
			}
		}
	}()

	go func() {
		<-c.ctx.Done()
		c.Flush()
		close(c.closed)
		return
	}()

	return c, nil
}

// Init _
func (c *Client) Init() error {
	storageFile := files.NewFile(c.storagePath)

	if !storageFile.IsExist() {
		if err := storageFile.Create(); err != nil {
			return errors.Wrap(err, "Creating storage file")
		}

		if err := c.Flush(); err != nil {
			return errors.Wrap(err, "Writing initial state to storage file")
		}
	} else {
		storageContent, err := storageFile.ReadAllContent()

		if err != nil {
			return errors.Wrap(err, "Reading storage file content")
		}

		storageFileIsNotEmpty := len(strings.Trim(storageContent, " ")) > 0

		if storageFileIsNotEmpty {
			err = json.Unmarshal([]byte(storageContent), &c.registry)

			if err != nil {
				return errors.Wrap(err, "Unmarshaling storage content")
			}
		}
	}

	return nil
}

// Flush _
func (c *Client) Flush() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.changedFromLastFlush {
		return nil
	}

	storageFile := files.NewFile(c.storagePath)

	err := storageFile.Create()

	if err != nil {
		return errors.Wrap(err, "Truncating storage file")
	}

	bRegistry, err := json.Marshal(c.registry)

	if err != nil {
		return errors.Wrap(err, "Marshaling registry to bytes")
	}

	bReader := bytes.NewReader(bRegistry)

	storageWriter, err := storageFile.WriteContent()

	_, err = io.Copy(storageWriter, bReader)

	if err != nil {
		return errors.Wrap(err, "Writing registry content to file")
	}

	c.changedFromLastFlush = false

	return nil
}

// Prune _
func (c *Client) Prune() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, cont := range c.registry.Data {
		_, extractionErr := c.extractContainer(cont)

		select {
		case <-c.ctx.Done():
			return nil
		default:
			if errors.Is(extractionErr, errExpired) || errors.Is(extractionErr, errEmptyValue) {
				delete(c.registry.Data, key)
			}
		}
	}

	return nil
}

// Get _
func (c *Client) Get(key string) ([]byte, error) {
	val, err := c.extractContainer(c.registry.Data[key])

	if err != nil {
		if errors.Is(err, errExpired) || errors.Is(err, errEmptyValue) {
			return nil, ukvs.ErrNotFound
		}

		return nil, err
	}

	return val, nil
}

// Set _
func (c *Client) Set(key string, val []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.registry.Data[key] = &Container{
		Val: val,
	}

	c.changedFromLastFlush = true

	return nil
}

// GetAll _
func (c *Client) GetAll(fctx context.Context) (chan []byte, chan error) {
	return c.FindAll(fctx, "*")
}

// FindAll _
func (c *Client) FindAll(fctx context.Context, pattern string) (chan []byte, chan error) {
	results := make(chan []byte, 1)
	failures := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(failures)

		for key, cont := range c.registry.Data {
			val, extractionErr := c.extractContainer(cont)

			keyMatches, matchErr := filepath.Match(pattern, key)

			if matchErr != nil {
				failures <- matchErr
				return
			}

			if !keyMatches || extractionErr != nil {
				select {
				case <-c.ctx.Done():
					return
				case <-fctx.Done():
					return
				default:
					continue
				}
			}

			select {
			case <-c.ctx.Done():
				return
			case <-fctx.Done():
				return
			default:
				results <- val
			}
		}
	}()

	return results, failures
}

// Destroy _
func (c *Client) Destroy(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.registry.Data, key)

	c.changedFromLastFlush = true

	return nil
}

// ExpireAt _
func (c *Client) ExpireAt(key string, t time.Time) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	cont := c.registry.Data[key]
	_, err := c.extractContainer(cont)

	if errors.Is(err, errEmptyValue) || errors.Is(err, errExpired) {
		return ukvs.ErrNotFound
	}

	if time.Now().After(t) {
		delete(c.registry.Data, key)
		return nil
	}

	cont.ExpiresAt = &t
	c.changedFromLastFlush = true

	return nil
}

func (c *Client) extractContainer(cont *Container) ([]byte, error) {
	if cont == nil || cont.Val == nil || len(cont.Val) == 0 {
		return nil, errEmptyValue
	}

	if cont.ExpiresAt != nil && time.Now().After(*cont.ExpiresAt) {
		return nil, errExpired
	}

	return cont.Val, nil
}

// Closed _
func (c *Client) Closed() <-chan struct{} {
	return c.closed
}
