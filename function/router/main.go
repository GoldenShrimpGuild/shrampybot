package router

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	DefaultResponseHeaders = ResponseHeaders{
		ContentType: "application/json",
	}

	ErrorMap = map[int]string{
		5:  "Unhandled exception occurred while routing.",
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

type GenericBody struct {
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

	action func(route *Route) Response
	Router *Router

	views map[string]any
}

type Router struct {
	ctx   *context.Context
	event *Event

	routes []Route
}

func NewRouter(ctx *context.Context, event *Event) Router {
	return Router{
		ctx:   ctx,
		event: event,
	}
}

func (r *Router) AddRoute(matchPath string, action func(route *Route) Response, requireAuth bool) *Route {
	q, _ := url.ParseQuery(r.event.RawQueryString)

	path := strings.Split(r.event.RawPath, "/")
	if len(path) > 1 {
		path = path[1:]
	}

	new_route := Route{
		match_path: matchPath,
		Body:       r.event.Body,
		Path:       path,
		Method:     r.event.RequestContext.Http.Method,
		Query:      q,
		action:     action,
		Router:     r,
	}

	r.routes = append(r.routes, new_route)

	return &new_route
}

func (r *Router) Route() Response {
	context := map[string]string{
		"environment": lambdacontext.FunctionName,
	}

	bodyBasic := map[string]any{
		"status": map[string]string{
			"msg": "Error processing request.",
		},
		"context": context,
	}

	if r.event.Headers.ContentType != "application/json" {
		return Response{
			Body:       bodyBasic,
			StatusCode: "400",
			Headers:    &DefaultResponseHeaders,
		}
	}

	if !r.event.CheckAuthorization() {
		bodyBasic["status"] = map[string]string{
			"msg": "Unauthorized access. Attempt logged.",
		}
		return Response{
			Body:       bodyBasic,
			StatusCode: "401",
			Headers:    &DefaultResponseHeaders,
		}
	}

	for i := 0; i < len(r.routes); i++ {
		if r.routes[i].Path[0] == r.routes[i].match_path {
			routeResp := r.routes[i].action(&r.routes[i])
			if routeResp.Body != nil {
				routeResp.Body["context"] = context
			} else {
				routeResp.Body = map[string]any{
					"context": context,
				}
			}
			return routeResp
		}
	}

	return Response{
		Body:       bodyBasic,
		StatusCode: "400",
		Headers:    &DefaultResponseHeaders,
	}
}

func (r *Router) ErrorBody(errorCode int, msg string) map[string]any {
	errCodes := maps.Keys(ErrorMap)
	errMsg := ""
	if slices.Contains(errCodes, errorCode) {
		errMsg = ErrorMap[errorCode]
	}

	body, err := json.Marshal(GenericBody{
		Status: &Status{
			Msg:       msg,
			ErrorMsg:  errMsg,
			ErrorCode: errorCode,
		},
		Count: 0,
		Data:  []any{},
	})
	if err != nil {
		return map[string]any{}
	}

	retmap := map[string]any{}
	err = json.Unmarshal(body, &retmap)
	if err != nil {
		return map[string]any{}
	}

	return retmap
}
