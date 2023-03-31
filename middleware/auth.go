package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		env := os.Getenv("ENVIROMENT")

		if env != "production" {
			err := godotenv.Load()
			if err != nil {
				fmt.Printf("Error loading environment variables: %s\n", err.Error())
				c.AbortWithStatusJSON(500, gin.H{
					"status": "failure",
					"error":  "Internal Server Error",
				})
				return
			}
		}

		apiKey := os.Getenv("MS_KEY")
		key := c.GetHeader("API-Key")
		openAiKey := c.GetHeader("OAI-KEY")

		if key != apiKey && len(openAiKey) > 0 {
			c.Next()
		}

		if key != apiKey && len(openAiKey) == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"status": "failure",
				"error":  "Unauthorized",
			})
			return
		}

		oaiKey := os.Getenv("OPENAI_API_KEY")
		c.Request.Header.Set("OAI-KEY", oaiKey)

		c.Next()
	}
}
