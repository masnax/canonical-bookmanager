package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/masnax/rest-api/collection"
	"github.com/masnax/rest-api/parser"
)

type bookCollectionHandler struct {
	sync.Mutex
	db *sql.DB
}

func NewBookCollectionHandler(db *sql.DB) *bookCollectionHandler {
	ch := &bookCollectionHandler{
		db: db,
	}
	http.Handle("/collections", ch)
	http.Handle("/collections/", ch)
	return ch
}

func (ch *bookCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer ch.Unlock()
	ch.Lock()

	_, err := parser.URLParser(r.URL)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: '%s'", r.URL.Path))
		return
	}
	switch r.Method {
	case "DELETE":
		ch.delete(w, r)
	case "GET":
		ch.get(w, r)
	case "POST":
		ch.post(w, r)
	default:
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid method: '%s' for path: '%s'", r.Method, r.URL.Path))
	}
}

func (ch *bookCollectionHandler) delete(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	stmt, err := ch.db.Prepare("DELETE from book_collection WHERE book_id=? AND collection_id=?")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	var bookCollection collection.BookCollection
	err = json.Unmarshal(body, &bookCollection)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}
	_, err = stmt.Exec(bookCollection.BookID, bookCollection.CollectionID)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *bookCollectionHandler) get(w http.ResponseWriter, r *http.Request) {
	rows, err := ch.db.Query("SELECT * from book_collection")
	defer rows.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to query database due to error: %v", err))
		return
	}
	bookCollections := []collection.BookCollection{}
	for rows.Next() {
		var bookCollection collection.BookCollection
		err := rows.Scan(&bookCollection.BookID, &bookCollection.CollectionID)
		if err != nil {
			parser.ErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unable to scan results: %v", err))
			return
		}
		bookCollections = append(bookCollections, bookCollection)
	}
	parser.JSONResponse(w, http.StatusOK, bookCollections)
}

func (ch *bookCollectionHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Malformed request body: %v", err))
		return
	}
	var bookCollection collection.BookCollection
	err = json.Unmarshal(body, &bookCollection)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}

	stmt, err := ch.db.Prepare("INSERT INTO book_collection (book_id, collection_id) VALUES (?, ?)")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(bookCollection.BookID, bookCollection.CollectionID)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}
