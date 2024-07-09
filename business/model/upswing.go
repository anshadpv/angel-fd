package model

type GenerateTokenResponse struct {
	AccessToken      string `json:"access_token,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	IDToken          string `json:"id_token,omitempty"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
	Scope            string `json:"scope,omitempty"`
	USWPC            string `json:"USW-PC,omitempty"`
}

type PCIRegistrationRequest struct {
	PartnerCustomerId string `json:"partnerCustomerId,omitempty"`
}

type PCIRegistrationResponse struct {
	ICI               string `json:"ici,omitempty"`
	GuestSessionToken string `json:"guestSessionToken,omitempty"`
}

type NetWorthResponse struct {
	TotalInvestedAmount    TotalInvestedAmount `json:"totalInvestedAmount"`
	TotalInterestEarned    TotalInterestEarned `json:"totalInterestEarned"`
	CurrentAmount          CurrentAmount       `json:"currentAmount"`
	ActiveTermDepositCount int                 `json:"activeTermDepositCount"`
}
type TotalInvestedAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
type TotalInterestEarned struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
type CurrentAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type DataIngestionRequest struct {
	Pan           string             `json:"pan,omitempty"`
	Fullname      *DataIngestionName `json:"fullName,omitempty"`
	MotherName    *DataIngestionName `json:"motherName,omitempty"`
	FatherName    *DataIngestionName `json:"fatherName,omitempty"`
	Gender        string             `json:"gender,omitempty"`
	Email         string             `json:"email,omitempty"`
	DateOfBirth   string             `json:"dateOfBirth,omitempty"`
	Occupation    string             `json:"occupation,omitempty"`
	Income        int                `json:"income,omitempty"`
	MaritalStatus string             `json:"maritalStatus,omitempty"`
	Address       []struct {
		AddressLine1 string `json:"addressLine1,omitempty"`
		AddressLine2 string `json:"addressLine2,omitempty"`
		AddressLine3 string `json:"addressLine3,omitempty"`
		PostalCode   string `json:"postalCode,omitempty"`
	} `json:"address,omitempty"`
	WithdrawalBank []struct {
		BankAccountNumber string `json:"bankAccountNumber,omitempty"`
		Ifsc              string `json:"ifsc,omitempty"`
	} `json:"withdrawalBank,omitempty"`
	Nominee []struct {
		FullName         string `json:"fullName,omitempty"`
		Relation         string `json:"relation,omitempty"`
		DateOfBirth      string `json:"dateOfBirth,omitempty"`
		PhoneNumber      string `json:"phoneNumber,omitempty"`
		Email            string `json:"email,omitempty"`
		IsAddressSimilar bool   `json:"isAddressSimilar,omitempty"`
		Address          struct {
			AddressLine1 string `json:"addressLine1,omitempty"`
			AddressLine2 string `json:"addressLine2,omitempty"`
			AddressLine3 string `json:"addressLine3,omitempty"`
			PostalCode   string `json:"postalCode,omitempty"`
		} `json:"address,omitempty"`
		GuardianInfo struct {
			FullName         string `json:"fullName,omitempty"`
			Relation         string `json:"relation,omitempty"`
			DateOfBirth      string `json:"dateOfBirth,omitempty"`
			PhoneNumber      string `json:"phoneNumber,omitempty"`
			IsAddressSimilar bool   `json:"isAddressSimilar,omitempty"`
			Address          struct {
				AddressLine1 string `json:"addressLine1,omitempty"`
				AddressLine2 string `json:"addressLine2,omitempty"`
				AddressLine3 string `json:"addressLine3,omitempty"`
				PostalCode   string `json:"postalCode,omitempty"`
			} `json:"address,omitempty"`
		} `json:"guardianInfo,omitempty"`
	} `json:"nominee,omitempty"`
}

type DataIngestionName struct {
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
}

type UpSwingWebhookEvent struct {
	Pci             string  `json:"pci,omitempty"`
	Fsi             string  `json:"fsi,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
	Tenure          string  `json:"tenure,omitempty"`
	JourneyID       string  `json:"journeyId,omitempty"`
	EventType       string  `json:"eventType,omitempty"`
	TermDepositType string  `json:"termDepositType,omitempty"`
	Reason          string  `json:"reason,omitempty"`
}

type PendingJourneyResponse struct {
	JourneyPending          bool `json:"journeyPending"`
	JourneyPendingOnPayment bool `json:"journeyPendingOnPayment"`
	JourneyPendingOnVkyc    bool `json:"journeyPendingOnVkyc"`
}
