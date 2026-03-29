package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	PokemonAPI  string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://pokemon:battle_secret@localhost:5432/pokemon_battle?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		PokemonAPI:  getEnv("POKEMON_API", "https://pokeapi.co"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
