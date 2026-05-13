package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/kohjunjie/kinji/bff/internal/model"
	"github.com/kohjunjie/kinji/bff/internal/repository"
)

const (
	monthLayout   = "2006-01"
	summaryMonths = 6
)

type SummaryService interface {
	GenerateMonthlySummary(ctx context.Context, userID string, from, to *time.Time) (*model.TransactionSummary, error)
}

type summaryService struct {
	repo repository.Repository
}

type monthMetrics struct {
	totalIncome float64
	totalSpent  float64
	netSavings  float64
	savingsRate float64
	byCategory  []model.CategorySpending
}

func NewSummaryService(repo repository.Repository) SummaryService {
	return &summaryService{repo: repo}
}

func (s *summaryService) GenerateMonthlySummary(ctx context.Context, userID string, from, to *time.Time) (*model.TransactionSummary, error) {
	if from == nil || to == nil {
		now := time.Now()
		t := now.AddDate(0, -(summaryMonths - 1), 0)
		from = &t
		to = &now
	}

	allTx, err := s.repo.ListRange(ctx, userID, from.Format(monthLayout), to.Format(monthLayout))
	if err != nil {
		return nil, err
	}

	transactionsByMonth := groupTransactionsByMonth(allTx)
	metrics := calcMetricsByMonth(transactionsByMonth)

	curMonthDateString := to.Format(monthLayout)
	prevMonthDateString := to.AddDate(0, -1, 0).Format(monthLayout)
	curMonthMetrics := metrics[curMonthDateString]
	prevMonthMetrics := metrics[prevMonthDateString]

	changeInSpending := percentageChange(curMonthMetrics.totalSpent, prevMonthMetrics.totalSpent)
	changeInNetSavings := percentageChange(curMonthMetrics.netSavings, prevMonthMetrics.netSavings)
	changeInSavingsRate := percentageChange(curMonthMetrics.savingsRate, prevMonthMetrics.savingsRate)

	topCategory := findTopSpendingCategory(curMonthMetrics.byCategory)

	monthlySummary := generateMonthlySummary(changeInSpending, len(transactionsByMonth[prevMonthDateString]) > 0, curMonthMetrics)
	trend := buildMonthlyTrend(to, metrics)
	spendingByDayOfWeek := computeSpendingByDayOfWeek(transactionsByMonth[curMonthDateString])
	categorySpendingChanges := computeBiggestSpendingChangeCategory(curMonthMetrics.byCategory, prevMonthMetrics.byCategory)
	topMerchants := findTopSpendingMerchants(transactionsByMonth[curMonthDateString])
	recentTransactions := recentTransactions(transactionsByMonth[curMonthDateString], 5)

	return &model.TransactionSummary{
		TotalIncome:        curMonthMetrics.totalIncome,
		TotalSpent:         model.ValueAndChange{Value: curMonthMetrics.totalSpent, Change: changeInSpending},
		NetSavings:         model.ValueAndChange{Value: curMonthMetrics.netSavings, Change: changeInNetSavings},
		SavingsRate:        model.ValueAndChange{Value: curMonthMetrics.savingsRate, Change: changeInSavingsRate},
		LastMonthSpent:     prevMonthMetrics.totalSpent,
		TopCategory:        topCategory,
		MonthlySummary:     monthlySummary,
		SpendingByCategory: curMonthMetrics.byCategory,
		MonthlyTrend:       trend,
		DailyTrend:         spendingByDayOfWeek,
		BiggestChanges:     categorySpendingChanges,
		TopMerchants:       topMerchants,
		RecentTransactions: recentTransactions,
	}, nil
}

func findTopSpendingCategory(categories []model.CategorySpending) *model.CategorySpending {
	var topCategory *model.CategorySpending
	if len(categories) > 0 {
		topCategory = &categories[0]
	}

	return topCategory
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

	if len(byCategory) > 5 {
		byCategory = byCategory[:5]
	}

	return monthMetrics{
		totalIncome: roundTo2Dp(income),
		totalSpent:  roundTo2Dp(spent),
		netSavings:  roundTo2Dp(net),
		savingsRate: roundTo2Dp(safeDivide(net, income) * 100),
		byCategory:  byCategory,
	}
}

