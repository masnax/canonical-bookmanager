package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	keys, err := parser.URLParser(r.URL)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: '%s'", r.URL.Path))
		return
	}
	if err := ch.validateUrl(keys, r.URL); err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest, err.Error())
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

func (ch *bookCollectionHandler) validateUrl(keys []string, url *url.URL) error {
	if len(keys) > 3 {
		return errors.New(fmt.Sprintf("Invalid path: '%s'", url.Path))
	}
	if len(keys) > 0 {
		key := keys[len(keys)-1]
		if len(key) != 0 {
			return errors.New(fmt.Sprintf("Invalid key: %s from path: %s", key, url.Path))
		}
	}
	return nil
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
	var data collection.BookCollectionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}
	_, err = stmt.Exec(data.BookID, data.CollectionID)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *bookCollectionHandler) get(w http.ResponseWriter, r *http.Request) {
	q := `SELECT collection.id, collection.collection, count(book.id) as size from collection 
LEFT JOIN (book, book_collection) on 
		book.id = book_collection.book_id 
		AND 
		collection.id = book_collection.collection_id 
		GROUP BY collection.id 
		ORDER BY size DESC`
	rows, err := ch.db.Query(q)
	defer rows.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to query database due to error: %v", err))
		return
	}

	bookCollections := []collection.BookCollection{}
	for rows.Next() {
		var bc collection.BookCollection
		err := rows.Scan(&bc.ID, &bc.Collection, &bc.Size)
		if err != nil {
			parser.ErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unable to scan results: %v", err))
			return
		}
		bookCollections = append(bookCollections, bc)
	}
	parser.JSONResponse(w, http.StatusOK, bookCollections)

	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to scan results: %v", err))
		return
	}
}

func (ch *bookCollectionHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Malformed request body: %v", err))
		return
	}
	var bc collection.BookCollectionData
	err = json.Unmarshal(body, &bc)
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
	_, err = stmt.Exec(bc.BookID, bc.CollectionID)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}
