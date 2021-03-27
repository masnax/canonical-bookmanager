package edit

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/masnax/canonical-bookmanager/cli/cmd/rest"
	"github.com/masnax/canonical-bookmanager/collection"
)

func EditCollection(sourceUrl string, path string, argPath string, name string) {
	url := sourceUrl + path + "/" + argPath
	collection := collection.Collection{Collection: name}

	bodyBytes, err := json.Marshal(collection)
	if err != nil {
		log.Printf("parsing error: %v", err)
	}
	reader := bytes.NewReader(bodyBytes)

	_, err = rest.MakeRequest(url, "PUT", reader)
	if err != nil {
		log.Printf("request error: %v", err)
	}
}
