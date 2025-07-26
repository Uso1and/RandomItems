CREATE TABLE IF NOT EXISTS drop_events (
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER REFERENCES users(id),
    item_id      INTEGER REFERENCES items(id),
    dropped_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_guaranteed BOOLEAN DEFAULT FALSE
);
