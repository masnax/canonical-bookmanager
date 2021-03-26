CREATE TABLE IF NOT EXISTS book_collection (
	book_id       INTEGER NOT NULL,
	collection_id INTEGER NOT NULL,
	PRIMARY KEY (book_id, collection_id)
	FOREIGN KEY (book_id) REFERENCES book(id),
	FOREIGN KEY (collection_id) REFERENCES collection(id)
);
