package admin

import "shrampybot/router"

type Collection struct {
	router.View
}

func NewCollection() *Collection {
	c := Collection{}
	return &c
}

func (c *Collection) Get(route *router.Route) *router.Response {
	body := router.GenericBody{}

	response := router.NewResponse(body, "200")
	return response
}

func (v *Collection) CallMethod(route *router.Route) *router.Response {
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

	return router.NewResponse(router.GenericBody{}, "500")
}
