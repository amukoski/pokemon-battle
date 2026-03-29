export interface PokemonStats {
    hp: number;
    attack: number;
    defense: number;
    special_attack: number;
    special_defense: number;
    speed: number;
}

export interface Pokemon {
    id: number;
    name: string;
    height: number;
    weight: number;
    types: string[];
    stats: PokemonStats;
    abilities: string[];
    sprite_url: string;
}

export interface BattleRound {
    category: string;
    pokemon1_val: number;
    pokemon2_val: number;
    advantage: string;
    points1: number;
    points2: number;
    weight: number;
}

export interface Battle {
    id: string;
    pokemon1: Pokemon;
    pokemon2: Pokemon;
    winner: string;
    battle_log: BattleRound[];
    pokemon1_score: number;
    pokemon2_score: number;
    created_at: string;
}
