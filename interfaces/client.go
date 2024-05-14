// interfaces/client.go
package interfaces

import (
	"github.com/stasshander/lead-management/types"
)

type ClientProvider interface {
	AddClient(c *types.Client) error
	GetClient(id string) (*types.Client, error)
	GetAllClients() []*types.Client
	IncrementLeadCount(clientID string) error
}
