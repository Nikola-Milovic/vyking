package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Nikola-Milovic/vyking-interview/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	store            domain.Store
	countryAPIClient domain.CountryAPIClient
}

func New(store domain.Store, countryAPIClient domain.CountryAPIClient) Service {
	return Service{
		store:            store,
		countryAPIClient: countryAPIClient,
	}
}

func (s Service) GetCountryPlayerStats(ctx context.Context, req domain.GetCountryPlayerStatsRequest) (domain.GetCountryPlayerStatsResponse, error) {
	query := domain.GetTopCountriesByPlayerActivityQuery{
		Limit: req.Limit,
	}

	result, err := s.store.GetTopCountriesByPlayerActivity(ctx, query)
	if err != nil {
		return domain.GetCountryPlayerStatsResponse{}, err
	}

	res := domain.GetCountryPlayerStatsResponse{
		Stats: make([]domain.CountryPlayerStatsWithInfo, 0, len(result.Stats)),
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	var mu sync.Mutex
	statsWithInfo := make([]domain.CountryPlayerStatsWithInfo, len(result.Stats))

	for i, stat := range result.Stats {
		g.Go(func() error {
			countryInfo, err := s.countryAPIClient.GetCountryInfo(ctx, stat.CountryCode)
			if err != nil {
				slog.Error("failed to fetch country info", "country_code", stat.CountryCode, "error", err)
				// Graceful degradation
				countryInfo = domain.CountryInfo{}
			}

			mu.Lock()
			statsWithInfo[i] = domain.CountryPlayerStatsWithInfo{
				CountryPlayerStats: stat,
				CountryInfo:        countryInfo,
			}
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return domain.GetCountryPlayerStatsResponse{}, fmt.Errorf("failed to wait: %w", err)
	}

	res.Stats = statsWithInfo

	return res, nil
}
