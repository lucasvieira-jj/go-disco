package pipeline

type Artist struct {
	Name     string `json:"name"`
	PageLink string `json:"page_link"`
}

func Cleaner(artists []Artist) []Artist {
	uniqueArtists := make([]Artist, 0)
	seen := make(map[string]bool)

	for _, artist := range artists {
		key := artist.Name + artist.PageLink
		if !seen[key] {
			uniqueArtists = append(uniqueArtists, artist)
			seen[key] = true
		}
	}

	return uniqueArtists
}
