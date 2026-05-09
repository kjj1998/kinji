package repository

import (
	"fmt"
	"strings"

	"github.com/kohjunjie/kinji/bff/internal/model"
)

const userPrefix = "USER#"

type dynamoTransaction struct {
	User string         `dynamodbav:"user"`
	Item string         `dynamodbav:"item"`
	Data map[string]any `dynamodbav:"data"`
}

func extractField[T any](data map[string]any, key string) (T, error) {
	v, ok := data[key].(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("invalid or missing %s", key)
	}
	return v, nil
}

func (d dynamoTransaction) ConvertToTransaction() (model.Transaction, error) {
	userID, ok := strings.CutPrefix(d.User, userPrefix)
	if !ok {
		return model.Transaction{}, fmt.Errorf("invalid user key: %s", d.User)
	}

	id, err := extractField[string](d.Data, "id")
	if err != nil {
		return model.Transaction{}, err
	}
	date, err := extractField[string](d.Data, "date")
	if err != nil {
		return model.Transaction{}, err
	}
	merchant, err := extractField[string](d.Data, "merchant")
	if err != nil {
		return model.Transaction{}, err
	}
	category, err := extractField[string](d.Data, "category")
	if err != nil {
		return model.Transaction{}, err
	}
	amount, err := extractField[float64](d.Data, "amount")
	if err != nil {
		return model.Transaction{}, err
	}

	tx := model.Transaction{
		UserID:   userID,
		ID:       id,
		Date:     date,
		Merchant: merchant,
		Category: model.Category(category),
		Amount:   amount,
	}
	if notes, ok := d.Data["notes"].(string); ok {
		tx.Notes = &notes
	}
	if split, ok := d.Data["split"].(float64); ok {
		tx.Split = &split
	}
	return tx, nil
}
