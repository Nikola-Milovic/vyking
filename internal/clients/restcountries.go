package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Nikola-Milovic/vyking-interview/internal/cache"
	"github.com/Nikola-Milovic/vyking-interview/internal/domain"
)

type RestCountriesClient struct {
	httpClient *http.Client
	cache      cache.Cache
	baseURL    string
	cacheTTL   time.Duration
}

func NewRestCountriesClient(cache cache.Cache, cacheTTL time.Duration) *RestCountriesClient {
	return &RestCountriesClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache:    cache,
		baseURL:  "https://restcountries.com/v3.1",
		cacheTTL: cacheTTL,
	}
}

type countryResponse struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Region  string   `json:"region"`
	Borders []string `json:"borders"`
}

func (c *RestCountriesClient) GetCountryInfo(ctx context.Context, countryCode string) (domain.CountryInfo, error) {
	cacheKey := strings.ToLower(countryCode)

	if cached, found := c.cache.Get(ctx, cacheKey); found {
		if info, ok := cached.(domain.CountryInfo); ok {
			slog.Debug("cache hit", slog.String("country_code", countryCode))
			return info, nil
		}
	}

	url := fmt.Sprintf("%s/alpha/%s", c.baseURL, countryCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.CountryInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.CountryInfo{}, fmt.Errorf("failed to fetch country info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.CountryInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var countries []countryResponse
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return domain.CountryInfo{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(countries) == 0 {
		return domain.CountryInfo{}, fmt.Errorf("no country found for code: %s", countryCode)
	}

	country := countries[0]
	info := domain.CountryInfo{
		Name:    country.Name.Common,
		Region:  country.Region,
		Borders: country.Borders,
	}

	if info.Borders == nil {
		info.Borders = []string{}
	}

	slog.Debug("got country info", slog.Any("info", info))

	c.cache.Set(ctx, cacheKey, info, c.cacheTTL)

	return info, nil
}
