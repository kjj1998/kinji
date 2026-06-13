package app

import (
	"context"
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/model"
)

// summaryTopN is the number of top merchants/categories fetched from the
// repository for the monthly summary.
const summaryTopN = 5

// SummaryService is the use-case API for a user's monthly spending summary.
type SummaryService interface {
	GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*model.MonthlySummary, error)
}

type summaryService struct {
	repo TransactionRepository
	calc model.SummaryCalculator
}

func NewSummaryService(repo TransactionRepository) SummaryService {
	return &summaryService{repo: repo, calc: model.NewSummaryCalculator()}
}

// GenerateMonthlySummary gathers the raw monthly data from the repository and
// hands it to the domain SummaryCalculator. It performs no calculation itself.
func (s *summaryService) GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*model.MonthlySummary, error) {
	currentMonth, err := s.repo.GetMonthlyTransactions(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get current month transactions for %s-%s, user id %s: %w", month, year, userId, err)
	}

	totalIncome, totalSpent, netSavings, lastMonthSpent, err := s.repo.GetTotalIncomeTotalSpentAndNetSavings(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get total income, total spent, net savings for %s-%s: %w", month, year, err)
	}

	topMerchants, err := s.repo.GetMonthlyTopMerchants(ctx, userId, month, year, summaryTopN)
	if err != nil {
		return nil, fmt.Errorf("get top merchants for %s-%s: %w", month, year, err)
	}

	topCategories, err := s.repo.GetMonthlyTopCategories(ctx, userId, month, year, summaryTopN)
	if err != nil {
		return nil, fmt.Errorf("get top categories for %s-%s: %w", month, year, err)
	}

	monthlyExpenses, err := s.repo.GetLastSixMonthsExpenses(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get last six month expenses from %s-%s: %w", month, year, err)
	}

	curCategorySpending, prevCategorySpending, err := s.repo.GetCategorySpendingForLastTwoMonths(ctx, userId, month, year)
	if err != nil {
		return nil, fmt.Errorf("get category spending for %s-%s: %w", month, year, err)
	}

	return s.calc.Calculate(model.SummaryInput{
		Month:                month,
		Year:                 year,
		CurrentMonth:         currentMonth,
		TotalIncome:          totalIncome,
		TotalSpent:           totalSpent,
		NetSavings:           netSavings,
		LastMonthSpent:       lastMonthSpent,
		TopMerchants:         topMerchants,
		TopCategories:        topCategories,
		MonthlyExpenses:      monthlyExpenses,
		CurCategorySpending:  curCategorySpending,
		PrevCategorySpending: prevCategorySpending,
	})
}
