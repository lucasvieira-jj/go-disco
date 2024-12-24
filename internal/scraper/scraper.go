package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/lucasvieira-jj/go-disco/config"
	"github.com/lucasvieira-jj/go-disco/internal/utils"
	"github.com/lucasvieira-jj/go-disco/models"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func (c *Scraper) RetrieveArtistList() []models.Artist {
	var artists []models.Artist
	var maxArtistReached bool

	c.collector.OnHTML(".card-artist-name > span", func(e *colly.HTMLElement) {
		if maxArtistReached {
			return
		}
		relativeLink := e.ChildAttr("a", "href")
		baseUrl, _ := url.Parse(config.BaseURL)

		absoluteLink := baseUrl.ResolveReference(&url.URL{Path: relativeLink}).String()

		artist := models.Artist{
			Name:     e.ChildText("a"),
			PageLink: absoluteLink,
		}

		artists = append(artists, artist)
		if len(artists) >= config.MaxArtists {
			maxArtistReached = true
		}

	})

	err := c.collector.Visit(config.GenreUrl)
	if err != nil {
		return nil
	}

	return artists
}

func (c *Scraper) ArtistDetails(artist models.Artist) []models.Artist {
	var artistDetails []models.Artist
	//Getting attributes from page for the artist

	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		membersText := e.ChildText("table.table_1fWaB tr:nth-child(3) td")
		members := utils.ExtractString(membersText)

		//TODO: Valid why the code doest return the href but return the website name
		websiteText := e.ChildAttrs("table.table_1fWaB tr:nth-child(2) td a", "href")

		artist := models.Artist{
			Name:     e.ChildText(".title_1q3xW"),
			Genre:    "Hip-Hop",
			Members:  members,
			Websites: websiteText,
			PageLink: artist.PageLink,
		}

		artistDetails = append(artistDetails, artist)
	})

	artistPage := artist.PageLink
	err := c.collector.Visit(artistPage)
	if err != nil {
		return nil
	}

	return artistDetails
}

func (c *Scraper) ArtistAlbums(artist models.Artist) models.Artist {
	var albums []models.Album

	client := &http.Client{Timeout: 60 * time.Second}

	filteredUrl := artist.PageLink + config.DiscographyFilter
	request, err := http.NewRequest("GET", filteredUrl, nil)
	if err != nil {
		log.Fatal("Error creating request: ", err)
		return models.Artist{}
	}

	request.Header.Add("User-Agent", RandomString(config.UserAgent))
	request.Cookies()
	request.Referer()

	response, err := client.Do(request)
	if err != nil {
		log.Println("Error getting response: ", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Error when closing request: ", response.Status)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Println("Error getting response: ", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Println("Error parsing response: ", err)
	}

	scriptTag := doc.Find("script#dsdata")

	data := scriptTag.Text()
	var dataFormatted map[string]interface{}

	// Unmarshal JSON into the map
	_ = json.Unmarshal([]byte(data), &dataFormatted)

	var maxAlbumReached bool
	// Assuming `dataFormatted` contains a key "data" that points to another object or array
	if dataInterface, exists := dataFormatted["data"]; exists {
		dataMap, ok := dataInterface.(map[string]interface{})
		if !ok {
			log.Println("Unexpected data type for 'data', expected a map[string]interface{}")
		}

		for key, releaseData := range dataMap {
			if maxAlbumReached {
				break
			}

			// Check if the key starts with "Release:" (you can customize this pattern as needed)
			if strings.HasPrefix(key, "Release:") {
				releaseDataMap, ok := releaseData.(map[string]interface{})
				if !ok {
					log.Println("Unexpected data type for release data, expected a map[string]interface{}")
					continue
				}

				album := models.Album{
					Name:     releaseDataMap["title"].(string),
					PageLink: config.BaseURL + strings.TrimPrefix(releaseDataMap["siteUrl"].(string), "/"),
				}

				albums = append(albums, album)

				if len(albums) == config.MaxAlbums {
					log.Println("\nMax number of albums: ", len(albums))
					maxAlbumReached = true
				}

			}
		}
	}

	artist.Albums = append(albums)

	return artist
}

func (c *Scraper) AlbumDetails(artist models.Artist) models.Artist {
	var albumDetails []models.Album

	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		album := models.Album{
			Name: e.ChildText(".title_1q3xW"),
			Year: 1997,
			//Record: ,
			//Styles:   ,
			//Tracks: ,
		}

		albumDetails = append(albumDetails, album)
	})

	for _, album := range artist.Albums {
		albumPage := album.PageLink
		err := c.collector.Visit(albumPage)
		if err != nil {
			return models.Artist{}
		}

	}

	return artist
}
