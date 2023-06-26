package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Webhook struct {
	Progress int `form:"progress" json:"progress" xml:"progress"  binding:"required"`
}

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/webhook", func(c *gin.Context) {
		var json Webhook
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"progress": json.Progress})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
