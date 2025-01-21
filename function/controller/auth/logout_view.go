package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LogoutView struct {
	router.View `tstype:",extends,required"`
}

type LogoutResponseBody struct {
	UserId string `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
}

func NewLogoutView() *LogoutView {
	c := LogoutView{}
	return &c
}

func (v *LogoutView) CallMethod(route *router.Route) *router.Response {
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

func (v *LogoutView) Post(route *router.Route) *router.Response {
	log.Println("Entered route: Auth.Logout.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	cookies, err := http.ParseCookie(route.Router.Event.Headers.Cookie)
	if err != nil {
		log.Printf("Issue parsing cookies from header: %v\n", err)
		response.StatusCode = "500"
		return response
	}
	var oldRefreshToken string
	for _, c := range cookies {
		if c.Name == "RefreshToken" {
			oldRefreshToken = c.Value
		}
	}
	if oldRefreshToken == "" {
		log.Println("No RefreshToken provided.")
		response.StatusCode = "500"
		return response
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	token := validateRefreshToken(oldRefreshToken)
	if token == nil || !token.Valid {
		response.StatusCode = "401"
		return response
	}
	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		response.StatusCode = "500"
		return response
	}

	oAuth, err := n.GetOAuth(claims["kid"].(string))
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	// Revoke UUID and save
	oAuth.RefreshUID = fmt.Sprintf("REVOKED:%v", oAuth.RefreshUID)
	err = n.PutOAuth(oAuth)
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	body := LogoutResponseBody{
		UserId: oAuth.Id,
		Status: "success",
	}
	bodyBytes, _ := json.Marshal(body)

	// RefreshToken in httponly cookie
	cookie := http.Cookie{
		Name:        "RefreshToken",
		Value:       "REVOKED",
		HttpOnly:    true,
		SameSite:    http.SameSiteNoneMode,
		Secure:      true,
		Partitioned: true,
		Expires:     time.Unix(0, 0),
	}
	response.Headers.SetCookie = cookie.String()

	response.Body = string(bodyBytes)

	response.StatusCode = "200"
	log.Println("Exiting route: Auth.Validate.Post")
	return response
}
