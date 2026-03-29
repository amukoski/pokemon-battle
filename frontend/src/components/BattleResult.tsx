import React from 'react';
import {Battle} from '../types';
import PokemonCard from './PokemonCard';

interface BattleResultProps {
    battle: Battle;
    onReset: () => void;
}

const BattleResult: React.FC<BattleResultProps> = ({battle, onReset}) => {
    return (
        <div className="battle-result">
            <div className="battle-cards">
                <PokemonCard
                    pokemon={battle.pokemon1}
                    isWinner={battle.winner === battle.pokemon1.name}
                    score={battle.pokemon1_score}
                />
                <div className="vs-divider">
                    <span>VS</span>
                </div>
                <PokemonCard
                    pokemon={battle.pokemon2}
                    isWinner={battle.winner === battle.pokemon2.name}
                    score={battle.pokemon2_score}
                />
            </div>

            <div className="battle-breakdown">
                <h3>Battle Breakdown</h3>
                <table>
                    <thead>
                    <tr>
                        <th>Category</th>
                        <th>{battle.pokemon1.name}</th>
                        <th>{battle.pokemon2.name}</th>
                        <th>Advantage</th>
                    </tr>
                    </thead>
                    <tbody>
                    {battle.battle_log.map((round) => (
                        <tr key={round.category}>
                            <td>{round.category}</td>
                            <td className={round.advantage === battle.pokemon1.name ? 'highlight' : ''}>
                                {round.pokemon1_val} ({round.points1.toFixed(3)})
                            </td>
                            <td className={round.advantage === battle.pokemon2.name ? 'highlight' : ''}>
                                {round.pokemon2_val} ({round.points2.toFixed(3)})
                            </td>
                            <td className={`advantage ${round.advantage === 'tie' ? 'tie' : ''}`}>
                                {round.advantage}
                            </td>
                        </tr>
                    ))}
                    <tr className="total-row">
                        <td><strong>Total</strong></td>
                        <td><strong>{battle.pokemon1_score.toFixed(2)}</strong></td>
                        <td><strong>{battle.pokemon2_score.toFixed(2)}</strong></td>
                        <td><strong>{battle.winner}</strong></td>
                    </tr>
                    </tbody>
                </table>
            </div>

            <button className="battle-again-btn" onClick={onReset}>
                Battle Again
            </button>
        </div>
    );
};

export default BattleResult;
