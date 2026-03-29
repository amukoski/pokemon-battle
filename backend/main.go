package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/pokemon-battle/backend/config"
	"github.com/pokemon-battle/backend/internal/cache"
	"github.com/pokemon-battle/backend/internal/client"
	"github.com/pokemon-battle/backend/internal/handler"
	"github.com/pokemon-battle/backend/internal/repository"
	"github.com/pokemon-battle/backend/internal/service"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		slog.ErrorContext(ctx, "failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.InfoContext(ctx, "connected to postgres")

	if err = runMigrations(db); err != nil {
		slog.ErrorContext(ctx, "failed to run migrations", "error", err)
	}
	slog.InfoContext(ctx, "migrations completed")

	redisCache := cache.NewRedisCache(cfg.RedisURL, 24*time.Hour)
	if err = redisCache.Ping(ctx); err != nil {
		slog.WarnContext(ctx, "redis unavailable: running without cache", "error", err)
	} else {
		slog.InfoContext(ctx, "connected to redis")
	}

	pokeClient := client.NewPokeAPIClient(cfg.PokemonAPI)
	pokemonRepo := repository.NewPostgresPokemonRepo(db)
	pokemonSvc := service.NewPokemonService(pokeClient, redisCache, pokemonRepo)
	battleRepo := repository.NewPostgresBattleRepo(db)
	battleSvc := service.NewBattleService(pokemonSvc, battleRepo)

	mux := http.NewServeMux()
	h := handler.New(battleSvc, pokemonSvc)
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	slog.InfoContext(ctx, "server starting on port", "port", cfg.Port)
	if err = server.ListenAndServe(); err != nil {
		slog.ErrorContext(ctx, "server error", "error", err)
	}

	slog.InfoContext(ctx, "server stopped")
}

func runMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("migrations/001_create_tables.sql")
	if err != nil {
		return fmt.Errorf("reading migration file: %w", err)
	}

	_, err = db.Exec(string(migration))
	return err
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
