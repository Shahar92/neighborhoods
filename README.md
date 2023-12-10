# neighborhoods


PostgreSQL Installation:
sudo apt update
sudo apt upgrade
sudo apt install postgresql postgresql-contrib 



-------  TODO: put this on endpoint url or on startup of system -------
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

