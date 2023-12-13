-- Query to filter neighborhoods by age range, distance, and sort by average income
SELECT *
FROM neighborhoods
WHERE average_age BETWEEN 20 AND 40
    AND distance_from_city_center <= 10
ORDER BY average_income DESC;

-- Query to calculate the average age of all neighborhoods
SELECT AVG(average_age) AS average_age
FROM neighborhoods;

-- Query to filter neighborhoods by city and sort by distance from the city center
SELECT *
FROM neighborhoods
WHERE city = 'Los Dos Caminos'
ORDER BY distance_from_city_center;
