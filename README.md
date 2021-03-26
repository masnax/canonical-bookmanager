# Book Manager

Simple REST API with CLI to manage books.

# Resources

- **net/http** for the API
- **go-mysql-driver/mysql** for connecting to MySQL database
- MySQL database running on a Docker container -- files included

# Functionality

- Add and edit books in the system
  - Books have a **title, author, published date, edition, description, and genre**
- Add and edit book collections
  - Collections have a **name, size, and creation date**

# Database Structure

- ## Table Books

# CLI

## Books

```bash
bm list                   # lists all books
bm add                    # creates a new book with the given information
bm remove                 # deletes an existing book
bm edit                   # opens config file to edit the given book
```

## Collections

```bash
bm collection list        # lists all collections
bm collection add         # adds a new collection
bm collection remove      # deletes an existing collection
bm collection edit        # edits an existing collection
```

# REST API

## Output Structure

```js
{
	"status-code": 200,
	"status": "Ok",
	"data": {}
}
```

## Supported Status Codes

Table goes here

## Filtering

Necessary as per requirements

## Endpoints

Guide
