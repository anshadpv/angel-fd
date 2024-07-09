package model

type Homepage struct {
	MostBought []Plan  `json:"mostBought"`
	AllFDS     []Plan  `json:"allFDs"`
	Journey    Journey `json:"journey"`
}

type Journey struct {
	Pending      bool         `json:"pending"`
	PendingState PendingState `json:"pendingState"`
}

type PendingState struct {
	Payment bool `json:"payment"`
	KYC     bool `json:"kyc"`
}
