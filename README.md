# Book Manager

Simple REST API with CLI to manage books.

# Resources

- `net/http` for the API
- `go-mysql-driver/mysql` for connecting to MySQL database
- `mitchellh/mapstructure` for parsing responses in CLI
- MySQL database running on a Docker container -- files included

# Functionality

- Add and edit books in the system
  - Books have a **title, author, published date, edition, description, and genre**
- Add and edit book collections
  - Collections have a **name, size, and creation date**

# Database Structure
- There are three tables: `book`, `collection`, and `book_collection`.
  - `book` holds information about all books 
  - `collection` holds information pertaining to a collection
  - `book_collection` associates books with collections
- SQL files are present in the `docker-files` directory

# CLI

## Books

```bash
go run cli/main.go list                   # lists all books
go run cli/main.go add                    # creates a new book with the given information
go run cli/main.go delete                 # deletes an existing book
go run cli/main.go edit                   # opens config file to edit the given book
```

## Collections

```bash
go run cli/main.go collection list        # lists all collections
go run cli/main.go collection new         # adds a new collection
go run cli/main.go collection delete      # deletes an existing collection
go run cli/main.go collection edit        # edits an existing collection
go run cli/main.go collection add         # adds a book to an existing collection
go run cli/main.go collection drop        # drops a book from an existing collection
```

## Flags

```bash
--filter "filter args"  # filters books on a given field      -- compatible with 'list', 'collection list --name'
--name collection-name  # shows all books for a collection    -- compatible with 'collection list'
--bid  book-id          # shows all collections for a book id -- compatible with 'collection list'
```

```bash
# 'edit' has its own set of flags to update an existing book:
--title
--author
--published
--edition
--description
--genre
```

# REST API

- `/books`
  - `/books/{id}`
- `/collections`
  - `/collections/manage/`
  - `/collections/manage/{id}`
  - `/collections/book/{id}`
  - `/collections/collection/{name}`

## Details

### `/books`
#### GET
-  returns list of all books
- Data:
```js
[
   {
      "id": 4,
      "title": "Title",
      "author": "FirstName LastName",
      "published": "2005-04-11",
      "edition": 1,
      "description": "Text",
      "genre": "horror"
    }
]
```
#### POST
- adds a new book to the list of all books
- Input:
```js
   {
      "id": 4,
      "title": "Title",
      "author": "FirstName LastName",
      "published": "2005-04-11",
      "edition": 1,
      "description": "Text",
      "genre": "horror"
    }
```

### `/books/{id}`
#### GET
- returns a book with the given id
- Data:
```js
   {
      "id": 4,
      "title": "Title",
      "author": "FirstName LastName",
      "published": "2005-04-11",
      "edition": 1,
      "description": "Text",
      "genre": "horror"
    }
```
#### PUT
- updates all book attributes for given id
- input fields are not mandatory (book information could be unknown), except published date, for formatting
- Input:
```js
   {
      "title": "Title",
      "author": "FirstName LastName",
      "published": "2005-04-11",
      "edition": 1,
      "description": "Text",
      "genre": "horror"
    }
```
#### DELETE
- deletes a book with the given id


### `/collections`
#### GET
- gets list of all collections and their size
- Data:
```js
[   
    {
      "id": 7,
      "collection": "Name",
      "size": 5
    }
]
```
#### POST
- adds a book to an existing collection
- Input:
```js
    {
      "book_id": 3,
      "collection_id": 1
    }
```
#### DELETE
- removes a book from a collection
```js
    {
      "book_id": 3,
      "collection_id": 1
    }
```
### `/collections/manage/`
#### POST
- adds a new collection
- Input:
```js
    {
      "id": 2,
      "collection": "Name"
    }
```
### `/collections/manage/{id}`
#### GET
- gets the collection with the given id
- Data:
```js
    {
      "id": 2,
      "collection": "Name"
    }
```
#### PUT
- updates collection name for the given id
- Input:
```js
    {
      "collection": "Name"
    }
```
#### DELETE
- deletes collection for the given id

### `/collections/book/{id}`
#### GET
- gets all collections that the book with the given id is part of
- Data:
```js
[
    {
      "id": 2,
      "collection": "Name"
    }
]
```

### `/collections/collection/{name}`
#### GET
- gets all books for the collection with the given name
```js
[
   {
      "id": 4,
      "title": "Title",
      "author": "FirstName LastName",
      "published": "2005-04-11",
      "edition": 1,
      "description": "Text",
      "genre": "horror"
    }
]
```


## Output Structure

```js
{
	"status-code": 200,
	"status": "Ok",
	"data": {}
}
```

## Filtering

- Filtering allows for filtering on a specific key for queries that return book results.
  - These endpoints are `/books/`, `/books/{id}` and `/collections/collection/{name}`
  - follows a format of `?filter=KEY+OP+VAL`
    - `KEY` is any field of a book
    - `OP`  is one of `[eq, ne, lt, gt, le, ne]`
    - `VAL` is a series of `+` delimited words representing the value of the field `KEY`
  - Example: `/books?filter=author+eq+max+asna`

