package admin

import (
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

	return router.NewResponse(router.GenericBody{}, "500")
}

func (c *CollectionView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Collection.Get")
	var response *router.Response
	var logins *[]map[string]any
	body := router.GenericBody{
		Count: 0,
		Data:  []any{},
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
	}

	// Fetch login names from our stored Twitch users
	logins, err = n.GetActiveTwitchLogins()
	if err != nil {
		log.Println("Could not get saved Twitch logins.")
	}

	// Munge users into displayable format
	for _, u := range *logins {
		body.Data = append(body.Data, u["login"])
	}
	body.Count = int64(len(*logins))
	response = router.NewResponse(body, "200")
	log.Println("Exited route: Collection.Get")
	return response
}

// A PATCH call will gather, assemble, and update all the user data required
// to do other Shrampy tasks. This is the linchpin of Shrampybot.
func (c *CollectionView) Patch(route *router.Route) *router.Response {
	log.Println("Entered route: Collection.Patch")
	var response *router.Response

	body := router.GenericBody{
		Count: 0,
		Data:  []any{},
	}

	users, err := getTwitchUsers()
	if err != nil || len(*users) == 0 {
		route.Router.ErrorBody(2, "")
		response = router.NewResponse(body, "500")
		log.Println("Exited route abnormally: Collection.Patch")
		return response
	}

	err = saveActiveTwitchUsers(users)
	if err != nil {
		route.Router.ErrorBody(3, "")
		response = router.NewResponse(body, "500")
		log.Println("Exited route abnormally: Collection.Patch")
		return response
	}

	// Munge users into displayable format
	for _, u := range *users {
		body.Data = append(body.Data, u["login"])
	}
	body.Count = int64(len(*users))
	response = router.NewResponse(body, "200")
	log.Println("Exited route: Collection.Patch")
	return response
}

func getTwitchUsers() (*[]map[string]string, error) {
	var err error
	var users *[]map[string]string
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
	users, err = th.GetUsers(&loginList)
	if err != nil {
		return users, err
	}

	return users, nil
}

func saveActiveTwitchUsers(users *[]map[string]string) error {
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
			if pastUser["id"] == newUser["id"] {
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			disableIds = append(disableIds, pastUser["id"].(string))
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
