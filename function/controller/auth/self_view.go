package auth

import (
	"encoding/json"
	"log"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v5"
)

type SelfView struct {
	router.View `tstype:",extends,required"`
}

type SelfResponseBody struct {
	discordgo.User `tstype:",extends,required"`
	Member         *discordgo.Member `json:"member,omitempty" tstype:"discordgo.Member"`
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

func (v *RefreshView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Self.Get")
	var err error
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	if !route.Router.Event.CheckAuthorizationJWT() {
		response.StatusCode = "403"
		return response
	}
	// Get token object, defined when CheckingAuthorizationJWT above
	token := route.Router.Event.Token
	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		response.StatusCode = "500"
		return response
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	dOAuth, err := n.GetDiscordOAuth(claims["kid"].(string))
	if err != nil {
		log.Printf("Could not get Discord OAuth record")
		response.StatusCode = "500"
		return response
	}

	d, err := discord.NewOAuthClient(dOAuth.AccessToken)
	if err != nil {
		log.Println("Could not create new Discord oauth client.")
		response.StatusCode = "403"
		return response
	}

	self, err := d.GetSelf()
	if err != nil || self.ID == "" {
		log.Println("Could not retrieve Discord user with new OAuth credentials.")
		response.StatusCode = "500"
		return response
	}

	dc, err := discord.NewBotClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	body := SelfResponseBody{}
	selfBytes, err := json.Marshal(self)
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	err = json.Unmarshal(selfBytes, &body)
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	member, err := dc.GetGuildMember(self.ID)
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	body.Member = member
	bodyBytes, _ := json.Marshal(body)
	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Self.Get")
	return response
}
