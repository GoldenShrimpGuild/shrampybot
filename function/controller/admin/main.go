package admin

import (
	"shrampybot/config"
	"shrampybot/router"
	"shrampybot/utility"
	"shrampybot/utility/nosqldb"

	"github.com/golang-jwt/jwt/v5"
)

func AdminController(route *router.Route) *router.Response {
	resp := &router.Response{
		Body:       route.Router.ErrorBody(7),
		StatusCode: "403",
		Headers:    &router.DefaultResponseHeaders,
	}
	if len(route.Path) < 2 {
		resp.Body = route.Router.ErrorBody(10)
		resp.StatusCode = "400"
		return resp
	}

	// if Token is valid (JWT auth), check if it has admin scope
	claims := route.Router.Event.Claims
	scopes := route.Router.Event.Scopes

	switch route.Path[1] {
	case "category":
		if utility.MatchScope(scopes, "admin:category") {
			c := NewCategoryView()
			return c.CallMethod(route)
		}
	case "collection":
		if utility.MatchScope(scopes, "admin:collection") {
			c := NewCollectionView()
			return c.CallMethod(route)
		}
	case "filter":
		if utility.MatchScope(scopes, "admin:filter") {
			c := NewFilterView()
			return c.CallMethod(route)
		}
	case "user":
		if utility.MatchScope(scopes, "admin:user") {
			c := NewUserView()
			return c.CallMethod(route)
		}
	case "token":
		if utility.MatchScope(scopes, "admin") {
			// Do not allow token management with a static token
			if claims["aud"].(string) != "static" {
				c := NewTokenView()
				return c.CallMethod(route)
			}
		}
	}

	return resp
}

func generateStaticToken(static *nosqldb.StaticTokenDatum) (string, error) {
	staticTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    config.BotName,
		"aud":    "static",
		"sub":    static.CreatorId,
		"kid":    static.Id,
		"iat":    static.CreatedAt.Unix(),
		"exp":    static.ExpiresAt.Unix(),
		"jti":    static.Id,
		"scopes": static.Scopes,
	})
	return staticTokenRaw.SignedString([]byte(static.SecretKey))
}
