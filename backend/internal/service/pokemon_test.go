package service

import (
	"context"
	"errors"
	"testing"

	"github.com/pokemon-battle/backend/internal/model"
)

// mockFetcher implements PokemonFetcher for testing.
type mockFetcher struct {
	pokemon *model.Pokemon
	err     error
	calls   int
}

func (m *mockFetcher) FetchPokemon(ctx context.Context, name string) (*model.Pokemon, error) {
	m.calls++
	return m.pokemon, m.err
}

// mockCache implements PokemonCache for testing.
type mockCache struct {
	store map[string]*model.Pokemon
}

func newMockCache() *mockCache {
	return &mockCache{store: make(map[string]*model.Pokemon)}
}

func (m *mockCache) Get(ctx context.Context, name string) (*model.Pokemon, error) {
	p, ok := m.store[name]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockCache) Set(ctx context.Context, name string, p *model.Pokemon) error {
	m.store[name] = p
	return nil
}

// mockPokemonRepo implements PokemonRepository for testing.
type mockPokemonRepo struct {
	store map[string]*model.Pokemon
}

func newMockPokemonRepo() *mockPokemonRepo {
	return &mockPokemonRepo{store: make(map[string]*model.Pokemon)}
}

func (m *mockPokemonRepo) Save(ctx context.Context, p *model.Pokemon) error {
	m.store[p.Name] = p
	return nil
}

func (m *mockPokemonRepo) GetByName(ctx context.Context, name string) (*model.Pokemon, error) {
	p, ok := m.store[name]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPokemonRepo) SearchNames(ctx context.Context, prefix string) ([]string, error) {
	var matches []string
	for name := range m.store {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			matches = append(matches, name)
		}
	}
	return matches, nil
}

func newTestService(fetcher *mockFetcher, cache *mockCache, repo *mockPokemonRepo) *PokemonService {
	return NewPokemonService(fetcher, cache, repo)
}

func TestGetPokemon_CacheMiss_DBMiss_FetchesFromAPI(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	result, err := svc.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", result.Name)
	}
	if fetcher.calls != 1 {
		t.Errorf("expected 1 API call, got %d", fetcher.calls)
	}
}

func TestGetPokemon_CacheHit_SkipsDBAndAPI(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	cache.store["pikachu"] = pikachu
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	result, err := svc.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", result.Name)
	}
	if fetcher.calls != 0 {
		t.Errorf("expected 0 API calls (cache hit), got %d", fetcher.calls)
	}
}

func TestGetPokemon_CacheMiss_DBHit_SkipsAPI(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	repo.store["pikachu"] = pikachu
	svc := newTestService(fetcher, cache, repo)

	result, err := svc.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", result.Name)
	}
	if fetcher.calls != 0 {
		t.Errorf("expected 0 API calls (DB hit), got %d", fetcher.calls)
	}
	// Should also populate cache
	if cache.store["pikachu"] == nil {
		t.Error("expected pokemon to be cached after DB hit")
	}
}

func TestGetPokemon_NameNormalization(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	_, err := svc.GetPokemon(context.Background(), "  PIKACHU  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetcher.calls != 1 {
		t.Errorf("expected 1 API call, got %d", fetcher.calls)
	}
}

func TestGetPokemon_EmptyName_ReturnsError(t *testing.T) {
	fetcher := &mockFetcher{}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	_, err := svc.GetPokemon(context.Background(), "  ")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestGetPokemon_APIError_PropagatesError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("api error")}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	_, err := svc.GetPokemon(context.Background(), "pikachu")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPokemon_APIFetch_SavesToDB(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	_, err := svc.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be saved to DB
	dbPokemon, _ := repo.GetByName(context.Background(), "pikachu")
	if dbPokemon == nil {
		t.Error("expected pokemon to be saved to DB after API fetch")
	}
}

func TestGetPokemon_APIFetch_SavestoCache(t *testing.T) {
	pikachu := &model.Pokemon{ID: 25, Name: "pikachu"}
	fetcher := &mockFetcher{pokemon: pikachu}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	svc := newTestService(fetcher, cache, repo)

	_, err := svc.GetPokemon(context.Background(), "pikachu")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cached, _ := cache.Get(context.Background(), "pikachu")
	if cached == nil {
		t.Error("expected pokemon to be stored in cache after fetch")
	}
}

func TestSearchNames_PrefixMatch_FromDB(t *testing.T) {
	fetcher := &mockFetcher{}
	cache := newMockCache()
	repo := newMockPokemonRepo()
	repo.store["pikachu"] = &model.Pokemon{Name: "pikachu"}
	repo.store["pidgey"] = &model.Pokemon{Name: "pidgey"}
	repo.store["pidgeotto"] = &model.Pokemon{Name: "pidgeotto"}
	repo.store["charizard"] = &model.Pokemon{Name: "charizard"}
	svc := newTestService(fetcher, cache, repo)

	results, err := svc.SearchNames(context.Background(), "pid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(results), results)
	}
}

func TestSearchNames_EmptyQuery(t *testing.T) {
	svc := newTestService(&mockFetcher{}, newMockCache(), newMockPokemonRepo())

	results, err := svc.SearchNames(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearchNames_CaseInsensitive(t *testing.T) {
	repo := newMockPokemonRepo()
	repo.store["pikachu"] = &model.Pokemon{Name: "pikachu"}
	repo.store["charizard"] = &model.Pokemon{Name: "charizard"}
	svc := newTestService(&mockFetcher{}, newMockCache(), repo)

	results, err := svc.SearchNames(context.Background(), "PIK")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0] != "pikachu" {
		t.Errorf("expected [pikachu], got %v", results)
	}
}

func TestSearchNames_OnlyPreviouslyPlayedPokemon(t *testing.T) {
	svc := newTestService(&mockFetcher{}, newMockCache(), newMockPokemonRepo())

	results, err := svc.SearchNames(context.Background(), "pik")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results when no pokemon in DB, got %d", len(results))
	}
}
