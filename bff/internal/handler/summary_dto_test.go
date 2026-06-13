package handler

import (
	"testing"
	"time"

	"github.com/kjj1998/kinji/bff/internal/model"
)

func TestToTransactionSummary_LabelsAndTruncation(t *testing.T) {
	changes := make([]model.CategorySpendingChange, 5) // more than maxBiggestChanges
	for i := range changes {
		changes[i] = model.CategorySpendingChange{Category: model.CategoryFood, Amount: 100 - i}
	}
	recent := make([]model.Transaction, 7) // more than maxRecentTransactions
	for i := range recent {
		recent[i] = model.Transaction{ID: string(rune('a' + i))}
	}

	in := &model.MonthlySummary{
		TotalIncome:      model.ValueAndChange[int]{Value: 1000, Change: 100},
		SavingsRate:      43.4,
		LastMonthSpent:   500,
		SummaryStatement: "you did fine",
		DailyTrend: []model.DaySpending{
			{Weekday: time.Monday, Amount: 500},
			{Weekday: time.Sunday, Amount: 90},
		},
		MonthlyTrend: []model.MonthSpending{
			{Month: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), Amount: 0},
			{Month: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), Amount: 1000},
		},
		BiggestChanges:     changes,
		RecentTransactions: recent,
	}

	out := ToTransactionSummary(in)

	if out.MonthlySummary != "you did fine" {
		t.Errorf("SummaryStatement should map to MonthlySummary, got %q", out.MonthlySummary)
	}
	if out.TotalIncome.Value != 1000 || out.TotalIncome.Change != 100 {
		t.Errorf("ValueAndChange not mapped: %+v", out.TotalIncome)
	}

	// weekday/month labels rendered by the mapper
	if out.DailyTrend[0].Date != "Mon" || out.DailyTrend[1].Date != "Sun" {
		t.Errorf("daily labels wrong: %+v", out.DailyTrend)
	}
	if out.MonthlyTrend[0].Date != "Jan" || out.MonthlyTrend[1].Date != "Jun" {
		t.Errorf("monthly labels wrong: %+v", out.MonthlyTrend)
	}

	// view truncation applied by the mapper
	if len(out.BiggestChanges) != maxBiggestChanges {
		t.Errorf("expected %d biggest changes, got %d", maxBiggestChanges, len(out.BiggestChanges))
	}
	if len(out.RecentTransactions) != maxRecentTransactions {
		t.Errorf("expected %d recent transactions, got %d", maxRecentTransactions, len(out.RecentTransactions))
	}
}

func TestToTransactionSummary_Nil(t *testing.T) {
	if ToTransactionSummary(nil) != nil {
		t.Error("expected nil for nil input")
	}
}
