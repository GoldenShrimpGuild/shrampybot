package router

import (
	"shrampybot/config"
	"strings"
)

/*
For now I'm just going to somewhat replicate the original
shrampybot behaviour, but this should be enhanced for security
down the road.  - Aria
*/

func (e *Event) CheckAuthorization() bool {
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

	return true
}
