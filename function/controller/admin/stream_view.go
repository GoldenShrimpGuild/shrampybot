package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"time"
)

type StreamView struct {
	router.View `tstype:",extends,required"`
}

type StreamStatusPutRequest struct {
	EndNow bool `json:"endNow,omitempty"`
}

type StreamPutResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
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

func (v *StreamView) Put(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Stream.Put")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	responseBody := StreamPutResponse{}

	if len(route.Path) < 4 {
		log.Printf("Insufficient path length: %v\n", len(route.Path))
		response.StatusCode = "400"
		return response
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	if route.Path[2] == "status" {
		stream, err := n.GetStream(route.Path[3])
		if err != nil {
			log.Printf("Specified stream %v not found.", route.Path[3])
		}

		requestBody := StreamStatusPutRequest{}
		err = json.Unmarshal([]byte(route.Body), &requestBody)
		if err != nil {
			log.Printf("Could not unmarshal body json: %v\n", err)
			response.StatusCode = "500"
			return response
		}

		responseBody.Id = stream.ID
		responseBody.Status = router.StatusUnknown.String()

		if requestBody.EndNow {
			if stream.EndedAt.After(time.Time{}) {
				log.Printf("Already ended stream %v at %v, do nothing.", stream.ID, stream.EndedAt.Format(time.RFC3339))
				responseBody.Status = router.StatusNotNeeded.String()
			} else {
				stream.EndedAt = time.Now()
				log.Printf("Setting stream %v end time to %v.", stream.ID, stream.EndedAt.Format(time.RFC3339))

				err = n.PutStream(stream)
				if err != nil {
					log.Printf("Failed to set end time on stream.")
					responseBody.Status = router.StatusFailure.String()
				} else {
					responseBody.Status = router.StatusSuccess.String()
				}
			}
		}
	} else {
		log.Printf("Invalid path: %v\n", route.Path[2])
		response.StatusCode = "400"
		return response
	}

	response.StatusCode = "200"
	bodyBytes, _ := json.Marshal(responseBody)
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Stream.Put")
	return response
}
