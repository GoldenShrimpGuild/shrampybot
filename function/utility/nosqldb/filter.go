package nosqldb

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	filterTableName = "filter"
)

type FilterDatum struct {
	Id              string `json:"id"`
	Keyword         string `json:"keyword"`
	CaseInsensitive bool   `json:"case_insensitive"`
	IsRegex         bool   `json:"is_regex"`
}

func (n *NoSqlDb) GetFilterKeyword(id string) (*FilterDatum, error) {
	var err error
	fullTableName := n.prefix + filterTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &FilterDatum{}, err
	}
	output := FilterDatum{}
	attributevalue.UnmarshalMap(result.Item, &output)
	rFilt := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rFilt)

	return &output, nil
}

// Fills the ID value in the passed-in filter reference if a match exists in the db
func (n *NoSqlDb) FillFilterIdIfAny(filter *FilterDatum) error {
	var results *[]map[string]any
	var err error

	fullTableName := n.prefix + filterTableName
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\"", fullTableName),
	)
	results, err = n.QueryDB(statement)
	if err != nil {
		return err
	}

	for _, rFilt := range *results {
		if rFilt["keyword"] == filter.Keyword {
			filter.Id = rFilt["id"].(string)
			return nil
		}
	}

	return nil
}

func (n *NoSqlDb) GetFilterKeywords() ([]*FilterDatum, error) {
	var err error
	fullTableName := n.prefix + filterTableName
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\"", fullTableName),
	)
	output := []*FilterDatum{}
	results, err := n.QueryDB(statement)
	if err != nil {
		return output, err
	}
	// Manual marshalling to expand tags
	for _, result := range *results {
		tempCat := FilterDatum{}
		tempBytes, _ := json.Marshal(result)
		json.Unmarshal(tempBytes, &tempCat)
		output = append(output, &tempCat)
	}

	return output, nil
}

func (n *NoSqlDb) PutFilterKeywords(filterKeywords []*FilterDatum) error {
	var err error
	fullTableName := n.prefix + filterTableName

	mapFilterKeywords := []map[string]any{}
	for _, fk := range filterKeywords {
		tempMap := map[string]any{}

		fkBytes, _ := json.Marshal(fk)
		err := json.Unmarshal(fkBytes, &tempMap)
		if err != nil {
			log.Println("Error unmarshaling FilterKeyword data into temporary map")
			return err
		}

		if fk.Id != "" {
			tempMap["id"] = fk.Id
		} else {
			tempMap["id"] = uuid.NewString()
		}

		mapFilterKeywords = append(mapFilterKeywords, tempMap)
	}

	// Iterate through batch chunks, 25 items at a time.
	for subList := range slices.Chunk(mapFilterKeywords, batchSize) {
		var writeReqs []types.WriteRequest

		for _, fk := range subList {
			var item map[string]types.AttributeValue
			item, err = attributevalue.MarshalMap(fk)
			if err != nil {
				log.Printf("Couldn't marshal filter keyword %v for batch writing because: %v\n", fk["keyword"], err)
				continue
			}

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
			log.Printf("Couldn't batch-submit categories to %v because: %v\n", fullTableName, err)
		}
	}
	return err
}

func (n *NoSqlDb) RemoveFilterKeyword(ids *[]string) error {
	var err error

	fullTableName := n.prefix + filterTableName

	for _, id := range *ids {
		keyMap := map[string]types.AttributeValue{}
		keyMap["id"] = &types.AttributeValueMemberS{Value: id}

		_, err = n.db.DeleteItem(n.ctx, &dynamodb.DeleteItemInput{
			Key:       keyMap,
			TableName: &fullTableName,
		})
		if err != nil {
			log.Printf("Couldn't delete filter keyword item id %v\n", id)
			continue
		}
	}

	return nil
}
