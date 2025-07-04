INSERT INTO players (name, email, country_code) VALUES
('Marko Marković', 'marko.markovic@example.com', 'RS'),
('Ana Petrović', 'ana.petrovic@example.com', 'RS'),
('Milan Jovanović', 'milan.jovanovic@example.com', 'RS'),
('Jelena Nikolić', 'jelena.nikolic@example.com', 'RS'),
('Stefan Stojanović', 'stefan.stojanovic@example.com', 'RS'),
('Milica Đorđević', 'milica.djordjevic@example.com', 'RS'),
('Nikola Stanković', 'nikola.stankovic@example.com', 'RS'),
('Tijana Milić', 'tijana.milic@example.com', 'RS'),
('Aleksandar Pavlović', 'aleksandar.pavlovic@example.com', 'RS'),
('Jovana Popović', 'jovana.popovic@example.com', 'RS'),

('Hans Mueller', 'hans.mueller@example.com', 'DE'),
('Anna Schmidt', 'anna.schmidt@example.com', 'DE'),
('Klaus Weber', 'klaus.weber@example.com', 'DE'),
('Emma Fischer', 'emma.fischer@example.com', 'DE'),
('Thomas Wagner', 'thomas.wagner@example.com', 'DE'),

('João Silva', 'joao.silva@example.com', 'BR'),
('Maria Santos', 'maria.santos@example.com', 'BR'),
('Pedro Oliveira', 'pedro.oliveira@example.com', 'BR'),
('Ana Costa', 'ana.costa@example.com', 'BR'),
('Carlos Rodrigues', 'carlos.rodrigues@example.com', 'BR'),
('Lucia Ferreira', 'lucia.ferreira@example.com', 'BR'),
('Rafael Almeida', 'rafael.almeida@example.com', 'BR'),

('James Smith', 'james.smith@example.com', 'UK'),
('Emma Johnson', 'emma.johnson@example.com', 'UK'),
('Oliver Brown', 'oliver.brown@example.com', 'UK'),
('Sophie Williams', 'sophie.williams@example.com', 'UK'),

('Carlos García', 'carlos.garcia@example.com', 'ES'),
('Maria López', 'maria.lopez@example.com', 'ES'),
('Antonio Martínez', 'antonio.martinez@example.com', 'ES');

INSERT INTO bets (player_id, amount, created_at) 
SELECT 
    p.id,
    ROUND(50 + (RAND() * 450), 2),
    DATE_SUB(NOW(), INTERVAL FLOOR(RAND() * 90) DAY)
FROM players p
CROSS JOIN (SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5) AS bet_count
WHERE RAND() < 0.8;

-- high value bets 
INSERT INTO bets (player_id, amount, created_at)
SELECT 
    p.id,
    ROUND(500 + (RAND() * 1500), 2),
    DATE_SUB(NOW(), INTERVAL FLOOR(RAND() * 30) DAY)
FROM players p
WHERE p.country_code IN ('RS', 'DE', 'BR')
AND RAND() < 0.3;
