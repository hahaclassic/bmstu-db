package main

import (
	"github.com/hahaclassic/databases/07_gorm/config"
	"github.com/hahaclassic/databases/07_gorm/internal/app"
)

func main() {
	app.Start(config.MustLoad())
}
