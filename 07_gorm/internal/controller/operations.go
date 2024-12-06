package controller

import (
	tableoutput "github.com/hahaclassic/databases/07_gorm/pkg/table"
	"github.com/jedib0t/go-pretty/table"
)

type Operation int

const (
	operationsStart Operation = iota

	BestExplicitTracks
	CountTracksByGenre
	AlbumsWithMaxTracks
	ArtistsWithReleasedAlbumYear
	UsersOlderThan

	TracksByGenre
	AlbumsWithTrackCounts
	AddUser
	UpdateUserName
	DeleteUser
	AlbumsByArtist

	ExportUsersToJSON
	ImportUsersFromJSON

	operationsEnd
	Exit Operation = 0
)

func (op Operation) String() string {
	switch op {
	case BestExplicitTracks:
		return "BestExplicitTracks"
	case CountTracksByGenre:
		return "CountTracksByGenre"
	case AlbumsWithMaxTracks:
		return "AlbumsWithMaxTracks"
	case ArtistsWithReleasedAlbumYear:
		return "ArtistsWithReleasedAlbumYear"
	case UsersOlderThan:
		return "UsersOlderThan"
	case TracksByGenre:
		return "TracksByGenre"
	case AlbumsWithTrackCounts:
		return "AlbumsWithTrackCounts"
	case AddUser:
		return "AddUser"
	case UpdateUserName:
		return "UpdateUserName"
	case DeleteUser:
		return "DeleteUser"
	case AlbumsByArtist:
		return "AlbumsByArtist"
	case ExportUsersToJSON:
		return "ExportUsersToJSON"
	case ImportUsersFromJSON:
		return "ImportUsersFromJSON"
	case Exit:
		return "Exit"
	default:
		return "UnknownOperation"
	}
}

func menu() {
	headers := []string{"#", "Operation"}
	rows := make([][]any, 0, int(operationsEnd-operationsStart-1))

	for i := operationsStart + 1; i < operationsEnd; i++ {
		rows = append(rows, []any{int(i), i.String()})
	}
	rows = append(rows, []any{int(Exit), Exit.String()})

	tableoutput.PrintTable(table.StyleDefault, headers, rows)
}
