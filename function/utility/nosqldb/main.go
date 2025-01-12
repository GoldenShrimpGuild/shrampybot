package nosqldb

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsC "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	batchSize       = 25
	ActiveUserField = "shrampybot_active"
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

func (n *NoSqlDb) DisableTwitchIds(ids *[]string) error {
	var err error
	fullTableName := n.prefix + "twitch_users"

	for subList := range slices.Chunk(*ids, batchSize) {
		var writeReqs []types.WriteRequest

		for _, id := range subList {
			item := map[string]types.AttributeValue{}
			item["id"] = &types.AttributeValueMemberS{Value: id}
			item["shrampybot_active"] = &types.AttributeValueMemberBOOL{Value: false}

			writeReqs = append(
				writeReqs,
				types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: item,
					},
				},
			)
		}
		_, err = n.db.BatchWriteItem(n.ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{fullTableName: writeReqs},
		})
		if err != nil {
			log.Printf("Couldn't batch-submit ids to %v because: %v\n", fullTableName, err)
		}
	}
	return nil
}

func (n *NoSqlDb) GetActiveTwitchIds() (*[]map[string]any, error) {
	var output *[]map[string]any
	var err error
	fullTableName := n.prefix + "twitch_users"
	statement := aws.String(
		fmt.Sprintf("SELECT id FROM \"%v\" WHERE shrampybot_active=true", fullTableName),
	)
	output, err = n.QueryDB(statement)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (n *NoSqlDb) GetActiveTwitchLogins() (*[]map[string]any, error) {
	var output *[]map[string]any
	var err error
	fullTableName := n.prefix + "twitch_users"
	statement := aws.String(
		fmt.Sprintf("SELECT login FROM \"%v\" WHERE shrampybot_active=true", fullTableName),
	)
	output, err = n.QueryDB(statement)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (n *NoSqlDb) GetActiveTwitchUsers() (*[]map[string]any, error) {
	var output *[]map[string]any
	var err error
	fullTableName := n.prefix + "twitch_users"
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\" WHERE shrampybot_active=true", fullTableName),
	)
	output, err = n.QueryDB(statement)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (n *NoSqlDb) PutTwitchUsers(users *[]map[string]string) error {
	var err error

	fullTableName := n.prefix + "twitch_users"

	for subList := range slices.Chunk(*users, batchSize) {
		var writeReqs []types.WriteRequest

		for _, user := range subList {
			var item map[string]types.AttributeValue
			item, err = attributevalue.MarshalMap(user)
			// Mark this list of users as active.
			if err != nil {
				log.Printf("Couldn't marshal user %v for batch writing because: %v\n", user["login"], err)
				continue
			}
			item["shrampybot_active"] = &types.AttributeValueMemberBOOL{Value: true}

			writeReqs = append(
				writeReqs,
				types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: item,
					},
				},
			)
		}

		_, err = n.db.BatchWriteItem(n.ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{fullTableName: writeReqs},
		})
		if err != nil {
			log.Printf("Couldn't batch-submit users to %v because: %v\n", fullTableName, err)
		}
	}
	return err
}
