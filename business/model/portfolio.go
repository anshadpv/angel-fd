package model

var (
	EmptyPortfolio = Portfolio{TotalActiveDeposits: 0, InvestedValue: 0, CurrentValue: 0, InterestEarned: 0, ReturnsValue: 0, ReturnsPercentage: 0}
)

type Portfolio struct {
	TotalActiveDeposits int     `json:"totalActiveDeposits"`
	InvestedValue       float64 `json:"investedValue"`
	CurrentValue        float64 `json:"currentValue"`
	InterestEarned      float64 `json:"interestEarned"`
	ReturnsValue        float64 `json:"returnsValue"`
	ReturnsPercentage   float64 `json:"returnsPercentage"`
}
