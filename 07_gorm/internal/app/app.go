package app

import (
	"context"
	"log/slog"

	"github.com/hahaclassic/databases/07_gorm/config"
	"github.com/hahaclassic/databases/07_gorm/internal/controller"
	"github.com/hahaclassic/databases/07_gorm/internal/storage/postgres"
)

func Start(cfg *config.Config) {
	ctx := context.Background()

	db, err := postgres.New(ctx, &cfg.Postres)
	if err != nil {
		slog.Error("POSTGRES", "err", err)
		return
	}

	controller.NewController(db).Start(ctx)
}
