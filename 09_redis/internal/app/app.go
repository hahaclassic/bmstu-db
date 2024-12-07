package app

import (
	"context"
	"log/slog"

	"github.com/hahaclassic/databases/09_redis/config"
	cacheredis "github.com/hahaclassic/databases/09_redis/internal/cache/redis"
	"github.com/hahaclassic/databases/09_redis/internal/controller"
	"github.com/hahaclassic/databases/09_redis/internal/service"
	"github.com/hahaclassic/databases/09_redis/internal/storage/postgres"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	db, err := postgres.New(ctx, &cfg.Postres)
	if err != nil {
		slog.Error("POSTGRES", "err", err)
		return
	}

	defer db.Close()

	cacheRedis, err := cacheredis.New(ctx, &cfg.Redis)
	if err != nil {
		slog.Error("REDIS", "err", err)
		return
	}

	s := service.New(cacheRedis, db)

	controller.New(s).Run(ctx)
}
