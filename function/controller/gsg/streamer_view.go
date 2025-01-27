package gsg

import (
	"shrampybot/router"
)

type StreamerView struct {
	router.View `tstype:",extends,required"`
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
