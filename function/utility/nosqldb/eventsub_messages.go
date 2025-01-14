package nosqldb

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EventsubMessageDatum struct {
	Id    string `json:"id"`
	Time  string `json:"time"` // RFC3339
	Type  string `json:"type"`
	Retry string `json:"retry"`
}

const (
	eventsubMessageTableName = "eventsub_messages"
)

func (n *NoSqlDb) GetEventsubMessage(id string) (*EventsubMessageDatum, error) {
	var err error
	fullTableName := n.prefix + eventsubMessageTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &EventsubMessageDatum{}, err
	}
	output := EventsubMessageDatum{}

	tempMap := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &tempMap)
	tempBytes, _ := json.Marshal(tempMap)
	json.Unmarshal(tempBytes, &output)

	return &output, nil
}

func (n *NoSqlDb) PutEventsubMessage(eventsub *EventsubMessageDatum) error {
	var err error
	fullTableName := n.prefix + eventsubMessageTableName

	esInterface := map[string]any{}
	esBytes, _ := json.Marshal(eventsub)
	json.Unmarshal(esBytes, &esInterface)

	var esItem map[string]types.AttributeValue
	esItem, err = attributevalue.MarshalMap(esInterface)
	if err != nil {
		log.Printf("Couldn't marshal item %v for writing because: %v\n", eventsub.Id, err)
		return err
	}

	_, err = n.db.PutItem(n.ctx, &dynamodb.PutItemInput{
		Item:      esItem,
		TableName: &fullTableName,
	})
	if err != nil {
		log.Printf("Couldn't record Eventsub Message: %v", err)
	}

	return err
}
