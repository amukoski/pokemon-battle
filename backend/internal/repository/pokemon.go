package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pokemon-battle/backend/internal/model"
)

type PostgresPokemonRepo struct {
	db *sql.DB
}

func NewPostgresPokemonRepo(db *sql.DB) *PostgresPokemonRepo {
	return &PostgresPokemonRepo{db: db}
}

func (r *PostgresPokemonRepo) Save(ctx context.Context, p *model.Pokemon) error {
	typesJSON, err := json.Marshal(p.Types)
	if err != nil {
		return fmt.Errorf("marshalling types: %w", err)
	}
	statsJSON, err := json.Marshal(p.Stats)
	if err != nil {
		return fmt.Errorf("marshalling stats: %w", err)
	}
	abilitiesJSON, err := json.Marshal(p.Abilities)
	if err != nil {
		return fmt.Errorf("marshalling abilities: %w", err)
	}

	query := `
		INSERT INTO pokemons (id, name, height, weight, types, stats, abilities, sprite_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			types = EXCLUDED.types,
			stats = EXCLUDED.stats,
			abilities = EXCLUDED.abilities,
			sprite_url = EXCLUDED.sprite_url`

	_, err = r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.Height, p.Weight,
		typesJSON, statsJSON, abilitiesJSON, p.SpriteURL,
	)
	if err != nil {
		return fmt.Errorf("saving pokemon: %w", err)
	}
	return nil
}

func (r *PostgresPokemonRepo) GetByName(ctx context.Context, name string) (*model.Pokemon, error) {
	query := `
		SELECT id, name, height, weight, types, stats, abilities, sprite_url
		FROM pokemons WHERE name = $1`

	var p model.Pokemon
	var typesJSON, statsJSON, abilitiesJSON []byte

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&p.ID, &p.Name, &p.Height, &p.Weight,
		&typesJSON, &statsJSON, &abilitiesJSON, &p.SpriteURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying pokemon: %w", err)
	}

	if err := json.Unmarshal(typesJSON, &p.Types); err != nil {
		return nil, fmt.Errorf("unmarshalling types: %w", err)
	}
	if err := json.Unmarshal(statsJSON, &p.Stats); err != nil {
		return nil, fmt.Errorf("unmarshalling stats: %w", err)
	}
	if err := json.Unmarshal(abilitiesJSON, &p.Abilities); err != nil {
		return nil, fmt.Errorf("unmarshalling abilities: %w", err)
	}

	return &p, nil
}

func (r *PostgresPokemonRepo) SearchNames(ctx context.Context, prefix string) ([]string, error) {
	query := `
		SELECT name FROM pokemons
		WHERE name LIKE $1
		ORDER BY name
		LIMIT 10`

	rows, err := r.db.QueryContext(ctx, query, prefix+"%")
	if err != nil {
		return nil, fmt.Errorf("searching pokemon names: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning pokemon name: %w", err)
		}
		names = append(names, name)
	}
	return names, rows.Err()
}
