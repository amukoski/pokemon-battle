package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPokemon_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/pokemon/pikachu" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 25,
			"name": "pikachu",
			"height": 4,
			"weight": 60,
			"types": [{"type": {"name": "electric"}}],
			"stats": [
				{"base_stat": 35, "stat": {"name": "hp"}},
				{"base_stat": 55, "stat": {"name": "attack"}},
				{"base_stat": 40, "stat": {"name": "defense"}},
				{"base_stat": 50, "stat": {"name": "special-attack"}},
				{"base_stat": 50, "stat": {"name": "special-defense"}},
				{"base_stat": 90, "stat": {"name": "speed"}}
			],
			"abilities": [{"ability": {"name": "static"}}],
			"sprites": {
				"front_default": "https://example.com/pikachu.png",
				"other": {"official-artwork": {"front_default": "https://example.com/pikachu-art.png"}}
			}
		}`))
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	pokemon, err := client.FetchPokemon(context.Background(), "pikachu")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pokemon.Name != "pikachu" {
		t.Errorf("expected name pikachu, got %s", pokemon.Name)
	}
	if pokemon.ID != 25 {
		t.Errorf("expected ID 25, got %d", pokemon.ID)
	}
	if len(pokemon.Types) != 1 || pokemon.Types[0] != "electric" {
		t.Errorf("expected types [electric], got %v", pokemon.Types)
	}
	if pokemon.Stats.HP != 35 {
		t.Errorf("expected HP 35, got %d", pokemon.Stats.HP)
	}
	if pokemon.Stats.Attack != 55 {
		t.Errorf("expected Attack 55, got %d", pokemon.Stats.Attack)
	}
	if pokemon.Stats.Speed != 90 {
		t.Errorf("expected Speed 90, got %d", pokemon.Stats.Speed)
	}
	if pokemon.SpriteURL != "https://example.com/pikachu-art.png" {
		t.Errorf("expected official artwork URL, got %s", pokemon.SpriteURL)
	}
	if len(pokemon.Abilities) != 1 || pokemon.Abilities[0] != "static" {
		t.Errorf("expected abilities [static], got %v", pokemon.Abilities)
	}
}

func TestFetchPokemon_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not Found`))
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	_, err := client.FetchPokemon(context.Background(), "notapokemon")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrPokemonNotFound) {
		t.Errorf("expected ErrPokemonNotFound, got %v", err)
	}
}

func TestFetchPokemon_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	_, err := client.FetchPokemon(context.Background(), "pikachu")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrAPIUnavailable) {
		t.Errorf("expected ErrAPIUnavailable, got %v", err)
	}
}

func TestFetchPokemon_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	_, err := client.FetchPokemon(context.Background(), "pikachu")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchPokemon_FallbackSprite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1, "name": "bulbasaur", "height": 7, "weight": 69,
			"types": [], "stats": [], "abilities": [],
			"sprites": {
				"front_default": "https://example.com/fallback.png",
				"other": {"official-artwork": {"front_default": ""}}
			}
		}`))
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	pokemon, err := client.FetchPokemon(context.Background(), "bulbasaur")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pokemon.SpriteURL != "https://example.com/fallback.png" {
		t.Errorf("expected fallback sprite, got %s", pokemon.SpriteURL)
	}
}
