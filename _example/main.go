package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Webhook struct {
	Progress int `form:"progress" json:"progress" xml:"progress"  binding:"required"`
}

func main() {
	router := gin.Default()

	router.POST("/webhook", func(c *gin.Context) {
		var json Webhook
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"progress": json.Progress})
	})

	router.POST("/webhook2", func(c *gin.Context) {
		var json Webhook
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if v, ok := c.Request.Header["X-Server-Token"]; ok {
			log.Println("show server token: ", v)
		}

		if v, ok := c.Request.Header["X-Data-Uuid"]; ok {
			log.Println("show data uuid: ", v)
		}

		c.JSON(http.StatusOK, gin.H{"progress": json.Progress})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
