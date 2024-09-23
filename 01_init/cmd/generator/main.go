package main

import (
	"flag"

	"github.com/hahaclassic/databases/01_init/config"
	"github.com/hahaclassic/databases/01_init/internal/generator"
)

func main() {
	conf, generatorConf := &config.Config{}, &config.GeneratorConfig{}

	flag.BoolVar(&generatorConf.DeleteCmd, "d", false, "Deletes all records in all tables")
	flag.IntVar(&generatorConf.RecordsPerTable, "c", 1000, "'-c N': generates N records for each table")
	flag.StringVar(&generatorConf.OutputCSV, "csv", "", "'-csv /path/to/your/file': output result to CSV")
	flag.Parse()

	// The database config is not needed to generate csv
	if generatorConf.OutputCSV == "" {
		conf = config.MustLoad()
	}

	conf.Generator = *generatorConf

	generator.Run(conf)
}
