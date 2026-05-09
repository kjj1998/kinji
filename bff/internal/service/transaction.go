package service

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/kohjunjie/kinji/bff/internal/model"
	"github.com/kohjunjie/kinji/bff/internal/repository"
)

type TransactionService interface {
	Summary(ctx context.Context, userID string) (*model.TransactionSummary, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

type monthMetrics struct {
	totalIncome float64
	totalSpent  float64
	netSavings  float64
	savingsRate float64
	byCategory  []model.CategorySpending
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Summary(ctx context.Context, userID string) (*model.TransactionSummary, error) {
	now := time.Now()
	to := now.Format("2006-01")
	from := now.AddDate(0, -5, 0).Format("2006-01")

	allTx, err := s.repo.ListRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	byMonth := make(map[string][]model.Transaction)
	for _, t := range allTx {
		month := t.Date[:7]
		byMonth[month] = append(byMonth[month], t)
	}

	metrics := make(map[string]monthMetrics, len(byMonth))
	for month, txs := range byMonth {
		metrics[month] = calcMonthMetrics(txs)
	}

	curMonth := now.Format("2006-01")
	prevMonth := now.AddDate(0, -1, 0).Format("2006-01")

	cur := metrics[curMonth]
	prev := metrics[prevMonth]

	spendChange := percentageChange(cur.totalSpent, prev.totalSpent)
	var delta *float64
	if len(byMonth[prevMonth]) > 0 {
		delta = &spendChange
	}
	var summary string
	if len(cur.byCategory) > 0 {
		summary = monthlySummary(delta, cur.byCategory[0].Category, cur.byCategory[0].Amount, cur.netSavings, cur.savingsRate)
	}

	trend := make([]model.DateSpending, 6)
	for i := range 6 {
		month := now.AddDate(0, -(5 - i), 0).Format("2006-01")
		trend[i] = model.DateSpending{
			Date:   month,
			Amount: metrics[month].totalSpent,
		}
	}

	spendingByDayOfWeek := spendingByDayOfWeek(byMonth[curMonth])
	categorySpendingChanges := categorySpendingChanges(cur.byCategory, prev.byCategory)
	topMerchants := topMerchants(byMonth[curMonth])
	recent := byMonth[curMonth]
	slices.Reverse(recent)

	return &model.TransactionSummary{
		TotalIncome:         cur.totalIncome,
		TotalSpent:          model.ValueAndChange{Value: cur.totalSpent, Change: spendChange},
		NetSavings:          model.ValueAndChange{Value: cur.netSavings, Change: percentageChange(cur.netSavings, prev.netSavings)},
		SavingsRate:         model.ValueAndChange{Value: cur.savingsRate, Change: percentageChange(cur.savingsRate, prev.savingsRate)},
		MonthlySummary:      summary,
		SpendingByCategory:  cur.byCategory,
		MonthlyTrend:        trend,
		SpendingByDayOfWeek: spendingByDayOfWeek,
		BiggestChanges:      categorySpendingChanges,
		TopMerchants:        topMerchants,
		RecentTransactions:  recent,
	}, nil
}

func calcMonthMetrics(txs []model.Transaction) monthMetrics {
	var income, spent float64
	catTotals := make(map[model.Category]float64)
	for _, t := range txs {
		if t.Category == model.CategoryIncome {
			income += t.Amount
		} else if isExpense(t) {
			spent += -t.Amount
			catTotals[t.Category] += -t.Amount
		}
	}
	net := income - spent

	byCategory := make([]model.CategorySpending, 0, len(catTotals))
	for cat, amount := range catTotals {
		byCategory = append(byCategory, model.CategorySpending{Category: cat, Amount: amount})
	}
	sortByAmountDesc(byCategory, func(c model.CategorySpending) float64 { return c.Amount })

	return monthMetrics{
		totalIncome: income,
		totalSpent:  spent,
		netSavings:  net,
		savingsRate: safeDivide(net, income) * 100,
		byCategory:  byCategory,
	}
}

// monthlySummary generates a human-readable summary sentence.
// spendDelta is nil when there is no previous month to compare against.
func monthlySummary(spendDelta *float64, topCategory model.Category, topCategoryAmount, netSavings, savingsRate float64) string {
	s := ""
	if spendDelta != nil {
		direction := "more"
		if *spendDelta < 0 {
			direction = "less"
		}
		s = fmt.Sprintf("You spent %.0f%% %s than last month. ", math.Abs(*spendDelta), direction)
	}
	return s + fmt.Sprintf(
		"Your biggest expense was %s at %s, and you saved %s (%.0f%% of income).",
		topCategory, formatAmount(topCategoryAmount), formatAmount(netSavings), savingsRate,
	)
}

func spendingByDayOfWeek(txs []model.Transaction) []model.DateSpending {
	totals := make(map[time.Weekday]float64)
	for _, t := range txs {
		if !isExpense(t) {
			continue
		}
		date, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			continue
		}
		totals[date.Weekday()] += -t.Amount
	}

	days := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday,
		time.Thursday, time.Friday, time.Saturday, time.Sunday,
	}

	result := make([]model.DateSpending, len(days))
	for i, day := range days {
		result[i] = model.DateSpending{
			Date:   day.String()[:3],
			Amount: totals[day],
		}
	}

	return result
}

func categorySpendingChanges(cur, prev []model.CategorySpending) []model.CategorySpendingChange {
	prevMap := make(map[model.Category]float64, len(prev))
	for _, c := range prev {
		prevMap[c.Category] = c.Amount
	}

	result := make([]model.CategorySpendingChange, len(cur))
	for i, c := range cur {
		prevAmount := prevMap[c.Category]
		result[i] = model.CategorySpendingChange{
			Category:         c.Category,
			Amount:           c.Amount,
			Change:           c.Amount - prevAmount,
			PercentageChange: int(percentageChange(c.Amount, prevAmount)),
		}
	}
	return result
}

func topMerchants(txs []model.Transaction) []model.Merchant {
	totals := make(map[string]float64)
	for _, t := range txs {
		if isExpense(t) {
			totals[t.Merchant] += -t.Amount
		}
	}

	result := make([]model.Merchant, 0, len(totals))
	for name, amount := range totals {
		result = append(result, model.Merchant{Name: name, Amount: amount})
	}
	sortByAmountDesc(result, func(m model.Merchant) float64 { return m.Amount })
	return result
}

func formatAmount(amount float64) string {
	digits := fmt.Sprintf("%d", int(math.Round(amount)))
	n := len(digits)
	var b strings.Builder
	b.Grow(1 + n + (n-1)/3)
	b.WriteByte('$')
	for i, ch := range digits {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	return b.String()
}

func percentageChange(current, previous float64) float64 {
	return safeDivide(current-previous, previous) * 100
}

func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func isExpense(t model.Transaction) bool {
	return t.Category != model.CategoryIncome && t.Amount < 0
}

func sortByAmountDesc[T any](s []T, amount func(T) float64) {
	slices.SortFunc(s, func(a, b T) int {
		if amount(b) > amount(a) {
			return 1
		}
		if amount(b) < amount(a) {
			return -1
		}
		return 0
	})
}
