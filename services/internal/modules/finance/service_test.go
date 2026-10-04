package finance

import (
	"testing"
	"time"
)

func TestAmountForMonth(t *testing.T) {
	amounts := []RecurringAmount{
		{EffectiveMonth: "2026-01", Amount: 5000},
		{EffectiveMonth: "2026-07", Amount: 5500},
	}
	cases := []struct {
		month string
		want  float64
		ok    bool
	}{
		{"2025-12", 0, false},
		{"2026-01", 5000, true},
		{"2026-06", 5000, true},
		{"2026-07", 5500, true},
		{"2027-03", 5500, true},
	}
	for _, c := range cases {
		got, ok := amountForMonth(amounts, c.month)
		if ok != c.ok || got.Amount != c.want {
			t.Errorf("amountForMonth(%s) = %v, %v; want %v, %v", c.month, got.Amount, ok, c.want, c.ok)
		}
	}
}

func TestActiveInMonth(t *testing.T) {
	end := "2026-06"
	r := Recurring{StartMonth: "2026-03", EndMonth: &end}
	for month, want := range map[string]bool{
		"2026-02": false, "2026-03": true, "2026-06": true, "2026-07": false,
	} {
		if got := activeInMonth(r, month); got != want {
			t.Errorf("activeInMonth(%s) = %v, want %v", month, got, want)
		}
	}
	open := Recurring{StartMonth: "2026-03"}
	if !activeInMonth(open, "2030-01") {
		t.Error("open-ended item should stay active")
	}
}

func TestClampDay(t *testing.T) {
	feb := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if got := clampDay(feb, 31); got.Day() != 28 {
		t.Errorf("clampDay(feb, 31) day = %d, want 28", got.Day())
	}
	oct := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	if got := clampDay(oct, 15); got.Day() != 15 || got.Month() != time.October {
		t.Errorf("clampDay(oct, 15) = %v", got)
	}
}

func TestFillDaily(t *testing.T) {
	month := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	upto := time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)
	got := fillDaily(month, upto, []DailyTotal{
		{Date: "2026-09-01", Amount: 10},
		{Date: "2026-09-03", Amount: 5},
		{Date: "2026-09-20", Amount: 99},
	})
	wantCumulative := []float64{10, 10, 15, 15}
	if len(got) != len(wantCumulative) {
		t.Fatalf("len = %d, want %d", len(got), len(wantCumulative))
	}
	for i, w := range wantCumulative {
		if got[i].Cumulative != w {
			t.Errorf("day %d cumulative = %v, want %v", i+1, got[i].Cumulative, w)
		}
	}
}

func TestWithDerived(t *testing.T) {
	got := withDerived(Totals{Fixed: 1808.04, Spend: 816.86, NonSalary: 182.76})
	if got.Gross < 2624.89 || got.Gross > 2624.91 {
		t.Errorf("gross = %v", got.Gross)
	}
	if got.Net < 2442.13 || got.Net > 2442.15 {
		t.Errorf("net = %v", got.Net)
	}
}

func TestParseMonth(t *testing.T) {
	if _, err := parseMonth("month", "2026-09"); err != nil {
		t.Errorf("valid month rejected: %v", err)
	}
	for _, bad := range []string{"", "2026-13", "2026-9-1", "sept"} {
		if _, err := parseMonth("month", bad); err == nil {
			t.Errorf("parseMonth(%q) should fail", bad)
		}
	}
}
