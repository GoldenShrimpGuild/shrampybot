package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"strconv"
)

type CurrentEventView struct {
	router.View `tstype:",extends,required"`
}

type CurrentEventPutRequestBody struct {
	EventId string `json:"eventId"`
}

type CurrentEventGetResponseBody struct {
	Status       string                    `json:"status,omitempty"`
	CurrentEvent nosqldb.CurrentEventDatum `json:"currentEvent,omitempty"`
}

type CurrentEventPutResponseBody struct {
	Status         string                    `json:"status,omitempty"`
	CurrentEvent   nosqldb.CurrentEventDatum `json:"currentEvent"`
	RetrievedEvent EventApiResponse          `json:"retrievedEvent"`
}

type CurrentEventDeleteResponseBody struct {
	Status string `json:"status,omitempty"`
}

func NewCurrentEventView() *CurrentEventView {
	c := CurrentEventView{}
	return &c
}

func (v *CurrentEventView) CallMethod(route *router.Route) *router.Response {
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

func (c *CurrentEventView) Get(route *router.Route) *router.Response {
	var currentEventIndex int64
	var err error
	log.Println("Entered route: Admin.CurrentEvent.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	response.StatusCode = "200"

	if len(route.Path) != 3 {
		currentEventIndex = 1
	} else {
		currentEventIndex, _ = strconv.ParseInt(route.Path[2], 10, 8)
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	responseBody := CurrentEventGetResponseBody{}
	responseBody.Status = router.StatusText[router.StatusUnknown]

	currentEvent, err := n.GetCurrentEvent(uint8(currentEventIndex))
	if err != nil {
		log.Printf("Could not retrieve current event by index %v: %v", currentEventIndex, err)
		responseBody.Status = router.StatusText[router.StatusFailure]
	} else {
		if currentEvent.EventId != "" {
			responseBody.Status = router.StatusText[router.StatusSuccess]
			responseBody.CurrentEvent = *currentEvent
		} else {
			responseBody.Status = router.StatusText[router.StatusNotFound]
		}
	}

	bodyBytes, _ := json.Marshal(responseBody)
	response.Body = string(bodyBytes)
	log.Println("Exited route: Admin.CurrentEvent.Get")
	return response
}

func (c *CurrentEventView) Put(route *router.Route) *router.Response {
	var err error
	var currentEventIndex int64
	log.Println("Entered route: Admin.CurrentEvent.Put")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	response.StatusCode = "200"

	if len(route.Path) != 3 {
		currentEventIndex = 1
	} else {
		currentEventIndex, _ = strconv.ParseInt(route.Path[2], 10, 8)
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	requestBody := CurrentEventPutRequestBody{}
	err = json.Unmarshal([]byte(route.Body), &requestBody)
	if err != nil {
		log.Printf("Could not unmarshal body json: %v\n", err)
		response.StatusCode = "500"
		return response
	}
	log.Printf("Requested Event ID is \"%v\"", requestBody.EventId)

	responseBody := CurrentEventPutResponseBody{}
	responseBody.Status = router.StatusText[router.StatusUnknown]
	retEvent, err := sendSignedCacheInvalidation(requestBody.EventId)
	responseBody.RetrievedEvent = *retEvent
	if err != nil {
		responseBody.Status = router.StatusText[router.StatusFailure]
		log.Printf("Error retrieving remote event with ID %v: %v", requestBody.EventId, err)
	} else {
		log.Printf("Retrieved event ID is \"%v\"", responseBody.RetrievedEvent.Meta.EventId)
		// Only set the current event if the retrieved EventId aligns
		if requestBody.EventId == responseBody.RetrievedEvent.Meta.EventId {
			currentEvent := nosqldb.CurrentEventDatum{
				EventId:     requestBody.EventId,
				Title:       responseBody.RetrievedEvent.Title,
				Description: responseBody.RetrievedEvent.Description,
			}
			if n.PutCurrentEvent(uint8(currentEventIndex), &currentEvent) == nil {
				responseBody.Status = router.StatusText[router.StatusSuccess]
				responseBody.CurrentEvent = currentEvent
			}
		}
	}

	// Invalidate the cache on the "current" event
	// Ignore all errors and just continue for now.
	_, _ = sendSignedCacheInvalidation("current")

	bodyBytes, _ := json.Marshal(responseBody)
	response.Body = string(bodyBytes)
	log.Println("Exited route: Admin.CurrentEvent.Put")
	return response
}

func (c *CurrentEventView) Delete(route *router.Route) *router.Response {
	var currentEventIndex int64
	var err error
	log.Println("Entered route: Admin.CurrentEvent.Delete")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	response.StatusCode = "200"

	if len(route.Path) != 3 {
		currentEventIndex = 1
	} else {
		currentEventIndex, _ = strconv.ParseInt(route.Path[2], 10, 8)
	}

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	responseBody := CurrentEventDeleteResponseBody{}
	responseBody.Status = router.StatusText[router.StatusUnknown]

	err = n.DeleteCurrentEvent(uint8(currentEventIndex))
	if err != nil {
		log.Printf("Could not delete current event by index %v: %v", currentEventIndex, err)
		responseBody.Status = router.StatusText[router.StatusFailure]
	} else {
		responseBody.Status = router.StatusText[router.StatusSuccess]
	}

	bodyBytes, _ := json.Marshal(responseBody)
	response.Body = string(bodyBytes)
	log.Println("Exited route: Admin.CurrentEvent.Delete")
	return response
}

// func retrieveEventApiResponse(eventId string) (*EventApiResponse, error) {
// 	apiResponse := EventApiResponse{}

// 	url := "https://" + config.EventApiHost + config.EventApiPath + eventId

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Printf("Error creating request to %v, %v", url, err)
// 		return &apiResponse, err
// 	}
// 	req.Host = config.EventApiHost

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Error making request: %v", err)
// 		return &apiResponse, err
// 	}
// 	defer resp.Body.Close()

// 	apiResponse.StatusCode = resp.StatusCode
// 	if resp.StatusCode != 200 {
// 		log.Printf("Request failed with status code: %v", resp.StatusCode)
// 		err = fmt.Errorf("request failed with status code %v", resp.StatusCode)
// 		return &apiResponse, err
// 	}

// 	body, err := io.ReadAll(io.Reader(resp.Body))
// 	if err != nil {
// 		fmt.Printf("Error reading response body: %v", err)
// 		return &apiResponse, err
// 	}
// 	err = json.Unmarshal(body, &apiResponse)
// 	if err != nil {
// 		fmt.Printf("Error unmarshalling remote body json: %v", err)
// 		return &apiResponse, err
// 	}

// 	apiResponse.StatusCode = resp.StatusCode
// 	return &apiResponse, nil
// }
