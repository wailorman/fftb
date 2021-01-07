package ukvs_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/files"

	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/ukvs"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/inmemory"
)

// StoresExample _
type StoresExample struct {
	driverName  string
	persistent  bool
	buildClient func(ctx context.Context) (ukvs.IStore, error)
}

var storesTbl []*StoresExample
var gctx, cancel = context.WithCancel(context.Background())

var clients = make(map[string]ukvs.IStore)

var localfileStoragePath = files.NewTempPath(
	fmt.Sprintf(
		"ukvs_test_%s/%s",
		uuid.New().String(),
		uuid.New().String(),
	),
).FullPath()

func buildLocalfileClient(ctx context.Context) (*localfile.Client, error) {
	lfFile := files.NewFile(localfileStoragePath)
	lfPath := lfFile.BuildPath()

	err := lfPath.Create()

	if err != nil {
		return nil, err
	}

	c, err := localfile.NewClient(ctx, lfFile.FullPath())

	if err != nil {
		return nil, err
	}

	if clients["localfile"] == nil {
		clients["localfile"] = c
	}

	return c, nil
}

func setup() {
	storesTbl = make([]*StoresExample, 0)

	storesTbl = append(
		storesTbl,
		&StoresExample{
			"inmemory",
			false,
			func(ctx context.Context) (ukvs.IStore, error) {
				c := inmemory.NewClient(ctx)

				if clients["inmemory"] == nil {
					clients["inmemory"] = c
				}

				return c, nil
			},
		},
	)

	storesTbl = append(
		storesTbl,
		&StoresExample{
			"localfile",
			true,
			func(ctx context.Context) (ukvs.IStore, error) { return buildLocalfileClient(ctx) },
		},
	)
}

func teardown() {
	cancel()
	for _, client := range clients {
		<-client.Closed()
	}

	err := files.NewPath(localfileStoragePath).Destroy()

	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test__GetSet(t *testing.T) {
	for _, tableItem := range storesTbl {
		store, err := tableItem.buildClient(gctx)

		assert.Nil(t, err)

		storedValue := []byte("value")

		err = store.Set("key", storedValue)

		assert.Nil(t, err, fmt.Sprintf("driver %s: Setting value error", tableItem.driverName))

		receivedVal, err := store.Get("key")

		assert.Nil(t, err, fmt.Sprintf("driver %s: Getting value error", tableItem.driverName))

		assert.Equal(t, storedValue, receivedVal)
	}
}

func Test__Destroy(t *testing.T) {
	for _, tableItem := range storesTbl {
		store, err := tableItem.buildClient(gctx)

		assert.Nil(t, err)

		err = store.Set("key", []byte("value"))

		assert.Nil(t, err, fmt.Sprintf("driver %s: Set error", tableItem.driverName))

		err = store.Destroy("key")
		assert.Nil(t, err, fmt.Sprintf("driver %s: Destroy error", tableItem.driverName))

		val, err := store.Get("key")
		assert.Nil(t, val, fmt.Sprintf("driver %s: Destroyed result", tableItem.driverName))
		assert.Equal(t, err, ukvs.ErrNotFound, fmt.Sprintf("driver %s: Not found destroyed", tableItem.driverName))

		err = store.Destroy("other_key111")
		assert.Nil(t, err, fmt.Sprintf("driver %s: Destroy non-existent key error", tableItem.driverName))
	}
}

func Test__GetAll(t *testing.T) {
	for _, tableItem := range storesTbl {
		store, err := tableItem.buildClient(gctx)

		assert.Nil(t, err)

		store.Set("key1", []byte("value"))
		store.Set("key2", []byte("value"))
		store.Set("key3", []byte("value"))

		ctx, cancel := context.WithCancel(context.Background())

		_, values, _ := store.GetAll(ctx)

		val1, ok1 := <-values
		cancel()
		time.Sleep(100 * time.Millisecond)
		val2, ok2 := <-values
		val3, ok3 := <-values

		assert.Equal(t, string([]byte("value")), string(val1), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, true, ok1, fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, string(""), string(val2), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, false, ok2, fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, string(""), string(val3), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, false, ok3, fmt.Sprintf("driver: %s", tableItem.driverName))
	}
}

func Test__FindAll(t *testing.T) {
	for _, tableItem := range storesTbl {
		store, err := tableItem.buildClient(gctx)

		assert.Nil(t, err)

		store.Set("orders/1", []byte("orders"))
		store.Set("orders/2", []byte("orders"))
		store.Set("orders/3", []byte("orders"))
		store.Set("orders/4", []byte("orders"))
		store.Set("tasks/1", []byte("tasks"))
		store.Set("tasks/2", []byte("tasks"))
		store.Set("tasks/3", []byte("tasks"))
		store.Set("tasks/4", []byte("tasks"))

		ctx, cancel := context.WithCancel(context.Background())

		_, values, _ := store.FindAll(ctx, "orders/*")

		val1, ok1 := <-values
		cancel()
		time.Sleep(100 * time.Millisecond)
		val2, ok2 := <-values
		val3, ok3 := <-values

		assert.Equal(t, string([]byte("orders")), string(val1), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, true, ok1, fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, string(""), string(val2), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, false, ok2, fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, string(""), string(val3), fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, false, ok3, fmt.Sprintf("driver: %s", tableItem.driverName))
	}
}

func Test__ExpireAt(t *testing.T) {
	for _, tableItem := range storesTbl {
		store, err := tableItem.buildClient(gctx)

		assert.Nil(t, err)

		store.Set("key1", []byte("value"))
		store.ExpireAt("key1", time.Now().Add(1*time.Minute))

		val, err := store.Get("key1")

		assert.Nil(t, err, "Not expired err", fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Equal(t, string([]byte("value")), string(val), fmt.Sprintf("driver: %s", tableItem.driverName))

		store.Set("key1", []byte("value"))
		store.ExpireAt("key1", time.Now().Add(100*time.Millisecond))
		time.Sleep(300 * time.Millisecond)

		val, err = store.Get("key1")

		assert.Equal(t, ukvs.ErrNotFound, err, fmt.Sprintf("driver: %s", tableItem.driverName))
		assert.Nil(t, val, fmt.Sprintf("driver: %s", tableItem.driverName))
	}
}

func Test__Persistence(t *testing.T) {
	for _, tableItem := range storesTbl {
		if !tableItem.persistent {
			continue
		}

		ctx1, cancel1 := context.WithCancel(context.TODO())
		store1, err := tableItem.buildClient(ctx1)
		assert.Nil(t, err)

		store1.Set("persistent_key", []byte("persistent_value"))
		cancel1()
		<-store1.Closed()

		ctx2 := context.TODO()
		store2, err := tableItem.buildClient(ctx2)
		assert.Nil(t, err)

		value2, _ := store2.Get("persistent_key")
		assert.Equal(t, []byte("persistent_value"), value2)
	}
}
