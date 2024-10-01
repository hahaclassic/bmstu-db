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
