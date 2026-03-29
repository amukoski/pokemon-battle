package service

import (
	"testing"

	"github.com/pokemon-battle/backend/internal/model"
)

func TestCalculateBattle_HigherOverallStatsWins(t *testing.T) {
	p1 := &model.Pokemon{
		Name: "strong",
		Stats: model.PokemonStats{
			HP: 100, Attack: 100, Defense: 100,
			SpecialAttack: 100, SpecialDefense: 100, Speed: 100,
		},
	}
	p2 := &model.Pokemon{
		Name: "weak",
		Stats: model.PokemonStats{
			HP: 50, Attack: 50, Defense: 50,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
		},
	}

	battle := CalculateBattle(p1, p2)

	if battle.Winner != "strong" {
		t.Errorf("expected winner 'strong', got '%s'", battle.Winner)
	}
	if battle.Pokemon1Score <= battle.Pokemon2Score {
		t.Errorf("expected p1 score > p2 score, got %.2f vs %.2f", battle.Pokemon1Score, battle.Pokemon2Score)
	}
}

func TestCalculateBattle_EqualStatsTieBreakBySpeed(t *testing.T) {
	p1 := &model.Pokemon{
		Name: "fast",
		Stats: model.PokemonStats{
			HP: 50, Attack: 50, Defense: 50,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 100,
		},
	}
	p2 := &model.Pokemon{
		Name: "slow",
		Stats: model.PokemonStats{
			HP: 50, Attack: 50, Defense: 50,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
		},
	}

	battle := CalculateBattle(p1, p2)

	// p1 has higher speed so gets more points in speed category
	// Even though other stats are equal, the speed difference tips it
	if battle.Winner != "fast" {
		t.Errorf("expected winner 'fast', got '%s'", battle.Winner)
	}
}

func TestCalculateBattle_CompletelyEqualTieBreakAlphabetical(t *testing.T) {
	stats := model.PokemonStats{
		HP: 50, Attack: 50, Defense: 50,
		SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
	}
	p1 := &model.Pokemon{Name: "alpha", Stats: stats}
	p2 := &model.Pokemon{Name: "beta", Stats: stats}

	battle := CalculateBattle(p1, p2)

	if battle.Winner != "alpha" {
		t.Errorf("expected winner 'alpha' (alphabetical tie-break), got '%s'", battle.Winner)
	}
	if battle.Pokemon1Score != battle.Pokemon2Score {
		t.Errorf("expected equal scores, got %.2f vs %.2f", battle.Pokemon1Score, battle.Pokemon2Score)
	}
}

func TestCalculateBattle_BattleLogHasSixRounds(t *testing.T) {
	stats := model.PokemonStats{
		HP: 50, Attack: 50, Defense: 50,
		SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
	}
	p1 := &model.Pokemon{Name: "a", Stats: stats}
	p2 := &model.Pokemon{Name: "b", Stats: stats}

	battle := CalculateBattle(p1, p2)

	if len(battle.BattleLog) != 6 {
		t.Errorf("expected 6 battle rounds, got %d", len(battle.BattleLog))
	}

	expectedCategories := []string{"HP", "Attack", "Defense", "Sp. Attack", "Sp. Defense", "Speed"}
	for i, round := range battle.BattleLog {
		if round.Category != expectedCategories[i] {
			t.Errorf("round %d: expected category '%s', got '%s'", i, expectedCategories[i], round.Category)
		}
	}
}

func TestCalculateBattle_ZeroStats(t *testing.T) {
	p1 := &model.Pokemon{
		Name: "zero1",
		Stats: model.PokemonStats{
			HP: 0, Attack: 0, Defense: 0,
			SpecialAttack: 0, SpecialDefense: 0, Speed: 0,
		},
	}
	p2 := &model.Pokemon{
		Name: "zero2",
		Stats: model.PokemonStats{
			HP: 0, Attack: 0, Defense: 0,
			SpecialAttack: 0, SpecialDefense: 0, Speed: 0,
		},
	}

	battle := CalculateBattle(p1, p2)

	// Should not panic, winner determined by alphabetical tie-break
	if battle.Winner != "zero1" {
		t.Errorf("expected winner 'zero1' (alphabetical tie-break), got '%s'", battle.Winner)
	}
	if battle.Pokemon1Score != 0 || battle.Pokemon2Score != 0 {
		t.Errorf("expected scores of 0, got %.2f vs %.2f", battle.Pokemon1Score, battle.Pokemon2Score)
	}
}

func TestCalculateBattle_ScoresAddUpToOne(t *testing.T) {
	p1 := &model.Pokemon{
		Name: "pikachu",
		Stats: model.PokemonStats{
			HP: 35, Attack: 55, Defense: 40,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 90,
		},
	}
	p2 := &model.Pokemon{
		Name: "charizard",
		Stats: model.PokemonStats{
			HP: 78, Attack: 84, Defense: 78,
			SpecialAttack: 109, SpecialDefense: 85, Speed: 100,
		},
	}

	battle := CalculateBattle(p1, p2)

	total := battle.Pokemon1Score + battle.Pokemon2Score
	if total < 0.99 || total > 1.01 {
		t.Errorf("expected scores to add up to ~1.0, got %.4f (%.2f + %.2f)",
			total, battle.Pokemon1Score, battle.Pokemon2Score)
	}
}

func TestCalculateBattle_AdvantageField(t *testing.T) {
	p1 := &model.Pokemon{
		Name: "attacker",
		Stats: model.PokemonStats{
			HP: 50, Attack: 100, Defense: 50,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
		},
	}
	p2 := &model.Pokemon{
		Name: "defender",
		Stats: model.PokemonStats{
			HP: 50, Attack: 50, Defense: 100,
			SpecialAttack: 50, SpecialDefense: 50, Speed: 50,
		},
	}

	battle := CalculateBattle(p1, p2)

	// Attack round should favor p1
	attackRound := battle.BattleLog[1]
	if attackRound.Advantage != "attacker" {
		t.Errorf("expected Attack advantage 'attacker', got '%s'", attackRound.Advantage)
	}

	// Defense round should favor p2
	defenseRound := battle.BattleLog[2]
	if defenseRound.Advantage != "defender" {
		t.Errorf("expected Defense advantage 'defender', got '%s'", defenseRound.Advantage)
	}

	// HP round should be tie
	hpRound := battle.BattleLog[0]
	if hpRound.Advantage != "tie" {
		t.Errorf("expected HP advantage 'tie', got '%s'", hpRound.Advantage)
	}
}
