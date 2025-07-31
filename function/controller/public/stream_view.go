package public

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"slices"
	"strings"
)

type StreamView struct {
	router.View `tstype:",extends,required"`
}

type StreamBody struct {
	router.GenericBodyDataFlat `tstype:",extends,required"`
	Data                       *[]nosqldb.StreamHistoryDatum `json:"data" tstype:"nosqldb.StreamHistoryDatum[]"`
}

func NewStreamView() *StreamView {
	c := StreamView{}
	return &c
}

func (v *StreamView) CallMethod(route *router.Route) *router.Response {
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

func (v *StreamView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Public.Stream.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	streams := []nosqldb.StreamHistoryDatum{}
	if len(route.Path) == 3 {
		// Get single result
		stream, err := n.GetStream(route.Path[2])
		if err != nil {
			log.Printf("Get stream [%v] failed.", route.Path[2])
			response.StatusCode = "500"
			return response
		}
		streams = append(streams, *stream)
	} else {
		// // Fetch active streams
		streamsRef, err := n.GetActiveStreams()
		if err != nil || streamsRef == nil {
			log.Println("Could not get active streams.")
			response.StatusCode = "500"
			return response
		}
		streams = *streamsRef
	}

	// Load requisite category listing
	categories, err := n.GetCategoryMap()
	if err != nil {
		log.Println("Could not retrieve categories for stream sorting.")
		response.StatusCode = "500"
		return response
	}
	lcCatStrings := []string{}
	for _, c := range *categories {
		lcCatStrings = append(lcCatStrings, strings.ToLower(c.TwitchCategory))
	}

	// Only include streams which were not filtered out by Shrampybot in this list
	outputStreams := []*nosqldb.StreamHistoryDatum{}
	for _, stream := range streams {
		if stream.ShrampybotFiltered {
			continue
		}

		if !slices.Contains(lcCatStrings, strings.ToLower(stream.GameName)) {
			continue
		}

		outputStreams = append(outputStreams, &stream)
	}
	streamBytes, _ := json.Marshal(outputStreams)

	body := map[string]any{}
	body["count"] = len(streams)
	bodyRef := []map[string]any{}
	json.Unmarshal(streamBytes, &bodyRef)
	body["data"] = bodyRef

	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)
	log.Println("Exited route: Public.Stream.Get")
	return response
}
