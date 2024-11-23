package main

import (
	"github.com/hahaclassic/databases/06_cli_app/config"
	"github.com/hahaclassic/databases/06_cli_app/internal/app"
)

func main() {
	app.Start(config.MustLoad())
}
