package nosqldb

import (
	"encoding/json"
	"shrampybot/utility"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	staticTokenTableName = "static_token"
)

type StaticTokenDatum struct {
	Id           string    `json:"id"`
	CreatorId    string    `json:"creator_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
	Scopes       string    `json:"scopes,omitempty"`
	Purpose      string    `json:"purpose"`
	SecretKey    string    `json:"-"`
	SecretKeyIV  string    `json:"secret_key_iv,omitempty"`
	SecretKeyEnc string    `json:"secret_key_enc,omitempty"`
}

func (n *NoSqlDb) GetStaticToken(id string) (*StaticTokenDatum, error) {
	var err error
	fullTableName := n.prefix + staticTokenTableName

	keyMap := map[string]types.AttributeValue{}
	keyMap["id"] = &types.AttributeValueMemberS{Value: id}

	result, err := n.db.GetItem(n.ctx, &dynamodb.GetItemInput{
		Key:       keyMap,
		TableName: &fullTableName,
	})
	if err != nil {
		return &StaticTokenDatum{}, err
	}
	output := StaticTokenDatum{}

	rStatic := map[string]any{}
	attributevalue.UnmarshalMap(result.Item, &rStatic)
	oBytes, _ := json.Marshal(rStatic)
	json.Unmarshal(oBytes, &output)

	if output.Id != "" && output.SecretKeyEnc != "" && output.SecretKeyIV != "" {
		// Decrypt secret values
		output.SecretKey, _ = utility.DecryptSecret(output.SecretKeyEnc, output.SecretKeyIV)
	}

	return &output, nil
}

func (n *NoSqlDb) PutStaticToken(static *StaticTokenDatum) error {
	var err error
	fullTableName := n.prefix + staticTokenTableName

	// Encrypt secret values to be stored
	static.SecretKeyEnc, static.SecretKeyIV, _ = utility.EncryptSecret(static.SecretKey)

	tempMap := map[string]string{}
	tempBytes, _ := json.Marshal(static)
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
