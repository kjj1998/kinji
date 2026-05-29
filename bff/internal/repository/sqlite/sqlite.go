package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kjj1998/kinji/bff/internal/models"
	_ "modernc.org/sqlite"
)

type Repository struct {
	client *sql.DB
}

func NewClient(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&journal_mode=WAL"

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, fmt.Errorf("opening sqlite database %q: %w", path, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging sqlite: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("applying schema: %w", err)
	}

	return db, nil
}

func NewRepository(client *sql.DB) *Repository {
	return &Repository{client: client}
}

func (d *Repository) getTransactionsWithinDateRange(
	ctx context.Context,
	userId, from, to string,
) ([]models.Transaction, error) {
	rows, err := d.client.QueryContext(ctx, getAllTransactionsWithinDateRange, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying transactions userId %q from %q to %q: %w", userId, from, to, err)
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Date, &t.Merchant,
			&t.Category, &t.Amount, &t.Direction, &t.Notes, &t.Split); err != nil {
			return nil, fmt.Errorf("scanning transaction: %w", err)
		}
		txs = append(txs, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating transaction rows: %w", err)
	}

	return txs, nil
}

func (d *Repository) getAvailabilites(
	ctx context.Context,
	userId string,
) ([]models.TransactionsAvailability, error) {
	rows, err := d.client.QueryContext(ctx, getMonthAndYearWhichTransactionsOccur, userId)
	if err != nil {
		return nil, fmt.Errorf("getting availabilities: %w", err)
	}
	defer rows.Close()

	byYear := map[int]*models.TransactionsAvailability{}
	var order []int
	for rows.Next() {
		var y, m int
		if err := rows.Scan(&y, &m); err != nil {
			return nil, fmt.Errorf("scanning availabilities: %w", err)
		}
		if _, ok := byYear[y]; !ok {
			byYear[y] = &models.TransactionsAvailability{Year: y}
			order = append(order, y)
		}
		byYear[y].Months = append(byYear[y].Months, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating availabilities rows: %w", err)
	}

	out := make([]models.TransactionsAvailability, 0, len(order))
	for _, y := range order {
		out = append(out, *byYear[y])
	}
	return out, nil
}

func (d *Repository) GetMonthlyTransactions(
	ctx context.Context,
	userId, month, year string,
) (models.Transactions, error) {
	from, to := GetMonthRangeDateStrings(month, year)
	txs, err := d.getTransactionsWithinDateRange(ctx, userId, from, to)
	if err != nil {
		return models.Transactions{}, fmt.Errorf("querying monthly transactions: %w", err)
	}

	avail, err := d.getAvailabilites(ctx, userId)
	if err != nil {
		return models.Transactions{}, fmt.Errorf("getting availabilities: %w", err)
	}

	return models.Transactions{Transactions: txs, Availabilities: avail}, nil
}

func (d *Repository) GetMonthlyTopMerchants(
	ctx context.Context,
	userId, month, year string,
	limit int,
) ([]models.Merchant, error) {
	from, to := GetMonthRangeDateStrings(month, year)
	rows, err := d.client.QueryContext(ctx, getTopSpendingMerchantsWithinDateRange, userId, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("querying top %d merchants: %w", limit, err)
	}
	defer rows.Close()

	var merchants []models.Merchant
	for rows.Next() {
		var m models.Merchant
		if err := rows.Scan(&m.Name, &m.Amount, &m.Category); err != nil {
			return nil, fmt.Errorf("scanning merchant: %w", err)
		}
		merchants = append(merchants, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating merchant rows: %w", err)
	}

	return merchants, nil
}

func (d *Repository) GetTotalIncomeTotalSpentAndNetSavings(
	ctx context.Context,
	userId, month, year string,
) (models.ValueAndChange[int], models.ValueAndChange[int], models.ValueAndChange[int], int, error) {
	curMonth, prevMonth, err := currentAndPreviousMonth(month, year)
	if err != nil {
		return models.ValueAndChange[int]{}, models.ValueAndChange[int]{}, models.ValueAndChange[int]{},
			0, fmt.Errorf("computing current and previous month: %w", err)
	}
	rows, err := d.client.QueryContext(
		ctx, getTotalIncomeTotalSpentAndNetSavingsForTwoMonths, userId, curMonth, prevMonth)
	if err != nil {
		return models.ValueAndChange[int]{}, models.ValueAndChange[int]{}, models.ValueAndChange[int]{},
			0, fmt.Errorf("querying income/spent/savings: %w", err)
	}
	defer rows.Close()

	type monthTotals struct {
		income, spend, saving int
	}
	totalsByMonth := make(map[string]monthTotals, 2)
	for rows.Next() {
		var month string
		var t monthTotals
		if err := rows.Scan(&month, &t.income, &t.spend, &t.saving); err != nil {
			return models.ValueAndChange[int]{}, models.ValueAndChange[int]{}, models.ValueAndChange[int]{},
				0, fmt.Errorf("scanning income/spent/savings: %w", err)
		}
		totalsByMonth[month] = t
	}

	if err := rows.Err(); err != nil {
		return models.ValueAndChange[int]{}, models.ValueAndChange[int]{}, models.ValueAndChange[int]{},
			0, fmt.Errorf("iterating income/spent/savings rows: %w", err)
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

	return models.NewValueAndChange(incomes),
		models.NewValueAndChange(spendings),
		models.NewValueAndChange(savings),
		prev.spend,
		nil
}

func (d *Repository) GetCategorySpendingForLastTwoMonths(
	ctx context.Context,
	userId, month, year string,
) (map[models.Category]int, map[models.Category]int, error) {
	curMonth, prevMonth, err := currentAndPreviousMonth(month, year)
	if err != nil {
		return nil, nil, fmt.Errorf("computing current and previous month: %w", err)
	}

	rows, err := d.client.QueryContext(ctx, getCategorySpendingForTwoMonths, userId, curMonth, prevMonth)
	if err != nil {
		return nil, nil, fmt.Errorf("querying category spending: %w", err)
	}
	defer rows.Close()

	cur := make(map[models.Category]int)
	prev := make(map[models.Category]int)
	for rows.Next() {
		var monthKey string
		var category models.Category
		var total int
		if err := rows.Scan(&monthKey, &category, &total); err != nil {
			return nil, nil, fmt.Errorf("scanning category spending: %w", err)
		}
		if monthKey == prevMonth {
			prev[category] = total
		} else {
			cur[category] = total
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterating category spending rows: %w", err)
	}

	return cur, prev, nil
}

func (d *Repository) GetMonthlyTopCategories(
	ctx context.Context,
	userId, month, year string,
	limit int,
) ([]models.CategorySpending, error) {
	from, to := GetMonthRangeDateStrings(month, year)
	rows, err := d.client.QueryContext(ctx, getTopSpendingCategoriesWithinDateRange, userId, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("querying top %d categories: %w", limit, err)
	}
	defer rows.Close()

	var categorySpendings []models.CategorySpending
	for rows.Next() {
		var cs models.CategorySpending
		if err := rows.Scan(&cs.Category, &cs.Amount); err != nil {
			return nil, fmt.Errorf("scanning category: %w", err)
		}
		categorySpendings = append(categorySpendings, cs)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating category rows: %w", err)
	}

	return categorySpendings, nil
}

func (d *Repository) GetLastSixMonthsExpenses(
	ctx context.Context,
	userId, month, year string,
) (map[string]int, error) {
	_, to := GetMonthRangeDateStrings(month, year)

	t, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, fmt.Errorf("parsing string %s into datetime: %w", to, err)
	}
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	from := firstOfMonth.AddDate(0, -5, 0).Format("2006-01-02")

	rows, err := d.client.QueryContext(ctx, getTotalMonthlyExpensesWithinDateRange, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying total monthly expenses between %s and %s: %w", from, to, err)
	}
	defer rows.Close()

	totalsByMonth := make(map[string]int)
	for rows.Next() {
		var month string
		var amount int
		if err := rows.Scan(&month, &amount); err != nil {
			return nil, fmt.Errorf("scanning monthly expense: %w", err)
		}
		totalsByMonth[month] = amount
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating monthly expense rows: %w", err)
	}

	return totalsByMonth, nil
}
