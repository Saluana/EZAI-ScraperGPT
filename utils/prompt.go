package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func GetNote(chunk string) (string, error) {

	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error loading environment variables: %s\n", err.Error())
		return "", err
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "I need you to create a bullet point list of short notes from key topics in the text the user has provided. Use '->' as the bullet point for each note. Please remove and ignore any unwanted text, such as things related to website cookies, website newletters, and website advertisements. Use the example text, and example notes as a guide.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Penguins are flightless birds that live in cold climates, primarily in Antarctica. They have adapted to their environment by developing thick feathers and a layer of fat, called blubber, to keep warm in the freezing water.",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "->Penguins are flightless birds. ->Live in cold climates, primarily in Antarctica. ->Have adapted to their environment with thick feathers and blubber. ->Blubber helps them stay warm in freezing water.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: chunk,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GetNotes(splitContent []string) (interface{}, error) {
	// Create a slice to store the results of each note
	notes := make([]string, len(splitContent))

	// Create a wait group with the number of chunks to process
	var wg sync.WaitGroup
	wg.Add(len(splitContent))

	// Launch a goroutine for each chunk to get its note concurrently
	for i, content := range splitContent {
		go func(i int, content string) {
			defer wg.Done()
			// Get the note for the content and store the result in the notes slice
			note, err := GetNote(content)

			if err != nil {
				log.Println(err)
				// if there is an error, return from the goroutine without assigning the result to the notes array
				return
			}

			notes[i] = note
		}(i, content)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Filter out any empty notes from the array
	filteredNotes := make([]string, 0)
	for _, n := range notes {
		if n != "" {
			filteredNotes = append(filteredNotes, n)
		}
	}

	// Parse the notes and remove any unwanted text
	parsedNotes := strings.Join(filteredNotes, "\n")
	parsedNotes = strings.ReplaceAll(parsedNotes, "->", "\n")
	parsedNotes = strings.ReplaceAll(parsedNotes, "\n\n", "\n")
	parsedNotes = strings.TrimSpace(parsedNotes)

	// Split the notes into individual strings and remove any leading/trailing whitespace
	completeNotes := make([]string, 0)
	for _, note := range strings.Split(parsedNotes, "\n") {
		note = strings.TrimSpace(note)
		if note != "" {
			completeNotes = append(completeNotes, note)
		}
	}

	// Return an error message if no notes were found
	if len(completeNotes) == 0 {
		return ErrorMessage{
			Status:  "failure",
			Message: "No notes found.",
		}, nil
	}

	// Return a success message with the complete notes
	return completeNotes, nil
}

func GetSummaryChunk(chunk string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading environment variables: %s\n", err.Error())
		return "", err
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Write a very short and concise summary of the text provided by the user. The summary should be no longer then a Tweet. Remove any unwanted text, such as things related to website cookies, website newletters, and website advertisements.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: chunk + "\n Short summary:",
				},
			},
			Temperature: 0.85,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GetSummary(splitContent []string) (string, error) {
	// Create a channel to receive results from the goroutines
	results := make(chan string)

	// Create a wait group with the number of chunks to summarize
	var wg sync.WaitGroup
	wg.Add(len(splitContent))

	// Launch a goroutine for each chunk to summarize it concurrently
	for _, chunk := range splitContent {
		go func(chunk string) {
			defer wg.Done()
			// Summarize the chunk and send the result to the results channel
			summary, err := GetSummaryChunk(chunk)
			if err != nil {
				log.Println(err)
				return
			}
			results <- summary
		}(chunk)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect the results from the channel
	summaryList := make([]string, 0)
	for summary := range results {
		if summary != "" {
			summaryList = append(summaryList, summary)
		}
	}

	// Concatenate the summaries into a single string and return it
	rawSummary := strings.Join(summaryList, "\n")
	completeSummary, err := GetSummaryChunk(rawSummary)

	if (err != nil) || (completeSummary == "") {
		completeSummary = rawSummary
	}

	if completeSummary == "" {
		return "", errors.New("Could not summarize article.")
	}
	return completeSummary, nil
}
