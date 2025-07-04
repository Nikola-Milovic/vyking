package http

import (
	"context"
	"net/http"

	"github.com/Nikola-Milovic/vyking-interview/internal/domain"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("Player Activity API")
	s.OpenAPISchema().SetDescription("This service provides player activity reports by country with enriched country information.")
	s.OpenAPISchema().SetVersion("v1.0.0")

	s.Wrap(
		gzip.Middleware,
	)

	s.Get("/country-player-stats", h.getCountryPlayerStats())
	s.Get("/health", h.health())

	s.Docs("/docs", swgui.New)

	mux.Handle("/", s)
}

type getCountryPlayerStatsInput struct {
	Limit int `query:"limit" default:"10" minimum:"1" maximum:"100" description:"Maximum number of countries to return"`
}

type getCountryPlayerStatsOutput struct {
	Stats []CountryPlayerStatsResponse `json:"stats"`
}

func (h *Handler) getCountryPlayerStats() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input getCountryPlayerStatsInput, output *getCountryPlayerStatsOutput) error {
		req := domain.GetCountryPlayerStatsRequest{
			Limit: input.Limit,
		}

		resp, err := h.service.GetCountryPlayerStats(ctx, req)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		output.Stats = make([]CountryPlayerStatsResponse, 0, len(resp.Stats))
		for _, stat := range resp.Stats {
			response := CountryPlayerStatsResponse{
				CountryCode:     stat.CountryCode,
				PlayerCount:     stat.PlayerCount,
				TotalBets:       stat.TotalBets,
				AvgBetPerPlayer: stat.AvgBetPerPlayer,
			}

			if !stat.CountryInfo.IsZero() {
				response.CountryInfo = &CountryInfo{
					Name:    stat.CountryInfo.Name,
					Region:  stat.CountryInfo.Region,
					Borders: stat.CountryInfo.Borders,
				}
			}

			output.Stats = append(output.Stats, response)
		}

		return nil
	})

	u.SetTitle("Get Country Player Statistics")
	u.SetDescription("Returns player activity statistics grouped by country with enriched country information")
	u.SetTags("Statistics")

	u.SetExpectedErrors(
		status.InvalidArgument,
		status.Internal,
		status.Unavailable,
	)

	return u
}

type healthOutput struct {
	Status string `json:"status"`
}

func (h *Handler) health() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *healthOutput) error {
		output.Status = "OK"
		return nil
	})

	u.SetTitle("Health Check")
	u.SetDescription("Simple health check endpoint")
	u.SetTags("System")

	u.SetExpectedErrors(
		status.Unavailable,
	)

	return u
}
