import React, {useState} from 'react';
import AutocompleteInput from './AutocompleteInput';

interface BattleFormProps {
    onBattle: (pokemon1: string, pokemon2: string) => void;
    loading: boolean;
}

const BattleForm: React.FC<BattleFormProps> = ({onBattle, loading}) => {
    const [pokemon1, setPokemon1] = useState('');
    const [pokemon2, setPokemon2] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (pokemon1.trim() && pokemon2.trim()) {
            onBattle(pokemon1.trim(), pokemon2.trim());
        }
    };

    return (
        <form className="battle-form" onSubmit={handleSubmit}>
            <div className="form-inputs">
                <AutocompleteInput
                    placeholder="First Pokemon"
                    value={pokemon1}
                    onChange={setPokemon1}
                    disabled={loading}
                />
                <span className="vs-badge">VS</span>
                <AutocompleteInput
                    placeholder="Second Pokemon"
                    value={pokemon2}
                    onChange={setPokemon2}
                    disabled={loading}
                />
            </div>
            <button type="submit" disabled={loading || !pokemon1.trim() || !pokemon2.trim()}>
                {loading ? 'Battling...' : 'Battle!'}
            </button>
        </form>
    );
};

export default BattleForm;
