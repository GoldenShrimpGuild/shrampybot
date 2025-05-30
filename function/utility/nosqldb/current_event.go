package nosqldb

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	currentEventTableName = "current_event"
)

type CurrentEventDatum struct {
	EventId     string `json:"eventId"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

func (n *NoSqlDb) GetCurrentEvent(index uint8) (*CurrentEventDatum, error) {
	var err error
	fullTableName := n.prefix + currentEventTableName
	currentEvent := CurrentEventDatum{}

	filt := expression.Key("id").Equal(expression.Value(strconv.Itoa(int(index))))
	expr, err := expression.NewBuilder().WithKeyCondition(filt).Build()
	if err == nil {
		results, err := n.QueryDBWithExpr(&fullTableName, &expr, nil)
		if err == nil && len(*results) > 0 {
			resultBytes, _ := json.Marshal((*results)[0])
			json.Unmarshal(resultBytes, &currentEvent)
		}
	}

	return &currentEvent, err
}

func (n *NoSqlDb) PutCurrentEvent(index uint8, currentEvent *CurrentEventDatum) error {
	var err error
	fullTableName := n.prefix + currentEventTableName

	tempMap := map[string]any{}
	tempBytes, _ := json.Marshal(currentEvent)
	tempMap["id"] = strconv.Itoa(int(index))
	json.Unmarshal(tempBytes, &tempMap)

	var item map[string]types.AttributeValue
	item, err = attributevalue.MarshalMap(tempMap)
	if err != nil {
		log.Printf("Couldn't marshal current event %v index %v for writing because: %v\n", currentEvent.EventId, index, err)
		return err
	}

	_, err = n.db.PutItem(n.ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &fullTableName,
	})
	if err != nil {
		log.Printf("Couldn't record current event %v index %v because: %v", currentEvent.EventId, index, err)
	}

	return err
}

func (n *NoSqlDb) DeleteCurrentEvent(index uint8) error {
	var err error
	fullTableName := n.prefix + currentEventTableName

	tempMap := map[string]any{}
	tempMap["id"] = strconv.Itoa(int(index))

	var item map[string]types.AttributeValue
	item, _ = attributevalue.MarshalMap(tempMap)

	_, err = n.db.DeleteItem(n.ctx, &dynamodb.DeleteItemInput{
		Key:       item,
		TableName: &fullTableName,
	})
	if err != nil {
		log.Printf("Couldn't delete current event index %v because: %v", index, err)
	}

	return err
}
