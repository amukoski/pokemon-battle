import React from 'react';
import {Pokemon} from '../types';

interface PokemonCardProps {
    pokemon: Pokemon;
    isWinner: boolean;
    score: number;
}

const typeColors: Record<string, string> = {
    normal: '#A8A878', fire: '#F08030', water: '#6890F0', electric: '#F8D030',
    grass: '#78C850', ice: '#98D8D8', fighting: '#C03028', poison: '#A040A0',
    ground: '#E0C068', flying: '#A890F0', psychic: '#F85888', bug: '#A8B820',
    rock: '#B8A038', ghost: '#705898', dragon: '#7038F8', dark: '#705848',
    steel: '#B8B8D0', fairy: '#EE99AC',
};

const statLabels: Record<string, string> = {
    hp: 'HP',
    attack: 'ATK',
    defense: 'DEF',
    special_attack: 'SP.ATK',
    special_defense: 'SP.DEF',
    speed: 'SPD',
};

const PokemonCard: React.FC<PokemonCardProps> = ({pokemon, isWinner, score}) => {
    const stats = pokemon.stats;
    const statEntries = [
        {key: 'hp', value: stats.hp},
        {key: 'attack', value: stats.attack},
        {key: 'defense', value: stats.defense},
        {key: 'special_attack', value: stats.special_attack},
        {key: 'special_defense', value: stats.special_defense},
        {key: 'speed', value: stats.speed},
    ];

    return (
        <div className={`pokemon-card ${isWinner ? 'winner' : ''}`}>
            {isWinner && <div className="winner-crown">Winner!</div>}
            <img src={pokemon.sprite_url} alt={pokemon.name} className="pokemon-sprite"/>
            <h2 className="pokemon-name">{pokemon.name}</h2>
            <p className="pokemon-id">#{pokemon.id}</p>
            <div className="pokemon-types">
                {pokemon.types.map((type_) => (
                    <span
                        key={type_}
                        className="type-badge"
                        style={{backgroundColor: typeColors[type_] || '#888'}}
                    >
            {type_}
          </span>
                ))}
            </div>
            <div className="pokemon-info">
                <span>Height: {(pokemon.height / 10).toFixed(1)}m</span>
                <span>Weight: {(pokemon.weight / 10).toFixed(1)}kg</span>
            </div>
            <div className="pokemon-abilities">
                <strong>Abilities:</strong> {pokemon.abilities.join(', ')}
            </div>
            <div className="pokemon-stats">
                {statEntries.map(({key, value}) => (
                    <div key={key} className="stat-row">
                        <span className="stat-label">{statLabels[key]}</span>
                        <div className="stat-bar-container">
                            <div
                                className="stat-bar"
                                style={{width: `${Math.min((value / 255) * 100, 100)}%`}}
                            />
                        </div>
                        <span className="stat-value">{value}</span>
                    </div>
                ))}
            </div>
            <div className="pokemon-score">Score: {score.toFixed(2)}</div>
        </div>
    );
};

export default PokemonCard;
