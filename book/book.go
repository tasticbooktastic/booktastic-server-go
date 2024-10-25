package book

import "booktastic-server-go/author"

type Book struct {
	ID      uint64          `json:"id"`
	Title   string          `json:"title"`
	Isbn13  string          `json:"isbn13"`
	Authors []author.Author `json:"authors"`
}
