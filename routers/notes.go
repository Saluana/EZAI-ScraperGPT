package routers

import (
	"ezai_scraper_api/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// create route

func NotesRouter(r *gin.Engine) {
	notesGroup := r.Group("/notes")

	notesGroup.POST("/", func(c *gin.Context) {
		type RequestBody struct {
			URL string `json:"url"`
		}

		var body RequestBody
		err := c.BindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failure",
				"message": "Problem getting website",
			})
			return
		}

		url := body.URL

		content, err := utils.ContentFinder(url)

		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"status":  "failure",
				"message": "Problem getting content",
			})
			return
		}

		text := strings.Join(content.Text, " ")
		chunks := utils.CreateContentChunks(text)
		// get the notes from the url
		notes, err := utils.GetNotes(chunks)

		// if there is an error, return the error
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "failure",
				"message": "Problem getting notes",
			})
			return
		}

		// return the notes
		c.JSON(200, gin.H{
			"status":  "success",
			"notes":   notes,
			"title":   content.Title,
			"url":     url,
			"message": "Successfully scraped notes from the url",
		})
	})

}
