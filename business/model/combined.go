package model

type CombinedResponse struct {
	NetWorthData   NetWorthResponse           `json:"portfolio"`
	PendingJourney PendingJourneyResponseTest `json:"pendingJourney"`
}

type PendingJourneyResponseTest struct {
	JourneyPending          bool `json:"pending"`
	JourneyPendingOnPayment bool `json:"payment"`
	JourneyPendingOnVkyc    bool `json:"kyc"`
}
