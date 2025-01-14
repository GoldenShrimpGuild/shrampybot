package nosqldb

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsC "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	batchSize = 25
)

type NoSqlDb struct {
	ctx    context.Context
	prefix string
	db     *dynamodb.Client
}

func NewClient() (*NoSqlDb, error) {
	n := NoSqlDb{}
	n.ctx = context.Background()
	n.prefix = lambdacontext.FunctionName + "."

	sdkConfig, err := awsC.LoadDefaultConfig(n.ctx)
	if err != nil {
		return &n, err
	}

	n.db = dynamodb.NewFromConfig(sdkConfig)
	return &n, nil
}

func (n *NoSqlDb) CreateTable(tableName string) error {
	fullTableName := n.prefix + tableName

	_, err := n.db.CreateTable(n.ctx, &dynamodb.CreateTableInput{
		TableName:   &fullTableName,
		BillingMode: "PAY_PER_REQUEST",
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeN,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *NoSqlDb) TableInfo(tableName string) (*dynamodb.DescribeTableOutput, error) {
	fullTableName := n.prefix + tableName

	info, err := n.db.DescribeTable(n.ctx, &dynamodb.DescribeTableInput{
		TableName: &fullTableName,
	})
	if err != nil {
		return &dynamodb.DescribeTableOutput{}, err
	}

	return info, nil
}

func (n *NoSqlDb) QueryDB(statement *string) (*[]map[string]any, error) {
	var output []map[string]any
	var nextToken *string

	for moreData := true; moreData; {
		result, err := n.db.ExecuteStatement(n.ctx, &dynamodb.ExecuteStatementInput{
			Statement: statement,
			Limit:     aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return &output, err
		}
		var pageOutput []map[string]any
		err = attributevalue.UnmarshalListOfMaps(result.Items, &pageOutput)
		if err != nil {
			return &output, err
		}
		output = append(output, pageOutput...)
		nextToken = result.NextToken
		moreData = nextToken != nil
	}

	return &output, nil
}

// Safe query using expression builder
// Query requires the table be indexed correctly to prevent a full scan
func (n *NoSqlDb) QueryDBWithExpr(tableName *string, expr *expression.Expression, indexName *string) (*[]map[string]any, error) {
	var output []map[string]any

	result, err := n.db.Query(n.ctx, &dynamodb.QueryInput{
		TableName:                 tableName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		IndexName:                 indexName,
	})
	if err != nil {
		return &output, err
	}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &output)
	if err != nil {
		return &output, err
	}

	return &output, nil
}

// Safe scan using expression builder
// Scan will always parse the entire table. Try to avoid.
func (n *NoSqlDb) ScanDBWithExpr(tableName *string, expr *expression.Expression, indexName *string) (*[]map[string]any, error) {
	var output []map[string]any

	result, err := n.db.Scan(n.ctx, &dynamodb.ScanInput{
		TableName:                 tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		IndexName:                 indexName,
	})
	if err != nil {
		return &output, err
	}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &output)
	if err != nil {
		return &output, err
	}

	return &output, nil
}
