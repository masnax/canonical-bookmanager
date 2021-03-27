package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/masnax/rest-api/db"
	"github.com/masnax/rest-api/handler"
)

func main() {

	db, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	handler.NewBookHandler(db)
	handler.NewCollectionHandler(db)
	handler.NewBookCollectionHandler(db)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
