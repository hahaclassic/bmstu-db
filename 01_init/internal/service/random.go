package service

import (
	"math/rand/v2"
	"time"
)

const (
	maxAlbumsPerArtist   = 5
	maxTracksPerAlbum    = 20
	maxPlaylistsPerUser  = 5
	maxTracksPerPlaylist = 25
)

var genres = []string{
	"Pop", "Rock", "Hip-hop", "Rap", "Electronic", "Jazz", "Blues", "Classical",
	"Reggae", "Metal", "Country", "Folk", "Soul", "R&B", "Alternative", "Punk",
	"Hardcore", "Ambient", "Funk", "Latin", "Disco", "Dance", "Techno", "Trance",
	"Dubstep", "Indie", "Gothic", "New Age", "Progressive", "Crossover", "Ska",
	"Acoustic", "Lounge", "Psychedelic Rock", "Hard Rock", "Traditional", "Synth-pop",
	"Alternative Hip-hop", "Chamber", "World", "Celtic", "Musical Theatre",
}

func randomGenre() string {
	return genres[rand.IntN(len(genres))]
}

var countries = []string{
	"Russia", "USA", "Canada", "Germany", "France", "United Kingdom", "Italy", "Spain", "China", "Japan",
	"Australia", "India", "Brazil", "Mexico", "South Korea", "Netherlands", "Turkey", "Sweden", "Norway",
	"Finland", "Denmark", "Belgium", "Switzerland", "Austria", "Ireland", "New Zealand", "South Africa",
	"Argentina", "Chile", "Colombia", "Saudi Arabia", "United Arab Emirates", "Singapore", "Malaysia",
	"Thailand", "Philippines", "Vietnam", "Israel", "Greece", "Portugal", "Czech Republic", "Hungary",
	"Romania", "Poland", "Ukraine", "Slovakia", "Slovenia", "Croatia", "Bulgaria", "Serbia",
}

func randomCountry() string {
	return countries[rand.IntN(len(countries))]
}

func randomBool() bool {
	return rand.IntN(2) != 0
}

func randomDates() (time.Time, time.Time, time.Time) {
	min := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)
	seconds := rand.Int64N(max.Unix() - min.Unix())

	birth := min.Add(time.Duration(seconds) * time.Second)

	min = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	max = time.Now()
	seconds = rand.Int64N(max.Unix() - min.Unix())

	registration := min.Add(time.Duration(seconds) * time.Second)

	min = time.Now()
	max = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	seconds = rand.Int64N(max.Unix() - min.Unix())

	premium := min.Add(time.Duration(seconds) * time.Second)

	return birth, registration, premium
}

func randomDateAfter(t time.Time) time.Time {
	seconds := rand.Int64N(time.Now().Unix() - t.Unix())

	return t.Add(time.Duration(seconds) * time.Second)
}
