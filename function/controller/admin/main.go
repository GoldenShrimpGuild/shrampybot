package admin

import (
	"log"
	"shrampybot/config"
	"shrampybot/connector/discord"
	"shrampybot/router"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

// Should these be moved to the template.yml?
const (
	discordAdminRole = ""
)

func AdminController(route *router.Route) *router.Response {
	resp := &router.Response{
		Body:       route.Router.ErrorBody(5),
		StatusCode: "500",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		resp.Body = route.Router.ErrorBody(10)
		resp.StatusCode = "400"
		return resp
	}

	dc, err := discord.NewBotClient()
	if err != nil {
		resp.StatusCode = "500"
		return resp
	}

	// Apply credential level check to all endpoints on this controller.

	// if Token is valid (JWT auth), check if it's an admin user
	token := route.Router.Event.Token
	if token != nil && token.Valid {
		claims, res := route.Router.Event.Token.Claims.(jwt.MapClaims)
		if !res {
			log.Println("Failed to map JWT claims.")
			resp.StatusCode = "500"
			return resp
		}
		if !userIsAdmin(claims["kid"].(string), dc) {
			log.Printf("User %v is not a valid admin on Discord.\n", claims["kid"].(string))
			resp.StatusCode = "403"
			return resp
		}
	}
	// static tokens are all admin level so we can skip the check for them

	switch route.Path[1] {
	case "category":
		c := NewCategoryView()
		return c.CallMethod(route)
	case "collection":
		c := NewCollectionView()
		return c.CallMethod(route)
	case "filter":
		c := NewFilterView()
		return c.CallMethod(route)
	case "user":
		c := NewUserView()
		return c.CallMethod(route)
	}

	return resp
}

func userIsAdmin(id string, dc *discord.BotClient) bool {
	membership, err := dc.GetGuildMember(id)
	if err != nil {
		return false
	}
	if slices.Contains(membership.Roles, config.DiscordAdminRole) {
		return true
	}

	return false
}

func userIsDev(id string, dc *discord.BotClient) bool {
	membership, err := dc.GetGuildMember(id)
	if err != nil {
		return false
	}
	if slices.Contains(membership.Roles, config.DiscordDevRole) {
		return true
	}

	return false
}
