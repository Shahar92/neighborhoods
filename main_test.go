// main_test.go
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func generateRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func TestGetNeighborhoods(t *testing.T) {

	// Generate Random numbers for testing
	ageRangeA := generateRandomNumber(18, 70)
	ageRangeB := generateRandomNumber(18, 70)
	maxDistance := generateRandomNumber(0, 100)
	minAge := ageRangeA
	maxAge := ageRangeB

	if ageRangeA > ageRangeB {
		minAge = ageRangeB
		maxAge = ageRangeA
	}

	requestURL := fmt.Sprintf("/neighborhoods?ageRange=%d,%d&maxDistance=%d&sortBy=average_age,asc", minAge, maxAge, maxDistance)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetNeighborhoods)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"
	actualContentType := rr.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v",
			actualContentType, expectedContentType)
	}

	// Parse the response body to check individual neighborhoods
	var neighborhoods []Neighborhood
	err = json.Unmarshal(rr.Body.Bytes(), &neighborhoods)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}

	// Check that neighborhoods are within the specified age range
	for _, n := range neighborhoods {
		if n.AverageAge < minAge || n.AverageAge > maxAge {
			t.Errorf("Neighborhood outside of age range: %v", n)
		}
	}

	// Check that neighborhoods are within the specified distance
	for _, n := range neighborhoods {
		if n.DistanceFromCityCenter > float64(maxDistance) {
			t.Errorf("Neighborhood outside of distance range: %v", n)
		}
	}

	// Check that neighborhoods are sorted by average age in ascending order
	for i := 0; i < len(neighborhoods)-1; i++ {
		if neighborhoods[i].AverageAge > neighborhoods[i+1].AverageAge {
			t.Errorf("Neighborhoods are not sorted correctly")
		}
	}
}
