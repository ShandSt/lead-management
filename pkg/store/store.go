package store

import (
	"errors"
	"sync"

	"github.com/stasshander/lead-management/types"
)

type InMemoryStore struct {
	sync.Mutex
	clients map[string]*types.Client
}

var instance *InMemoryStore
var once sync.Once

func GetInstance() *InMemoryStore {
	once.Do(func() {
		instance = &InMemoryStore{clients: make(map[string]*types.Client)}
	})
	return instance
}

func (s *InMemoryStore) AddClient(c *types.Client) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		return errors.New("client cannot be nil")
	}
	if _, exists := s.clients[c.ID]; exists {
		return errors.New("client already exists")
	}

	s.clients[c.ID] = c
	return nil
}

func (s *InMemoryStore) GetClient(id string) (*types.Client, error) {
	s.Lock()
	defer s.Unlock()
	client, exists := s.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (s *InMemoryStore) GetAllClients() []*types.Client {
	s.Lock()
	defer s.Unlock()

	allClients := make([]*types.Client, 0, len(s.clients))
	for _, client := range s.clients {
		allClients = append(allClients, client)
	}

	return allClients
}

func (s *InMemoryStore) IncrementLeadCount(clientID string) error {
	s.Lock()
	defer s.Unlock()

	client, exists := s.clients[clientID]
	if !exists {
		return errors.New("client not found")
	}

	if client.LeadCount >= client.Capacity {
		return errors.New("client has reached maximum capacity")
	}

	client.LeadCount++
	return nil
}
