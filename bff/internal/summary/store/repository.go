package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
	sharedstore "github.com/kjj1998/kinji/bff/internal/shared/store"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
	"github.com/kjj1998/kinji/bff/internal/summary/service"
)

// Repository is the sqlite-backed implementation of the summary repository.
type Repository struct {
	client *sql.DB
}

// compile-time check that Repository satisfies the application port.
var _ service.TransactionRepository = (*Repository)(nil)

func NewRepository(client *sql.DB) *Repository {
	return &Repository{client: client}
}

func (d *Repository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	from, to := shared.GetMonthRangeDateStrings(month, year)
	return sharedstore.TransactionsInDateRange(ctx, d.client, userId, from, to)
}

func (d *Repository) GetMonthlyTopMerchants(
	ctx context.Context,
	userId, month, year string,
	limit int,
) ([]domain.MerchantSpending, error) {
	from, to := shared.GetMonthRangeDateStrings(month, year)
	merchants, err := database.QueryRows(ctx, d.client, getTopSpendingMerchantsWithinDateRange,
		func(r *sql.Rows) (domain.MerchantSpending, error) {
			var m domain.MerchantSpending
			return m, r.Scan(&m.Name, &m.Amount, &m.Category)
		}, userId, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("querying top %d merchants: %w", limit, err)
	}
	return merchants, nil
}

func (d *Repository) GetTotalIncomeTotalSpentAndNetSavings(
	ctx context.Context,
	userId, month, year string,
) (domain.ValueAndChange[int], domain.ValueAndChange[int], domain.ValueAndChange[int], int, error) {
	zero := domain.ValueAndChange[int]{}
	curMonth, prevMonth, err := currentAndPreviousMonth(month, year)
	if err != nil {
		return zero, zero, zero, 0, err
	}

	type monthTotals struct {
		month                 string
		income, spend, saving int
	}
	rows, err := database.QueryRows(ctx, d.client, getTotalIncomeTotalSpentAndNetSavingsForTwoMonths,
		func(r *sql.Rows) (monthTotals, error) {
			var t monthTotals
			return t, r.Scan(&t.month, &t.income, &t.spend, &t.saving)
		}, userId, curMonth, prevMonth)
	if err != nil {
		return zero, zero, zero, 0, fmt.Errorf("querying income/spent/savings: %w", err)
	}

	totalsByMonth := make(map[string]monthTotals, len(rows))
	for _, t := range rows {
		totalsByMonth[t.month] = t
	}

	cur := totalsByMonth[curMonth]
	prev, hasPrev := totalsByMonth[prevMonth]

	incomes := []int{cur.income}
	spendings := []int{cur.spend}
	savings := []int{cur.saving}
	if hasPrev {
		incomes = append(incomes, prev.income)
		spendings = append(spendings, prev.spend)
		savings = append(savings, prev.saving)
	}

	return domain.NewValueAndChange(incomes),
		domain.NewValueAndChange(spendings),
		domain.NewValueAndChange(savings),
		prev.spend,
		nil
}

func (d *Repository) GetCategorySpendingForLastTwoMonths(
	ctx context.Context,
	userId, month, year string,
) (map[shared.Category]int, map[shared.Category]int, error) {
	curMonth, prevMonth, err := currentAndPreviousMonth(month, year)
	if err != nil {
		return nil, nil, err
	}

	type categoryRow struct {
		month    string
		category shared.Category
		total    int
	}
	rows, err := database.QueryRows(ctx, d.client, getCategorySpendingForTwoMonths,
		func(r *sql.Rows) (categoryRow, error) {
			var c categoryRow
			return c, r.Scan(&c.month, &c.category, &c.total)
		}, userId, curMonth, prevMonth)
	if err != nil {
		return nil, nil, fmt.Errorf("querying category spending: %w", err)
	}

	cur := make(map[shared.Category]int)
	prev := make(map[shared.Category]int)
	for _, c := range rows {
		if c.month == prevMonth {
			prev[c.category] = c.total
		} else {
			cur[c.category] = c.total
		}
	}

	return cur, prev, nil
}

func (d *Repository) GetMonthlyTopCategories(
	ctx context.Context,
	userId, month, year string,
	limit int,
) ([]domain.CategorySpending, error) {
	from, to := shared.GetMonthRangeDateStrings(month, year)
	categorySpendings, err := database.QueryRows(ctx, d.client, getTopSpendingCategoriesWithinDateRange,
		func(r *sql.Rows) (domain.CategorySpending, error) {
			var cs domain.CategorySpending
			return cs, r.Scan(&cs.Category, &cs.Amount)
		}, userId, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("querying top %d categories: %w", limit, err)
	}
	return categorySpendings, nil
}

func (d *Repository) GetLastSixMonthsExpenses(
	ctx context.Context,
	userId, month, year string,
) (map[string]int, error) {
	_, to := shared.GetMonthRangeDateStrings(month, year)

	t, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, fmt.Errorf("parsing string %s into datetime: %w", to, err)
	}
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	from := firstOfMonth.AddDate(0, -5, 0).Format("2006-01-02")

	type monthExpense struct {
		month  string
		amount int
	}
	rows, err := database.QueryRows(ctx, d.client, getTotalMonthlyExpensesWithinDateRange,
		func(r *sql.Rows) (monthExpense, error) {
			var m monthExpense
			return m, r.Scan(&m.month, &m.amount)
		}, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying total monthly expenses between %s and %s: %w", from, to, err)
	}

	totalsByMonth := make(map[string]int, len(rows))
	for _, m := range rows {
		totalsByMonth[m.month] = m.amount
	}

	return totalsByMonth, nil
}
