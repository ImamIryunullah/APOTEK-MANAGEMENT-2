package routes

import (
    "apotek-management/controllers"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    api := router.Group("/api")
    {
        api.GET("/stok", controllers.GetAllStok)
        api.POST("/stok", controllers.CreateStok)
        api.PUT("/stok/:id", controllers.UpdateStok)
        api.DELETE("/stok/:id", controllers.DeleteStok)
        api.GET("/stok/search", controllers.SearchStok)
        api.GET("/stok/filter", controllers.FilterStok)
        api.GET("/stok/summary", controllers.GetStokSummary)
    }
}
