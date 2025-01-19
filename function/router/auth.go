package router

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"log"
	"shrampybot/config"
	"shrampybot/utility/nosqldb"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (e *Event) CheckAuthorizationStatic() bool {
	log.Println("Checking Bearer Authorization (Static)...")
	if e.Headers.Authorization == "" {
		return false
	}

	bearer := strings.Split(e.Headers.Authorization, " ")
	if len(bearer) < 2 || strings.ToLower(bearer[0]) != "bearer" {
		return false
	}

	if bearer[1] != config.GsgAdminToken {
		return false
	}

	log.Println("Bearer Auth Succeeded.")
	return true
}

func (e *Event) CheckAuthorizationJWT() bool {
	var oAuth *nosqldb.OAuthDatum
	var claims jwt.MapClaims

	log.Println("Checking Bearer Authorization (JWT)...")
	if e.Headers.Authorization == "" {
		return false
	}

	bearer := strings.Split(e.Headers.Authorization, " ")
	if len(bearer) < 2 || strings.ToLower(bearer[0]) != "bearer" {
		return false
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		return false
	}

	token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
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
			log.Printf("Could not retrieve OAuth detail for kid %v\n", claims["kid"])
			return nil, err
		}

		return []byte(oAuth.SecretKey), nil
	})
	if err != nil {
		log.Printf("Signature check for JWT failed: %v\n", err)
		return false
	}
	if !token.Valid {
		return false
	}

	claims, res := token.Claims.(jwt.MapClaims)
	if !res {
		return false
	}
	if claims["iss"] != config.BotName {
		return false
	}
	if claims["sub"] != "access" {
		return false
	}
	if time.Unix(int64(claims["exp"].(float64)), 0).Before(time.Now()) {
		return false
	}

	e.Token = token
	return true
}

func (e *Event) calculateSHA256Signature() string {
	combined := e.Headers.TwitchEventsubMessageId +
		e.Headers.TwitchEventsubMessageTimestamp +
		e.Body

	h := hmac.New(sha256.New, []byte(config.TwitchEventSecret))
	_, err := h.Write([]byte(combined))
	if err != nil {
		return ""
	}
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

// Function for handling Twitch Webhook notification authorization
func (e *Event) CheckTwitchAuthorization() bool {
	log.Println("Checking Twitch Webhook Authorization...")
	if config.TwitchEventSecret == "" ||
		e.Headers.TwitchEventsubMessageId == "" ||
		e.Headers.TwitchEventsubMessageTimestamp == "" ||
		e.Headers.TwitchEventsubMessageSignature == "" {
		log.Println("Blank values in Twitch headers!")
		return false
	}

	log.Println("Comparing SHA256 signature values...")
	our_digest := []byte(e.calculateSHA256Signature())
	their_digest := []byte(e.Headers.TwitchEventsubMessageSignature)

	log.Printf("Their digest: %v\n", their_digest)
	log.Printf("Our digest  : %v\n", our_digest)

	// Use constant time comparison functions
	if subtle.ConstantTimeEq(int32(len(our_digest)), int32(len(their_digest))) == 0 {
		return false
	}
	if subtle.ConstantTimeCompare(our_digest, their_digest) == 0 {
		return false
	}

	log.Println("Twitch Auth Succeeded.")
	return true
}
