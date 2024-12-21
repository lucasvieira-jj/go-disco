package models

type Artist struct {
	Name     string   `json:"name"`
	Genre    string   `json:"genre"`
	Members  []string `json:"members"`
	Websites []string `json:"websites"`
	Albums   []Album  `json:"albums"`
	PageLink string   `json:"page_link"`
}

type Album struct {
	Name   string   `json:"name"`
	Year   string   `json:"year"`
	Record []string `json:"Record"`
	Styles []string `json:"styles"`
	Tracks []Track  `json:"tracks"`
}

type Track struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}
