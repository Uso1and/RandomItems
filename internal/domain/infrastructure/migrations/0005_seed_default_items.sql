-- Убедимся, что таблица существует
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    base_chance FLOAT NOT NULL,
    min_pity INTEGER NOT NULL
);


INSERT INTO items (id, name, rarity, base_chance, min_pity) VALUES
(1, 'Малый зелье здоровья', 'common', 0.5, 0),
(2, 'Средний меч', 'uncommon', 0.35, 0),
(3, 'Редкий посох маны', 'rare', 0.25, 15),
(4, 'Эпические сапоги скорости', 'epic', 0.15, 30),
(5, 'Легендарный меч дракона', 'legendary', 0.05, 50)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    rarity = EXCLUDED.rarity,
    base_chance = EXCLUDED.base_chance,
    min_pity = EXCLUDED.min_pity;

