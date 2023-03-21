package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type SuccessMessage struct {
	Status string   `json:"status"`
	Title  string   `json:"title"`
	Text   []string `json:"text"`
}

type ErrorMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// fetchHtml retrieves the HTML content of a given URL
func fetchHtml(url string) (string, error) {
	// Send HTTP GET request to the given URL
	resp, err := http.Get(url)
	if err != nil {
		// Return an empty string and the error if the request failed
		return "", err
	}
	// Make sure the response body is closed when the function exits
	defer resp.Body.Close()

	// Check if the response status code is not OK (200)
	if resp.StatusCode != 200 {
		// Return an error if the status code is not OK
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body into a byte slice
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Return an empty string and the error if reading the response body failed
		return "", err
	}

	// Convert the byte slice to a string and return it
	return string(bodyBytes), nil
}

func ContentFinder(URI string) (SuccessMessage, error) {
	// Fetch the HTML content from the given URI
	html, err := fetchHtml(URI)
	if err != nil {
		log.Println(err)
		// Return an error message if fetching the web page failed
		return SuccessMessage{}, err
	}

	// Parse the HTML content into a goquery document object
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println(err)
		// Return an error message if parsing the HTML document failed
		return SuccessMessage{}, err
	}

	// Find the page title
	title := doc.Find("title").Text()

	// Find the main content of the article by searching for elements with certain class names
	content := make([]string, 0)
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		classNames, _ := s.Attr("class")
		classArr := strings.Split(classNames, " ")
		for _, class := range classArr {
			if strings.Contains(class, "article") || strings.Contains(class, "post") || strings.Contains(class, "content") {
				// Find all <p> elements inside the matched element and add their text to the content slice
				s.Find("p").Each(func(i int, p *goquery.Selection) {
					text := p.Text()
					if len(strings.Fields(text)) > 20 {
						content = append(content, text)
					}
				})
				// Stop searching for content once a matching element is found
				break
			}
		}
	})

	// Return an error message if no content was found
	if len(content) == 0 {
		return SuccessMessage{}, err
	}

	// Return a success message with the page title and content
	return SuccessMessage{
		Status: "success",
		Title:  title,
		Text:   content,
	}, nil
}

func ConcurrentContentFinder(URIs []string) ([]SuccessMessage, error) {
	// Create a slice to store the results of each web page
	results := make([]SuccessMessage, len(URIs))

	// Create a wait group with the number of web pages to fetch
	var wg sync.WaitGroup
	wg.Add(len(URIs))

	// Launch a goroutine for each web page to fetch its content concurrently
	for i, uri := range URIs {
		go func(i int, uri string) {
			defer wg.Done()
			// Fetch the content of the web page and store the result in the results slice
			res, err := ContentFinder(uri)

			if err != nil {
				log.Println(err)
				// if there is an error, return from the goroutine without assigning the result to the results array
				return
			}

			results[i] = res
		}(i, uri)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Filter out any empty results from the array
	filteredResults := make([]SuccessMessage, 0)
	for _, r := range results {
		if r.Status != "" {
			filteredResults = append(filteredResults, r)
		}
	}

	// Return the filtered results slice
	return filteredResults, nil
}
