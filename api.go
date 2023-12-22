package main

import (
	// "net/http"
	"github.com/gin-gonic/gin"
)

type Controller struct{}
func (t *Controller) Default(c *gin.Context) {
	c.JSON(200, gin.H{"response": "here it is"})
}

func APIValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("secret")

		if secret != "ooga_booga" {
			c.JSON(401, gin.H{"response": "bro tf up with ur api key 🤣"})
			c.Abort()
			
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.Use(APIValidation())
	{
		controller := new(Controller)
		v1.GET("/default", controller.Default)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"response": "not found"})
	})

	router.Run("localhost:8080")
}
