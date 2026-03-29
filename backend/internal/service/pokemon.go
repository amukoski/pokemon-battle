package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/pokemon-battle/backend/internal/model"
)

type PokemonFetcher interface {
	FetchPokemon(ctx context.Context, name string) (*model.Pokemon, error)
}

type PokemonCache interface {
	Get(ctx context.Context, name string) (*model.Pokemon, error)
	Set(ctx context.Context, name string, p *model.Pokemon) error
}

type PokemonRepository interface {
	Save(ctx context.Context, p *model.Pokemon) error
	GetByName(ctx context.Context, name string) (*model.Pokemon, error)
	SearchNames(ctx context.Context, prefix string) ([]string, error)
}

type PokemonService struct {
	client PokemonFetcher
	cache  PokemonCache
	repo   PokemonRepository
}

func NewPokemonService(client PokemonFetcher, cache PokemonCache, repo PokemonRepository) *PokemonService {
	return &PokemonService{
		client: client,
		cache:  cache,
		repo:   repo,
	}
}

func (s *PokemonService) GetPokemon(ctx context.Context, name string) (*model.Pokemon, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, fmt.Errorf("pokemon name cannot be empty")
	}

	if cached, err := s.cache.Get(ctx, name); err == nil && cached != nil {
		return cached, nil
	}

	if dbPokemon, err := s.repo.GetByName(ctx, name); err == nil && dbPokemon != nil {
		_ = s.cache.Set(ctx, name, dbPokemon)
		return dbPokemon, nil
	}

	pokemon, err := s.client.FetchPokemon(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, pokemon); err != nil {
		slog.WarnContext(ctx, "failed to save pokemon to db", "name", name, "error", err)
	}

	_ = s.cache.Set(ctx, name, pokemon)

	return pokemon, nil
}

func (s *PokemonService) SearchNames(ctx context.Context, query string) ([]string, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []string{}, nil
	}

	names, err := s.repo.SearchNames(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("searching pokemon names: %w", err)
	}
	if names == nil {
		return []string{}, nil
	}
	return names, nil
}
