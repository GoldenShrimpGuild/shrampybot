package gsg

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"strings"

	"github.com/litui/helix/v3"
)

type StreamerView struct {
	router.View `tstype:",extends,required"`
}

type StreamerResponseBody struct {
	Count int           `json:"count"`
	Data  []*helix.User `json:"data" tstype:"helix.User[]"`
}

func NewStreamerView() *StreamerView {
	c := StreamerView{}
	return &c
}

func (v *StreamerView) CallMethod(route *router.Route) *router.Response {
	switch route.Method {
	case "GET":
		return v.Get(route)
	case "POST":
		return v.Post(route)
	case "PUT":
		return v.Put(route)
	case "PATCH":
		return v.Patch(route)
	case "DELETE":
		return v.Delete(route)
	}

	return router.NewResponse(router.GenericBodyDataFlat{}, "500")
}

func (v *StreamerView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: GSG.Streamer.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	var streamers []*nosqldb.TwitchUserDatum

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	log.Printf("Query map: %v\n", route.Query)
	log.Printf("Query string raw: %v\n", route.Router.Event.RawQueryString)

	logins := route.Query["login"]

	if len(logins) > 0 {
		loginIdMap, err := n.GetTwitchLoginIdMap()
		if err != nil {
			log.Printf("Could not get Twitch Login<->Id map: %v\n", err)
			response.StatusCode = "500"
			return response
		}

		for _, login := range logins {
			if loginIdMap[strings.ToLower(login)] != "" {
				user, err := n.GetTwitchUser(loginIdMap[strings.ToLower(login)])
				if err != nil {
					log.Printf("Could not retrieve Twitch user %v: %v\n", login, err)
					response.StatusCode = "500"
					return response
				}

				streamers = append(streamers, user)
			}
		}
	} else {
		// Fetch login names from our stored Twitch users
		streamers, err = n.GetTwitchUsers()
		if err != nil {
			log.Println("Could not get saved Twitch logins.")
			response.StatusCode = "500"
			return response
		}
	}

	respBody := StreamerResponseBody{}
	streamerBytes, _ := json.Marshal(streamers)
	json.Unmarshal(streamerBytes, &respBody.Data)
	respBody.Count = len(respBody.Data)

	bodyBytes, _ := json.Marshal(respBody)
	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: GSG.Streamer.Get")
	return response
}
