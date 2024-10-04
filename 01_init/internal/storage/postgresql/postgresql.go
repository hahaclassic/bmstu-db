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
		return fmt.Errorf("%w: %w", storage.ErrCreateArtist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateAlbum(ctx context.Context, album *models.Album) error {
	query := `INSERT INTO albums (id, title, release_date, label, genre) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, album.ID, album.Title, album.ReleaseDate, album.Label, album.Genre)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateAlbum, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateTrack(ctx context.Context, track *models.Track) error {
	query := `INSERT INTO tracks (id, name, order_in_album, album_id, explicit, duration, genre, stream_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.db.Exec(ctx, query, track.ID, track.Name,
		track.OrderInAlbum, track.AlbumID, track.Explicit,
		track.Duration, track.Genre, track.StreamCount)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	query := `INSERT INTO playlists (id, title, description, private, rating) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, playlist.ID, playlist.Title, playlist.Description, playlist.Private, playlist.Rating)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreatePlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, registration_date, birth_date, premium, premium_expiration) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(ctx, query, user.ID, user.Name, user.RegistrationDate, user.BirthDate, user.Premium, user.PremiumExpiration)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateUser, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddPlaylist(ctx context.Context, userPlaylist *models.UserPlaylist) error {
	query := `INSERT INTO user_playlists (playlist_id, user_id, is_favorite, access_level) VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(ctx, query, userPlaylist.ID, userPlaylist.UserID, userPlaylist.IsFavorite, userPlaylist.AccessLevel)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddTrackToPlaylist(ctx context.Context, track *models.PlaylistTrack) error {
	query := `INSERT INTO playlist_tracks (track_id, playlist_id, track_order) VALUES ($1, $2, `

	var err error

	if track.TrackOrder == -1 {
		query += `(SELECT COALESCE(MAX(track_order), 0) + 1 FROM playlist_tracks WHERE playlist_id = $2))`

		_, err = s.db.Exec(ctx, query, track.ID, track.PlaylistID)
	} else {
		query += `$3)`

		_, err = s.db.Exec(ctx, query, track.ID, track.PlaylistID, track.TrackOrder)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddTrackToPlaylist, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddArtistTrack(ctx context.Context, trackID uuid.UUID, artistID uuid.UUID) error {
	query := `INSERT INTO tracks_by_artists (track_id, artist_id) VALUES ($1, $2)`

	if _, err := s.db.Exec(ctx, query, trackID, artistID); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddArtistTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) AddArtistAlbum(ctx context.Context, albumID uuid.UUID, artistID uuid.UUID) error {
	query := `INSERT INTO albums_by_artists (album_id, artist_id) VALUES ($1, $2)`

	if _, err := s.db.Exec(ctx, query, albumID, artistID); err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddArtistTrack, err)
	}

	return nil
}

func (s *MusicServiceStorage) DeleteAll(ctx context.Context) error {
	tables := []string{
		"artists",
		"albums",
		"tracks",
		"playlists",
		"users",
		"playlist_tracks",
		"user_playlists",
		"tracks_by_artists",
		"albums_by_artists",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)

		if _, err := s.db.Exec(ctx, query); err != nil {
			return fmt.Errorf("%w: %v", storage.ErrDeleteAll, err)
		}
	}

	return nil
}
