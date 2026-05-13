package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kohjunjie/kinji/bff/internal/model"
)

type Repository struct {
	client *dynamodb.Client
	table  string
}

func NewClient(endpoint, region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return client, nil
}

func NewRepository(client *dynamodb.Client, table string) *Repository {
	return &Repository{client: client, table: table}
}

func (d *Repository) queryTransactions(ctx context.Context, keyEx expression.KeyConditionBuilder) ([]model.Transaction, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression: %w", err)
	}

	paginator := dynamodb.NewQueryPaginator(d.client, &dynamodb.QueryInput{
		TableName:                 aws.String(d.table),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	var transactions []model.Transaction
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("query page: %w", err)
		}
		var page []dynamoTransaction
		if err := attributevalue.UnmarshalListOfMaps(response.Items, &page); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}
		for _, tx := range page {
			t, err := tx.ConvertToTransaction()
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, t)
		}
	}
	return transactions, nil
}

func (d *Repository) List(ctx context.Context, userID string, month string, year string) ([]model.Transaction, error) {
	pk := expression.Key("user").Equal(expression.Value(fmt.Sprintf("USER#%s", userID)))
	var keyEx expression.KeyConditionBuilder
	if month != "" && year != "" {
		keyEx = pk.And(expression.Key("item").Between(
			expression.Value(fmt.Sprintf("TX#%s-%s-01", year, month)),
			expression.Value(fmt.Sprintf("TX#%s-%s-31~", year, month)),
		))
	} else {
		keyEx = pk
	}
	return d.queryTransactions(ctx, keyEx)
}

func (d *Repository) ListRange(ctx context.Context, userID, from, to string) ([]model.Transaction, error) {
	pk := expression.Key("user").Equal(expression.Value(fmt.Sprintf("USER#%s", userID)))
	keyEx := pk.And(expression.Key("item").Between(
		expression.Value(fmt.Sprintf("TX#%s-01", from)),
		expression.Value(fmt.Sprintf("TX#%s-31~", to)),
	))
	return d.queryTransactions(ctx, keyEx)
}

func (d *Repository) Create(ctx context.Context, tx model.Transaction) error {
	// TODO
	return nil
}

func (d *Repository) Update(ctx context.Context, tx model.Transaction) error {
	// TODO
	return nil
}

func (d *Repository) Delete(ctx context.Context, userID, id string) error {
	// TODO
	return nil
}
