package admin

import (
	"encoding/json"
	"log"
	"shrampybot/router"
	"shrampybot/utility/nosqldb"
)

type CategoryView struct {
	router.View
}

type CategoryBody struct {
	router.GenericBodyDataFlat
	Data *[]nosqldb.CategoryDatum `json:"data"`
}

func NewCategoryView() *CategoryView {
	c := CategoryView{}
	return &c
}

func (v *CategoryView) CallMethod(route *router.Route) *router.Response {
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

func (v *CategoryView) Get(route *router.Route) *router.Response {
	log.Println("Entered route: Admin.Category.Get")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	categories := []nosqldb.CategoryDatum{}
	if len(route.Path) == 3 {
		// Get single result
		category, err := n.GetCategory(route.Path[2])
		if err != nil {
			log.Println("Get category failed.")
			response.StatusCode = "500"
			return response
		}
		categories = append(categories, *category)
	} else {
		// Fetch categories
		categoriesRef, err := n.GetCategoryMap()
		if err != nil || categoriesRef == nil {
			log.Println("Could not get category map.")
			response.StatusCode = "500"
			return response
		}
		categories = *categoriesRef
	}

	catBytes, _ := json.Marshal(categories)

	body := map[string]any{}
	body["count"] = len(categories)
	bodyRef := []map[string]any{}
	json.Unmarshal(catBytes, &bodyRef)
	body["data"] = bodyRef

	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Category.Get")
	return response
}

func (v *CategoryView) Post(route *router.Route) *router.Response {
	var err error
	log.Println("Entered route: Admin.Category.Post")
	response := &router.Response{}
	response.Headers = &router.DefaultResponseHeaders

	// Parse submitted category data
	requestBody := CategoryBody{}
	json.Unmarshal([]byte(route.Body), &requestBody)

	// Instantiate DynamoDB
	n, err := nosqldb.NewClient()
	if err != nil {
		log.Println("Could not instantiate dynamodb.")
		response.StatusCode = "500"
		return response
	}

	// Get existing list so we can remove defunct entries.
	existingCategories, err := n.GetCategoryMap()
	if err != nil {
		log.Println("Could not retrieve category_map.")
		response.StatusCode = "500"
		return response
	}

	// Make list of defunct entries
	removeIds := []string{}
	for _, extCat := range *existingCategories {
		foundMatch := false

		for elem, newCat := range *requestBody.Data {
			if extCat.TwitchCategory == newCat.TwitchCategory {
				foundMatch = true
				(*requestBody.Data)[elem].Id = extCat.Id
			}
		}

		if !foundMatch {
			removeIds = append(removeIds, extCat.Id)
		}
	}

	// Remove defunct entries
	err = n.RemoveCategory(&removeIds)
	if err != nil {
		log.Println("Could not remove categories.")
		response.StatusCode = "500"
		return response
	}

	// Add/update category map
	err = n.PutCategories(requestBody.Data)
	if err != nil {
		log.Println("Could not save categories.")
		response.StatusCode = "500"
		return response
	}

	// Re-fetch categories
	categories, err := n.GetCategoryMap()
	if err != nil {
		log.Println("Could not get saved category_map.")
		response.StatusCode = "500"
		return response
	}
	catBytes, _ := json.Marshal(categories)

	body := map[string]any{}
	body["count"] = len(*categories)
	bodyRef := []map[string]any{}
	json.Unmarshal(catBytes, &bodyRef)
	body["data"] = bodyRef

	bodyBytes, _ := json.Marshal(body)

	response.StatusCode = "200"
	response.Body = string(bodyBytes)

	log.Println("Exited route: Admin.Category.Post")
	return response
}
