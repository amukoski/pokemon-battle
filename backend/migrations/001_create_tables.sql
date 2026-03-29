CREATE TABLE IF NOT EXISTS pokemons (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    height      INTEGER NOT NULL,
    weight      INTEGER NOT NULL,
    types       JSONB NOT NULL,
    stats       JSONB NOT NULL,
    abilities   JSONB NOT NULL,
    sprite_url  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pokemons_name ON pokemons(name);
CREATE INDEX IF NOT EXISTS idx_pokemons_name_prefix ON pokemons(name varchar_pattern_ops);

CREATE TABLE IF NOT EXISTS battles (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pokemon1_id    INTEGER NOT NULL REFERENCES pokemons(id),
    pokemon2_id    INTEGER NOT NULL REFERENCES pokemons(id),
    winner         VARCHAR(100) NOT NULL,
    battle_log     JSONB NOT NULL,
    pokemon1_score DECIMAL(6,2) NOT NULL,
    pokemon2_score DECIMAL(6,2) NOT NULL,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_battles_created_at ON battles(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_battles_winner ON battles(winner);
