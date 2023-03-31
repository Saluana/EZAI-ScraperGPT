package main

import (
	"ezai_scraper_api/middleware"
	"ezai_scraper_api/routers"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENVIROMENT")

	if env != "production" {
		err := godotenv.Load()

		if err != nil {
			fmt.Printf("Error loading environment variables: %s\n", err.Error())
			return
		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8484"
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	/*
		config := cors.DefaultConfig()
		config.AllowHeaders = []string{"OAI-KEY", "Content-Type", "API-Key"}
		config.AllowAllOrigins = true
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowCredentials = true

		// Register auth router
		r.Use(cors.New(config))
	*/
	r.Use(CORSMiddleware())
	r.Use(middleware.RateLimit())
	r.Use(middleware.Auth())
	routers.NotesRouter(r)
	routers.SummaryRouter(r)

	r.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, OAI-KEY, API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
