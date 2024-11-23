package controller

type Operation int

const (
	CountArtists Operation = iota + 1
	TracksAdditionalInfo
	BestTracks
	TablesNames
	CountArtistTracks
	ArtistAlbums
	FreePremiumForPensioners
	CurrentUser
	CreateReviewsTable
	AddReview
)

func (op Operation) String() string {
	switch op {
	case CountArtists:
		return "CountArtists"
	case TracksAdditionalInfo:
		return "TracksAdditionalInfo"
	case BestTracks:
		return "BestTracks"
	case TablesNames:
		return "TablesNames"
	case CountArtistTracks:
		return "CountArtistTracks"
	case ArtistAlbums:
		return "ArtistAlbums"
	case FreePremiumForPensioners:
		return "FreePremiumForPensioners"
	case CurrentUser:
		return "CurrentUser"
	case CreateReviewsTable:
		return "CreateReviewsTable"
	case AddReview:
		return "AddReview"
	default:
		return "UnknownOperation"
	}
}
