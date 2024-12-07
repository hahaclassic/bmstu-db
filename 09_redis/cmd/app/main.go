package main

import (
	"github.com/hahaclassic/databases/09_redis/config"
	"github.com/hahaclassic/databases/09_redis/internal/app"
)

func main() {
	app.Run(config.MustLoad())
}
