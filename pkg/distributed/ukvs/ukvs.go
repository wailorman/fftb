package ukvs

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound _
var ErrNotFound = errors.New("Not found")

// IStore _
type IStore interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte) error
	GetAll(ctx context.Context) (chan []byte, chan error)
	FindAll(ctx context.Context, pattern string) (chan []byte, chan error)
	Destroy(key string) error
	ExpireAt(key string, time time.Time) error
	Closed() <-chan struct{}
}
