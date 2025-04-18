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
	helix.User `tstype:",extends,required"`
	// The active flag refers to membership in either the GSG team or tagged presence
	// on the GSG mastodon instance
	ShrampybotArtistName       string `json:"shrampybot_artist_name"`
	ShrampybotActive           bool   `json:"shrampybot_active"`
	ShrampybotLocation         string `json:"shrampybot_location,omitempty"`
	ShrampybotBirthMonth       int    `json:"shrampybot_birth_month"`
	ShrampybotBirthDay         int    `json:"shrampybot_birth_day"`
	ShrampybotEmail            string `json:"shrampybot_email,omitempty"`
	ShrampybotOverlayNow       string `json:"shrampybot_overlay_now,omitempty"`
	ShrampybotOverlayNext      string `json:"shrampybot_overlay_next,omitempty"`
	ShrampybotOverlayLater     string `json:"shrampybot_overlay_later,omitempty"`
	ShrampybotGroupsTeam       bool   `json:"shrampybot_groups_team"`
	ShrampybotGroupsMainPage   bool   `json:"shrampybot_groups_main_page"`
	ShrampybotGroupsMDLRMonday bool   `json:"shrampybot_groups_mdlr_monday"`
	ShrampybotGroupsClub       bool   `json:"shrampybot_groups_club"`
	MastodonUserId             string `json:"mastodon_user_id,omitempty"`
	DiscordUserId              string `json:"discord_user_id,omitempty"`
	DiscordUsername            string `json:"discord_username,omitempty"`
	BlueskyUserId              string `json:"bluesky_user_id,omitempty"`
	BlueskyUsername            string `json:"bluesky_username,omitempty"`
	YoutubeUserId              string `json:"youtube_user_id,omitempty"`
	YoutubeUsername            string `json:"youtube_username,omitempty"`
	XUserId                    string `json:"x_user_id,omitempty"`
	XUsername                  string `json:"x_username,omitempty"`
	InstagramUserId            string `json:"instagram_user_id,omitempty"`
	InstagramUserName          string `json:"instagram_username,omitempty"`
	SteamUserId                string `json:"steam_user_id,omitempty"`
	SteamUsername              string `json:"steam_username,omitempty"`
	FacebookUserId             string `json:"facebook_user_id,omitempty"`
	FacebookUsername           string `json:"facebook_username,omitempty"`
	GithubUserId               string `json:"github_user_id,omitempty"`
	GithubUsername             string `json:"github_username,omitempty"`
	TiktokUserId               string `json:"tiktok_user_id,omitempty"`
	TiktokUsername             string `json:"tiktok_username,omitempty"`
	SpotifyUserId              string `json:"spotify_user_id,omitempty"`
	SpotifyUsername            string `json:"spotify_username,omitempty"`
}

func (n *NoSqlDb) DisableTwitchUsers(ids *[]string) error {
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

func (n *NoSqlDb) GetTwitchIdLoginMap() (map[string]string, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT id,login FROM \"%v\"", fullTableName),
	)
	output := map[string]string{}
	result, err := n.QueryDB(statement)
	if err != nil {
		return output, err
	}

	for _, tu := range *result {
		tempId := objx.New(tu).Get("id")
		tempLogin := objx.New(tu).Get("login")
		output[tempId.Str()] = tempLogin.Str()
	}

	return output, nil
}

func (n *NoSqlDb) GetTwitchLoginIdMap() (map[string]string, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT id,login FROM \"%v\"", fullTableName),
	)
	output := map[string]string{}
	result, err := n.QueryDB(statement)
	if err != nil {
		return output, err
	}

	for _, tu := range *result {
		tempId := objx.New(tu).Get("id")
		tempLogin := objx.New(tu).Get("login")
		output[tempLogin.Str()] = tempId.Str()
	}

	return output, nil
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

func (n *NoSqlDb) GetTwitchUsers() ([]*TwitchUserDatum, error) {
	var err error
	fullTableName := n.prefix + twitchUsersTableName
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\"", fullTableName),
	)
	output := []*TwitchUserDatum{}
	results, err := n.QueryDB(statement)
	if err != nil {
		return output, err
	}
	// Manual marshalling to expand tags
	for _, result := range *results {
		tempCat := TwitchUserDatum{}
		tempBytes, _ := json.Marshal(result)
		json.Unmarshal(tempBytes, &tempCat)
		output = append(output, &tempCat)
	}

	return output, nil
}

func (n *NoSqlDb) PutTwitchUser(user *TwitchUserDatum) error {
	var err error

	fullTableName := n.prefix + twitchUsersTableName

	tempMap := map[string]any{}
	tempBytes, _ := json.Marshal(user)
	json.Unmarshal(tempBytes, &tempMap)

	var item map[string]types.AttributeValue
	item, err = attributevalue.MarshalMap(tempMap)
	if err != nil {
		log.Printf("Couldn't marshal twitch user %v for writing because: %v\n", user.ID, err)
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

func (n *NoSqlDb) PutTwitchUsers(users []*TwitchUserDatum) error {
	var err error

	fullTableName := n.prefix + twitchUsersTableName

	for subList := range slices.Chunk(users, batchSize) {
		var writeReqs []types.WriteRequest

		for _, user := range subList {
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
