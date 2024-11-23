package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/hahaclassic/databases/06_cli_app/config"
	"github.com/hahaclassic/databases/06_cli_app/internal/storage"
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
	slog.Info("CONNECTION CLOSED")
}
