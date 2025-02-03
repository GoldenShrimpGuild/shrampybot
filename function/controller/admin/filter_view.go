package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
	"strings"
)

type FilterView struct {
	router.View `tstype:",extends,required"`
}

type FilterBody struct {
	router.GenericBodyDataFlat `tstype:",extends,required"`
	Data                       []*nosqldb.FilterDatum `json:"data"`
}

func NewFilterView() *FilterView {
	c := FilterView{}
	return &c
}

func (v *FilterView) CallMethod(route *router.Route) *router.Response {
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

func (v *FilterView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.Filter.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	filterItems, err := n.GetFilterKeywords()
	if err != nil {
		log.Println("Could not load filter keywords from db.")
		response.StatusCode = "500"
		return response
	}

	respBody := FilterBody{}
	respBody.Count = len(filterItems)
	respBody.Data = filterItems

	bodyBytes, err := json.Marshal(respBody)
	if err != nil {
		log.Println("Could not marshal body bytes from body json.")
		response.StatusCode = "500"
		return response
	}

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Filter.Get")
	return response
}

func (v *FilterView) Put(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Filter.Put")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	requestBody := nosqldb.FilterDatum{}
	err = json.Unmarshal([]byte(route.Body), &requestBody)
	if err != nil {
		log.Printf("Could not unmarshal body json: %v\n", err)
		response.StatusCode = "500"
		return response
	}

	requestBody.Id = strings.Trim(requestBody.Id, " ")

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	err = n.FillFilterIdIfAny(&requestBody)
	if err != nil {
		log.Println("Error searching for keyword.")
		response.StatusCode = "500"
		return response
	}

	kwList := append([]*nosqldb.FilterDatum{}, &requestBody)
	err = n.PutFilterKeywords(kwList)
	if err != nil {
		log.Println("Could not save filters.")
		response.StatusCode = "500"
		return response
	}

	body := FilterBody{}
	body.Count = 1
	body.Data = append(body.Data, &requestBody)
	bodyBytes, _ := json.Marshal(body)

	response.Body = string(bodyBytes)
	response.StatusCode = "200"

	log.Println("Exited route: Admin.Filter.Put")
	return response
}

func (v *FilterView) Post(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Filter.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Parse submitted filter data
	requestBody := FilterBody{}
	json.Unmarshal([]byte(route.Body), &requestBody)

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Get existing list so we can remove defunct entries.
	existingFilterKeywords, err := n.GetFilterKeywords()
	if err != nil {
		log.Println("Could not retrieve filter keywords.")
		response.StatusCode = "500"
		return response
	}

	// Make list of defunct entries
	removeIds := []string{}
	for _, extFK := range existingFilterKeywords {
		foundMatch := false

		for elem, newFK := range requestBody.Data {
			if extFK.Keyword == newFK.Keyword {
				foundMatch = true
				(requestBody.Data)[elem].Id = extFK.Id
			}
		}

		if !foundMatch {
			removeIds = append(removeIds, extFK.Id)
		}
	}

	// Remove defunct entries
	err = n.RemoveFilterKeyword(&removeIds)
	if err != nil {
		log.Println("Could not remove filter keyword.")
		response.StatusCode = "500"
		return response
	}

	// Add/update category map
	err = n.PutFilterKeywords(requestBody.Data)
	if err != nil {
		log.Println("Could not save filter keywords.")
		response.StatusCode = "500"
		return response
	}

	// Re-fetch filter keywords
	filterKeywords, err := n.GetFilterKeywords()
	if err != nil {
		log.Println("Could not retrieve filter keywords.")
		response.StatusCode = "500"
		return response
	}
	fkBytes, _ := json.Marshal(filterKeywords)

	body := map[string]any{}
	body["count"] = len(filterKeywords)
	bodyRef := []map[string]any{}
	json.Unmarshal(fkBytes, &bodyRef)
	body["data"] = bodyRef

	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Filter.Post")
	return response
}

func (v *FilterView) Delete(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Filter.Delete")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	if len(route.Path) == 3 {
		ids := append([]string{}, route.Path[2])

		err := n.RemoveFilterKeyword(&ids)
		if err != nil {
			log.Println("Remove filter failed.")
			response.StatusCode = "500"
			return response
		}
	} else {
		log.Println("No ID specified.")
		response.StatusCode = "400"
		return response
	}

	response.StatusCode = "200"
	log.Println("Exited route: Admin.Filter.Delete")
	return response
}
