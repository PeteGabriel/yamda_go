package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) serve() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		//The "" and 0 indicate that the
		// log.Logger instance should not use a prefix or any flags.
		ErrorLog: log.New(app.logger, "", 0),
	}

	shutdown := make(chan error)

	go func() {
		//contains os signals to handle graceful shutdown
		quit := make(chan os.Signal, 1) //buffered channel

		//specify which signals to handle
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit //block until signal arrives

		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//we indicate that shutdown is in progress and will finish in 5 seconds.
		//basically we instruct our server to stop accepting any new HTTP requests,
		//and give any in-flight requests a "grace period" of 5 seconds to complete
		//before the application is terminated.
		shutdown <- srv.Shutdown(ctx)

	}()

	// Likewise log a "starting server" message.
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.Env,
	})

	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for the shutdown to complete.
	err = <-shutdown
	if err != nil {
		return err
	}

	// Log the "server stopped" message.
	app.logger.PrintInfo("server stopped", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
