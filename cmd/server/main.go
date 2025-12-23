package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"go-connect-todo/internal/handler"
	"go-connect-todo/internal/repository"
	"go-connect-todo/internal/server"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		slog.Error(fmt.Sprintf("Error processing envconfig: %s", err.Error()))
		os.Exit(1)
	}

	// connect DB
	db, err := connectDB(c.DatabaseURL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	txnRepo := repository.NewTransactionRepo(db)
	txnHandler := handler.NewTransaction(txnRepo)

	r := repository.New()
	h := handler.NewTodo(r)

	s := server.New(h, txnHandler)

	slog.Info("Server listening on :8080")
	if err := s.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func connectDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return db, nil
}
