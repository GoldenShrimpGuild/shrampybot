package main

type ResponseHeaders struct {
	// SetCookie   string `json:"Set-Cookie"`
	ContentType string `json:"Content-Type"`
}

type Response struct {
	Body       string          `json:"body"`
	StatusCode string          `json:"statusCode"`
	Headers    ResponseHeaders `json:"headers"`
}
