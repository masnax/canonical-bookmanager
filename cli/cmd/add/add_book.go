package add

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/masnax/canonical-bookmanager/book"
	"github.com/masnax/canonical-bookmanager/cli/cmd/rest"
)

func AddBook(sourceUrl string, path string, args []string) {
	url := sourceUrl + path
	edition, err := strconv.Atoi(args[3])
	if err != nil {
		log.Printf("edition must be integer")
		return
	}
	if _, err := time.Parse("2006-01-02", args[2]); err != nil {
		log.Println("published date must be of form Y-M-D")
		return
	}
	book := book.Book{
		Title:       args[0],
		Author:      args[1],
		Published:   args[2],
		Edition:     edition,
		Description: args[4],
		Genre:       args[5],
	}

	bodyBytes, err := json.Marshal(book)
	if err != nil {
		log.Printf("parsing error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "POST", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}
