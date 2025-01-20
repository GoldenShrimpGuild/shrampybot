package auth

import (
	"fmt"
	"log"
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
		c := NewRefreshView()
		return c.CallMethod(route)
	case "validate":
		// Validate discord oAuth and produce new access & refresh tokens
		c := NewValidateView()
		return c.CallMethod(route)
	case "self":
		c := NewSelfView()
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
	return accessTokenRaw.SignedString([]byte(oauth.SecretKey))
}

func generateRefreshToken(oauth *nosqldb.OAuthDatum) (string, error) {
	// Generate jwt refreshToken
	refreshTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": config.BotName,
		"sub": "refresh",
		"kid": oauth.Id,
		"iat": time.Now().Unix(),
		// Refresh token lasts for 2 weeks
		"exp": time.Now().Add(336 * time.Hour).Unix(),
		"jti": oauth.RefreshUID,
	})
	return refreshTokenRaw.SignedString([]byte(oauth.SecretKey))
}

func validateRefreshToken(refreshToken string) *jwt.Token {
	var oAuth *nosqldb.OAuthDatum
	var claims jwt.MapClaims

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		return nil
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		_, res := token.Method.(*jwt.SigningMethodHMAC)
		if !res {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Use kid claim as key to retrieve the OAuth secret from DynamoDB
		claims, res := token.Claims.(jwt.MapClaims)
		if !res {
			return nil, fmt.Errorf("could not retrieve token claims")
		}
		oAuth, err = n.GetOAuth(claims["kid"].(string))
		if err != nil {
			return nil, err
		}

		return []byte(oAuth.SecretKey), nil
	})
	if err != nil {
		log.Println("Signature check for JWT failed.")
		return nil
	}
	if !token.Valid {
		return nil
	}

	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		return nil
	}
	if claims["iss"] != config.BotName {
		return nil
	}
	if claims["sub"] != "refresh" {
		return nil
	}

	// Check jti for a match to UUID
	if oAuth.RefreshUID != claims["jti"] {
		log.Println("Invalid UUID for RefreshToken.")
		return nil
	}

	if time.Unix(int64(claims["exp"].(float64)), 0).Before(time.Now()) {
		return nil
	}

	return token
}
