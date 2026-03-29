package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/pokemon-battle/backend/internal/model"
)

type BattleRepository interface {
	Save(ctx context.Context, b *model.Battle) error
	GetByID(ctx context.Context, id string) (*model.Battle, error)
	List(ctx context.Context, limit, offset int) ([]model.Battle, error)
}

type BattleService struct {
	pokemonSvc *PokemonService
	repo       BattleRepository
}

func NewBattleService(pokemonSvc *PokemonService, repo BattleRepository) *BattleService {
	return &BattleService{
		pokemonSvc: pokemonSvc,
		repo:       repo,
	}
}

func (s *BattleService) ExecuteBattle(ctx context.Context, name1, name2 string) (*model.Battle, error) {
	name1 = strings.ToLower(strings.TrimSpace(name1))
	name2 = strings.ToLower(strings.TrimSpace(name2))

	if name1 == "" || name2 == "" {
		return nil, fmt.Errorf("both pokemon names are required")
	}
	if name1 == name2 {
		return nil, fmt.Errorf("cannot battle a pokemon against itself")
	}

	var p1, p2 *model.Pokemon
	var err1, err2 error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		p1, err1 = s.pokemonSvc.GetPokemon(ctx, name1)
	}()
	go func() {
		defer wg.Done()
		p2, err2 = s.pokemonSvc.GetPokemon(ctx, name2)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, fmt.Errorf("fetching %s: %w", name1, err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("fetching %s: %w", name2, err2)
	}

	battle := CalculateBattle(p1, p2)

	if err := s.repo.Save(ctx, battle); err != nil {
		return nil, fmt.Errorf("saving battle: %w", err)
	}

	return battle, nil
}

func (s *BattleService) GetBattle(ctx context.Context, id string) (*model.Battle, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BattleService) ListBattles(ctx context.Context, limit, offset int) ([]model.Battle, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(ctx, limit, offset)
}

func CalculateBattle(p1, p2 *model.Pokemon) *model.Battle {
	type statCategory struct {
		name   string
		p1Val  int
		p2Val  int
		weight float64
	}

	categories := []statCategory{
		{"HP", p1.Stats.HP, p2.Stats.HP, 0.15},
		{"Attack", p1.Stats.Attack, p2.Stats.Attack, 0.20},
		{"Defense", p1.Stats.Defense, p2.Stats.Defense, 0.15},
		{"Sp. Attack", p1.Stats.SpecialAttack, p2.Stats.SpecialAttack, 0.20},
		{"Sp. Defense", p1.Stats.SpecialDefense, p2.Stats.SpecialDefense, 0.15},
		{"Speed", p1.Stats.Speed, p2.Stats.Speed, 0.15},
	}

	var score1, score2 float64
	var battleLog []model.BattleRound

	for _, cat := range categories {
		total := float64(cat.p1Val + cat.p2Val)
		var s1, s2 float64
		if total > 0 {
			s1 = (float64(cat.p1Val) / total) * cat.weight
			s2 = (float64(cat.p2Val) / total) * cat.weight
		}
		score1 += s1
		score2 += s2

		advantage := "tie"
		if cat.p1Val > cat.p2Val {
			advantage = p1.Name
		} else if cat.p2Val > cat.p1Val {
			advantage = p2.Name
		}

		battleLog = append(battleLog, model.BattleRound{
			Category:    cat.name,
			Pokemon1Val: cat.p1Val,
			Pokemon2Val: cat.p2Val,
			Advantage:   advantage,
			Points1:     math.Round(s1*1000) / 1000,
			Points2:     math.Round(s2*1000) / 1000,
			Weight:      cat.weight,
		})
	}

	score1 = math.Round(score1*100) / 100
	score2 = math.Round(score2*100) / 100

	winner := determineWinner(p1, p2, score1, score2)

	return &model.Battle{
		Pokemon1:      *p1,
		Pokemon2:      *p2,
		Winner:        winner,
		BattleLog:     battleLog,
		Pokemon1Score: score1,
		Pokemon2Score: score2,
	}
}

func determineWinner(p1, p2 *model.Pokemon, score1, score2 float64) string {
	if score1 > score2 {
		return p1.Name
	}
	if score2 > score1 {
		return p2.Name
	}
	// Tie-breaker: Speed
	if p1.Stats.Speed > p2.Stats.Speed {
		return p1.Name
	}
	if p2.Stats.Speed > p1.Stats.Speed {
		return p2.Name
	}
	// Final tie-breaker: alphabetical
	if p1.Name < p2.Name {
		return p1.Name
	}
	return p2.Name
}
