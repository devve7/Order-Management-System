package product

import (
	"context"
	"errors"
	"sync"
)

var errCacheMiss = errors.New("cache miss")

type mockCache struct {
	mu   sync.RWMutex
	data map[string]string

	getErr    error
	setErr    error
	deleteErr error

	getCalls    []string
	setCalls    []string
	deleteCalls []string
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string]string),
	}
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getCalls = append(m.getCalls, key)

	if m.getErr != nil {
		return "", m.getErr
	}

	value, ok := m.data[key]
	if !ok {
		return "", errCacheMiss
	}

	return value, nil
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.setCalls = append(m.setCalls, key)

	if m.setErr != nil {
		return m.setErr
	}

	m.data[key] = value
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.deleteCalls = append(m.deleteCalls, key)

	if m.deleteErr != nil {
		return m.deleteErr
	}

	delete(m.data, key)
	return nil
}
