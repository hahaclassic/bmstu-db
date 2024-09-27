package postgresql

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/config"
	"github.com/hahaclassic/databases/01_init/internal/models"
	"github.com/hahaclassic/databases/01_init/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MusicServiceStorage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, config *config.PostgresConfig) (*MusicServiceStorage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, net.JoinHostPort(config.Host, config.Port), config.DB, config.SSLMode)

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	if dbpool.Ping(ctx) != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	return &MusicServiceStorage{dbpool}, nil
}

func (s *MusicServiceStorage) Close() {
	s.db.Close()
}

func (s *MusicServiceStorage) CreateArtist(ctx context.Context, artist *models.Artist) error {
	query := `INSERT INTO artists (id, name, genre, country, debut_year) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, artist.ID, artist.Name, artist.Genre, artist.Country, artist.DebutYear)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateArtist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateAlbum(ctx context.Context, album *models.Album) error {
	query := `INSERT INTO albums (id, title, release_date, label, genre) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, album.ID, album.Title, album.ReleaseDate, album.Label, album.Genre)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateAlbum, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateTrack(ctx context.Context, track *models.Track) error {
	query := `INSERT INTO tracks (id, name, explicit, duration, genre, stream_count) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(ctx, query, track.ID, track.Name, track.Explicit, track.Duration, track.Genre, track.StreamCount)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	query := `INSERT INTO playlists (id, title, description, private, last_updated, rating) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(ctx, query, playlist.ID, playlist.Title, playlist.Description, playlist.Private, playlist.LastUpdated, playlist.Rating)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreatePlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, registration_date, birth_date, premium, premium_expiration) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(ctx, query, user.ID, user.Name, user.RegistrationDate, user.BirthDate, user.Premium, user.PremiumExpiration)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrCreateUser, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddPlaylist(ctx context.Context, userPlaylist *models.UserPlaylist) error {
	query := `INSERT INTO user_playlists (playlist_id, user_id, is_favorite, access_level) VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(ctx, query, userPlaylist.PlaylistID, userPlaylist.UserID, userPlaylist.IsFavorite, userPlaylist.AccessLevel)
	if err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToPlaylist(ctx context.Context, trackID uuid.UUID, playlistID uuid.UUID) error {
	query := `INSERT INTO playlist_tracks (track_id, playlist_id, date_added, track_order) VALUES ($1, $2, NOW(), 
		(SELECT COALESCE(MAX(track_order), 0) + 1 FROM playlist_tracks WHERE playlist_id = $2))`

	if _, err := s.db.Exec(ctx, query, trackID, playlistID); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddTrackToPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToAlbum(ctx context.Context, trackID uuid.UUID, albumID uuid.UUID) error {
	query := `INSERT INTO album_tracks (track_id, album_id, track_order) VALUES ($1, $2, 
		(SELECT COALESCE(MAX(track_order), 0) + 1 FROM album_tracks WHERE album_id = $2))`

	if _, err := s.db.Exec(ctx, query, trackID, albumID); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddTrackToAlbum, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddArtistTrack(ctx context.Context, trackID uuid.UUID, artistID uuid.UUID) error {
	query := `INSERT INTO tracks_by_artists (track_id, album_id) VALUES ($1, $2)`

	if _, err := s.db.Exec(ctx, query, trackID, artistID); err != nil {
		return fmt.Errorf("%w: %v", storage.ErrAddArtistTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) DeleteAll(ctx context.Context) error {
	tables := []string{
		"user_playlists",
		"playlist_tracks",
		"album_tracks",
		"tracks_by_artists",
		"albums_by_artists",
		"tracks",
		"albums",
		"artists",
		"users",
		"playlists",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)

		if _, err := s.db.Exec(ctx, query); err != nil {
			return fmt.Errorf("%w: %v", storage.ErrDeleteAll, err)
		}
	}

	return nil
}
