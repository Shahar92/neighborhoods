package main

import (
	"encoding/json"
	"io/ioutil"
)

// ReadNeighborhoodsFromFile reads neighborhoods from a JSON file and inserts them into the database
func ReadNeighborhoodsFromFile(filePath string) ([]Neighborhood, error) {

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into a slice of Neighborhood objects
	var neighborhoods []Neighborhood
	err = json.Unmarshal(jsonData, &neighborhoods)
	if err != nil {
		return nil, err
	}

	return neighborhoods, nil
}
