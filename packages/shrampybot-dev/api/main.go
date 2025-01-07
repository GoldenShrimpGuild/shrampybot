package main

import (
	"context"
)

func Main(ctx context.Context, event Event) Response {
	resp := Response{
		Body:       event.Http.Method,
		StatusCode: "200",
		Headers: ResponseHeaders{
			ContentType: "application/json",
		},
	}
	return resp
}
