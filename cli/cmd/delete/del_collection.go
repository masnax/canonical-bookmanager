package delete

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/masnax/rest-api/cli/cmd/rest"
	"github.com/masnax/rest-api/collection"
)

func DelCollection(sourceUrl string, path string, args []string) {
	url := sourceUrl + path
	if _, err := strconv.Atoi(args[0]); err != nil {
		log.Printf("expected numerical id as input")
		return
	}
	url += "/" + args[0]
	_, err := rest.MakeRequest(url, "DELETE", nil)
	if err != nil {
		log.Printf("error from request: %v", err)
	}
}

func RemoveFromCollection(sourceUrl string, path string, bookId string, collectionId string) {
	url := sourceUrl + path + "/"
	bid, err := strconv.Atoi(bookId)
	if err != nil {

		log.Printf("expected integer book id")
		return
	}
	cid, err := strconv.Atoi(collectionId)
	if err != nil {
		log.Printf("expected integer collection id")
		return
	}
	collection := collection.BookCollectionData{BookID: bid, CollectionID: cid}

	bodyBytes, err := json.Marshal(collection)
	if err != nil {
		log.Printf("parsing error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "DELETE", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}
