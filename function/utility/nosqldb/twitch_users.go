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
	"github.com/litui/helix/v3"
	"github.com/stretchr/objx"
)

const (
	activeUserField      = "shrampybot_active"
	twitchUsersTableName = "twitch_users"
)

type TwitchUserDatum struct {
	helix.User       `tstype:",extends,required"`
	ShrampybotActive bool   `json:"shrampybot_active,omitempty"`
	MastodonUserId   string `json:"mastodon_user_id,omitempty"`
	DiscordUserId    string `json:"discord_user_id,omitempty"`
	BlueskyUserId    string `json:"bluesky_user_id,omitempty"`
}

func (n *NoSqlDb) DisableTwitchIds(ids *[]string) error {
	var err error
	fullTableName := n.prefix + twitchUsersTableName

	for subList := range slices.Chunk(*ids, batchSize) {
		var writeReqs []types.WriteRequest

		for _, id := range subList {
			item := map[string]types.AttributeValue{}
			item["id"] = &types.AttributeValueMemberS{Value: id}
			item[activeUserField] = &types.AttributeValueMemberBOOL{Value: false}

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

func (n *NoSqlDb) GetTwitchUser(id string) (*TwitchUserDatum, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		fmt.Printf("Could not retrieve Twitch user %v: %v\n", id, err)
		return &TwitchUserDatum{}, err
	}
	output := TwitchUserDatum{}

	tempRes := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &tempRes)
	tempBytes, _ := json.Marshal(tempRes)
	json.Unmarshal(tempBytes, &output)

	return &output, nil
}

func (n *NoSqlDb) GetActiveTwitchIds() (*[]string, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT id FROM \"%v\" WHERE shrampybot_active=true", fullTableName),
	)
	output := []string{}
	result, err := n.QueryDB(statement)
	if err != nil {
		return &output, err
	}

	for _, id := range *result {
		tempId := objx.New(id).Get("id")
		output = append(output, tempId.Str())
	}

	return &output, nil
}

func (n *NoSqlDb) GetActiveTwitchLogins() (*[]string, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT login FROM \"%v\" WHERE %v=true", fullTableName, activeUserField),
	)
	output := []string{}
	result, err := n.QueryDB(statement)
	if err != nil {
		return &output, err
	}

	for _, id := range *result {
		tempId := objx.New(id).Get("login")
		output = append(output, tempId.Str())
	}

	return &output, nil
}

func (n *NoSqlDb) GetActiveTwitchUsers() (*[]TwitchUserDatum, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\" WHERE %v=true", fullTableName, activeUserField),
	)
	output := []TwitchUserDatum{}
	results, err := n.QueryDB(statement)
	if err != nil {
		return &output, err
	}
	// Manual marshalling to expand tags
	for _, result := range *results {
		tempCat := TwitchUserDatum{}
		tempBytes, _ := json.Marshal(result)
		json.Unmarshal(tempBytes, &tempCat)
		output = append(output, tempCat)
	}

	return &output, nil
}

func (n *NoSqlDb) PutTwitchUsers(users *[]TwitchUserDatum) error {
	var err error

	fullTableName := n.prefix + twitchUsersTableName

	for subList := range slices.Chunk(*users, batchSize) {
		var writeReqs []types.WriteRequest

		for _, user := range subList {
			// Mark this list of users as active.
			user.ShrampybotActive = true

			userInterface := map[string]any{}
			userBytes, _ := json.Marshal(user)
			json.Unmarshal(userBytes, &userInterface)

			var item map[string]types.AttributeValue
			item, err = attributevalue.MarshalMap(userInterface)
			if err != nil {
				log.Printf("Couldn't marshal user %v for batch writing because: %v\n", user.Login, err)
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
			log.Printf("Couldn't batch-submit users to %v because: %v\n", fullTableName, err)
		}
	}
	return err
}
