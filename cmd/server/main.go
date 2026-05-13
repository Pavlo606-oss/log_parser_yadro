package main

import (
	"database/sql"
	"log"
	"net/http"
	"repo/internal/config"
	"repo/internal/handler"
	"repo/internal/repository"
	"repo/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const BaseURL = ":8080"

func main() {
	c := config.Load()

	db, err := sql.Open("pgx", c.DSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := repository.NewRepository(db)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/parse", h.PostLog)
	mux.HandleFunc("GET /api/v1/topology/{log_id}", h.GetNodesTopology)
	mux.HandleFunc("GET /api/v1/node/{node_id}", h.GetNodeDetail)
	mux.HandleFunc("GET /api/v1/port/{node_id}", h.GetPorts)
	mux.HandleFunc("GET /api/v1/log/{log_id}", h.GetLogMeta)

	log.Fatal(http.ListenAndServe(BaseURL, mux))
}
