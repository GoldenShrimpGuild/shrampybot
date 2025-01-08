package router

type EventHttpHeaders struct {
	Accept          string `json:"accept"`
	AcceptEncoding  string `json:"accept-encoding"`
	ContentType     string `json:"content-type"`
	Host            string `json:"host"`
	UserAgent       string `json:"user-agent"`
	XForwardedFor   string `json:"x-forwarded-for"`
	XForwardedProto string `json:"x-forwarded-proto"`
	XRequestId      string `json:"x-request-id"`
}

type EventHttp struct {
	HttpHeaders     *EventHttpHeaders `json:"headers"`
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	Body            string            `json:"body"`
	QueryString     string            `json:"queryString"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type Event struct {
	Http *EventHttp `json:"http"`
	Name string     `json:"name"`
	Foo  string     `json:"foo"`
}
