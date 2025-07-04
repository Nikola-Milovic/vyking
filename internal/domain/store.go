package domain

import "context"

type Store interface {
	GetTopCountriesByPlayerActivity(ctx context.Context, query GetTopCountriesByPlayerActivityQuery) (*GetTopCountriesByPlayerActivityResult, error)
}

type (
	GetTopCountriesByPlayerActivityQuery struct {
		Limit int
	}
	GetTopCountriesByPlayerActivityResult struct {
		Stats []CountryPlayerStats
	}
)
