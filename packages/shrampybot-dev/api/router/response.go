package router

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
	Body       map[string]any   `json:"body"`
	StatusCode string           `json:"statusCode"`
	Headers    *ResponseHeaders `json:"headers"`
	// IsBase64Encoded bool             `json:"isBase64Encoded"`
}
