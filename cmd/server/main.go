package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"repo/internal/config"
	"repo/internal/handler"
	"repo/internal/repository"
	"repo/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	setupLogger(cfg.LogLevel)

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	r := repository.NewRepository(db)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/parse", h.PostLog)
	mux.HandleFunc("GET /api/v1/topology/{log_id}", h.GetNodesTopology)
	mux.HandleFunc("GET /api/v1/node/{node_id}", h.GetNodeDetail)
	mux.HandleFunc("GET /api/v1/port/{node_id}", h.GetPorts)
	mux.HandleFunc("GET /api/v1/log/{log_id}", h.GetLogMeta)

	addr := ":" + cfg.Port
	slog.Info("server started", "addr", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func setupLogger(logLevel string) {
	var level slog.Level

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}
