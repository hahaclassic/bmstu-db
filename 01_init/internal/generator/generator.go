package generator

import (
	"context"
	"log/slog"

	"github.com/hahaclassic/databases/01_init/config"
	"github.com/hahaclassic/databases/01_init/internal/service"
	"github.com/hahaclassic/databases/01_init/internal/storage"
	"github.com/hahaclassic/databases/01_init/internal/storage/csv"
	"github.com/hahaclassic/databases/01_init/internal/storage/postgresql"
)

func Run(conf *config.Config) {
	ctx := context.Background()

	var (
		musicStorage storage.MusicServiceStorage
		err          error
	)

	if conf.Generator.OutputCSV == "" {
		musicStorage, err = postgresql.New(ctx, &conf.Postres)
	} else {
		musicStorage, err = csv.New(conf.Generator.OutputCSV)
	}

	if err != nil {
		slog.Error("", "[ERR]", err)

		return
	}

	defer musicStorage.Close()

	musicService := service.New(musicStorage)

	switch {
	case conf.Generator.DeleteCmd:
		err = musicService.DeleteAll(ctx)

	default:
		err = musicService.Generate(ctx, conf.Generator.RecordsPerTable)
	}

	if err != nil {
		slog.Error("", "[ERR]", err)
	}
}
