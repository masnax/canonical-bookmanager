package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/masnax/rest-api/book"
	"github.com/masnax/rest-api/parser"
)

type bookHandler struct {
	sync.Mutex
	db *sql.DB
}

func NewBookHandler(db *sql.DB) *bookHandler {
	bh := &bookHandler{
		db: db,
	}
	http.Handle("/books", bh)
	http.Handle("/books/", bh)
	return bh
}

func (bh *bookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer bh.Unlock()
	bh.Lock()
	key, err := parser.URLParser(r.URL)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: '%s'", r.URL.Path))
		return
	}
	switch r.Method {
	case "PUT":
		bh.put(w, r, key)
	case "DELETE":
		bh.delete(w, r, key)
	case "GET":
		bh.get(w, r, key)
	case "POST":
		bh.post(w, r)
	default:
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid method: '%s' for path: '%s'", r.Method, r.URL.Path))
	}
}

func (bh *bookHandler) get(w http.ResponseWriter, r *http.Request, key string) {
	q := "SELECT * from book"
	if len(key) > 0 {
		q += " WHERE id = " + key
	}
	rows, err := bh.db.Query(q)
	defer rows.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to query database due to error: %v", err))
		return
	}
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
}

func (bh *bookHandler) post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Malformed request body: %v", err))
		return
	}
	var book book.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}

	stmt, err := bh.db.Prepare("INSERT INTO book " +
		"(title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(book.Title, book.Author, book.Published_date, book.Edition, book.Description, book.Genre)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (bh *bookHandler) put(w http.ResponseWriter, r *http.Request, key string) {
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
	var book book.Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected non-JSON request: %v", err))
		return
	}

	stmt, err := bh.db.Prepare("UPDATE book SET " +
		"title=?, author=?, published_date=?, edition=?, description=?, genre=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("Unable to update database: %v", err))
		return
	}
	_, err = stmt.Exec(book.Title, book.Author, book.Published_date, book.Edition, book.Description, book.Genre, key)
	if err != nil {
		parser.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	parser.JSONResponse(w, http.StatusOK, nil)
}

func (bh *bookHandler) delete(w http.ResponseWriter, r *http.Request, key string) {
	if len(key) == 0 {
		parser.ErrorResponse(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid path: %s", r.URL.Path))
		return
	}
	stmt, err := bh.db.Prepare("DELETE from book WHERE id=?")
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