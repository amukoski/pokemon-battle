package model

type Pokemon struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Height    int          `json:"height"`
	Weight    int          `json:"weight"`
	Types     []string     `json:"types"`
	Stats     PokemonStats `json:"stats"`
	Abilities []string     `json:"abilities"`
	SpriteURL string       `json:"sprite_url"`
}

type PokemonStats struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}
