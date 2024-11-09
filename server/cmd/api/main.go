package main

import (
	"flag"
	"fmt"
	"github.com/Avixph/receipt-processor-challenge/server/internal/data"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// String containing the application version number.
const version = "1.0.0"

// Config struct holding all the configuration settings for the
// application (network port, current operating environment
// (development, staging, production, etc.)).
type config struct {
	port int
	env  string
}

// Application struct holding the dependencies for the HTTP
// handlers, helpers, and middleware.
type application struct {
	config config
	logger *slog.Logger
	store  data.Stores
}

func main() {
	// Instance of the config struct.
	var cfg config

	// Read values from the port and env command-line flags into
	// the config struct. We default to port '8080' and 'development'
	// environment if no corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Structured logger that writes log entries to the standard out stream.
	lgr := slog.New(slog.NewTextHandler(os.Stdout, nil))

	str := data.NewStores()

	// Instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: lgr,
		store:  str,
	}

	// HTTP server that listens on the port provided in the config struct,
	// uses the serverMux as the handler, timeout settings (idle, read and write)
	// and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(lgr.Handler(), slog.LevelError),
	}

	// Start HTTP serer.
	lgr.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	lgr.Error(err.Error())
	os.Exit(1)
}
