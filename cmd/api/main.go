package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	app := &Application{
		config:   cfg,
		logger:   logger,
		provider: provider.New(cfg),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		//The "" and 0 indicate that the
		// log.Logger instance should not use a prefix or any flags.
		ErrorLog: log.New(logger, "", 0),
	}
	// Start the HTTP server.
	logMsg := fmt.Sprintf("Starting server on port %s", app.config.Port)
	logger.PrintInfo(logMsg, map[string]string{
		"addr": srv.Addr,
		"env":  cfg.Env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