func generateMonthlySummary(spendChange float64, hasPrevMonth bool, cur monthMetrics) string {
	if len(cur.byCategory) == 0 {
		return ""
	}
	suffix := fmt.Sprintf(
		"Your biggest expense was %s at $%.2f, and you saved $%.2f (%.2f%% of income).",
		cur.byCategory[0].Category, cur.byCategory[0].Amount, cur.netSavings, cur.savingsRate,
	)
	if !hasPrevMonth {
		return suffix
	}
	direction := "more"
	if spendChange < 0 {
		direction = "less"
	}
	return fmt.Sprintf("You spent %.0f%% %s than last month. ", math.Abs(spendChange), direction) + suffix
}

func computeSpendingByDayOfWeek(txs []model.Transaction) []model.DateSpending {
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
			Amount: roundTo2Dp(totals[day]),
		}
	}

	return result
}

func computeBiggestSpendingChangeCategory(cur, prev []model.CategorySpending) []model.CategorySpendingChange {
	prevMap := make(map[model.Category]float64, len(prev))
	for _, c := range prev {
		prevMap[c.Category] = c.Amount
	}

	result := make([]model.CategorySpendingChange, len(cur))
	for i, c := range cur {
		prevAmount := prevMap[c.Category]
		result[i] = model.CategorySpendingChange{
			Category:         c.Category,
			Amount:           roundTo2Dp(c.Amount),
			Change:           roundTo2Dp(c.Amount - prevAmount),
			PercentageChange: int(percentageChange(c.Amount, prevAmount)),
		}
	}

	sortByAmountDesc(result, func(c model.CategorySpendingChange) float64 { return math.Abs(float64(c.PercentageChange)) })

	if len(result) < 3 {
		return result
	}

	return result[:3]
}

func findTopSpendingMerchants(txs []model.Transaction) []model.Merchant {
	totals := make(map[string]model.CategorySpending)
	for _, t := range txs {
		if isExpense(t) {
			entry := totals[t.Merchant]
			entry.Category = t.Category
			entry.Amount += -t.Amount
			totals[t.Merchant] = entry
		}
	}

	result := make([]model.Merchant, 0, len(totals))
	for name, spending := range totals {
		result = append(result, model.Merchant{Name: name, Amount: roundTo2Dp(spending.Amount), Category: spending.Category})
	}
	sortByAmountDesc(result, func(m model.Merchant) float64 { return m.Amount })
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func recentTransactions(txs []model.Transaction, n int) []model.Transaction {
	if n > len(txs) {
		n = len(txs)
	}
	out := make([]model.Transaction, n)
	for i := range n {
		out[i] = txs[len(txs)-1-i]
	}
	return out
}

func buildMonthlyTrend(to *time.Time, metrics map[string]monthMetrics) []model.DateSpending {
	toMonth := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, to.Location())
	trend := make([]model.DateSpending, summaryMonths)
	for i := range summaryMonths {
		t := toMonth.AddDate(0, -(summaryMonths - 1 - i), 0)
		month := t.Format(monthLayout)
		trend[i] = model.DateSpending{
			Date:   t.Format("Jan"),
			Amount: roundTo2Dp(metrics[month].totalSpent),
		}
	}
	return trend
}

func calcMetricsByMonth(transactionsByMonth map[string][]model.Transaction) map[string]monthMetrics {
	metrics := make(map[string]monthMetrics, len(transactionsByMonth))
	for month, txs := range transactionsByMonth {
		metrics[month] = calcMonthMetrics(txs)
	}
	return metrics
}

func groupTransactionsByMonth(txs []model.Transaction) map[string][]model.Transaction {
	byMonth := make(map[string][]model.Transaction)
	for _, t := range txs {
		if len(t.Date) < 7 {
			continue
		}
		month := t.Date[:7]
		byMonth[month] = append(byMonth[month], t)
	}
	return byMonth
}

func isExpense(t model.Transaction) bool {
	return t.Category != model.CategoryIncome && t.Amount < 0
}
