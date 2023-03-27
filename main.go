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

	// Register auth router
	r.Use(middleware.RateLimit())
	r.Use(middleware.Auth())
	routers.NotesRouter(r)
	routers.SummaryRouter(r)

	r.Run(":" + port)
}
