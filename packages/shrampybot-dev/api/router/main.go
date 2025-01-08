package router

import (
	"context"
	"net/url"
	"strings"
)

var (
	errorCodes = map[int]string{
		5:  "Unhandled exception occurred while routing: %s",
		10: "No arguments to route.",
		11: "No route provided.",
		12: "Invalid route: %s.",
		13: "No applicable methods or invalid command.",
		14: "Authentication failed.",
		15: "Key not found in input json: %s.",
		16: "Corrupt json found on input.",
		17: "Duplicate Twitch message ID; possible replay attack.",
	}
)

type Route struct {
	match_path string

	body   string
	path   []string
	method string
	query  url.Values

	action func(route *Route) Response
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

func (r *Router) AddRoute(match_path string, action func(route *Route) Response) {
	q, _ := url.ParseQuery(r.event.Http.QueryString)

	new_route := Route{
		match_path: match_path,
		body:       r.event.Http.Body,
		path:       strings.Split(r.event.Http.Path, "/")[1:],
		method:     r.event.Http.Method,
		query:      q,
		action:     action,
	}

	r.routes = append(r.routes, new_route)
}

func (r *Router) Route() Response {
	if r.event.Http.HttpHeaders.ContentType != "application/json" {
		return Response{
			Body:       "",
			StatusCode: "400",
			// Headers: ResponseHeaders{
			// 	ContentType: "application/json",
			// },
		}
	}

	for i := 0; i < len(r.routes); i++ {
		if r.routes[i].path[0] == r.routes[i].match_path {
			return r.routes[i].action(&r.routes[i])
		}
	}

	return Response{
		Body:       "",
		StatusCode: "404",
		// Headers: ResponseHeaders{
		// 	ContentType: "application/json",
		// },
	}
}
