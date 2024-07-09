package entity

type WebhookEvent struct {
	ClientCode    string
	Vendor        string
	TrackingId    string
	EventType     string
	Institution   string
	Type          string
	Amount        float64
	TenureMonths  int
	TenureDays    int
	FailureReason string
	CreatedBy     string
	UpdatedBy     string
}
