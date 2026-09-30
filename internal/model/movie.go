package model

type Movie struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Year     int     `json:"year"`
	Genres   []Genre `json:"genres"`
	Director string  `json:"director"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
