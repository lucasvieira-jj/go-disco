package config

import (
	"encoding/json"
	"github.com/lucasvieira-jj/go-disco/models"
)

func JsonConverter(list []models.Artist) string {
	convertedArtists, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return ""
	}

	return string(convertedArtists)
}
