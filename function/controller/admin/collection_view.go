package admin

import (
	"encoding/json"
	"log"
	"shrampybot/connector/mastodon"
	"shrampybot/connector/twitch"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"slices"
	"sort"
)

type CollectionView struct {
	router.View `tstype:",extends,required"`
}

func NewCollectionView() *CollectionView {
	c := CollectionView{}
	return &c
}

func (v *CollectionView) CallMethod(route *router.Route) *router.Response {
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

func (c *CollectionView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.Collection.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	var logins *[]string

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Fetch login names from our stored Twitch users
	logins, err = n.GetActiveTwitchLogins()
	if err != nil {
		log.Println("Could not get saved Twitch logins.")
		response.StatusCode = "500"
		return response
	}

	body := map[string]any{}
	body["count"] = len(*logins)
	body["data"] = logins
	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Collection.Get")
	return response
}

// A PATCH call will gather, assemble, and update all the user data required
// to do other Shrampy tasks. This is the linchpin of Shrampybot.
func (c *CollectionView) Patch(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.Collection.Patch")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	teamUsers, err := getTwitchUsers()
	if err != nil || len(*teamUsers) == 0 {
		log.Println("Exited route abnormally: Collection.Patch")
		response.StatusCode = "500"
		return response
	}

	storedUsers, err := n.GetTwitchUsers()
	if err != nil || len(*storedUsers) == 0 {
		log.Println("Exited route abnormally: Collection.Patch")
		response.StatusCode = "500"
		return response
	}

	intersectUsers := []*nosqldb.TwitchUserDatum{}
	extraUsers := []*nosqldb.TwitchUserDatum{}
	diffUsers := []*nosqldb.TwitchUserDatum{}

	for _, u := range *teamUsers {
		foundStoredUser := false

		for _, su := range *storedUsers {
			if u.ID == su.ID {
				foundStoredUser = true
				// Update changeable fields from Twitch
				su.Login = u.Login
				su.DisplayName = u.DisplayName
				su.Description = u.Description
				su.OfflineImageURL = u.OfflineImageURL
				su.ProfileImageURL = u.ProfileImageURL
				su.ViewCount = u.ViewCount
				su.MastodonUserId = u.MastodonUserId
				su.ShrampybotActive = true

				intersectUsers = append(intersectUsers, &su)
				break
			}
		}

		if !foundStoredUser {
			u.ShrampybotActive = true
			extraUsers = append(extraUsers, &u)
		}
	}

	for _, su := range *storedUsers {
		foundTeamUser := false

		for _, tu := range *teamUsers {
			if tu.ID == su.ID {
				foundTeamUser = true
			}
		}

		if !foundTeamUser {
			su.ShrampybotActive = false
			diffUsers = append(diffUsers, &su)
		}
	}

	err = n.PutTwitchUsers(intersectUsers)
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	err = n.PutTwitchUsers(diffUsers)
	if err != nil {
		response.StatusCode = "500"
		return response
	}
	err = n.PutTwitchUsers(extraUsers)
	if err != nil {
		response.StatusCode = "500"
		return response
	}

	activeUsers := append(intersectUsers, extraUsers...)

	body := map[string]any{}
	body["count"] = len(activeUsers)
	data := []string{}

	// Munge users into displayable format
	for _, u := range activeUsers {
		data = append(data, u.Login)
	}
	body["data"] = data
	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Collection.Patch")
	return response
}

func getTwitchUsers() (*[]nosqldb.TwitchUserDatum, error) {
	var err error
	var loginList []string

	// connect to Twitch
	th, _ := twitch.NewClient()
	// connect to Mastodon
	mh, _ := mastodon.NewClient()

	// Parallelize assembling the logins list from multiple sources
	chTh := make(chan string)
	chMh := make(chan map[string]string)
	go th.GetTeamMemberLoginsThreaded(chTh)
	go mh.GetMappedTwitchLoginsThreaded(chMh)

	mastodonMap := <-chMh
	for t := range mastodonMap {
		loginList = append(loginList, t)
	}

	for login := range chTh {
		loginList = append(loginList, login)
	}
	sort.Strings(loginList)
	loginList = slices.Compact(loginList)

	// Fetch user records for logins on our list
	output := []nosqldb.TwitchUserDatum{}
	users, err := th.GetUsers(&loginList)
	if err != nil {
		return &output, err
	}

	userBytes, _ := json.Marshal(users)
	json.Unmarshal(userBytes, &output)

	// Add mastodon IDs to the retrieved users
	for i, u := range output {
		output[i].MastodonUserId = mastodonMap[u.Login]
	}

	return &output, nil
}
