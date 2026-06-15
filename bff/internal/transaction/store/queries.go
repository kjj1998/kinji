package store

const getMonthAndYearWhichTransactionsOccur = `
	SELECT CAST(substr(date,1,4) AS INTEGER), CAST(substr(date,6,2) AS INTEGER)
	FROM transactions
	WHERE user_id = ?
	GROUP BY 1, 2
	ORDER BY 1, 2`

const saveTransactions = `
	INSERT INTO transactions (id, user_id, date, merchant, category, amount, direction, notes, split)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
