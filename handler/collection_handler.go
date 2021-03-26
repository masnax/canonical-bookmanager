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

type collectionHandler struct {
	sync.Mutex
	db *sql.DB
}

func NewCollectionHandler(db *sql.DB) *collectionHandler {
	ch := &collectionHandler{
		db: db,
	}
	http.Handle("/collections/manage", ch)
	http.Handle("/collections/manage/", ch)
	return ch
}

func (ch *collectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (ch *collectionHandler) put(w http.ResponseWriter, r *http.Request, key string) {
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
	_, err = stmt.Exec(collection.Name, key)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (ch *collectionHandler) delete(w http.ResponseWriter, r *http.Request, key string) {
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

func (ch *collectionHandler) get(w http.ResponseWriter, r *http.Request, key string) {
	q := "SELECT * from collection"
	if len(key) > 0 {
		q += " WHERE id = " + key
	}
	rows, err := ch.db.Query(q)
	defer rows.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to query database due to error: %v", err))
		return
	}
	collections := []collection.Collection{}
	for rows.Next() {
		var collection collection.Collection
		err := rows.Scan(&collection.ID, &collection.Name)
		if err != nil {
			parser.ErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unable to scan results: %v", err))
			return
		}
		collections = append(collections, collection)
	}
	parser.JSONResponse(w, http.StatusOK, collections)
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
	_, err = stmt.Exec(collection.Name)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}
