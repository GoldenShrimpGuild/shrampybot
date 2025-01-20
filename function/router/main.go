package router

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	allowedOrigins = []string{
		"http://localhost:5173",
		"https://shrampybot.github.io",
	}

	DefaultResponseHeaders = ResponseHeaders{
		// AccessControlAllowOrigin:      "http://localhost:5173",
		AccessControlAllowMethods:     "GET, HEAD, POST, PUT, DELETE",
		AccessControlAllowCredentials: "true",
		AccessControlAllowHeaders:     "Content-Type, Authorization",
		ContentType:                   "application/json",
		Vary:                          "Origin",
	}

	ErrorMap = map[int]string{
		1:  "Wrong content type.",
		2:  "Data retrieval error.",
		3:  "Data storage error.",
		4:  "Database connection error.",
		5:  "Unhandled exception occurred while routing.",
		6:  "Request data parsing error.",
		10: "No arguments to route.",
		11: "No route provided.",
		12: "Invalid route.",
		13: "No applicable methods or invalid command.",
		14: "Authentication failed.",
		15: "Key not found in input json.",
		16: "Corrupt json found on input.",
		17: "Duplicate Twitch message ID; possible replay attack.",
	}
)

type GenericBodyDataFlat struct {
	Status *Status `json:"status,omitempty"`
	Count  int64   `json:"count,omitempty"`
	Data   []any   `json:"data,omitempty"`
}

type Status struct {
	Msg       string `json:"msg"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}

type Route struct {
	match_path string

	Body   string
	Method string
	Query  url.Values
	Path   []string

	action func(route *Route) *Response
	Router *Router

	RequireAuth bool
}

type Router struct {
	Event *Event

	ctx    *context.Context
	routes []Route
}

func NewRouter(ctx *context.Context, event *Event) Router {
	return Router{
		ctx:   ctx,
		Event: event,
	}
}

func (r *Router) AddRoute(matchPath string, action func(route *Route) *Response, requireAuth bool) *Route {
	q, _ := url.ParseQuery(r.Event.RawQueryString)

	path := strings.Split(r.Event.RawPath, "/")
	if len(path) > 1 {
		path = path[1:]
	}

	new_route := Route{
		match_path:  matchPath,
		Body:        r.Event.Body,
		Path:        path,
		Method:      r.Event.RequestContext.Http.Method,
		Query:       q,
		action:      action,
		Router:      r,
		RequireAuth: requireAuth,
	}

	r.routes = append(r.routes, new_route)

	return &new_route
}

func (r *Router) Route() *Response {
	log.Println("Started routing...")
	context := map[string]string{
		"environment": lambdacontext.FunctionName,
	}

	// Dynamically set allowed CORS origin in default headers based on allow list
	if slices.Contains(allowedOrigins, r.Event.Headers.Origin) {
		DefaultResponseHeaders.AccessControlAllowOrigin = r.Event.Headers.Origin
	}

	// if r.Event.Headers.ContentType != "application/json" {
	// 	return &Response{
	// 		Body:       r.ErrorBody(1),
	// 		StatusCode: "400",
	// 		Headers:    &DefaultResponseHeaders,
	// 	}
	// }

	for i := 0; i < len(r.routes); i++ {
		if r.routes[i].Path[0] == r.routes[i].match_path {
			if r.routes[i].RequireAuth {
				log.Printf("Authentication required for endpoint %v\n", r.routes[i].match_path)

				if !r.Event.CheckAuthorizationJWT() {
					log.Println("JWT authentication failed, trying static...")

					if !r.Event.CheckAuthorizationStatic() {
						log.Println("Static authentication failed!")

						return &Response{
							Body:       r.ErrorBody(14),
							StatusCode: "401",
							Headers:    &DefaultResponseHeaders,
						}
					}
				}

				log.Println("Authentication succeeded!")
			} else {
				log.Printf("No authentication required for endpoint %v\n", r.routes[i].match_path)
			}

			responseBody := map[string]any{}
			// Continue routing post-authentication or lack thereof.
			routeResp := r.routes[i].action(&r.routes[i])
			log.Printf("Body post-route: %v", routeResp.Body)

			// Fill in a nice happy json body if there is no existing body data
			if routeResp.Body == "" {
				responseBody["context"] = context
				json.Unmarshal([]byte(routeResp.Body), &responseBody)
				responseBytes, _ := json.Marshal(responseBody)
				routeResp.Body = string(responseBytes)
			}
			return routeResp
		}
	}

	// Failed to route
	return &Response{
		Body:       r.ErrorBody(12),
		StatusCode: "400",
		Headers:    &DefaultResponseHeaders,
	}
}

func (r *Router) ErrorBody(errorCode int) string {
	errCodes := maps.Keys(ErrorMap)
	errMsg := ""
	if slices.Contains(errCodes, errorCode) {
		errMsg = ErrorMap[errorCode]
	}

	body, err := json.Marshal(GenericBodyDataFlat{
		Status: &Status{
			ErrorMsg:  errMsg,
			ErrorCode: errorCode,
		},
		Count: 0,
		Data:  []any{},
	})
	if err != nil {
		return ""
	}

	return string(body)
}
