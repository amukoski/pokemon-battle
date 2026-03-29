package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pokemon-battle/backend/internal/model"
)

type PostgresBattleRepo struct {
	db *sql.DB
}

func NewPostgresBattleRepo(db *sql.DB) *PostgresBattleRepo {
	return &PostgresBattleRepo{db: db}
}

func (r *PostgresBattleRepo) Save(ctx context.Context, b *model.Battle) error {
	logJSON, err := json.Marshal(b.BattleLog)
	if err != nil {
		return fmt.Errorf("marshalling battle log: %w", err)
	}

	query := `
		INSERT INTO battles (pokemon1_id, pokemon2_id, winner, battle_log, pokemon1_score, pokemon2_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		b.Pokemon1.ID, b.Pokemon2.ID, b.Winner, logJSON, b.Pokemon1Score, b.Pokemon2Score,
	).Scan(&b.ID, &b.CreatedAt)
}

func (r *PostgresBattleRepo) GetByID(ctx context.Context, id string) (*model.Battle, error) {
	query := `
		SELECT
			b.id, b.winner, b.battle_log, b.pokemon1_score, b.pokemon2_score, b.created_at,
			p1.id, p1.name, p1.height, p1.weight, p1.types, p1.stats, p1.abilities, p1.sprite_url,
			p2.id, p2.name, p2.height, p2.weight, p2.types, p2.stats, p2.abilities, p2.sprite_url
		FROM battles b
		JOIN pokemons p1 ON b.pokemon1_id = p1.id
		JOIN pokemons p2 ON b.pokemon2_id = p2.id
		WHERE b.id = $1`

	var battle model.Battle
	var logJSON, p1Types, p1Stats, p1Abilities, p2Types, p2Stats, p2Abilities []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&battle.ID, &battle.Winner, &logJSON,
		&battle.Pokemon1Score, &battle.Pokemon2Score, &battle.CreatedAt,
		&battle.Pokemon1.ID, &battle.Pokemon1.Name, &battle.Pokemon1.Height, &battle.Pokemon1.Weight,
		&p1Types, &p1Stats, &p1Abilities, &battle.Pokemon1.SpriteURL,
		&battle.Pokemon2.ID, &battle.Pokemon2.Name, &battle.Pokemon2.Height, &battle.Pokemon2.Weight,
		&p2Types, &p2Stats, &p2Abilities, &battle.Pokemon2.SpriteURL,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("battle not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("querying battle: %w", err)
	}

	if err := unmarshalPokemonJSON(&battle.Pokemon1, p1Types, p1Stats, p1Abilities); err != nil {
		return nil, fmt.Errorf("unmarshalling pokemon1: %w", err)
	}
	if err := unmarshalPokemonJSON(&battle.Pokemon2, p2Types, p2Stats, p2Abilities); err != nil {
		return nil, fmt.Errorf("unmarshalling pokemon2: %w", err)
	}
	if err := json.Unmarshal(logJSON, &battle.BattleLog); err != nil {
		return nil, fmt.Errorf("unmarshalling battle log: %w", err)
	}

	return &battle, nil
}

func (r *PostgresBattleRepo) List(ctx context.Context, limit, offset int) ([]model.Battle, error) {
	query := `
		SELECT
			b.id, b.winner, b.battle_log, b.pokemon1_score, b.pokemon2_score, b.created_at,
			p1.id, p1.name, p1.height, p1.weight, p1.types, p1.stats, p1.abilities, p1.sprite_url,
			p2.id, p2.name, p2.height, p2.weight, p2.types, p2.stats, p2.abilities, p2.sprite_url
		FROM battles b
		JOIN pokemons p1 ON b.pokemon1_id = p1.id
		JOIN pokemons p2 ON b.pokemon2_id = p2.id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying battles: %w", err)
	}
	defer rows.Close()

	var battles []model.Battle
	for rows.Next() {
		var b model.Battle
		var logJSON, p1Types, p1Stats, p1Abilities, p2Types, p2Stats, p2Abilities []byte

		if err := rows.Scan(
			&b.ID, &b.Winner, &logJSON,
			&b.Pokemon1Score, &b.Pokemon2Score, &b.CreatedAt,
			&b.Pokemon1.ID, &b.Pokemon1.Name, &b.Pokemon1.Height, &b.Pokemon1.Weight,
			&p1Types, &p1Stats, &p1Abilities, &b.Pokemon1.SpriteURL,
			&b.Pokemon2.ID, &b.Pokemon2.Name, &b.Pokemon2.Height, &b.Pokemon2.Weight,
			&p2Types, &p2Stats, &p2Abilities, &b.Pokemon2.SpriteURL,
		); err != nil {
			return nil, fmt.Errorf("scanning battle row: %w", err)
		}

		if err := unmarshalPokemonJSON(&b.Pokemon1, p1Types, p1Stats, p1Abilities); err != nil {
			return nil, fmt.Errorf("unmarshalling pokemon1: %w", err)
		}
		if err := unmarshalPokemonJSON(&b.Pokemon2, p2Types, p2Stats, p2Abilities); err != nil {
			return nil, fmt.Errorf("unmarshalling pokemon2: %w", err)
		}
		if err := json.Unmarshal(logJSON, &b.BattleLog); err != nil {
			return nil, fmt.Errorf("unmarshalling battle log: %w", err)
		}

		battles = append(battles, b)
	}

	return battles, rows.Err()
}

func unmarshalPokemonJSON(p *model.Pokemon, types, stats, abilities []byte) error {
	if err := json.Unmarshal(types, &p.Types); err != nil {
		return fmt.Errorf("types: %w", err)
	}
	if err := json.Unmarshal(stats, &p.Stats); err != nil {
		return fmt.Errorf("stats: %w", err)
	}
	if err := json.Unmarshal(abilities, &p.Abilities); err != nil {
		return fmt.Errorf("abilities: %w", err)
	}
	return nil
}
