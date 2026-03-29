import React from 'react';
import {Battle} from '../types';

interface BattleHistoryProps {
    battles: Battle[];
}

const BattleHistory: React.FC<BattleHistoryProps> = ({battles}) => {
    if (battles.length === 0) {
        return null;
    }

    return (
        <div className="battle-history">
            <h3>Battle History</h3>
            <table>
                <thead>
                <tr>
                    <th>Date</th>
                    <th>Match</th>
                    <th>Score</th>
                    <th>Winner</th>
                </tr>
                </thead>
                <tbody>
                {battles.map((battle) => (
                    <tr key={battle.id}>
                        <td>{new Date(battle.created_at).toLocaleString()}</td>
                        <td>
                            {battle.pokemon1.name} vs {battle.pokemon2.name}
                        </td>
                        <td>
                            {battle.pokemon1_score.toFixed(2)} - {battle.pokemon2_score.toFixed(2)}
                        </td>
                        <td className="winner-cell">{battle.winner}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default BattleHistory;
