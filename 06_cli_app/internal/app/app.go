package app

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hahaclassic/databases/06_cli_app/config"
	"github.com/hahaclassic/databases/06_cli_app/internal/controller"
	"github.com/hahaclassic/databases/06_cli_app/internal/storage/postgres"
)

func Start(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.New(ctx, &cfg.Postres)
	if err != nil {
		slog.Error("POSTGRES", "err", err)
	}
	defer db.Close()

	ctrl := controller.NewController(db)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctrl.Start(ctx)
	}()

	wg.Wait()
	fmt.Println("Graceful shutdown.")
}
