package models

type Artist struct {
	Name     string   `json:"name"`
	Genre    string   `json:"genre"`
	Members  []string `json:"members"`
	Websites []string `json:"websites"`
	Albums   []Album  `json:"albums"`
	PageLink string   `json:"page_link"`
}
