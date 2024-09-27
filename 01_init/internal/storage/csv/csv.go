package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
)

const (
	artistsFileName        = "artists.csv"
	albumsFileName         = "albums.csv"
	tracksFileName         = "tracks.csv"
	playlistsFileName      = "playlists.csv"
	usersFileName          = "users.csv"
	usersPlaylistsFileName = "user_playlists.csv"
	playlistTracksFileName = "playlist_tracks.csv"
	albumTracksFileName    = "album_tracks.csv"
)

type MusicServiceStorage struct {
	artistWriter        *csv.Writer
	albumWriter         *csv.Writer
	trackWriter         *csv.Writer
	playlistWriter      *csv.Writer
	userWriter          *csv.Writer
	userPlaylistWriter  *csv.Writer
	playlistTrackWriter *csv.Writer
	albumTrackWriter    *csv.Writer
}

func New(pathToFolder string) (*MusicServiceStorage, error) {
	artistFile, err := os.Create(filepath.Join(pathToFolder, artistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create artists.csv: %w", err)
	}

	albumFile, err := os.Create(filepath.Join(pathToFolder, albumsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create albums.csv: %w", err)
	}

	trackFile, err := os.Create(filepath.Join(pathToFolder, tracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create tracks.csv: %w", err)
	}

	playlistFile, err := os.Create(filepath.Join(pathToFolder, playlistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create playlists.csv: %w", err)
	}

	userFile, err := os.Create(filepath.Join(pathToFolder, usersFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create users.csv: %w", err)
	}

	userPlaylistFile, err := os.Create(filepath.Join(pathToFolder, usersPlaylistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create user_playlists.csv: %w", err)
	}

	playlistTrackFile, err := os.Create(filepath.Join(pathToFolder, playlistTracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create playlist_tracks.csv: %w", err)
	}

	albumTrackFile, err := os.Create(filepath.Join(pathToFolder, albumTracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create album_tracks.csv: %w", err)
	}

	return &MusicServiceStorage{
		artistWriter:        csv.NewWriter(artistFile),
		albumWriter:         csv.NewWriter(albumFile),
		trackWriter:         csv.NewWriter(trackFile),
		playlistWriter:      csv.NewWriter(playlistFile),
		userWriter:          csv.NewWriter(userFile),
		userPlaylistWriter:  csv.NewWriter(userPlaylistFile),
		playlistTrackWriter: csv.NewWriter(playlistTrackFile),
		albumTrackWriter:    csv.NewWriter(albumTrackFile),
	}, nil
}

func (s *MusicServiceStorage) Close() {
	s.artistWriter.Flush()
	s.albumWriter.Flush()
	s.trackWriter.Flush()
	s.playlistWriter.Flush()
	s.userWriter.Flush()
	s.userPlaylistWriter.Flush()
	s.playlistTrackWriter.Flush()
	s.albumTrackWriter.Flush()
}

func (s *MusicServiceStorage) CreateArtist(ctx context.Context, artist *models.Artist) error {
	record := []string{artist.ID.String(), artist.Name, artist.Genre, artist.Country, fmt.Sprintf("%d", artist.DebutYear)}
	if err := s.artistWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateArtist, err)
	}
	return nil
}

func (s *MusicServiceStorage) CreateAlbum(ctx context.Context, album *models.Album) error {
	record := []string{album.ID.String(), album.Title, album.ReleaseDate.Format("2006-01-02"), album.Label, album.Genre}
	if err := s.albumWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateAlbum, err)
	}
	return nil
}

func (s *MusicServiceStorage) CreateTrack(ctx context.Context, track *models.Track) error {
	record := []string{track.ID.String(), track.Name, fmt.Sprintf("%t", track.Explicit), fmt.Sprintf("%d", track.Duration), track.Genre, fmt.Sprintf("%d", track.StreamCount)}
	if err := s.trackWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateTrack, err)
	}
	return nil
}

func (s *MusicServiceStorage) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	record := []string{playlist.ID.String(), playlist.Title, playlist.Description,
		fmt.Sprintf("%t", playlist.Private), playlist.LastUpdated.Format(time.RFC3339), fmt.Sprintf("%d", playlist.Rating)}
	if err := s.playlistWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreatePlaylist, err)
	}
	return nil
}

func (s *MusicServiceStorage) CreateUser(ctx context.Context, user *models.User) error {
	record := []string{user.ID.String(), user.Name, user.RegistrationDate.Format(time.RFC3339),
		user.BirthDate.Format("2006-01-02"), fmt.Sprintf("%t", user.Premium), user.PremiumExpiration.Format(time.RFC3339)}
	if err := s.userWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateUser, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddPlaylist(ctx context.Context, userPlaylist *models.UserPlaylist) error {
	record := []string{userPlaylist.PlaylistID.String(), userPlaylist.UserID.String(),
		fmt.Sprintf("%t", userPlaylist.IsFavorite), fmt.Sprintf("%d", userPlaylist.AccessLevel)}
	if err := s.userPlaylistWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToPlaylist(ctx context.Context, trackID uuid.UUID, playlistID uuid.UUID) error {
	record := []string{trackID.String(), playlistID.String(), time.Now().Format(time.RFC3339), "1"} // No data about last track in playlist
	if err := s.playlistTrackWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddTrackToPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToAlbum(ctx context.Context, trackID uuid.UUID, albumID uuid.UUID) error {
	record := []string{trackID.String(), albumID.String(), "1"} // No data about last track in album
	if err := s.albumTrackWriter.Write(record); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddTrackToAlbum, err)
	}
	return nil
}

func (s *MusicServiceStorage) DeleteAll(ctx context.Context) error {
	files := []string{
		"artists.csv",
		"albums.csv",
		"tracks.csv",
		"playlists.csv",
		"users.csv",
		"user_playlists.csv",
		"playlist_tracks.csv",
		"album_tracks.csv",
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("%w: %v", storage.ErrDeleteAll, err)
		}
	}

	return nil
}
