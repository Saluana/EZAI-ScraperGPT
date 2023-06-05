package main

import (
	"ezai_scraper_api/routers"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
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
		port = "8080"
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"} //Change this when ready for production
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "OAI-KEY", "API-Key"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	r.Use(cors.New(corsConfig))
	routers.NotesRouter(r)
	routers.SummaryRouter(r)

	r.Run(":" + port)
}
