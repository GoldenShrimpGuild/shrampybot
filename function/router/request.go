package router

// const (
// 	EVENTSUB_MESSAGE_ID        = "twitch-eventsub-message-id"
// 	EVENTSUB_MESSAGE_TIMESTAMP = "twitch-eventsub-message-timestamp"
// 	EVENTSUB_MESSAGE_SIGNATURE = "twitch-eventsub-message-signature"
// 	EVENTSUB_MESSAGE_TYPE      = "twitch-eventsub-message-type"
// )

type Headers struct {
	Accept                            string `json:"accept"`
	AcceptEncoding                    string `json:"accept-encoding"`
	Authorization                     string `json:"authorization"`
	ContentLength                     string `json:"content-length"`
	ContentType                       string `json:"content-type"`
	Host                              string `json:"host"`
	UserAgent                         string `json:"user-agent"`
	XForwardedFor                     string `json:"x-forwarded-for"`
	XForwardedPort                    string `json:"x-forwarded-port"`
	XForwardedProto                   string `json:"x-forwarded-proto"`
	XRequestId                        string `json:"x-request-id"`
	TwitchEventsubMessageId           string `json:"twitch-eventsub-message-id"`
	TwitchEventsubMessageRetry        string `json:"twitch-eventsub-message-retry"`
	TwitchEventsubMessageTimestamp    string `json:"twitch-eventsub-message-timestamp"`
	TwitchEventsubMessageSignature    string `json:"twitch-eventsub-message-signature"`
	TwitchEventsubMessageType         string `json:"twitch-eventsub-message-type"`
	TwitchEventsubSubscriptionType    string `json:"twitch-eventsub-subscription-type"`
	TwitchEventsubSubscriptionVersion string `json:"twitch-eventsub-subscription-version"`
}

type Http struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	SourceIp  string `json:"sourceIp"`
	UserAgent string `json:"userAgent"`
}

type RequestContext struct {
	AccountId    string `json:"accountId"`
	Time         int64  `json:"timeEpoch"`
	DomainPrefix string `json:"domainPrefix"`
	RequestId    string `json:"requestId"`
	DomainName   string `json:"domainName"`
	Http         *Http  `json:"http"`
	ApiId        string `json:"apiId"`
}

type Event struct {
	Headers         *Headers        `json:"headers"`
	IsBase64Encoded bool            `json:"isBase64Encoded"`
	RawPath         string          `json:"rawPath"`
	RequestContext  *RequestContext `json:"requestContext"`
	Name            string          `json:"name"`
	Body            string          `json:"body"`
	RawQueryString  string          `json:"rawQueryString"`
}
