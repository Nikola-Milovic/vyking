DROP PROCEDURE IF EXISTS GetTopCountriesByPlayerActivity;

CREATE PROCEDURE GetTopCountriesByPlayerActivity(IN limit_count INT)
BEGIN
    SELECT 
        p.country_code AS country_code,
        COUNT(DISTINCT p.id) AS player_count,
        COALESCE(SUM(b.amount), 0) AS total_bets,
        COALESCE(SUM(b.amount) / COUNT(DISTINCT p.id), 0) AS avg_bet_per_player
    FROM 
        players p
    LEFT JOIN 
        bets b ON p.id = b.player_id
    GROUP BY 
        p.country_code
    ORDER BY 
        player_count DESC,
        total_bets DESC
    LIMIT limit_count;
END;