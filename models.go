// models.go
package main

type Neighborhood struct {
	ID                          int     `json:"id"`
	Name                        string  `json:"neigborhood"` // Note the struct tag here
	City                        string  `json:"city"`
	State                       *string `json:"state"`
	AverageAge                  int     `json:"average age"`
	DistanceFromCityCenter      float64 `json:"distance from city center"`
	AverageIncome               int     `json:"average income"`
	PublicTransportAvailability string  `json:"public transport availability"`
	Latitude                    float64 `json:"latitude"`
	Longitude                   float64 `json:"longitude"`
}
