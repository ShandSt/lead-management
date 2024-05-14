package client

import (
	"github.com/stasshander/lead-management/interfaces"
	"github.com/stasshander/lead-management/types"
)

type Service struct {
	store interfaces.ClientProvider
}

func NewService(store interfaces.ClientProvider) *Service {
	return &Service{store: store}
}

func (s *Service) AddClient(c *types.Client) error {
	return s.store.AddClient(c)
}

func (s *Service) GetClient(id string) (*types.Client, error) {
	return s.store.GetClient(id)
}

func (s *Service) GetAllClients() []*types.Client {
	return s.store.GetAllClients()
}

func (s *Service) IncrementLeadCount(clientID string) error {
	return s.store.IncrementLeadCount(clientID)
}
