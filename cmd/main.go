package main

import (
	"fmt"
	"github.com/lucasvieira-jj/go-disco/internal/scraper"
	"github.com/lucasvieira-jj/go-disco/models"
	"time"
)

var maxRetries int = 1000
var waitTimeExecution = time.Second

func main() {
	fmt.Println("Starting the scraper")
	scraperService := scraper.NewScraperCollector()

	var artistList []models.Artist
	var artistDetails []models.Artist

	for i := 0; i < maxRetries; i++ {

		artistList = scraperService.RetrieveArtistList()
		if len(artistList) > 0 {
			fmt.Println("Artists Retrieved Successfully")
			break
		}

		if i < maxRetries-1 {
			time.Sleep(waitTimeExecution)
			waitTimeExecution *= 2
		} else {
			fmt.Println("Max retries reached")
		}
	}
	fmt.Println("Retrieved artists")

	for _, artist := range artistList {
		for i := 0; i < maxRetries; i++ {
			artistDetails = scraperService.ArtistDetails(artist)
			if artistDetails != nil {
				fmt.Printf("Artists details Retrieved Successfully: %v", artistDetails)
				break
			}

			if i < maxRetries-1 {
				time.Sleep(waitTimeExecution)
				waitTimeExecution *= 2
			} else {
				fmt.Println("Max retries reached")
			}
		}
	}

	fmt.Println("Artists Detailers Retrieved")

	for _, artist := range artistList {
		for i := 0; i < maxRetries; i++ {
			albumList := scraperService.ArtistAlbums(artist)
			if albumList.Albums != nil {
				fmt.Printf("Album details Retrieved Successfully: %v", albumList)
				break
			}

			if i < maxRetries-1 {
				time.Sleep(waitTimeExecution)
				waitTimeExecution *= 2
			} else {
				fmt.Println("Max retries reached")
			}
		}
	}

	fmt.Println("Artists Albums Retrieved")

	for _, artist := range artistList {
		for i := 0; i < maxRetries; i++ {
			artistDetailed := scraperService.AlbumDetails(artist)
			fmt.Printf("Album details Retrieved Successfully: %v", artistDetailed)

			if i < maxRetries-1 {
				time.Sleep(waitTimeExecution)
				waitTimeExecution *= 2
			} else {
				fmt.Println("Max retries reached")
			}
		}
	}

	fmt.Println("Artists All details Retrieved")

}
