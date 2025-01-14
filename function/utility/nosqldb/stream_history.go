package nosqldb

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/litui/helix/v3"
)

const (
	streamHistoryTableName = "stream_history"
)

type StreamHistoryDatum struct {
	helix.Stream
	DiscordPostId   string    `json:"discord_post_id,omitempty"`
	DiscordPostUrl  string    `json:"discord_post_url,omitempty"`
	MastodonPostId  string    `json:"mastodon_post_id,omitempty"`
	MastodonPostUrl string    `json:"mastodon_post_url,omitempty"`
	BlueskyPostId   string    `json:"bluesky_post_id,omitempty"`
	BlueskyPostUrl  string    `json:"bluesky_post_url,omitempty"`
	EndedAt         time.Time `json:"ended_at,omitempty"`
}

func (n *NoSqlDb) GetStream(id string) (*StreamHistoryDatum, error) {
	var err error
	fullTableName := n.prefix + streamHistoryTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &StreamHistoryDatum{}, err
	}
	output := StreamHistoryDatum{}

	rStream := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rStream)
	oBytes, _ := json.Marshal(rStream)
	json.Unmarshal(oBytes, &output)

	if rStream["tag_ids"] != nil {
		json.Unmarshal([]byte(rStream["tag_ids"].(string)), &output.TagIDs)
	}
	if rStream["tags"] != nil {
		json.Unmarshal([]byte(rStream["tags"].(string)), &output.Tags)
	}

	return &output, nil
}

func (n *NoSqlDb) GetLatestStreamByUserId(user_id string) (*StreamHistoryDatum, error) {
	var err error
	fullTableName := n.prefix + streamHistoryTableName
	indexName := fullTableName + ".user_id-index"

	filt := expression.Key("user_id").Equal(expression.Value(user_id))
	expr, err := expression.NewBuilder().WithKeyCondition(filt).Build()
	if err != nil {
		return &StreamHistoryDatum{}, err
	}

	result, err := n.QueryDBWithExpr(&fullTableName, &expr, &indexName)
	if err != nil {
		return &StreamHistoryDatum{}, err
	}

	output := StreamHistoryDatum{}
	output.StartedAt = time.Unix(0, 0)
	for _, res := range *result {
		startedAt, err := time.Parse(time.RFC3339, res["started_at"].(string))
		if err != nil {
			continue
		}

		if startedAt.After(output.StartedAt) {
			oBytes, _ := json.Marshal(res)
			json.Unmarshal(oBytes, &output)

			if res["tag_ids"] != nil {
				json.Unmarshal([]byte(res["tag_ids"].(string)), &output.TagIDs)
			}
			if res["tags"] != nil {
				json.Unmarshal([]byte(res["tags"].(string)), &output.Tags)
			}
		}
	}

	return &output, nil
}

func (n *NoSqlDb) PutStream(stream *StreamHistoryDatum) error {
	var err error

	fullTableName := n.prefix + streamHistoryTableName

	tempMap := map[string]string{}
	tempBytes, _ := json.Marshal(stream)
	json.Unmarshal(tempBytes, &tempMap)

	tagIds, _ := json.Marshal(stream.TagIDs)
	tempMap["tag_ids"] = string(tagIds)
	tags, _ := json.Marshal(stream.Tags)
	tempMap["tags"] = string(tags)

	var item map[string]types.AttributeValue
	item, err = attributevalue.MarshalMap(tempMap)
	if err != nil {
		log.Printf("Couldn't marshal stream %v for writing because: %v\n", stream.ID, err)
		return err
	}

	_, err = n.db.PutItem(n.ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &fullTableName,
	})
	if err != nil {
		log.Printf("Couldn't record Stream: %v", err)
	}

	return err
}
