package entity

type PendingJourneyEntity struct {
	ClientCode    string
	Provider      string
	Pending       bool
	Payment       bool
	KYC           bool
	CreatedBy     string
	UpdatedBy     string
	InvalidClient bool
	ApiError      string
}
