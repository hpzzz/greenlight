package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// shutdownError receives any errors returned by graceful Shutdown() func
	shutdownError := make(chan error)

	go func() {
		// Quit is a channel which carres os.Signal values
		quit := make(chan os.Signal, 1)

		// Listen for SIGINT and SIGTERM signals and relay them to quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from quit channel. This code blocks untill a signal is received
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		// Wait untill background goroutines have finished, then return nil to indicate that shutdown completed w/o issues
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Shutdown() immedietaly causes listenAndServe to get ErrServerClosed
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	// Part below happens if listenAndServe receives ErrServerClosed from Shutdown()

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
