package finance

import (
	"errors"
	"github.com/hungphan1911/tapestry/services/internal/core"
	"math"
	"testing"
)

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		in      string
		want    float64
		wantErr bool
	}{
		{"12", 12, false},
		{"=12.5+3+4", 19.5, false},
		{" 10 + 2.5 * 2 ", 15, false},
		{"(1+2)*3", 9, false},
		{"-5+2", -3, false},
		{"10/4", 2.5, false},
		{"2*-3", -6, false},
		{"", 0, true},
		{"1/0", 0, true},
		{"1+", 0, true},
		{"(1+2", 0, true},
		{"abc", 0, true},
		{"1 2", 0, true},
		{"1..2", 0, true},
	}
	for _, tt := range tests {
		got, err := EvaluateExpression(tt.in)
		if tt.wantErr {
			if err == nil || !errors.Is(err, core.ErrInvalidInput) {
				t.Errorf("%q: expected core.ErrInvalidInput, got %v", tt.in, err)
			}
			continue
		}
		if err != nil || math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("%q: got %v, %v; want %v", tt.in, got, err, tt.want)
		}
	}
}
