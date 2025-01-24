package admin

import (
	"log"
	"shrampybot/router"
	"strings"
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

	// if Token is valid (JWT auth), check if it has admin scope
	token := route.Router.Event.Token
	if token != nil && token.Valid {
		claims := route.Router.Event.Claims
		scopes := route.Router.Event.Scopes
		foundAdminScope := false
		for _, scope := range scopes {
			subScope := strings.Split(scope, ":")
			if subScope[0] == "admin" {
				foundAdminScope = true
				break
			}
		}
		if !foundAdminScope {
			log.Printf("%v token for user %v does not have the admin scope.\n", claims["aud"].(string), claims["sub"].(string))
			resp.StatusCode = "403"
			return resp
		}
	}

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
	case "token":
		c := NewTokenView()
		return c.CallMethod(route)
	}

	return resp
}
