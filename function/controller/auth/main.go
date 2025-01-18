package auth

import (
	"shrampybot/config"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthController(route *router.Route) *router.Response {
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
	case "refresh":
		// Request a refreshed set of tokens
	case "validate":
		// Validate discord oAuth and produce new access & refresh tokens
		c := NewValidateView()
		return c.CallMethod(route)
	}

	return resp
}

func generateAccessToken(oauth *nosqldb.OAuthDatum) (string, error) {
	// Generate jwt accessToken
	accessTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.BotName,
		"sub": "access",
		"kid": oauth.Id,
		"iat": time.Now().Unix(),
		// Access token lasts for 10 minutes
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})
	return accessTokenRaw.SignedString(oauth.SecretKey)
}

func generateRefreshToken(oauth *nosqldb.OAuthDatum) (string, error) {
	// Generate jwt refreshToken
	refreshTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.BotName,
		"sub": "refresh",
		"kid": oauth.Id,
		"iat": time.Now().Unix(),
		// Refresh token lasts for 2 weeks
		"exp": time.Now().Add(336 * time.Hour),
		"jti": oauth.RefreshUID,
	})
	return refreshTokenRaw.SignedString(oauth.SecretKey)
}
