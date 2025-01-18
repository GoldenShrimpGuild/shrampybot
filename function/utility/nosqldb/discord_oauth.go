package nosqldb

import (
	"encoding/json"
	"shrampybot/utility"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	discordOAuthTableName = "discord_oauth"
)

type DiscordOAuthDatum struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
	// Raw, unencrypted value of access token; never gets stored
	AccessToken string
	// Encrypted value of access token
	AccessTokenIV  string `json:"access_token_iv,omitempty"`
	AccessTokenEnc string `json:"access_token_enc,omitempty"`
	TokenType      string `json:"token_type,omitempty"`
	ExpiresIn      int    `json:"expires_in,omitempty"`
	// Raw, unencrypted value of refresh token; never gets stored
	RefreshToken string
	// Encrypted value of refresh token
	RefreshTokenIV  string `json:"refresh_token_iv,omitempty"`
	RefreshTokenEnc string `json:"refresh_token_enc,omitempty"`

	Scope string `json:"scope,omitempty"`
}

func (n *NoSqlDb) GetDiscordOAuth(id string) (*DiscordOAuthDatum, error) {
	var err error
	fullTableName := n.prefix + discordOAuthTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &DiscordOAuthDatum{}, err
	}
	output := DiscordOAuthDatum{}

	rOAuth := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rOAuth)
	oBytes, _ := json.Marshal(rOAuth)
	json.Unmarshal(oBytes, &output)

	// Decrypt secret values
	output.AccessToken, _ = utility.DecryptSecret(output.AccessTokenEnc, output.AccessTokenIV)
	output.RefreshToken, _ = utility.DecryptSecret(output.RefreshTokenEnc, output.RefreshTokenIV)

	return &output, nil
}

func (n *NoSqlDb) PutDiscordOAuth(oauth *DiscordOAuthDatum) error {
	var err error
	fullTableName := n.prefix + discordOAuthTableName

	// Encrypt secret values to be stored
	oauth.AccessTokenEnc, oauth.AccessTokenIV, _ = utility.EncryptSecret(oauth.AccessToken)
	oauth.RefreshTokenEnc, oauth.RefreshTokenIV, _ = utility.EncryptSecret(oauth.RefreshToken)

	tempMap := map[string]string{}
	tempBytes, _ := json.Marshal(oauth)
	json.Unmarshal(tempBytes, &tempMap)

	var item map[string]types.AttributeValue
	item, err = attributevalue.MarshalMap(tempMap)
	if err != nil {
		return err
	}

	_, err = n.db.PutItem(n.ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &fullTableName,
	})

	return err
}
