CREATE TABLE IF NOT EXISTS players (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    country_code VARCHAR(2) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_players_country_code ON players(country_code);
CREATE INDEX idx_players_email ON players(email);

CREATE TABLE IF NOT EXISTS bets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    player_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

CREATE INDEX idx_bets_player_id ON bets(player_id);
CREATE INDEX idx_bets_created_at ON bets(created_at);
