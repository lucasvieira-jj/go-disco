package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lucasvieira-jj/go-disco/config"
	"math/rand"
	"net/http"
)

type Artist struct {
	Name     string `json:"name"`
	PageLink string `json:"page_link"`
}

type Scraper struct {
	collector *colly.Collector
}

func NewScraperCollector() *Scraper {
	c := colly.NewCollector(
		colly.AllowURLRevisit())

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString(config.UserAgent))

		fmt.Println("User-Agent: ", r.Headers.Get("User-Agent"))
	})

	return &Scraper{collector: c}
}

func RandomString(userAgentList []string) string {
	randomIndex := rand.Intn(len(userAgentList))
	return userAgentList[randomIndex]
}

func (c *Scraper) RetrieveArtistList() string {
	var artists []Artist

	c.collector.OnHTML(".card-artist-name > span", func(e *colly.HTMLElement) {

		artist := Artist{
			Name:     e.ChildText("a"),
			PageLink: e.ChildAttr("a", "href"),
		}
		artists = append(artists, artist)
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		if err != nil && r.StatusCode == http.StatusForbidden {
			fmt.Println("Forbidden error: 403")
		} else if err != nil {
			fmt.Println("Error:", err)
		}
	})

	err := c.collector.Visit(config.BaseURL)
	if err != nil {
		return ""
	}

	convertedArtists, err := json.MarshalIndent(artists, "", "  ")
	if err != nil {
		return ""
	}

	return string(convertedArtists)
}
