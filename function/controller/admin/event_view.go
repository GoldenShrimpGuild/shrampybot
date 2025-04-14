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
	Meta EventApiResponseMeta `json:"meta"`
}

type EventApiResponseMeta struct {
	ApiVersion  string `json:"apiVersion"`
	Environment string `json:"environment"`
	EventId     string `json:"eventId"`
	GeneratedOn string `json:"generatedOn"`
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

	if len(route.Path) != 3 {
		log.Println("No ID specified.")
		response.StatusCode = "400"
		return response
	}

	err := sendSignedCacheInvalidation(route.Path[2])
	if err != nil {
		log.Printf("Error invalidating cache: %v", err)
		response.StatusCode = "500"
		return response
	}

	log.Printf("Successfully invalidated cache for event %v", route.Path[2])
	response.StatusCode = "200"

	log.Println("Exited route: Admin.Event.Get")
	return response
}

// Reference code examples:
// https://github.com/aws-samples/sigv4-signing-examples/blob/main/sdk/golang/main.go
// https://gist.github.com/secretorange/905b4811300d7c96c71fa9c6d115ee24

func sendSignedCacheInvalidation(eventId string) error {
	ctx := context.Background()
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
		return err
	}

	req.Host = config.EventApiHost
	req.Header.Add("Cache-Control", "max-age=0")

	signer := v4.NewSigner()
	if err = signer.SignHTTP(ctx, creds, req, payloadHash, config.EventApiService, config.EventApiRegion, time.Now()); err != nil {
		log.Printf("Error signing request: %v", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		err = errors.New("request failed with non-200 status code")
		return err
	}

	warning := resp.Header.Get("Warning")
	if strings.Contains(warning, "199") {
		err = errors.New("error 199 received from remote url")
		return err
	}

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return err
	}
	parsedBody := EventApiResponse{}
	json.Unmarshal(body, &parsedBody)

	if eventId != parsedBody.Meta.EventId {
		err = errors.New("no matching eventId found in response")
		return err
	}

	return nil
}
