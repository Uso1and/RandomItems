CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL UNIQUE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    pity_counter  INTEGER DEFAULT 0
);
