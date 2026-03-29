package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pokemon-battle/backend/internal/model"
)

var (
	ErrPokemonNotFound = errors.New("pokemon not found")
	ErrAPIUnavailable  = errors.New("pokemon API unavailable")
)

type PokeAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokeAPIClient(baseURL string) *PokeAPIClient {
	return &PokeAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *PokeAPIClient) FetchPokemon(ctx context.Context, name string) (*model.Pokemon, error) {
	url := fmt.Sprintf("%s/api/v2/pokemon/%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", ErrPokemonNotFound, name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrAPIUnavailable, resp.StatusCode)
	}

	var raw pokeAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return raw.toDomain(), nil
}

type pokeAPIResponse struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Height    int              `json:"height"`
	Weight    int              `json:"weight"`
	Types     []pokeAPIType    `json:"types"`
	Stats     []pokeAPIStat    `json:"stats"`
	Abilities []pokeAPIAbility `json:"abilities"`
	Sprites   pokeAPISprites   `json:"sprites"`
}

type pokeAPIType struct {
	Type pokeAPINameField `json:"type"`
}

type pokeAPIStat struct {
	BaseStat int              `json:"base_stat"`
	Stat     pokeAPINameField `json:"stat"`
}

type pokeAPIAbility struct {
	Ability pokeAPINameField `json:"ability"`
}

type pokeAPINameField struct {
	Name string `json:"name"`
}

type pokeAPISprites struct {
	FrontDefault string              `json:"front_default"`
	Other        pokeAPISpritesOther `json:"other"`
}

type pokeAPISpritesOther struct {
	OfficialArtwork pokeAPIOfficialArtwork `json:"official-artwork"`
}

type pokeAPIOfficialArtwork struct {
	FrontDefault string `json:"front_default"`
}

func (r *pokeAPIResponse) toDomain() *model.Pokemon {
	p := &model.Pokemon{
		ID:     r.ID,
		Name:   r.Name,
		Height: r.Height,
		Weight: r.Weight,
	}

	for _, t := range r.Types {
		p.Types = append(p.Types, t.Type.Name)
	}

	for _, a := range r.Abilities {
		p.Abilities = append(p.Abilities, a.Ability.Name)
	}

	if r.Sprites.Other.OfficialArtwork.FrontDefault != "" {
		p.SpriteURL = r.Sprites.Other.OfficialArtwork.FrontDefault
	} else {
		p.SpriteURL = r.Sprites.FrontDefault
	}

	for _, s := range r.Stats {
		switch s.Stat.Name {
		case "hp":
			p.Stats.HP = s.BaseStat
		case "attack":
			p.Stats.Attack = s.BaseStat
		case "defense":
			p.Stats.Defense = s.BaseStat
		case "special-attack":
			p.Stats.SpecialAttack = s.BaseStat
		case "special-defense":
			p.Stats.SpecialDefense = s.BaseStat
		case "speed":
			p.Stats.Speed = s.BaseStat
		}
	}

	return p
}
