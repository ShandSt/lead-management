package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/lead-management/interfaces"
	"github.com/stasshander/lead-management/internal/api"
	"github.com/stasshander/lead-management/internal/client"
	"github.com/stasshander/lead-management/types"
	"github.com/stretchr/testify/assert"
)

func SetupRouter(mockProvider interfaces.ClientProvider) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	clientService := client.NewService(mockProvider)
	api.RegisterRoutes(router, clientService)
	return router
}

type MockClientProvider struct {
	clients []*types.Client
}

func (mcp *MockClientProvider) AddClient(c *types.Client) error {
	if c.Name == "" {
		return errors.New("client name cannot be empty")
	}
	// Simulate adding the client to the "database"
	mcp.clients = append(mcp.clients, c)
	return nil
}

func (mcp *MockClientProvider) NewClient(c *types.Client) error {
	return mcp.AddClient(c)
}

func (mcp *MockClientProvider) GetAllClients() []*types.Client {
	return mcp.clients
}

func (mcp *MockClientProvider) GetClient(id string) (*types.Client, error) {
	for _, client := range mcp.clients {
		if client.ID == id {
			return client, nil
		}
	}
	return nil, errors.New("client not found")
}

// IncrementLeadCount simulates incrementing the lead count of a client.
func (mcp *MockClientProvider) IncrementLeadCount(clientID string) error {
	for _, client := range mcp.clients {
		if client.ID == clientID {
			client.LeadCount++
			return nil
		}
	}
	return errors.New("No suitable client found")
}

func TestCreateClientHandler(t *testing.T) {
	mockProvider := &MockClientProvider{}
	router := SetupRouter(mockProvider)

	today := time.Now()
	startTime := types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)}
	endTime := types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 20, 0, 0, 0, time.UTC)}

	client := &types.Client{
		ID:   "7",
		Name: "Client RR2343",
		WorkingHours: types.WorkingHours{
			Start: startTime,
			End:   endTime,
		},
		Priority:  5,
		LeadCount: 0,
		Capacity:  10,
	}

	body, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Failed to marshal test client: %v", err)
	}

	req := httptest.NewRequest("POST", "/clients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected HTTP status 201 to be returned")
}

func TestGetAllClientsHandler(t *testing.T) {
	today := time.Now()
	startTime := types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)}
	endTime := types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 20, 0, 0, 0, time.UTC)}

	mockProvider := &MockClientProvider{clients: []*types.Client{
		{
			ID:   "7",
			Name: "Client RR2343",
			WorkingHours: types.WorkingHours{
				Start: startTime,
				End:   endTime,
			},
			Priority:  5,
			LeadCount: 0,
			Capacity:  10,
		},
	}}
	router := SetupRouter(mockProvider)

	req := httptest.NewRequest("GET", "/clients", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 to be returned")
	assert.NotEmpty(t, w.Body.String(), "Expected non-empty response body")
}

// SetupClients creates a mock provider with predefined clients.
func SetupClients() *MockClientProvider {
	today := time.Now()
	return &MockClientProvider{clients: []*types.Client{
		{
			ID:   "6",
			Name: "Client RR",
			WorkingHours: types.WorkingHours{
				Start: types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 9, 0, 0, 0, time.UTC)},
				End:   types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 23, 0, 0, 0, time.UTC)},
			},
			Priority:  15,
			LeadCount: 5,
			Capacity:  8,
		},
		{
			ID:   "7",
			Name: "Client XYZ",
			WorkingHours: types.WorkingHours{
				Start: types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 8, 0, 0, 0, time.UTC)},
				End:   types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 22, 0, 0, 0, time.UTC)},
			},
			Priority:  10,
			LeadCount: 3,
			Capacity:  12,
		},
		{
			ID:   "8",
			Name: "Client ABC",
			WorkingHours: types.WorkingHours{
				Start: types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 7, 0, 0, 0, time.UTC)},
				End:   types.CustomTime{Time: time.Date(today.Year(), today.Month(), today.Day(), 21, 0, 0, 0, time.UTC)},
			},
			Priority:  20,
			LeadCount: 0,
			Capacity:  15,
		},
	}}
}

func TestCreateManyClientHandler(t *testing.T) {
	mockProvider := SetupClients()
	router := SetupRouter(mockProvider)

	for _, client := range mockProvider.clients {
		req := httptest.NewRequest("GET", "/clients/"+client.ID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 to be returned for client ID "+client.ID)
		assert.Contains(t, w.Body.String(), client.Name, "Response body should contain the client data for "+client.Name)
	}
}

func TestAssignLeadHandler(t *testing.T) {
	mockProvider := SetupClients()
	router := SetupRouter(mockProvider)

	req := httptest.NewRequest("GET", "/clients/assignLead", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 to be returned")

	responseBody := w.Body.String()

	expectedID := "\"id\":\"8\""
	expectedName := "\"name\":\"Client ABC\""
	expectedLeadCount := "\"leadCount\":1"

	assert.Contains(t, responseBody, "Lead assigned", "Response body should contain the success message 'Lead assigned'")
	assert.Contains(t, responseBody, expectedID, "Response should include the client's ID")
	assert.Contains(t, responseBody, expectedName, "Response should include the client's name")
	assert.Contains(t, responseBody, expectedLeadCount, "Response should reflect the correct lead count")
}
