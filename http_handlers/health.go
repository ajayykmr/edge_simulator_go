package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	log.Println("Health check endpoint hit")
	c.JSON(200, gin.H{"status": "ok"})
}
