package domain

import "time"

type Player struct {
	ID          int
	Name        string
	Email       string
	CountryCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Bet struct {
	ID        int
	PlayerID  int
	Amount    float64
	CreatedAt time.Time
}

type CountryPlayerStats struct {
	CountryCode     string
	PlayerCount     int
	TotalBets       float64
	AvgBetPerPlayer float64
}

type CountryInfo struct {
	Name    string
	Region  string
	Borders []string
}

func (c CountryInfo) IsZero() bool {
	return c.Name == "" && c.Region == ""
}

type CountryPlayerStatsWithInfo struct {
	CountryPlayerStats
	CountryInfo CountryInfo
}
