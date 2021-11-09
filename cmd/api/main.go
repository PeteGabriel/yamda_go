package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"yamda_go/internal/config"
	"yamda_go/internal/data/provider"
	"yamda_go/internal/services"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	cfg, err := config.New("./debug.env")
	if err != nil {
		logger.Fatal(err)
	}
	app := &Application{
		config: cfg,
		log: logger,
		movieSvc: services.New(provider.New(cfg)),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", app.config.Host, app.config.Port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// Start the HTTP server.
	logger.Printf("starting %s server on %s", app.config.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}