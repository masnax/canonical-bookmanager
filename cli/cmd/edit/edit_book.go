package edit

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/masnax/rest-api/book"
	"github.com/masnax/rest-api/cli/cmd/rest"
)

func EditBook(sourceUrl string, path string, argPath string, title string,
	author string, date string, edition int, description string, genre string) {
	url := sourceUrl + path + "/" + argPath

	if _, err := time.Parse("2006-01-02", date); err != nil {
		log.Println("published date must be of form Y-M-D")
		return
	}
	book := book.Book{
		Title:       title,
		Author:      author,
		Published:   date,
		Edition:     edition,
		Description: description,
		Genre:       genre,
	}

	bodyBytes, err := json.Marshal(book)
	if err != nil {
		log.Printf("parsing error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "PUT", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}
