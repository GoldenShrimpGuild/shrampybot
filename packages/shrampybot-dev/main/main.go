package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Event struct {
	Name string `json:"name"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

type ResponseHeaders struct {
	// SetCookie   string `json:"Set-Cookie"`
	ContentType string `json:"Content-Type"`
}

type Response struct {
	Body       string          `json:"body"`
	StatusCode string          `json:"statusCode"`
	Headers    ResponseHeaders `json:"headers"`
}

func Main(ctx context.Context, event Event) Response {
	router := httprouter.New()
	router.GET("/", Index)

	return Response{
		Body:       "{}",
		StatusCode: "200",
		Headers: ResponseHeaders{
			ContentType: "application/json",
		},
	}
}
