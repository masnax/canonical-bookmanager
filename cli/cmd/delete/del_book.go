package delete

import (
	"log"
	"strconv"

	"github.com/masnax/rest-api/cli/cmd/rest"
)

func DelBook(sourceUrl string, path string, args []string) {
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
