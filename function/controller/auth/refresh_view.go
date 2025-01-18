package auth

import (
	"shrampybot/router"
)

type RefreshView struct {
	router.View
}

type RefreshBody struct {
	router.GenericBodyDataFlat
}

func NewRefreshView() *RefreshView {
	c := RefreshView{}
	return &c
}

func (v *RefreshView) CallMethod(route *router.Route) *router.Response {
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

// func (v *RefreshView) Post(route *router.Route) *router.Response {
// 	log.Println("Entered route: Auth.Refresh.Post")

// 	log.Println("Exiting route: Auth.Refresh.Post")
// 	// return resp
// }
