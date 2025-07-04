package http

type CountryPlayerStatsResponse struct {
	CountryCode     string       `json:"country_code" description:"ISO 3166-1 alpha-2 country code"`
	PlayerCount     int          `json:"player_count" description:"Number of active players in this country"`
	TotalBets       float64      `json:"total_bets" description:"Total amount of bets placed by players from this country"`
	AvgBetPerPlayer float64      `json:"avg_bet_per_player" description:"Average bet amount per player in this country"`
	CountryInfo     *CountryInfo `json:"country_info" description:"Additional information about the country"`
}

type CountryInfo struct {
	Name    string   `json:"name" description:"Common name of the country"`
	Region  string   `json:"region" description:"Region where the country is located"`
	Borders []string `json:"borders" description:"List of ISO 3166-1 alpha-3 codes of bordering countries"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
