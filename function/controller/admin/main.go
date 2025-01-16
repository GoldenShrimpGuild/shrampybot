package admin

import (
	"shrampybot/router"
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
	}

	return resp
}
