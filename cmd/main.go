package main

import (
	"fmt"
	"github.com/lucasvieira-jj/go-disco/internal/scraper"
	"time"
)

var maxRetries int = 10
var waitTimeExecution = 5 * time.Second

func main() {
	fmt.Println("Starting the scraper")
	scraperService := scraper.NewScraperCollector()

	for i := 0; i < maxRetries; i++ {

		artists := scraperService.RetrieveArtistList()

		if len(artists) > 0 {
			fmt.Printf("Artists Retrieved Successfully: %v\n", artists)
			break
		}

		if i < maxRetries-1 {
			time.Sleep(waitTimeExecution)
			waitTimeExecution *= 2
		} else {
			fmt.Println("Max retries reached")
		}
	}

	//Clean duplicates

}
