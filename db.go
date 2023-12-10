// db.go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass1234"
	dbname   = "mytestdb"
)

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", connStr)
}

func initDB() error {
	// Check if there are already records in the neighborhoods table
	hasRecords, err := HasRecordsInNeighborhoods()
	if err != nil {
		return err
	}

	// If there are records, return an error or handle it as needed
	if hasRecords {
		return errors.New("neighborhoods table already has records")
	}

	neighborhoods, err := ReadNeighborhoodsFromFile("./neighborhoods_data.json")
	if err != nil {
		return err
	}

	// Insert each neighborhood into the database
	_, err = InsertNeighborhoods(neighborhoods)
	if err != nil {
		return err
	}
	return nil
}

func InsertNeighborhoods(neighborhoods []Neighborhood) ([]Neighborhood, error) {
	// Build the VALUES clause for the bulk insert
	values := make([]interface{}, 0, len(neighborhoods)*9)
	for _, n := range neighborhoods {
		values = append(values,
			n.Name,
			n.City,
			n.State,
			n.AverageAge,
			n.DistanceFromCityCenter,
			n.AverageIncome,
			n.PublicTransportAvailability,
			n.Latitude,
			n.Longitude,
		)
	}

	// Build the SQL query with multiple sets of values
	query := `
		INSERT INTO neighborhoods (
			name, city, state, average_age, distance_from_city_center, average_income,
			public_transport_availability, latitude, longitude
		) VALUES
	`

	// Add multiple sets of values to the query
	for i := 0; i < len(neighborhoods); i++ {
		query += "($" + strconv.Itoa(i*9+1) + ", $" + strconv.Itoa(i*9+2) + ", $" + strconv.Itoa(i*9+3) + ", $" +
			strconv.Itoa(i*9+4) + ", $" + strconv.Itoa(i*9+5) + ", $" + strconv.Itoa(i*9+6) + ", $" + strconv.Itoa(i*9+7) +
			", $" + strconv.Itoa(i*9+8) + ", $" + strconv.Itoa(i*9+9) + "),"
	}

	// Remove the trailing comma
	query = query[:len(query)-1]

	// Execute the bulk insert query
	_, err := db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	// Fetch the updated list of neighborhoods from the database
	updatedNeighborhoods, err := QueryNeighborhoods("SELECT * FROM neighborhoods")
	if err != nil {
		return nil, err
	}

	return updatedNeighborhoods, nil
}

// QueryNeighborhoods fetches neighborhoods from the database
func QueryNeighborhoods(sqlQuery string) ([]Neighborhood, error) {

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var neighborhoods []Neighborhood
	for rows.Next() {
		var neighborhood Neighborhood
		err := rows.Scan(
			&neighborhood.ID,
			&neighborhood.Name,
			&neighborhood.City,
			&neighborhood.State,
			&neighborhood.AverageAge,
			&neighborhood.DistanceFromCityCenter,
			&neighborhood.AverageIncome,
			&neighborhood.PublicTransportAvailability,
			&neighborhood.Latitude,
			&neighborhood.Longitude,
		)
		if err != nil {
			return nil, err
		}
		neighborhoods = append(neighborhoods, neighborhood)
	}

	return neighborhoods, nil
}

// HasRecordsInNeighborhoods checks if there is at least one record in the neighborhoods table
func HasRecordsInNeighborhoods() (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM neighborhoods").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// func InsertNeighborhood(neighborhood Neighborhood) error {
// 	query := `
// 		INSERT INTO neighborhoods (
// 			name, city, state, average_age, distance_from_city_center, average_income,
// 			public_transport_availability, latitude, longitude
// 		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
// 	`

// 	_, err := db.Exec(query,
// 		neighborhood.Name,
// 		neighborhood.City,
// 		neighborhood.State,
// 		neighborhood.AverageAge,
// 		neighborhood.DistanceFromCityCenter,
// 		neighborhood.AverageIncome,
// 		neighborhood.PublicTransportAvailability,
// 		neighborhood.Latitude,
// 		neighborhood.Longitude,
// 	)

// 	return err
// }
