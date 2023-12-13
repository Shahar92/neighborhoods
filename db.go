// db.go
package main

import (
	"database/sql"
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

// doesTableExist checks if a table with the given name exists in the database
func doesTableExist(tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM   information_schema.tables
			WHERE  table_name = $1
		)
	`

	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func buildNeighborhoodTable() error {
	_, err := db.Exec(`
		CREATE TABLE neighborhoods (
			ID SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			city VARCHAR(255) NOT NULL,
			state VARCHAR(255),
			average_age INTEGER NOT NULL,
			distance_from_city_center FLOAT NOT NULL,
			average_income INTEGER NOT NULL,
			public_transport_availability VARCHAR(255) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL
		);
	`)

	if err != nil {
		return err
	}

	// Add indexes to improve query performance
	_, err = db.Exec(`
		CREATE INDEX idx_average_age ON neighborhoods (average_age);
		CREATE INDEX idx_distance_from_city_center ON neighborhoods (distance_from_city_center);
		CREATE INDEX idx_average_income ON neighborhoods (average_income);
		CREATE INDEX idx_latitude_longitude ON neighborhoods (latitude, longitude);
	`)

	return err
}

func initDB() error {

	// Check if the neighborhoods table exists
	tableExists, err := doesTableExist("neighborhoods")
	if err != nil {
		return err
	}

	if tableExists {
		return nil
	}

	if err = buildNeighborhoodTable(); err != nil {
		return err
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
