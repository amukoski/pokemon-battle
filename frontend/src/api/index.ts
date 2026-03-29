import {Battle} from '../types';

const API_BASE = process.env.REACT_APP_API_URL || '';

async function handleResponse<T>(response: Response): Promise<T> {
    const data = await response.json();
    if (!response.ok) {
        throw new Error(data.error || "API Error")
    }
    return data as T;
}

export async function executeBattle(pokemon1: string, pokemon2: string): Promise<Battle> {
    const response = await fetch(`${API_BASE}/api/battle`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({pokemon1, pokemon2}),
    });
    return handleResponse<Battle>(response);
}

export async function getBattleHistory(): Promise<Battle[]> {
    const response = await fetch(`${API_BASE}/api/battles?limit=20`);
    return handleResponse<Battle[]>(response);
}

export async function searchPokemonNames(query: string): Promise<string[]> {
    if (!query.trim()) return [];
    const response = await fetch(`${API_BASE}/api/pokemon-names?q=${encodeURIComponent(query)}`);
    return handleResponse<string[]>(response);
}
