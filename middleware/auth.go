package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading environment variables: %s\n", err.Error())
			c.AbortWithStatusJSON(500, gin.H{
				"status": "failure",
				"error":  "Internal Server Error",
			})
			return
		}

		apiKey := os.Getenv("MS_KEY")
		key := c.GetHeader("API-Key")
		if key != apiKey {
			c.AbortWithStatusJSON(401, gin.H{
				"status": "failure",
				"error":  "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
