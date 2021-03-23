package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sort"
	"regexp"
	"github.com/gorilla/mux"
)

// Text defines structure of the request
type Text struct {
	Text string `json:"text"`
}

// Result defines structure of result
type Result struct {
	TextLength *TextLength      `json:"textLength"`
	WordCount  int              `json:"wordCount"`
	CharCount  []map[string]int `json:"charCount"`
}

// TextLength defines structure of textlength
type TextLength struct {
	WithSpaces    int `json:"withSpaces"`
	WithoutSpaces int `json:"withoutSpaces"`
}

// Analyze text
func analyzeText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var text Text
	var result Result
	_ = json.NewDecoder(r.Body).Decode(&text)

	result.WordCount = GetWordCount(text.Text)
	result.TextLength = GetTextLength(text.Text)
	result.CharCount = GetCharCount(text.Text)

	json.NewEncoder(w).Encode(result)
}

// GetWordCount return word count
func GetWordCount(text string) int {
	words := strings.Fields(text)
	return len(words)
}

// GetTextLength return text lenght with and without spaces
func GetTextLength(text string) *TextLength {
	wos := strings.ReplaceAll(text, " ", "")
	return &TextLength{WithSpaces: len(text), WithoutSpaces: len(wos)}
}

// GetCharCount return charater count in alphabetical order
func GetCharCount(text string) []map[string]int {

	// Make a Regex to say we only want letters
    reg, err := regexp.Compile("[^a-zA-Z]+")
    if err != nil {
        log.Fatal(err)
    }
	// change to lowercase and remove everything that is not alphabet
    cleanText := reg.ReplaceAllString(strings.ToLower(text), "")
	counter := make(map[string]int)
	for _, c := range cleanText {
			counter[string(c)]++
	}

	keys := make([]string, 0, len(counter))

	for k := range counter {
		keys = append(keys, k)
	}
	
	// sorting the keys
	sort.Strings(keys)

	result := make([]map[string]int, 0)

	for _, k := range keys {
		result = append(result, map[string]int{ k : counter[k] } )
	}

	return result
}

// Main function
func main() {
	// Init router
	r := mux.NewRouter()

	r.HandleFunc("/analyze", analyzeText).Methods("POST")

	// Start server
	log.Printf("Server started at port 8000")
	
	log.Fatal(http.ListenAndServe(":8000", r))
}
