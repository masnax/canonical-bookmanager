package book

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Published   string `json:"published"`
	Edition     int    `json:"edition"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
}
