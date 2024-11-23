package models

import "time"

type TracksAdditionalInfo struct {
	Name        string
	Album       string
	Artist      string
	ReleaseDate time.Time
	StreamCount int
}

type RankedTrack struct {
	Name        string
	StreamCount int
	Rank        int
}
