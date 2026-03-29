import React, {useEffect, useState} from 'react';
import {Battle} from './types';
import {executeBattle, getBattleHistory} from './api';
import BattleForm from './components/BattleForm';
import BattleResult from './components/BattleResult';
import BattleHistory from './components/BattleHistory';
import './styles/index.css';

const App: React.FC = () => {
    const [battle, setBattle] = useState<Battle | null>(null);
    const [history, setHistory] = useState<Battle[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const loadHistory = async () => {
        try {
            const battles = await getBattleHistory();
            setHistory(battles);
        } catch (err: any) {
            console.error(err);
        }
    };

    useEffect(() => {
        loadHistory()
    }, []);

    const handleBattle = async (pokemon1: string, pokemon2: string) => {
        setLoading(true);
        setError(null);
        setBattle(null);

        try {
            const result = await executeBattle(pokemon1, pokemon2);
            setBattle(result);
            loadHistory();
        } catch (err: any) {
            setError(err.message || 'An unexpected error occurred');
        } finally {
            setLoading(false);
        }
    };

    const handleReset = () => {
        setBattle(null);
        setError(null);
    };

    return (
        <div className="app">
            <header className="app-header">
                <h1>Pokemon Battle</h1>
                <p>Enter two Pokemon names to simulate a battle!</p>
            </header>

            <main>
                <BattleForm onBattle={handleBattle} loading={loading}/>

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}

                {battle && <BattleResult battle={battle} onReset={handleReset}/>}

                <BattleHistory battles={history}/>
            </main>
        </div>
    );
};

export default App;
