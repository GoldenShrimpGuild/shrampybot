package public

import (
	"fmt"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"sort"
	"strings"
)

type MultiView struct {
	router.View `tstype:",extends,required"`
}

type MultiBody struct {
	router.GenericBodyDataFlat `tstype:",extends,required"`
	Data                       *[]nosqldb.StreamHistoryDatum `json:"data" tstype:"nosqldb.StreamHistoryDatum[]"`
}

func NewMultiView() *MultiView {
	c := MultiView{}
	return &c
}

func (v *MultiView) CallMethod(route *router.Route) *router.Response {
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

func (v *MultiView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Public.Multi.Get")
	response := &router.Response{}
	response.Headers = &router.ResponseHeaders{
		ContentType: "text/plain; charset=utf-8",
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Fetch active streams
	streams, err := n.GetActiveStreams()
	if err != nil || streams == nil {
		log.Println("Could not get active streams.")
		response.StatusCode = "500"
		return response
	}

	validCategories, _ := n.GetCategoryMap()

	streamerNames := []string{}
	for _, stream := range *streams {
		hasValidCategory := false

		for _, c := range *validCategories {
			if stream.GameName == c.TwitchCategory && c.Id != "" {
				hasValidCategory = true
				break
			}
		}
		if !hasValidCategory {
			continue
		}

		if len(route.Path) == 3 && route.Path[2] != "" {
			if !strings.Contains(strings.ToLower(stream.Title), strings.ToLower(route.Path[2])) {
				continue
			}
		}

		if !stream.ShrampybotFiltered {
			streamerNames = append(streamerNames, stream.UserLogin)
		}
	}
	sort.Strings(streamerNames)
	body := strings.Join(streamerNames, "/")

	response.StatusCode = "301"
	response.Headers.Location = fmt.Sprintf("https://www.multitwitch.tv/%s", body)
	response.Body = body + "\n"
	log.Println("Exited route: Public.Multi.Get")
	return response
}
