package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lucasvieira-jj/go-disco/config"
	"github.com/lucasvieira-jj/go-disco/models"
	"math/rand"
	"net/http"
	"net/url"
)

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
	var artists []models.Artist

	c.collector.OnHTML(".card-artist-name > span", func(e *colly.HTMLElement) {
		relativeLink := e.ChildAttr("a", "href")
		baseUrl, _ := url.Parse(config.BaseURL) // A base URL

		absoluteLink := baseUrl.ResolveReference(&url.URL{Path: relativeLink}).String()

		artist := models.Artist{
			Name:     e.ChildText("a"),
			PageLink: absoluteLink,
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

	err := c.collector.Visit(config.GenreUrl)
	if err != nil {
		return ""
	}

	convertedArtists, err := json.MarshalIndent(artists, "", "  ")
	if err != nil {
		return ""
	}

	return string(convertedArtists)
}
