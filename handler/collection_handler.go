package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/masnax/rest-api/book"
	"github.com/masnax/rest-api/collection"
	"github.com/masnax/rest-api/parser"
)

type collectionHandler struct {
	sync.Mutex
	db *sql.DB
}

func NewCollectionHandler(db *sql.DB) *collectionHandler {
	ch := &collectionHandler{
		db: db,
	}
	http.Handle("/collections/book", ch)
	http.Handle("/collections/collection", ch)
	http.Handle("/collections/manage", ch)
	http.Handle("/collections/book/", ch)
	http.Handle("/collections/collection/", ch)
	http.Handle("/collections/manage/", ch)
	return ch
}

func (ch *collectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer ch.Unlock()
	ch.Lock()

	keys, err := parser.URLParser(r.URL)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: '%s'", r.URL.Path))
		return
	}
	if err = ch.validateUrl(keys, r.URL); err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	lastKey := keys[len(keys)-1]
	formKey := keys[len(keys)-2]
	switch r.Method {
	case "PUT":
		ch.put(w, r, lastKey, formKey)
	case "DELETE":
		ch.delete(w, r, lastKey, formKey)
	case "GET":
		ch.get(w, r, lastKey, formKey)
	case "POST":
		ch.post(w, r)
	default:
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid method: '%s' for path: '%s'", r.Method, r.URL.Path))
	}
}

func (ch *collectionHandler) validateUrl(keys []string, url *url.URL) error {
	if len(keys) != 4 {
		return errors.New(fmt.Sprintf("Invalid path: '%s'", url.Path))
	}
	lastKey := keys[len(keys)-1]
	formKey := keys[len(keys)-2]
	if formKey == "manage" || formKey == "book" {
		if len(lastKey) > 0 {
			if _, err := strconv.Atoi(lastKey); err != nil {
				return errors.New(fmt.Sprintf("Invalid key: %s from path: %s", lastKey, url.Path))
			}
		}
	} else if formKey == "collection" {
		if len(lastKey) == 0 {
			return errors.New(fmt.Sprintf("Invalid key: %s from path: %s", lastKey, url.Path))
		}
	}
	return nil
}

func (ch *collectionHandler) put(w http.ResponseWriter, r *http.Request, lastKey string, formKey string) {
	if len(lastKey) == 0 || formKey != "manage" {
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
	var collection collection.Collection
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
	_, err = stmt.Exec(collection.Collection, lastKey)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *collectionHandler) delete(w http.ResponseWriter, r *http.Request, lastKey string, formKey string) {
	if len(lastKey) == 0 || formKey != "manage" {
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
	_, err = stmt.Exec(lastKey)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *collectionHandler) get(w http.ResponseWriter, r *http.Request, lastKey string, formKey string) {
	if formKey != "book" && formKey != "collection" && len(lastKey) == 0 {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid keys: %s, %s from path: %s", lastKey, formKey, r.URL.Path))
		return
	}

	common := `bc.book_id = book.id 
	AND 
	collection.id = bc.collection_id
	WHERE `
	var q string
	if formKey == "collection" {
		q = `SELECT book.* from book 
		JOIN (book_collection as bc, collection) ON ` + common + `collection.name = "` + lastKey + `"`
	} else {
		q = `SELECT collection.* from collection 
JOIN (book_collection as bc, book) ON ` + common + `book.id = ` + lastKey
	}

	rows, err := ch.db.Query(q)
	defer rows.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to query database due to error: %v", err))
		return
	}

	if formKey == "collection" {
		books := []book.Book{}
		for rows.Next() {
			var book book.Book
			err := rows.Scan(&book.Id, &book.Title,
				&book.Author, &book.Published_date, &book.Edition, &book.Description, &book.Genre)
			if err != nil {
				parser.ErrorResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("Unable to scan results: %v", err))
				return
			}
			books = append(books, book)
		}
		parser.JSONResponse(w, http.StatusOK, books)
	} else {
		collections := []collection.Collection{}
		for rows.Next() {
			var collection collection.Collection
			err := rows.Scan(&collection.ID, &collection.Collection)
			if err != nil {
				parser.ErrorResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("Unable to scan results: %v", err))
				return
			}
			collections = append(collections, collection)
		}
		parser.JSONResponse(w, http.StatusOK, collections)
	}
}

func (ch *collectionHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Malformed request body: %v", err))
		return
	}
	var collection collection.Collection
	err = json.Unmarshal(body, &collection)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}

	stmt, err := ch.db.Prepare("INSERT INTO collection (name) VALUES (?)")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(collection.Collection)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}
