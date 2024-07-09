package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var webhookService WebhookService

func init() {
	webhookService = &webhookServiceImpl{}
}

type pattern struct {
	Pattern string
	Days    int
	Months  int
}

func TestExtractTenure(t *testing.T) {
	patterns := []pattern{{Pattern: "3m3d", Days: 3, Months: 3}, {Pattern: "0m0d", Days: 0, Months: 0}, {Pattern: "0m548d", Days: 0, Months: 18}, {Pattern: "0m730d", Days: 0, Months: 24}, {Pattern: "0m1826d", Days: 0, Months: 60}, {Pattern: "40m20d", Days: 20, Months: 40}, {Pattern: "11m11d", Days: 11, Months: 11}, {Pattern: "22m0d", Days: 0, Months: 22}}
	for _, p := range patterns {
		t.Logf("tesitng pattern: %s", p.Pattern)
		months, days := webhookService.extractTenure(p.Pattern)
		assert.Equal(t, p.Months, months, "Invalid Months")
		assert.Equal(t, p.Days, days, "Invalid Days")
	}
}
