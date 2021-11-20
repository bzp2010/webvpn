package storage

import "context"

// Keypair is data key-value pairs
type Keypair struct {
	Key   string
	Value string
}

// Storage is abstract interface for database
type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	List(ctx context.Context, key string) ([]Keypair, error)
	Create(ctx context.Context, key, val string) error
	Update(ctx context.Context, key, val string) error
	BatchDelete(ctx context.Context, keys []string) error
	Watch(ctx context.Context, key string) <-chan WatchEvent
}

// EventType is the type of data change
type EventType string

// WatchEvent is the record of datasource's data change
type WatchEvent struct {
	Keypair
	Error    error
	Canceled bool
	Type     EventType
}
