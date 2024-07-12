package model

type HomeInfo struct {
	MostBought []Plan  `json:"mostBought"`
	AllFDS     []Plan  `json:"allFDs"`
	Journey    Journey `json:"journey"`
	FAQs       []FAQ   `json:"faqs"`
}
