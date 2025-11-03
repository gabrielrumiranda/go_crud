package main

import (
	"go_crud/config"
	"go_crud/handlers"
	"go_crud/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	defer config.CloseDB()

	productRepo := repository.NewProductRepository(config.DB)
	productHandler := handlers.NewProductHandler(productRepo)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is working",
		})
	})

	api := router.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
