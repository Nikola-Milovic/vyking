package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Nikola-Milovic/vyking-interview/internal/domain"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	if db == nil {
		panic("db is nil")
	}

	return &Store{db: db}
}

func (s *Store) GetTopCountriesByPlayerActivity(ctx context.Context, q domain.GetTopCountriesByPlayerActivityQuery) (*domain.GetTopCountriesByPlayerActivityResult, error) {
	query := "CALL GetTopCountriesByPlayerActivity(?)"

	rows, err := s.db.QueryContext(ctx, query, q.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute stored procedure: %w", err)
	}
	defer rows.Close()

	var stats []domain.CountryPlayerStats
	for rows.Next() {
		var stat domain.CountryPlayerStats
		err := rows.Scan(
			&stat.CountryCode,
			&stat.PlayerCount,
			&stat.TotalBets,
			&stat.AvgBetPerPlayer,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &domain.GetTopCountriesByPlayerActivityResult{
		Stats: stats,
	}, nil
}
