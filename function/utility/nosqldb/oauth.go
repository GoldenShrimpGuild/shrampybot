package nosqldb

import (
	"encoding/json"
	"log"
	"shrampybot/utility"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	oAuthTableName = "oauth"
)

type OAuthDatum struct {
	Id string `json:"id"`
	// Raw, unencrypted value of secret key; never gets stored
	SecretKey    string `json:"-"`
	SecretKeyIV  string `json:"secret_key_iv,omitempty"`
	SecretKeyEnc string `json:"secret_key_enc,omitempty"`
	RefreshUID   string `json:"refresh_uid,omitempty"`
}

func (n *NoSqlDb) GetOAuth(id string) (*OAuthDatum, error) {
	var err error
	fullTableName := n.prefix + oAuthTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &OAuthDatum{}, err
	}
	output := OAuthDatum{}

	rOAuth := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rOAuth)
	oBytes, _ := json.Marshal(rOAuth)
	json.Unmarshal(oBytes, &output)

	if output.Id != "" && output.SecretKeyEnc != "" && output.SecretKeyIV != "" {
		// Decrypt secret values
		output.SecretKey, _ = utility.DecryptSecret(output.SecretKeyEnc, output.SecretKeyIV)
	}

	return &output, nil
}

func (n *NoSqlDb) PutOAuth(oauth *OAuthDatum) error {
	var err error
	fullTableName := n.prefix + oAuthTableName

	// Encrypt secret values to be stored
	oauth.SecretKeyEnc, oauth.SecretKeyIV, err = utility.EncryptSecret(oauth.SecretKey)
	if err != nil {
		log.Printf("Could not encrypt oauth secret key: %v\n", err)
		return err
	}

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
