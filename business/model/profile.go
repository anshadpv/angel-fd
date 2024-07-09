package model

type ProfileRequest struct {
	ClientCode string `json:"clientCode"`
}

type User struct {
	SubBrokerTag string `json:"subBrokerTag" binding:"required"`
	UserType     string `json:"userType" binding:"required"`
	Arn          string `json:"arn" binding:"required"`
	Euin         string `json:"euin" binding:"required"`
	IsARNExpired string `json:"isARNExpired" binding:"required"`
}
type UserSubBrokerDetails struct {
	SubBrokerTag string `json:"subBrokerTag" binding:"required"`
	UserType     string `json:"userType" binding:"required"`
	Arn          string `json:"arn" binding:"required"`
	Euin         string `json:"euin" binding:"required"`
	IsARNExpired string `json:"isARNExpired" binding:"required"`
}

type UserProfileInfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}
type ProfileResponse struct {
	Status    string             `json:"status"`
	Message   string             `json:"message,omitempty"`
	ErrorCode string             `json:"errorCode,omitempty"`
	Data      UserProfileDetails `json:"data,omitempty"`
}
type UserInfoRequest struct {
	ClientCode string `json:"ClientCode"`
}

type UserProfileResponse struct {
	Status     string             `json:"status"`
	StatusCode int                `json:"statusCode"`
	Message    string             `json:"message"`
	ErrorCode  string             `json:"error_code"`
	Data       UserProfileDetails `json:"data"`
}

type UserProfileDetails struct {
	ApplicationNo        string               `json:"applicationNo"`
	ClientID             string               `json:"clientId"`
	Pan                  string               `json:"pan"`
	CountryCode          string               `json:"countryCode"`
	Mobile               string               `json:"mobile"`
	Active               bool                 `json:"active"`
	ClientAccountType    string               `json:"clientAccountType"`
	ClientAccountSubType string               `json:"clientAccountSubType"`
	UserType             string               `json:"userType"`
	DpNumber             string               `json:"dpNumber"`
	PoaStatus            string               `json:"poaStatus"`
	BrokerageType        string               `json:"brokerageType"`
	ReferralCode         string               `json:"referralCode"`
	CreationDate         string               `json:"creationDate"`
	ModificationDate     string               `json:"modificationDate"`
	DpcFlag              bool                 `json:"dpcFlag"`
	ClientDetails        ClientDetails        `json:"clientDetails"`
	ActiveSegments       ActiveSegments       `json:"activeSegments"`
	BankDetails          []BankDetails        `json:"bankDetails"`
	DpDetails            DpDetails            `json:"dpDetails"`
	BrokerageDetails     BrokerageDetails     `json:"brokerageDetails"`
	Nominee              Nominee              `json:"nominee"`
	UserSubBrokerDetails UserSubBrokerDetails `json:"subBrokerDetails"`
}

type ClientDetails struct {
	FullName              string                `json:"fullName"`
	FirstName             string                `json:"firstName"`
	MiddleName            string                `json:"middleName"`
	LastName              string                `json:"lastName"`
	FatherName            string                `json:"fatherName"`
	Email                 string                `json:"email"`
	Branch                string                `json:"branch"`
	SubBroker             string                `json:"subBroker"`
	ShortName             string                `json:"shortName"`
	Birthdate             string                `json:"birthdate"`
	Gender                string                `json:"gender"`
	TradingPlatform       string                `json:"tradingPlatform"`
	DcName                string                `json:"dcName"`
	Address               string                `json:"address"`
	Zip                   string                `json:"zip"`
	RmDealerPhone         string                `json:"rmDealerPhone"`
	PermanentAddress      PermanentAddress      `json:"permanentAddress"`
	CorrespondenceAddress CorrespondenceAddress `json:"correspondenceAddress"`
}
type PermanentAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Pincode      string `json:"pincode"`
}
type CorrespondenceAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	Pincode      string `json:"pincode"`
}
type ActiveSegments struct {
	Equity     bool `json:"equity"`
	Futures    bool `json:"futures"`
	Currency   bool `json:"currency"`
	Commodity  bool `json:"commodity"`
	MutualFund bool `json:"mutualFund"`
	Ipo        bool `json:"ipo"`
	Nsecm      bool `json:"NSECM"`
	Bsecm      bool `json:"BSECM"`
	Nsefo      bool `json:"NSEFO"`
	Nsx        bool `json:"NSX"`
	Bsx        bool `json:"BSX"`
	Mcdx       bool `json:"MCDX"`
	Ncdx       bool `json:"NCDX"`
}
type BankDetails struct {
	BankName     string `json:"bankName"`
	BranchName   string `json:"branchName"`
	AccNO        string `json:"accNO"`
	IfscCode     string `json:"ifscCode"`
	MicrCode     string `json:"micrCode"`
	IsNetBanking bool   `json:"isNetBanking"`
	IsDefalutID  bool   `json:"isDefalutID"`
	BankLogo     string `json:"bankLogo"`
}
type DpDetails struct {
	DpIDNo      string `json:"dpIdNo"`
	BoAccountNo string `json:"boAccountNo"`
	DpType      string `json:"dpType"`
	DpName      string `json:"dpName"`
}
type BrokerageDetails struct {
	RevisedTariff    interface{} `json:"revisedTariff"`
	OldFormat        interface{} `json:"oldFormat"`
	Version          string      `json:"version"`
	BrokerageProduct string      `json:"brokerageProduct"`
	SchemeCode       string      `json:"schemeCode"`
}
type Nominee struct {
	Name                    string `json:"name"`
	RelationshipWithNominee string `json:"relationshipWithNominee"`
	NomineeDOB              string `json:"nomineeDOB"`
	GuardianName            string `json:"guardianName"`
	Address1                string `json:"address1"`
	Address2                string `json:"address2"`
	Address3                string `json:"address3"`
	City                    string `json:"city"`
	State                   string `json:"state"`
	Country                 string `json:"country"`
	Pincode                 string `json:"pincode"`
}

type PartnerDetailsResponse struct {
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	ErrorCode string         `json:"error_code"`
	Data      PartnerDetails `json:"data"`
}
type PartnerDetails struct {
	PartnerID          string             `json:"partnerId"`
	PartnerType        string             `json:"partnerType"`
	Active             bool               `json:"active"`
	ArnDetails         ArnDetails         `json:"arnDetails"`
	BankDetails        []BankDetails      `json:"bankDetails"`
	PartnerDetailsInfo PartnerDetailsInfo `json:"partnerDetails"`
	CountryCode        string             `json:"countryCode"`
	CreationDate       string             `json:"creationDate"`
	Mobile             string             `json:"mobile"`
	Pan                string             `json:"pan"`
	Email              string             `json:"email"`
	ModificationDate   string             `json:"modificationDate"`
	Category           string             `json:"category"`
}
type ArnDetails struct {
	Arn       string `json:"arn"`
	ValidFrom string `json:"validFrom"`
	ValidTo   string `json:"validTo"`
	Euin      string `json:"euin"`
}
type PartnerDetailsInfo struct {
	AccountType string `json:"accountType"`
	Branch      string `json:"branch"`
	FullName    string `json:"fullName"`
	Gender      string `json:"gender"`
	ParentTag   string `json:"parentTag"`
}

type PaymentHandler struct {
	AccountNo   string `json:"accountNo"`
	BankName    string `json:"bankName"`
	PaymentMode string `json:"paymentMode"`
}
