package nosqldb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	filterTableName = "filter"
)

type FilterDatum struct {
	Id      string `json:"id"`
	Keyword string `json:"keyword"`
}

func (n *NoSqlDb) GetFilterKeywords() (*[]FilterDatum, error) {
	var err error
	fullTableName := n.prefix + filterTableName
	statement := aws.String(
		fmt.Sprintf("SELECT * FROM \"%v\"", fullTableName),
	)
	output := []FilterDatum{}
	results, err := n.QueryDB(statement)
	if err != nil {
		return &output, err
	}
	// Manual marshalling to expand tags
	for _, result := range *results {
		tempCat := FilterDatum{}
		tempBytes, _ := json.Marshal(result)
		json.Unmarshal(tempBytes, &tempCat)
		output = append(output, tempCat)
	}

	return &output, nil
}
