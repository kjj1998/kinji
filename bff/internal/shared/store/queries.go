package store

const transactionsInDateRange = `
	SELECT id, user_id, date, merchant, category, amount, direction, notes, split
	FROM transactions
	WHERE user_id = ? AND date >= ? AND date <= ?
	ORDER BY date DESC`
