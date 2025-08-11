package utility

import "strings"

var (
	ValidStaticTokenScopes = []string{
		"login",
		"dev",
		"gsg",
		"gsg:streamer",
		"admin",
		"admin:categories",
		"admin:collection",
		"admin:events",
		"admin:filters",
		"admin:stream",
		"admin:tokens",
		"admin:users",
	}
)

func MatchScope(userScopes []string, targetScope string) bool {
	tsParts := strings.Split(targetScope, ":")
	if len(tsParts) == 0 {
		return false
	}

	for _, userScope := range userScopes {
		usParts := strings.Split(userScope, ":")
		matchedLevel := -1

		for i := 0; i < len(tsParts); i++ {

			if len(usParts) == i+1 && len(tsParts) >= i+1 {
				if usParts[i] == tsParts[i] {
					matchedLevel = i
				} else {
					break
				}
			}
		}

		if matchedLevel+1 == len(usParts) {
			return true
		}
	}
	return false
}
