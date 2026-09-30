package model

import (
	"time"
)

type Movie struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
	Year  int    `db:"year"`
	// Genres   []Genre   `db:"genres"`
	Director  string    `db:"director"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
