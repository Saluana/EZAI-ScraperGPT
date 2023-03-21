package utils

import (
	"strings"
)

// CreateContentChunks takes in a string of text and returns a slice of strings,
// where each string in the slice is a chunk of the original text containing up
// to 500 words (except for the last chunk, which can contain up to 1000 words).
func CreateContentChunks(text string) []string {
	// Split the input text into words using the Fields function from the strings package.
	words := strings.Fields(text)

	// Set the chunk size to 500.
	chunkSize := 500

	// Calculate the number of chunks needed based on the total number of words in the input text.
	numChunks := (len(words) + chunkSize - 1) / chunkSize

	// Create an empty slice to hold the chunks.
	chunks := make([]string, 0, numChunks)

	// Iterate over the words, appending each chunk to the chunks slice using the Join function
	// from the strings package to join the words together with spaces.
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))

		// Check if this is the last chunk and if it is larger than 500 words. If so, replace the last
		// chunk with a new chunk containing the remaining words.
		if len(chunks) == numChunks && len(chunks[len(chunks)-1]) > 1000 {
			chunks[len(chunks)-1] = strings.Join(words[i:], " ")
		}
	}

	// Return the chunks slice.
	return chunks
}
