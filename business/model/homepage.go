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

type PlanTest struct {
	Fsi           string  `json:"fsi"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	InterestRate  float64 `json:"interestRate"`
	LockinMonths  int     `json:"lockinMonths"`
	ImageURL      string  `json:"imageUrl"`
	Description   string  `json:"description"`
	InsuredAmount int     `json:"insuredAmount"`
}
