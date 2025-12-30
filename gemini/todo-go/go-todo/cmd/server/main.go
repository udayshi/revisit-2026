package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gemini/go-todo/internal/config"
	httpHandler "github.com/gemini/go-todo/internal/http"
	"github.com/gemini/go-todo/internal/storage/sqlite"
	"github.com/gemini/go-todo/internal/todo"
	"github.com/gemini/go-todo/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("starting server", "addr", cfg.HTTPAddr)

	if err := os.MkdirAll(filepath.Dir(cfg.SQLiteDSN), 0755); err != nil {
		log.Error("failed to create data directory", "error", err)
		os.Exit(1)
	}

	repo, err := sqlite.NewRepo(cfg.SQLiteDSN)
	if err != nil {
		log.Error("failed to create repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	service := todo.NewService(repo)
	handler := httpHandler.NewHandler(service, log)

	r := chi.NewRouter()
	r.Use(httpHandler.Cors(cfg.CORSAllowed))
	r.Use(httpHandler.RequestLogger(log))
	r.Use(httpHandler.PanicRecoverer(log, handler))
	handler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server exited properly")
}
