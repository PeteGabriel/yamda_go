package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {

				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")

				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter *rate.Limiter
		//with the last time seen we can run a goroutine to clean up the mapping.
		lastSeen time.Time // last time this client connected
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	//clean up the clients map every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, err)
			return
		}

		mu.Lock()

		// Check to see if the IP address already exists in the map. If it doesn't, then
		// initialize a new rate limiter and add the IP address and limiter to the map.
		if _, found := clients[ip]; !found {
			nl := rate.NewLimiter(rate.Limit(app.config.LimiterRPS), app.config.LimiterBurst)
			clients[ip] = &client{limiter: nl}
		}

		// Update the last seen time for the client.
		clients[ip].lastSeen = time.Now()

		//if the request is not allowed, respond with rate limit error code
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w)
			return
		}

		// Notice that we DON'T use defer to unlock the mutex, as that would mean
		// that the mutex isn't unlocked until all the handlers downstream of this
		// middleware have also returned.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
