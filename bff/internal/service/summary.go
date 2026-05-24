package service

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/kjj1998/kinji/bff/internal/model"
	"github.com/kjj1998/kinji/bff/internal/repository"
)

const (
	monthLayout   = "2006-01"
	summaryMonths = 6
)

type SummaryService interface {
	GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*model.TransactionSummary, error)
}

type summaryService struct {
	repo repository.Repository
}

func NewSummaryService(repo repository.Repository) SummaryService {
	return &summaryService{repo: repo}
}

func (s *summaryService) GenerateMonthlySummary(
	ctx context.Context,
	userId string,
	month, year string,
) (*model.TransactionSummary, error) {
	curMonthTransactions, err := s.repo.GetMonthlyTransactions(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get current month transactions for %s-%s, user id %s: %w", month, year, userId, err)
	}

	totalIncome, totalSpent, netSavings, lastMonthSpent, err := s.repo.GetTotalIncomeTotalSpentAndNetSavings(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get total income, total spent, net savings for %s-%s: %w", month, year, err)
	}
	savingsRate := roundTo2Dp(safeDivide(netSavings.Value, totalIncome.Value) * 100)

	topMerchants, err := s.repo.GetMonthlyTopMerchants(ctx, userId, month, year, 5)
	if err != nil {
		return nil, fmt.Errorf("get top merchants for %s-%s: %w", month, year, err)
	}

	topCategories, err := s.repo.GetMonthlyTopCategories(ctx, userId, month, year, 5)
	if err != nil {
		return nil, fmt.Errorf("get top categories for %s-%s: %w", month, year, err)
	}

	monthlyExpenses, err := s.repo.GetLastSixMonthsExpenses(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get last six month expenses from %s-%s: %w", month, year, err)
	}
	monthlyTrend, err := buildMonthlyTrend(month, year, monthlyExpenses)
	if err != nil {
		return nil, err
	}

	changeInSpending := totalSpent.Value - lastMonthSpent
	topTransaction := getTopTransaction(curMonthTransactions.Transactions)
	monthlySummary := generateMonthlySummary(
		float64(changeInSpending), topTransaction, lastMonthSpent > 0, netSavings.Value, savingsRate)

	dailySpendngTrend := computeDailySpendingTrend(curMonthTransactions.Transactions)

	curCategorySpending, prevCategorySpending, err := s.repo.GetCategorySpendingForLastTwoMonths(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get category spending for %s-%s: %w", month, year, err)
	}
	categorySpendingChanges := computeCategoriesWithBiggestSpendingChange(
		curCategorySpending, prevCategorySpending)

	recentTransactions := recentTransactions(curMonthTransactions.Transactions, 5)

	return &model.TransactionSummary{
		TotalIncome:        totalIncome,
		TotalSpent:         totalSpent,
		NetSavings:         netSavings,
		LastMonthSpent:     lastMonthSpent,
		SavingsRate:        savingsRate,
		TopMerchants:       topMerchants,
		TopCategories:      topCategories,
		MonthlyTrend:       monthlyTrend,
		MonthlySummary:     monthlySummary,
		DailyTrend:         dailySpendngTrend,
		BiggestChanges:     categorySpendingChanges,
		RecentTransactions: recentTransactions,
	}, nil
}

func generateMonthlySummary(
	difference float64,
	topTransaction *model.Transaction,
	hasPrevMonth bool,
	netSavings int,
	savingsRate float64,
) string {
	if topTransaction == nil {
		return ""
	}

	suffix := fmt.Sprintf(
		"Your biggest expense was %s at $%.2f, and you saved $%.2f (%.2f%% of income).",
		topTransaction.Category,
		float64(topTransaction.Amount)/100,
		float64(netSavings)/100,
		savingsRate,
	)
	if !hasPrevMonth {
		return suffix
	}
	direction := "more"
	if difference < 0 {
		direction = "less"
	}
	return fmt.Sprintf("You spent %.0f%% %s than last month. ", math.Abs(roundTo2Dp(difference/100)), direction) + suffix
}

func computeDailySpendingTrend(txs []model.Transaction) []model.DateSpending {
	totals := make(map[time.Weekday]int)
	for _, t := range txs {
		if t.Direction == model.Inflow {
			continue
		}
		date, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			continue
		}
		totals[date.Weekday()] += t.Amount
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

func computeCategoriesWithBiggestSpendingChange(cur, prev map[model.Category]int) []model.CategorySpendingChange {
	categories := make(map[model.Category]struct{}, len(cur)+len(prev))
	for cat := range cur {
		categories[cat] = struct{}{}
	}
	for cat := range prev {
		categories[cat] = struct{}{}
	}

	result := make([]model.CategorySpendingChange, 0, len(categories))
	for cat := range categories {
		curAmount := cur[cat]
		prevAmount := prev[cat]
		result = append(result, model.CategorySpendingChange{
			Category:         cat,
			Amount:           curAmount,
			Change:           curAmount - prevAmount,
			PercentageChange: int(percentageChange(curAmount, prevAmount)),
			IsNew:            prevAmount == 0,
		})
	}

	if len(prev) == 0 {
		// No baseline: "biggest movers" is undefined, fall back to biggest spenders.
		sortByAmountDesc(result, func(c model.CategorySpendingChange) int {
			return c.Amount
		})
	} else {
		sortByAmountDesc(result, func(c model.CategorySpendingChange) int {
			if c.PercentageChange < 0 {
				return -c.PercentageChange
			}
			return c.PercentageChange
		})
	}

	if len(result) < 3 {
		return result
	}

	return result[:3] // return top 3 categories with biggest spending changes
}

func recentTransactions(txs []model.Transaction, n int) []model.Transaction {
	copy := slices.Clone(txs)

	slices.SortFunc(copy, func(a, b model.Transaction) int {
		return cmp.Compare(b.Date, a.Date)
	})

	if n > len(copy) {
		n = len(copy)
	}
	return copy[:n]
}

func buildMonthlyTrend(month, year string, monthlyExpenses map[string]int) ([]model.DateSpending, error) {
	toMonth, err := time.Parse("2006-01", year+"-"+month)
	if err != nil {
		return nil, fmt.Errorf("parse %s-%s: %w", year, month, err)
	}

	trend := make([]model.DateSpending, summaryMonths)
	for i := range summaryMonths {
		t := toMonth.AddDate(0, -(summaryMonths - 1 - i), 0)
		month := t.Format(monthLayout)
		trend[i] = model.DateSpending{
			Date:   t.Format("Jan"),
			Amount: monthlyExpenses[month],
		}
	}
	return trend, nil
}

func getTopTransaction(transactions []model.Transaction) *model.Transaction {
	if len(transactions) == 0 {
		return nil
	}

	var top *model.Transaction
	for i, t := range transactions {
		if t.Direction != model.Outflow {
			continue
		}
		if top == nil || t.Amount > top.Amount {
			top = &transactions[i]
		}
	}
	return top
}
