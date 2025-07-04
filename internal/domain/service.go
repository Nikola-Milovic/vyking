package domain

import "context"

type Service interface {
	GetCountryPlayerStats(ctx context.Context, req GetCountryPlayerStatsRequest) (GetCountryPlayerStatsResponse, error)
}

type (
	GetCountryPlayerStatsRequest struct {
		Limit int
	}
	GetCountryPlayerStatsResponse struct {
		Stats []CountryPlayerStatsWithInfo
	}
)
