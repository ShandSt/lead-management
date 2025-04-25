package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/stasshander/lead-management/types"
)

var (
	ErrClientNil      = errors.New("client cannot be nil")
	ErrClientExists   = errors.New("client already exists")
	ErrClientNotFound = errors.New("client not found")
	ErrMaxCapacity    = errors.New("client has reached maximum capacity")
)

type InMemoryStore struct {
	mu      sync.RWMutex
	clients map[string]*types.Client
}

var instance *InMemoryStore
var once sync.Once

func GetInstance() *InMemoryStore {
	once.Do(func() {
		instance = &InMemoryStore{
			clients: make(map[string]*types.Client),
		}
	})
	return instance
}

func (s *InMemoryStore) AddClient(c *types.Client) error {
	if c == nil {
		return ErrClientNil
	}

	if c.ID == "" {
		return errors.New("client ID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[c.ID]; exists {
		return ErrClientExists
	}

	s.clients[c.ID] = c
	return nil
}

func (s *InMemoryStore) GetClient(id string) (*types.Client, error) {
	if id == "" {
		return nil, errors.New("client ID cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[id]
	if !exists {
		return nil, ErrClientNotFound
	}

	// Return a copy to prevent race conditions
	clientCopy := *client
	return &clientCopy, nil
}

func (s *InMemoryStore) GetAllClients() []*types.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	allClients := make([]*types.Client, 0, len(s.clients))
	for _, client := range s.clients {
		// Create a copy to prevent race conditions
		clientCopy := *client
		allClients = append(allClients, &clientCopy)
	}

	return allClients
}

func (s *InMemoryStore) IncrementLeadCount(clientID string) error {
	if clientID == "" {
		return errors.New("client ID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	client, exists := s.clients[clientID]
	if !exists {
		return fmt.Errorf("increment lead count failed: %w", ErrClientNotFound)
	}

	if client.LeadCount >= client.Capacity {
		return ErrMaxCapacity
	}

	client.LeadCount++
	return nil
}
