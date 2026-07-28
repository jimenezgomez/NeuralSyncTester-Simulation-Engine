package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// A generic batch insert function signature.
// Takes a DB handle, a context, and the batch of queued items.
type InsertFunc[T any] func(ctx context.Context, db *sql.DB, items []T) error

type DBManager[T any] struct {
	db          *sql.DB
	insertFunc  InsertFunc[T]
	buffer      []T
	mu          sync.Mutex
	flushTicker *time.Ticker
	flushChan   chan struct{}

	maxBatchSize  int
	flushInterval time.Duration
}

// NewDBManager creates a new manager for type T
func NewDBManager[T any](db *sql.DB, insertFunc InsertFunc[T], maxBatchSize int, flushInterval time.Duration) *DBManager[T] {
	m := &DBManager[T]{
		db:            db,
		insertFunc:    insertFunc,
		buffer:        make([]T, 0, maxBatchSize),
		maxBatchSize:  maxBatchSize,
		flushInterval: flushInterval,
		flushChan:     make(chan struct{}, 1),
		flushTicker:   time.NewTicker(flushInterval),
	}

	go m.loop()
	return m
}

// Add queues an item for batch insertion
func (m *DBManager[T]) Add(item T) {
	m.mu.Lock()
	m.buffer = append(m.buffer, item)
	full := len(m.buffer) >= m.maxBatchSize
	m.mu.Unlock()

	if full {
		// non-blocking signal
		select {
		case m.flushChan <- struct{}{}:
		default:
		}
	}
}

// loop flushes periodically or on demand
func (m *DBManager[T]) loop() {
	for {
		select {
		case <-m.flushTicker.C:
			m.Flush(context.Background())
		case <-m.flushChan:
			m.Flush(context.Background())
		}
	}
}

// Flush writes the buffer to the DB
func (m *DBManager[T]) Flush(ctx context.Context) {
	m.mu.Lock()
	if len(m.buffer) == 0 {
		m.mu.Unlock()
		return
	}
	batch := m.buffer
	m.buffer = make([]T, 0, m.maxBatchSize)
	m.mu.Unlock()
	if err := m.insertFunc(ctx, m.db, batch); err != nil {
		// TODO: retry policy
		fmt.Printf("ERROR INSERTING INTO DB! (%d items lost): %v\n", len(batch), err)
	}
}

// Close stops the ticker and flushes remaining items
func (m *DBManager[T]) Close(ctx context.Context) {
	m.flushTicker.Stop()
	m.Flush(ctx)
}
