package main

type EventHttpHeaders struct {
	Accept          string `json:"accept"`
	AcceptEncoding  string `json:"accept-encoding"`
	UserAgent       string `json:"user-agent"`
	XForwardedFor   string `json:"x-forwarded-for"`
	XForwardedProto string `json:"x-forwarded-proto"`
	XRequestId      string `json:"x-request-id"`
}

type EventHttp struct {
	HttpHeaders *EventHttpHeaders `json:"headers"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
}

type Event struct {
	Http *EventHttp `json:"http"`
	Name string     `json:"name"`
}
