package model

type Plans struct {
	Plans []Plan `json:"plans"`
}

type FsiPlans struct {
	Plans                  []Plan                 `json:"plans"`
	MaxInterestRate        float64                `json:"maxInterestRate"`
	MinInvestment          int                    `json:"minInvestment"`
	InsuredAmount          int                    `json:"insuredAmount"`
	CompareFsi             string                 `json:"compareFsi"`
	CompareFsiName         string                 `json:"compareFsiName"`
	CompareFsiInterestRate float64                `json:"compareFsiInterestRate"`
	CompareFsiImageUrl     string                 `json:"compareFsiImageUrl"`
	About                  map[string]interface{} `json:"about"`
	Calculator             interface{}            `json:"calculator"`
}

type Plan struct {
	Fsi           string  `json:"fsi"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TenureYears   int     `json:"tenureYears"`
	TenureMonths  int     `json:"tenureMonths"`
	TenureDays    int     `json:"tenureDays"`
	InterestRate  float64 `json:"interestRate"`
	LockinMonths  int     `json:"lockinMonths"`
	WomenBenefit  float64 `json:"womenBenefit"`
	SeniorCitizen float64 `json:"seniorCitizen"`
	ImageURL      string  `json:"imageUrl"`
	Description   string  `json:"description"`
	InsuredAmount int     `json:"insuredAmount"`
}
