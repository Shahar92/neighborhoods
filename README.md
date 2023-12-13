# neighborhoods

here is a link for the API documentation: https://app.swaggerhub.com/apis/SAMOYAL92/neighborhoods/1.0.0

Pre-Request: 
    1. PostgreSQL 
    2. Golang

additional info:
{
   PostgreSQL Installation:
   sudo apt update
   sudo apt upgrade
   sudo apt install postgresql postgresql-contrib 
}

all you need to do in order to run the server is to type the command: "go run ."

in order to check the test case run: "go test"

Performance Considerations:
    A. indexes are created on fields used in WHERE clauses for efficient querying.
    B. (caching) for frequently requested data: I will use library like Redis
    C. In general also we can: Optimize database queries. 


Scalability and Deployment:
   A. Scalability: using load balancers for high user loads. Such as: Nginx or AWS, GCP, Azure Load balancers.
   B. Deployment: Containerize my application using Docker for easier deployment and scaling


*For the init DB I will use: migrate with commands inside of a makefile