package add

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/masnax/canonical-bookmanager/cli/cmd/rest"
	"github.com/masnax/canonical-bookmanager/collection"
)

func AddNewCollection(sourceUrl string, path string, args []string) {
	url := sourceUrl + path + "/"
	collection := collection.Collection{Collection: args[0]}

	bodyBytes, err := json.Marshal(collection)
	if err != nil {
		log.Printf("parsing error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "POST", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}

func AddToCollection(sourceUrl string, path string, bookId string, collectionId string) {
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

	_, err = rest.MakeRequest(url, "POST", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}
