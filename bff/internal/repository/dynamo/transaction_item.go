package dynamo

import (
	"fmt"
	"strings"

	"github.com/kohjunjie/kinji/bff/internal/model"
)

const userPrefix = "USER#"

type dynamoTransactionData struct {
	ID        string `dynamodbav:"id"`
	Date      string `dynamodbav:"date"`
	Merchant  string `dynamodbav:"merchant"`
	Category  string `dynamodbav:"category"`
	Direction string `dynamodbav:"direction"`
	Amount    int    `dynamodbav:"amount"`
	Notes     string `dynamodbav:"notes"`
	Split     int    `dynamodbav:"split"`
}

type dynamoTransaction struct {
	Key  string                `dynamodbav:"key"`
	Item string                `dynamodbav:"item"`
	Data dynamoTransactionData `dynamodbav:"data"`
}

func (d dynamoTransaction) ConvertToTransaction() (model.Transaction, error) {
	userId, ok := strings.CutPrefix(d.Key, userPrefix)
	if !ok {
		return model.Transaction{}, fmt.Errorf("invalid key: %s", d.Key)
	}

	return model.Transaction{
		UserID:    userId,
		ID:        d.Data.ID,
		Date:      d.Data.Date,
		Merchant:  d.Data.Merchant,
		Category:  model.Category(d.Data.Category),
		Amount:    d.Data.Amount,
		Direction: model.Direction(d.Data.Direction),
		Notes:     d.Data.Notes,
		Split:     d.Data.Split,
	}, nil
}
