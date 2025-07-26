-- Сначала обновляем min_pity для существующих записей (если они есть)
UPDATE items SET min_pity = 0 WHERE id = 1;  -- Обычный предмет - всегда доступен
UPDATE items SET min_pity = 20 WHERE id = 3; -- Редкий
UPDATE items SET min_pity = 50 WHERE id = 4; -- Эпический
UPDATE items SET min_pity = 100 WHERE id = 5; -- Легендарный

-- Затем добавляем новые записи (если их нет)
INSERT INTO items (name, rarity, base_chance, min_pity) VALUES
('Малый зелье здоровья', 'common', 0.3, 0),
('Средний меч', 'uncommon', 0.25, 0),
('Редкий посох маны', 'rare', 0.2, 20),
('Эпические сапоги скорости', 'epic', 0.15, 50),
('Легендарный меч дракона', 'legendary', 0.1, 100)
ON CONFLICT (id) DO NOTHING;
