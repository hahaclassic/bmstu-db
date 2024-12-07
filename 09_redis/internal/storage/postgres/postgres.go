package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/hahaclassic/databases/09_redis/config"
	"github.com/hahaclassic/databases/09_redis/internal/models"
	"github.com/hahaclassic/databases/09_redis/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/rand"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, config *config.PostgresConfig) (*Storage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, net.JoinHostPort(config.Host, config.Port), config.DB, config.SSLMode)

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	if dbpool.Ping(ctx) != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	return &Storage{dbpool}, nil
}

func (s *Storage) Close() {
	s.db.Close()
	slog.Info("CONNECTION CLOSED")
}

func (s *Storage) Top10MostStreamedTracks(ctx context.Context) (tracks []*models.Track, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", storage.ErrTop10MostStreamedTracks, err)
		}
	}()

	query := `
        select t.id, t.name, t.album, 
		t.explicit, t.duration, t.genre, t.stream_count
        from temp_tracks t
        order by t.stream_count desc
        limit 10`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		track := &models.Track{}
		if err := rows.Scan(&track.ID, &track.Name, &track.Album, &track.Explicit,
			&track.Duration, &track.Genre, &track.StreamCount); err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return tracks, nil
}

func (s *Storage) AddTrack(ctx context.Context, track *models.Track) error {
	query := `
        INSERT INTO temp_tracks (id, name, album, explicit, duration, genre, stream_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(ctx, query, track.ID, track.Name, track.Album, track.Explicit,
		track.Duration, track.Genre, track.StreamCount)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrAddTracks, err)
	}

	return nil
}

func (s *Storage) DeleteRandomTrack(ctx context.Context) error {
	query := `
        DELETE FROM temp_tracks
        WHERE id = (
            SELECT id FROM temp_tracks
            ORDER BY RANDOM()
            LIMIT 1
        )`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDeleteRandomTrack, err)
	}

	return nil
}

func (s *Storage) UpdateRandomTrackStreams(ctx context.Context) error {
	query := `
        UPDATE temp_tracks
        SET stream_count = stream_count + $1
        WHERE id = (
            SELECT id FROM temp_tracks
            ORDER BY RANDOM()
            LIMIT 1
        )
    `

	randomIncrement := 1000 + rand.Intn(9000)
	_, err := s.db.Exec(ctx, query, randomIncrement)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrUpdateRandomTrackStreams, err)
	}

	return nil
}
