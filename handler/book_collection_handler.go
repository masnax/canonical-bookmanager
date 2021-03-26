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

	key, err := parser.URLParser(r.URL)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: '%s'", r.URL.Path))
		return
	}
	switch r.Method {
	case "PUT":
		ch.put(w, r, key)
	case "DELETE":
		ch.delete(w, r, key)
	case "GET":
		ch.get(w, r, key)
	case "POST":
		ch.post(w, r)
	default:
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid method: '%s' for path: '%s'", r.Method, r.URL.Path))
	}
}

func (ch *bookCollectionHandler) put(w http.ResponseWriter, r *http.Request, key string) {
	if len(key) == 0 {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: %s", r.URL.Path))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Malformed request body: %v", err))
		return
	}
	var bookCollection book_collection.BookCollection
	err = json.Unmarshal(body, &collection)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}

	stmt, err := ch.db.Prepare("UPDATE collection SET title=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(collection.Name, key)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *bookCollectionHandler) delete(w http.ResponseWriter, r *http.Request, key string) {
	if len(key) == 0 {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: %s", r.URL.Path))
		return
	}
	stmt, err := ch.db.Prepare("DELETE from collection WHERE id=?")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(key)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *bookCollectionHandler) get(w http.ResponseWriter, r *http.Request, key string) {
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
		err := rows.Scan(&bookCollection.ID, &bookCollection.BookID, &bookCollection.CollectionID)
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
