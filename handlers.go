// handlers.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetNeighborhoods handles the GET request to fetch neighborhoods based on specified criteria.
// It supports filtering and sorting based on parameters in the query string.
//
// Endpoint: GET /neighborhoods
//
// Parameters:
//   - ageRange: Filter neighborhoods by age range. Format: "minAge,maxAge".
//   - maxDistance: Filter neighborhoods by maximum distance from the city center in kilometers.
//   - sortBy: Sort neighborhoods by a specific field. Possible values: "name", "average_age", "distance_from_city_center", "average_income" , etc...
//
// Example Request:
//
//	GET /neighborhoods?ageRange=20,40&maxDistance=10&sortBy=name,asc
//
// Example Response:
//
//	[
//	  {
//	    "name": "Example Neighborhood",
//	    "city": "Example City",
//	    "state": "EX",
//	    "average_age": 30,
//	    "distance_from_city_center": 5.2,
//	    "average_income": 50000,
//	    "public_transport_availability": "high",
//	    "latitude": 40.7128,
//	    "longitude": -74.0060
//	  },...
//	]
func GetNeighborhoods(w http.ResponseWriter, r *http.Request) {

	// Fetch neighborhoods from the database
	neighborhoodsQuery, err := createNeighborhoodsQueryFromParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	neighborhoods, err := QueryNeighborhoods(neighborhoodsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neighborhoods)
}

func createNeighborhoodsQueryFromParams(params url.Values) (string, error) {
	// Parse query parameters
	ageRange := params.Get("ageRange")
	maxDistance := params.Get("maxDistance")
	sortBy := params.Get("sortBy")

	// Build the SQL query based on parameters
	sqlQuery := "SELECT * FROM neighborhoods WHERE 1=1"

	// Age range filter
	if ageRange != "" {

		ageRangeValues := strings.Split(ageRange, ",")
		if len(ageRangeValues) != 2 {
			return "", errors.New("invalid age range format, should contain 2 params: min_age, max_age")
		}

		minAge, err := strconv.Atoi(ageRangeValues[0])
		if err != nil {
			return "", errors.New("invalid min age range")
		}

		maxAge, err := strconv.Atoi(ageRangeValues[1])
		if err != nil {
			return "", errors.New("invalid max age range")
		}

		sqlQuery += fmt.Sprintf(" AND average_age BETWEEN %d AND %d", minAge, maxAge)
	}

	// Max distance filter
	if maxDistance != "" {
		maxDistanceValue, err := strconv.ParseFloat(maxDistance, 64)
		if err != nil {
			return "", errors.New("invalid max age range")
		}
		sqlQuery += fmt.Sprintf(" AND distance_from_city_center <= %f", maxDistanceValue)
	}

	// Sorting
	if sortBy != "" {
		sortByValues := strings.Split(sortBy, ",")
		if len(sortByValues) != 2 || !isExists([]string{"asc", "desc"}, strings.ToLower(sortByValues[1])) {
			return "", errors.New("invalid age sortBy format, should contain 2 params: field, order_direction(ASC/DESC)")
		}
		cols := []string{"name", "city", "state", "average_age", "distance_from_city_center", "average_income", "public_transport_availability", "latitude", "longitude"}
		if !isExists(cols, sortByValues[0]) {
			return "", fmt.Errorf("invalid field relevant ones can be: [%s]", strings.Join(cols, ","))
		}

		sqlQuery += fmt.Sprintf(" ORDER BY %s %s", sortByValues[0], sortByValues[1])
	}

	return sqlQuery, nil

}
