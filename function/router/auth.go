package router

import (
	"crypto"
	"crypto/hmac"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"shrampybot/config"
	"strings"
)

/*
For now I'm just going to somewhat replicate the original
shrampybot behaviour, but this should be enhanced for security
down the road.  - Aria
*/

func (e *Event) CheckAuthorization() bool {
	log.Println("Checking Bearer Authorization...")
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

func (e *Event) calculateSHA256Signature() string {
	combined := e.Headers.TwitchEventsubMessageId +
		e.Headers.TwitchEventsubMessageTimestamp +
		e.Body

	h := hmac.New(crypto.SHA256.HashFunc().New, []byte(config.TwitchEventSecret))
	return "sha256=" + hex.EncodeToString(h.Sum([]byte(combined)))
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
