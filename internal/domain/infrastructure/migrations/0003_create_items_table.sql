CREATE TABLE IF NOT EXISTS items (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    rarity      VARCHAR(20) NOT NULL,
    base_chance FLOAT NOT NULL,
    min_pity    INTEGER NOT NULL
);
