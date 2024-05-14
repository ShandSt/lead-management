package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stasshander/lead-management/internal/api"
	"github.com/stasshander/lead-management/internal/client"
	"github.com/stasshander/lead-management/pkg/store"
)

func main() {
	r := gin.Default()
	store := store.GetInstance()
	clientService := client.NewService(store)
	api.RegisterRoutes(r, clientService)
	r.Run(":8001")
}
