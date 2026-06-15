package store

const getTopSpendingMerchantsWithinDateRange = `
	SELECT merchant, SUM(amount) AS total, category
	FROM transactions
	WHERE user_id = ? AND date >= ? AND date <= ? AND direction = 'OUTFLOW'
	GROUP BY merchant
	ORDER BY total DESC
	LIMIT ?`

const getTotalIncomeTotalSpentAndNetSavingsForTwoMonths = `
	SELECT
		strftime('%Y-%m', date) AS month,
		COALESCE(SUM(CASE WHEN category = 'Income'                           THEN amount END), 0) AS total_income,
		COALESCE(SUM(CASE WHEN direction = 'OUTFLOW' AND category != 'Income' THEN amount END), 0) AS total_spent,
		COALESCE(SUM(CASE WHEN category = 'Income'                           THEN amount END), 0)
		- COALESCE(SUM(CASE WHEN direction = 'OUTFLOW' AND category != 'Income' THEN amount END), 0) AS net_savings
	FROM transactions
	WHERE user_id = ? AND strftime('%Y-%m', date) IN (?, ?)
	GROUP BY month
	ORDER BY month DESC;`

const getCategorySpendingForTwoMonths = `
	SELECT strftime('%Y-%m', date) AS month, category, SUM(amount) AS total
	FROM transactions
	WHERE user_id = ? AND strftime('%Y-%m', date) IN (?, ?)
		AND direction = 'OUTFLOW' AND category != 'Income'
	GROUP BY month, category;`

const getTopSpendingCategoriesWithinDateRange = `
	SELECT category, SUM(amount) as total
	FROM transactions
	WHERE user_id = ? AND date >= ? AND date <= ? AND direction = 'OUTFLOW'
	GROUP by category
	ORDER BY total DESC
	LIMIT ?`

const getTotalMonthlyExpensesWithinDateRange = `
	SELECT strftime('%Y-%m', date) AS month, SUM(amount) AS total
	FROM transactions
	WHERE user_id = ? AND date >= ? AND date <= ? AND direction = 'OUTFLOW'
	GROUP BY month
	ORDER BY month ASC`
