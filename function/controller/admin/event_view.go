package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"shrampybot/config"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"strings"
	"time"

	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type EventView struct {
	router.View `tstype:",extends,required"`
}

type EventApiResponse struct {
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Meta        EventApiResponseMeta `json:"meta"`
	StatusCode  int                  `json:"_statusCode,omitempty"`
}

type EventApiResponseMeta struct {
	ApiVersion  string `json:"apiVersion,omitempty"`
	Environment string `json:"environment,omitempty"`
	EventId     string `json:"eventId,omitempty"`
	GeneratedOn string `json:"generatedOn,omitempty"`
}

type EventGetResponseBody struct {
	Status         string                               `json:"status,omitempty"`
	EventResponse  *EventApiResponse                    `json:"eventResponse,omitempty"`
	IsCurrentEvent bool                                 `json:"isCurrentEvent"`
	CurrentEvent   *CurrentEventGetEventGetResponseBody `json:"currentEvent,omitempty"`
}

type CurrentEventGetEventGetResponseBody struct {
	Status        string            `json:"status,omitempty"`
	EventResponse *EventApiResponse `json:"eventResponse,omitempty"`
}

func NewEventView() *EventView {
	c := EventView{}
	return &c
}

func (v *EventView) CallMethod(route *router.Route) *router.Response {
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

func (c *EventView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.Event.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders
	response.StatusCode = "200"

	if len(route.Path) != 3 {
		log.Println("No ID specified.")
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

	inputEvent := strings.TrimSpace(route.Path[2])

	responseBody := EventGetResponseBody{}
	responseBody.Status = router.StatusText[router.StatusUnknown]
	responseBody.IsCurrentEvent = false
	responseBody.CurrentEvent = nil
	responseBody.EventResponse, err = sendSignedCacheInvalidation(inputEvent)
	if err != nil {
		log.Printf("Could not invalidate cache for event %v: %v", inputEvent, err)
		responseBody.Status = router.StatusText[router.StatusFailure]
	} else {
		log.Printf("Successfully invalidated cache for event %v", route.Path[2])

		responseBody.Status = router.StatusText[router.StatusSuccess]
		// Defaulting to event 1 based on current single-current-event behaviour
		currentEvent, err := n.GetCurrentEvent(1)
		if err != nil {
			log.Printf("Could not retrieve current event by index %v: %v", 1, err)
		} else {
			if currentEvent.EventId == inputEvent {
				responseBody.IsCurrentEvent = true
				responseBody.CurrentEvent = &CurrentEventGetEventGetResponseBody{}
				responseBody.CurrentEvent.EventResponse, err = sendSignedCacheInvalidation("current")
				if err != nil {
					responseBody.CurrentEvent.Status = router.StatusText[router.StatusFailure]
					log.Printf("Failed to reset current event cache: %v", err)
				} else {
					responseBody.CurrentEvent.Status = router.StatusText[router.StatusSuccess]
					log.Printf("Successfully invalidated cache for current event.")
				}
			} else {
				log.Printf("Event %v is not current event (%v), nothing to be done.", inputEvent, currentEvent.EventId)
			}
		}
	}

	bodyBytes, _ := json.Marshal(responseBody)
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Event.Get")
	return response
}

// Reference code examples:
// https://github.com/aws-samples/sigv4-signing-examples/blob/main/sdk/golang/main.go
// https://gist.github.com/secretorange/905b4811300d7c96c71fa9c6d115ee24

func sendSignedCacheInvalidation(eventId string) (*EventApiResponse, error) {
	ctx := context.Background()
	apiResponse := EventApiResponse{}
	payloadHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	creds := aws.Credentials{
		AccessKeyID:     config.AwsAccessKeyId,
		SecretAccessKey: config.AwsSecretAccessKey,
		SessionToken:    config.AwsSessionToken,
	}

	url := "https://" + config.EventApiHost + config.EventApiPath + eventId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request to %v, %v", url, err)
		return &apiResponse, err
	}

	req.Host = config.EventApiHost
	req.Header.Add("Cache-Control", "max-age=0")

	signer := v4.NewSigner()
	if err = signer.SignHTTP(ctx, creds, req, payloadHash, config.EventApiService, config.EventApiRegion, time.Now()); err != nil {
		log.Printf("Error signing request: %v", err)
		return &apiResponse, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return &apiResponse, err
	}
	defer resp.Body.Close()

	apiResponse.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		err = fmt.Errorf("request failed with status code %v", resp.StatusCode)
		return &apiResponse, err
	}

	warning := resp.Header.Get("Warning")
	if strings.Contains(warning, "199") {
		apiResponse.StatusCode = 199
		err = errors.New("error 199 received from remote url")
		return &apiResponse, err
	}

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return &apiResponse, err
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling remote body json: %v", err)
		return &apiResponse, err
	}

	apiResponse.StatusCode = resp.StatusCode
	return &apiResponse, err
}
