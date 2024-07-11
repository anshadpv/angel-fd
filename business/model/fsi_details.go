package model

type FsiStruct struct {
	Plans      []FsiDetailPlans       `json:"plans"`
	About      map[string]interface{} `json:"about"`
	Calculator interface{}            `json:"calculator"`
	FAQs       []FAQ                  `json:"faqs"`
}

type FsiDetailPlans struct {
	Fsi           string  `json:"fsi"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	InterestRate  float64 `json:"interestRate"`
	LockinMonths  int     `json:"lockinMonths"`
	WomenBenefit  float64 `json:"womenBenefit"`
	SeniorCitizen float64 `json:"seniorCitizen"`
	ImageURL      string  `json:"imageUrl"`
	Description   string  `json:"description"`
	InsuredAmount int     `json:"insuredAmount"`
}

type FsiDetailsList struct {
	FsiDetailsList []FsiStruct `json:"fsiList"`
}
