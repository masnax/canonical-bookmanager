package collection

type BookCollectionData struct {
	BookID       int `json:"book_id"`
	CollectionID int `json:"collection_id"`
}

type BookCollection struct {
	Collection string `json:"collection"`
	Size       int    `json:"size"`
}
