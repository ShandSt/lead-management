package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stasshander/lead-management/internal/client"
)

func RegisterRoutes(router *gin.Engine, clientService *client.Service) {
	router.POST("/clients", func(c *gin.Context) {
		createClientHandler(c, clientService)
	})
	router.GET("/clients", func(c *gin.Context) {
		getAllClientsHandler(c, clientService)
	})
	router.GET("/clients/assignLead", func(c *gin.Context) {
		assignLeadHandler(c, clientService)
	})
	router.GET("/clients/:id", func(c *gin.Context) {
		getClientHandler(c, clientService)
	})
}
