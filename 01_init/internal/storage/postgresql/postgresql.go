package postgresql

import (
	"context"
	"fmt"
	"net"

	"github.com/hahaclassic/bmstu-db/01_init/config"
	"github.com/hahaclassic/bmstu-db/01_init/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MusicServiceStorage struct {
	pool *pgxpool.Pool
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
