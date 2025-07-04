package domain

import "context"

type CountryAPIClient interface {
	GetCountryInfo(ctx context.Context, countryCode string) (CountryInfo, error)
}
