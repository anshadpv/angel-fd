package model

type FsiList struct {
	FsiList []FsiDetails `json:"fsiList"`
}

type FsiDetails struct {
	Fsi          string  `json:"fsi"`
	Name         string  `json:"name"`
	ImageUrl     string  `json:"imageUrl"`
	InterestRate float64 `json:"interestRate"`
}

type CompareFSIDBDetails struct {
	FSI                  string  `json:"fsi"`
	Name                 string  `json:"name"`
	TenureYears          int     `json:"tenureYears"`
	TenureMonths         int     `json:"tenureMonths"`
	TenureDays           int     `json:"tenureDays"`
	InterestRate         float64 `json:"interestRate"`
	MinDeposit           int     `json:"minDeposit"`
	SeniorCitizenBenefit bool    `json:"seniorCitizenBenefit"`
	BankAccount          string  `json:"bankAccount"`
	InsuredAmount        int     `json:"insuredAmount"`
	ImageURL             string  `json:"imageUrl"`
}
type CompareFSIDetails struct {
	FSI                  string             `json:"fsi"`
	Name                 string             `json:"name"`
	YearlyInterestRate   YearlyInterestRate `json:"yearlyInterestRate"`
	MinDeposit           int                `json:"minDeposit"`
	SeniorCitizenBenefit bool               `json:"seniorCitizenBenefit"`
	BankAccount          string             `json:"bankAccount"`
	InsuredAmount        int                `json:"insuredAmount"`
	ImageURL             string             `json:"imageUrl"`
}

type YearlyInterestRate struct {
	ZeroToOne   float64 `json:"0_to_1Y"`
	OneToTwo    float64 `json:"1_to_2Y"`
	TwoToThree  float64 `json:"2_to_3Y"`
	ThreeToFour float64 `json:"3_to_4Y"`
	FourToFive  float64 `json:"4_to_5Y"`
	FiveToSix   float64 `json:"5_to_6Y"`
}
