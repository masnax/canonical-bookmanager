package book

type Book struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	Published_date string `json:"published_date"`
	Edition        int    `json:"edition"`
	Description    string `json:"description"`
	Genre          string `json:"genre"`
}
