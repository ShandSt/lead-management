package api

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/lead-management/internal/client"
	"github.com/stasshander/lead-management/types"
)

func createClientHandler(c *gin.Context, clientService *client.Service) {
	var newClient types.Client
	if err := c.ShouldBindJSON(&newClient); err != nil {
		fmt.Printf("Error binding JSON: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	if err := clientService.AddClient(&newClient); err != nil {
		fmt.Printf("Error adding client: %+v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Client added", "client": newClient})
}

// getAllClientsHandler handles the retrieval of all clients.
func getAllClientsHandler(c *gin.Context, clientService *client.Service) {
	clients := clientService.GetAllClients()
	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

// getClientHandler handles the retrieval of a specific client by ID.
func getClientHandler(c *gin.Context, clientService *client.Service) {
	id := c.Param("id")
	client, err := clientService.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"client": client})
}

// assignLeadHandler assigns a lead to an appropriate client based on priority and availability
func assignLeadHandler(c *gin.Context, clientService *client.Service) {
	clients := clientService.GetAllClients()

	selectedClient := SelectClientForLead(clients)
	if selectedClient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No suitable client found"})
		return
	}

	if err := clientService.IncrementLeadCount(selectedClient.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign lead"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lead assigned", "client": selectedClient})
}

// SelectClientForLead returns the most appropriate client for a lead based on their priority and availability
func SelectClientForLead(clients []*types.Client) *types.Client {
	now := time.Now()
	customNow := types.CustomTime{Time: now}
	availableClients := filterAvailableClients(clients, customNow)

	if len(availableClients) == 0 {
		return nil
	}

	// Sort clients by priority (higher priority first)
	// and then by load factor (leadCount/capacity)
	sort.Slice(availableClients, func(i, j int) bool {
		if availableClients[i].Priority == availableClients[j].Priority {
			// Calculate load factor as float to avoid integer division issues
			loadFactorI := float64(availableClients[i].LeadCount) / float64(availableClients[i].Capacity)
			loadFactorJ := float64(availableClients[j].LeadCount) / float64(availableClients[j].Capacity)
			return loadFactorI < loadFactorJ
		}
		return availableClients[i].Priority > availableClients[j].Priority
	})

	// Return the first client that has available capacity
	for _, client := range availableClients {
		if client.LeadCount < client.Capacity {
			return client
		}
	}
	return nil
}

func filterAvailableClients(clients []*types.Client, now types.CustomTime) []*types.Client {
	var available []*types.Client
	for _, client := range clients {
		if client.Capacity > 0 && client.WorkingHours.Contains(now) {
			available = append(available, client)
		} else {
			fmt.Printf("Client %s not available at %s (capacity: %d)\n",
				client.Name, now.Time.Format("15:04"), client.Capacity)
		}
	}
	return available
}
