package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	//read the value of the port and env command line flags into the config struc. We default to using the port numbe 40000 and the environment "development" if no corresponding flags are provided
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "developmnent", "Environment (development | staging | production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//Declare an instance of the application struct, containing the config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests to the healthcheckHandler method (which we will create in a moment)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %d", cfg.env, cfg.port)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
