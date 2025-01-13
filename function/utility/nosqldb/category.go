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

type CategoryDatum struct {
	Id             string   `json:"id,omitempty"`
	TwitchCategory string   `json:"twitch_category"`
	MastodonTags   []string `json:"mastodon_tags,omitempty"`
	BlueskyTags    []string `json:"bluesky_tags,omitempty"`
}

func (n *NoSqlDb) GetCategory(id string) (*CategoryDatum, error) {
	var err error
	fullTableName := n.prefix + "category_map"

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &CategoryDatum{}, err
	}
	output := CategoryDatum{}
	rCat := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rCat)
	output.Id = rCat["id"].(string)
	output.TwitchCategory = rCat["twitch_category"].(string)
	json.Unmarshal([]byte(rCat["mastodon_tags"].(string)), &output.MastodonTags)
	json.Unmarshal([]byte(rCat["bluesky_tags"].(string)), &output.BlueskyTags)

	return &output, nil
}

func (n *NoSqlDb) GetCategoryByName(name string) (*CategoryDatum, error) {
	var err error
	fullTableName := n.prefix + "category_map"

	keyMap := map[string]types.AttributeValue{}
	keyMap["twitch_category"] = &types.AttributeValueMemberS{Value: name}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &CategoryDatum{}, err
	}
	output := CategoryDatum{}
	rCat := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rCat)
	attributevalue.UnmarshalMap(result.Item, &output)

	json.Unmarshal([]byte(rCat["mastodon_tags"].(string)), &output.MastodonTags)
	json.Unmarshal([]byte(rCat["bluesky_tags"].(string)), &output.BlueskyTags)

	return &output, nil
}

func (n *NoSqlDb) GetCategoryMap() (*[]CategoryDatum, error) {
	var results *[]map[string]any
	var err error
	output := []CategoryDatum{}

	fullTableName := n.prefix + "category_map"
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\"", fullTableName),
	)
	results, err = n.QueryDB(statement)
	if err != nil {
		return &output, err
	}

	// Manual marshalling to expand tags
	for _, rCat := range *results {
		tempCat := CategoryDatum{}
		tempCat.Id = rCat["id"].(string)
		tempCat.TwitchCategory = rCat["twitch_category"].(string)

		json.Unmarshal([]byte(rCat["mastodon_tags"].(string)), &tempCat.MastodonTags)
		json.Unmarshal([]byte(rCat["bluesky_tags"].(string)), &tempCat.BlueskyTags)

		output = append(output, tempCat)
	}

	return &output, nil
}

func (n *NoSqlDb) PutCategories(categories *[]CategoryDatum) error {
	var err error

	fullTableName := n.prefix + "category_map"

	// Manually remap categories to flatten tags fields
	mapCategories := []map[string]string{}
	for _, c := range *categories {
		tempMap := map[string]string{}
		if c.Id != "" {
			tempMap["id"] = c.Id
		} else {
			tempMap["id"] = uuid.NewString()
		}
		tempMap["twitch_category"] = c.TwitchCategory
		mTags, _ := json.Marshal(c.MastodonTags)
		tempMap["mastodon_tags"] = string(mTags)
		bTags, _ := json.Marshal(c.BlueskyTags)
		tempMap["bluesky_tags"] = string(bTags)

		mapCategories = append(mapCategories, tempMap)
	}

	// Iterate through batch chunks, 25 items at a time.
	for subList := range slices.Chunk(mapCategories, batchSize) {
		var writeReqs []types.WriteRequest

		for _, cat := range subList {
			log.Printf("Cat: %v\n", cat)
			var item map[string]types.AttributeValue
			item, err = attributevalue.MarshalMap(cat)
			if err != nil {
				log.Printf("Couldn't marshal category %v for batch writing because: %v\n", cat["twitch_category"], err)
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

func (n *NoSqlDb) RemoveCategory(ids *[]string) error {
	var err error

	fullTableName := n.prefix + "category_map"

	for _, id := range *ids {
		keyMap := map[string]types.AttributeValue{}
		keyMap["id"] = &types.AttributeValueMemberS{Value: id}

		_, err = n.db.DeleteItem(n.ctx, &dynamodb.DeleteItemInput{
			Key:       keyMap,
			TableName: &fullTableName,
		})
		if err != nil {
			log.Printf("Couldn't delete category_map item id %v\n", id)
			continue
		}
	}

	return nil
}
