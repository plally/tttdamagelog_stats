package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/plally/damagelog_parser/internal/charts"
	"github.com/plally/damagelog_parser/internal/config"
	"github.com/plally/damagelog_parser/internal/dal"
	"github.com/plally/damagelog_parser/internal/damagelog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	cfg := config.Get()
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := dal.New(pool)

	m, err := migrate.New("file://db/migrations", cfg.PostgresURL)
	if err != nil {
		panic(err)
	}

	slog.Info("running migrations")
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			slog.With("error", err).Info("No change")
		} else {
			panic(err)
		}
	}

	authGroup := r.Group(nil)
	authGroup.Use(middleware.BasicAuth("gmod", cfg.BasicAuth))

	authGroup.Post("/intake/damagelog/round", func(w http.ResponseWriter, r *http.Request) {
		err := damagelog.ProcessData(r.Context(), r.Body, queries)
		if err != nil {
			slog.With("error", err).Error("failed to process data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	r.Get("/stats/charts", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting charts")
		steamid64 := r.URL.Query().Get("id")
		first := r.URL.Query().Get("first")

		if steamid64 == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if first == "true" {
			chart, err := charts.GetRandomChart(r.Context(), steamid64, queries)
			if err != nil {
				slog.With("error", err).Error("failed to get random chart")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]charts.Chart{chart})
			return
		}

		charts, err := charts.GetDefaultCharts(r.Context(), steamid64, queries)
		if err != nil {
			slog.With("error", err).Error("failed to get charts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(charts)
	})

	slog.With("port", 8080).Info("starting server")
	http.ListenAndServe(":8080", r)
}

type Data struct {
	Damagelog json.RawMessage `json:"damagelog,omitempty"`
	Date      string          `json:"date,omitempty"`
	Day       string          `json:"day,omitempty"`
	ID        string          `json:"id,omitempty"`
	Map       string          `json:"map,omitempty"`
	Month     string          `json:"month,omitempty"`
	Round     string          `json:"round,omitempty"`
	Year      string          `json:"year,omitempty"`
}

func IDToType(i int) string {
	switch i {
	case 0:
		return "Damage"
	default:
		return "Unknown"
	}
}
