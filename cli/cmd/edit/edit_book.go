package edit

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/masnax/rest-api/book"
	"github.com/masnax/rest-api/cli/cmd/rest"
)

func EditBook(sourceUrl string, path string, args []string) {
	url := sourceUrl + path
	edition, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatal(err)
	}
	book := book.Book{
		Title:          args[0],
		Author:         args[1],
		Published_date: args[2],
		Edition:        edition,
		Description:    args[4],
		Genre:          args[5],
	}

	bodyBytes, err := json.Marshal(book)
	if err != nil {
		log.Printf("some error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "POST", reader)
	if err != nil {
		log.Printf("some error: %v", err)
	}
}
