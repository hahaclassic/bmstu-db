package generator

import (
	"context"
	"log/slog"

	"github.com/hahaclassic/bmstu-db/01_init/config"
	"github.com/hahaclassic/bmstu-db/01_init/internal/service"
)

func Run(conf *config.Config) {
	musicStorage := storage.New(conf.Postres)
	musicService := service.New(musicStorage)

	var err error

	switch {
	case conf.Generator.DeleteCmd:
		err = musicService.DeleteAll(context.Background())

	case conf.Generator.OutputCSV != "":
		err = musicService.GenerateCSV(context.Background(), conf.Generator.OutputCSV, conf.Generator.RecordsPerTable)

	default:
		err = musicService.Generate(context.Background(), conf.Generator.RecordsPerTable)
	}

	if err != nil {
		slog.Error("[ERR]", err)
	}
}
