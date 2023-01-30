package main

import (
	"os"
	"yamda_go/internal/config"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/jsonlog"
)

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	cfg, err := config.New("./debug.env")

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	prov := provider.NewMovieProvider(cfg, logger)

	app := &Application{
		config:   cfg,
		logger:   logger,
		provider: prov,
	}

	if err = app.serve(); err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
