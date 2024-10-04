package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	playlistTracksFileName = "playlist_tracks.csv"
	usersPlaylistsFileName = "user_playlists.csv"
	artistTracksFileName   = "tracks_by_artists.csv"
	artistAlbumFileName    = "albums_by_artists.csv"
)

type MusicServiceStorage struct {
	mu           *sync.Mutex
	writers      map[string]*csv.Writer
	files        map[*os.File]struct{}
	pathToFolder string
}

func New(pathToFolder string) (m *MusicServiceStorage, err error) {
	files := make(map[*os.File]struct{})

	defer func() {
		if err != nil {
			for file := range files {
				file.Close()
			}
		}
	}()

	artistFile, err := os.Create(filepath.Join(pathToFolder, artistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", artistsFileName, err)
	}
	files[artistFile] = struct{}{}

	albumFile, err := os.Create(filepath.Join(pathToFolder, albumsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", albumsFileName, err)
	}
	files[albumFile] = struct{}{}

	trackFile, err := os.Create(filepath.Join(pathToFolder, tracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", tracksFileName, err)
	}
	files[trackFile] = struct{}{}

	playlistFile, err := os.Create(filepath.Join(pathToFolder, playlistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", playlistsFileName, err)
	}
	files[playlistFile] = struct{}{}

	userFile, err := os.Create(filepath.Join(pathToFolder, usersFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", usersFileName, err)
	}
	files[userFile] = struct{}{}

	userPlaylistFile, err := os.Create(filepath.Join(pathToFolder, usersPlaylistsFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", usersPlaylistsFileName, err)
	}
	files[userPlaylistFile] = struct{}{}

	playlistTrackFile, err := os.Create(filepath.Join(pathToFolder, playlistTracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", playlistTracksFileName, err)
	}
	files[playlistTrackFile] = struct{}{}

	artistTrackFile, err := os.Create(filepath.Join(pathToFolder, artistTracksFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", artistTracksFileName, err)
	}
	files[artistTrackFile] = struct{}{}

	artistAlbumFile, err := os.Create(filepath.Join(pathToFolder, artistAlbumFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", artistAlbumFileName, err)
	}
	files[artistAlbumFile] = struct{}{}

	storage := &MusicServiceStorage{
		pathToFolder: pathToFolder,
		mu:           &sync.Mutex{},
		writers: map[string]*csv.Writer{
			artistsFileName:        csv.NewWriter(artistFile),
			albumsFileName:         csv.NewWriter(albumFile),
			tracksFileName:         csv.NewWriter(trackFile),
			playlistsFileName:      csv.NewWriter(playlistFile),
			usersFileName:          csv.NewWriter(userFile),
			usersPlaylistsFileName: csv.NewWriter(userPlaylistFile),
			playlistTracksFileName: csv.NewWriter(playlistTrackFile),
			artistTracksFileName:   csv.NewWriter(artistTrackFile),
			artistAlbumFileName:    csv.NewWriter(artistAlbumFile),
		},
		files: files,
	}

	if err = storage.writeHeaders(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *MusicServiceStorage) writeHeaders() error {
	headers := map[string][]string{
		artistsFileName:        {"id", "name", "genre", "country", "debut_year"},
		albumsFileName:         {"id", "title", "release_date", "label", "genre"},
		tracksFileName:         {"id", "name", "order_in_album", "album_id", "explicit", "duration", "genre", "stream_count"},
		playlistsFileName:      {"id", "title", "description", "private", "last_updated", "rating"},
		usersFileName:          {"id", "name", "registration_date", "birth_date", "premium", "premium_expiration"},
		usersPlaylistsFileName: {"playlist_id", "user_id", "is_favorite", "access_level"},
		playlistTracksFileName: {"track_id", "playlist_id", "date_added", "track_order"},
		artistTracksFileName:   {"track_id", "artist_id"},
		artistAlbumFileName:    {"album_id", "artist_id"},
	}

	for filename, header := range headers {
		if err := s.writers[filename].Write(header); err != nil {
			return err
		}
	}

	return nil
}

func (s *MusicServiceStorage) Close() {
	for _, w := range s.writers {
		w.Flush()
	}

	for file := range s.files {
		err := file.Close()
		if err != nil {
			slog.Error("error while closing file", "error", err)
		}
	}
}

func (s *MusicServiceStorage) CreateArtist(ctx context.Context, artist *models.Artist) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{artist.ID.String(), artist.Name, artist.Genre,
		artist.Country, fmt.Sprintf("%d", artist.DebutYear)}

	if err := s.writers[artistsFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateArtist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateAlbum(ctx context.Context, album *models.Album) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{album.ID.String(), album.Title, album.ReleaseDate.Format("2006-01-02"), album.Label, album.Genre}

	if err := s.writers[albumsFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateAlbum, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateTrack(ctx context.Context, track *models.Track) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{track.ID.String(), track.Name, fmt.Sprintf("%d", track.OrderInAlbum),
		track.AlbumID.String(), fmt.Sprintf("%t", track.Explicit), fmt.Sprintf("%d", track.Duration),
		track.Genre, fmt.Sprintf("%d", track.StreamCount)}

	if err := s.writers[tracksFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{playlist.ID.String(), playlist.Title, playlist.Description,
		fmt.Sprintf("%t", playlist.Private), playlist.LastUpdated.Format(time.RFC3339),
		fmt.Sprintf("%d", playlist.Rating)}

	if err := s.writers[playlistsFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreatePlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateUser(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{user.ID.String(), user.Name, user.RegistrationDate.Format(time.RFC3339),
		user.BirthDate.Format("2006-01-02"), fmt.Sprintf("%t", user.Premium),
		user.PremiumExpiration.Format(time.RFC3339)}

	if err := s.writers[usersFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateUser, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddPlaylist(ctx context.Context, userPlaylist *models.UserPlaylist) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{userPlaylist.ID.String(), userPlaylist.UserID.String(),
		fmt.Sprintf("%t", userPlaylist.IsFavorite), fmt.Sprintf("%d", userPlaylist.AccessLevel)}

	if err := s.writers[usersPlaylistsFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToPlaylist(ctx context.Context, track *models.PlaylistTrack) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{track.ID.String(), track.PlaylistID.String(),
		time.Now().Format(time.RFC3339), strconv.Itoa(track.TrackOrder)}

	if err := s.writers[playlistTracksFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddTrackToPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddArtistTrack(ctx context.Context, trackID uuid.UUID, artistID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{trackID.String(), artistID.String()}

	if err := s.writers[artistTracksFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddArtistTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddArtistAlbum(ctx context.Context, albumID uuid.UUID, artistID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := []string{albumID.String(), artistID.String()}

	if err := s.writers[artistAlbumFileName].Write(record); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddArtistAlbum, err)
	}

	return nil
}

func (s *MusicServiceStorage) DeleteAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range s.writers {
		w.Flush()
	}

	var err error

	for file := range s.files {
		err = file.Close()
		if err != nil {
			break
		}
	}

	return err
}
