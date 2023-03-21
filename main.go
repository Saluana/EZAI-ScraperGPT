package main

import (
	"ezai-scraper-api/routers"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error loading environment variables: %s\n", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8484"
	}

	r := gin.Default()

	// Register auth router
	routers.NotesRouter(r)
	routers.SummaryRouter(r)

	r.Run(":" + port)
}
