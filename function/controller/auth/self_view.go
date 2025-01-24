package auth

import (
	"encoding/json"
	"log"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"

	"github.com/bwmarrin/discordgo"
)

type SelfView struct {
	router.View `tstype:",extends,required"`
}

type SelfResponseBody struct {
	discordgo.User `tstype:",extends,required"`
	Member         *discordgo.Member           `json:"member,omitempty" tstype:"discordgo.Member"`
	Connections    []*discordgo.UserConnection `json:"connections,omitempty" tstype:"discordgo.UserConnection"`
}

func NewSelfView() *SelfView {
	c := SelfView{}
	return &c
}

func (v *SelfView) CallMethod(route *router.Route) *router.Response {
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
	case "OPTIONS":
		return v.Options(route)
	}

	return router.NewResponse(router.GenericBodyDataFlat{}, "500")
}

func (v *SelfView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Self.Get")
	var err error
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	if !route.Router.Event.CheckAuthorizationJWT() {
		log.Println("Failed JWT Auth check.")
		response.StatusCode = "401"
		return response
	}
	// Get token object, defined when CheckingAuthorizationJWT above
	// token := route.Router.Event.Token
	claims := route.Router.Event.Claims

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	dOAuth, err := n.GetDiscordOAuth(claims["sub"].(string))
	if err != nil {
		log.Printf("Could not get Discord OAuth record")
		response.StatusCode = "500"
		return response
	}
	d, err := discord.NewOAuthClient(dOAuth)
	if err != nil {
		log.Println("Could not create new Discord oauth client.")
		response.StatusCode = "500"
		return response
	}
	if dOAuth.Refreshed {
		err = n.PutDiscordOAuth(dOAuth)
		if err != nil {
			log.Println("Couldn't save refreshed discord oauth.")
			// Continue on anyway, for now.
		}
	}

	self, err := d.GetSelf()
	if err != nil || self.ID == "" {
		log.Println("Could not retrieve Discord user with new OAuth credentials.")
		response.StatusCode = "500"
		return response
	}

	// // Update discord connections in table whenever self is called.
	// err = mapDiscordConnections(self.ID, self.Username, n, d)
	// if err != nil {
	// 	log.Println("Could not map Discord connections.")
	// }

	dc, err := discord.NewBotClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	body := SelfResponseBody{}
	selfBytes, err := json.Marshal(self)
	if err != nil {
		log.Println("Failed to marshal self JSON")
		response.StatusCode = "500"
		return response
	}
	err = json.Unmarshal(selfBytes, &body)
	if err != nil {
		log.Println("Failed to unmarshal self JSON")
		response.StatusCode = "500"
		return response
	}

	member, err := dc.GetGuildMember(self.ID)
	if err != nil {
		log.Println("Failed to get Discord guild membership")
	}
	body.Member = member

	connections, err := d.GetConnections()
	if err != nil {
		log.Println("Failed to get Discord connections")
	}
	body.Connections = connections

	bodyBytes, _ := json.Marshal(body)
	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Self.Get")
	return response
}
