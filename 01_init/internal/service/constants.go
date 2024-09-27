package service

import "math/rand/v2"

const (
	albumsPerArtist   = 5
	tracksPerAlbum    = 5
	playlistsPerUser  = 5
	tracksPerPlaylist = 10
)

var genres = []string{
	"Pop", "Rock", "Hip-hop", "Rap", "Electronic", "Jazz", "Blues", "Classical",
	"Reggae", "Metal", "Country", "Folk", "Soul", "R&B", "Alternative", "Punk",
	"Hardcore", "Ambient", "Funk", "Latin", "Disco", "Dance", "Techno", "Trance",
	"Dubstep", "Indie", "Gothic", "New Age", "Progressive", "Crossover", "Ska",
	"Acoustic", "Lounge", "Psychedelic Rock", "Hard Rock", "Traditional", "Synth-pop",
	"Alternative Hip-hop", "Chamber", "World", "Celtic", "Musical Theatre",
}

func randGenre() string {
	return genres[rand.IntN(len(genres))]
}
