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
	router.View
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

	users, err := getTwitchUsers()
	if err != nil || len(*users) == 0 {
		log.Println("Exited route abnormally: Collection.Patch")
		response.StatusCode = "500"
		return response
	}

	err = saveActiveTwitchUsers(users)
	if err != nil {
		log.Println("Exited route abnormally: Collection.Patch")
		response.StatusCode = "500"
		return response
	}

	body := map[string]any{}
	body["count"] = len(*users)
	data := []string{}

	// Munge users into displayable format
	for _, u := range *users {
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
	ch_th := make(chan string)
	ch_mh := make(chan string)
	go th.GetTeamMemberLoginsThreaded(ch_th)
	go mh.GetMappedTwitchLoginsThreaded(ch_mh)

	for login := range ch_mh {
		loginList = append(loginList, login)
	}
	for login := range ch_th {
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

	return &output, nil
}

func saveActiveTwitchUsers(users *[]nosqldb.TwitchUserDatum) error {
	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		return err
	}

	pastActive, err := n.GetActiveTwitchIds()
	if err != nil {
		return err
	}
	disableIds := []string{}
	for _, pastUser := range *pastActive {
		foundMatch := false

		for _, newUser := range *users {
			if pastUser == newUser.ID {
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			disableIds = append(disableIds, pastUser)
		}
	}

	// Disable removed (eg: no longer shrampy) Twitch Ids
	err = n.DisableTwitchIds(&disableIds)
	if err != nil {
		return err
	}

	// Update/add twitch users and mark them active
	err = n.PutTwitchUsers(users)
	if err != nil {
		return err
	}

	return nil
}
