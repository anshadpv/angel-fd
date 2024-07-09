package entity

type PortfolioEntity struct {
	ClientCode          string
	TotalActiveDeposits int
	Provider            string
	InvestedValue       float64
	CurrentValue        float64
	InterestEarned      float64
	ReturnsValue        float64
	ReturnsPercentage   float64
	CreatedBy           string
	UpdatedBy           string
	InvalidClient       bool
	ApiError            string
}
