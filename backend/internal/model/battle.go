package model

import "time"

type Battle struct {
	ID            string        `json:"id"`
	Pokemon1      Pokemon       `json:"pokemon1"`
	Pokemon2      Pokemon       `json:"pokemon2"`
	Winner        string        `json:"winner"`
	BattleLog     []BattleRound `json:"battle_log"`
	Pokemon1Score float64       `json:"pokemon1_score"`
	Pokemon2Score float64       `json:"pokemon2_score"`
	CreatedAt     time.Time     `json:"created_at"`
}

type BattleRound struct {
	Category    string  `json:"category"`
	Pokemon1Val int     `json:"pokemon1_val"`
	Pokemon2Val int     `json:"pokemon2_val"`
	Advantage   string  `json:"advantage"`
	Points1     float64 `json:"points1"`
	Points2     float64 `json:"points2"`
	Weight      float64 `json:"weight"`
}

type BattleRequest struct {
	Pokemon1 string `json:"pokemon1"`
	Pokemon2 string `json:"pokemon2"`
}
