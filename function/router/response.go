package router

import (
	"encoding/json"
)

type ResponseHeaders struct {
	// SetCookie   string `json:"Set-Cookie"`
	ContentType string `json:"Content-Type"`
}

// type ResponseStatus struct {
// 	Msg       string `json:"msg"`
// 	ErrorCode int    `json:"errorCode"`
// }

// type ResponseBody struct {
// 	Data map[string]interface{} `json:"data"`
// }

type Response struct {
	Body       string           `json:"body"`
	StatusCode string           `json:"statusCode"`
	Headers    *ResponseHeaders `json:"headers"`
}

type AWSResponse struct {
	Body            string           `json:"body"`
	StatusCode      string           `json:"statusCode"`
	Headers         *ResponseHeaders `json:"headers"`
	IsBase64Encoded string           `json:"isBase64Encoded"`
}

func NewResponse(body GenericBodyDataFlat, statusCode string) *Response {
	response := Response{
		StatusCode: statusCode,
		Headers:    &DefaultResponseHeaders,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return &response
	}
	json.Unmarshal(bodyBytes, &response.Body)
	return &response
}

func (r *Response) FormatAWS() AWSResponse {
	return AWSResponse{
		Body:            r.Body,
		StatusCode:      r.StatusCode,
		Headers:         r.Headers,
		IsBase64Encoded: "false",
	}
}
