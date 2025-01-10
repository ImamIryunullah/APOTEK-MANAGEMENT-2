package main

import (
	"apotek-management/config"
	"apotek-management/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDB()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
