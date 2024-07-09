package model

type FAQResponse struct {
	Status string `json:"status" example:"success"`
	FAQs   []FAQ  `json:"faqs"`
}

type FAQ struct {
	ContentType string `json:"contentType" example:"json"`
	Question    string `json:"question" example:"When will my portfolio will get updated?"`
	Content     string `json:"content" example:"It will take 4 to 5 business days after realisation of funds"`
}
