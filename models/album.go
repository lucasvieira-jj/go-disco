package models

type Album struct {
	Name     string   `json:"name"`
	Year     int16    `json:"year"`
	Record   []string `json:"Record"`
	Styles   []string `json:"styles"`
	Tracks   []Track  `json:"tracks"`
	PageLink string   `json:"page_link"`
}
