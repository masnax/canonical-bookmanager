CREATE TABLE IF NOT EXISTS book (
	id							INTEGER AUTO_INCREMENT PRIMARY KEY,
	title						VARCHAR(255) NOT NULL,
	author					VARCHAR(255) NOT NULL,
	published_date  Date NOT NULL,
	edition					INTEGER NOT NULL,
	description     TEXT,
	genre						VARCHAR(255) NOT NULL
)
