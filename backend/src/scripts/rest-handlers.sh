# Base URL
BASE_URL="http://localhost:4040"

# Apiary Operations
curl -X POST "${BASE_URL}/apiary" \
  -H "Content-Type: application/json" \
  -d '{"location": "Test Location", "manager_id": 1, "establishment_date": "2023-01-01T15:04:05Z"}'
curl -X GET "${BASE_URL}/apiary/1"
curl -X PUT "${BASE_URL}/apiary" \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "location": "Updated Location", "manager_id": 1, "establishment_date": "2023-01-01"}'
curl -X GET "${BASE_URL}/apiaries"
curl -X DELETE "${BASE_URL}/apiary/1"

# Hive Operations
curl -X GET "${BASE_URL}/hives"
curl -X POST "${BASE_URL}/hive" \
  -H "Content-Type: application/json" \
  -d '{"apiary_id": 1, "name": "Test Hive", "hive_type": "Langstroth", "current_status": "Active"}'
curl -X POST "${BASE_URL}/hive" \
  -H "Content-Type: application/json" \
  -d '{"apiary_id": 1, "name": "Test Second Hive", "hive_type": "Langstroth", "current_status": "Inactive"}'
curl -X PUT "${BASE_URL}/hive" \
  -H "Content-Type: application/json" \
  -d '{"hive_id": 1, "apiary_id": 1, "name": "Updated Hive", "hive_type": "Langstroth", "current_status": "Inactive", "installation_date": "2023-01-01T15:04:05Z"}'
curl -X DELETE "${BASE_URL}/hive/1"
curl -X GET "${BASE_URL}/apiary/1/hives"

# Bee Community Operations
curl -X GET "${BASE_URL}/bee-communities"
curl -X POST "${BASE_URL}/bee-community" \
  -H "Content-Type: application/json" \
  -d '{"hive_id": 2, "queen_age": 2, "population_estimate": 100, "health_status": "healthy"}'
curl -X POST "${BASE_URL}/bee-community" \
  -H "Content-Type: application/json" \
  -d '{"hive_id": 2, "queen_age": 10, "population_estimate": 100, "health_status": "healthy"}'
curl -X PUT "${BASE_URL}/bee-community" \
  -H "Content-Type: application/json" \
  -d '{"community_id": 1, "hive_id": 1, "queen_age": 11, "population_estimate": 1200, "health_status": "Thriving"}'
curl -X DELETE "${BASE_URL}/bee-community/1"
curl -X GET "${BASE_URL}/hive/1/bee-communities"


# Honey Harvest Operations

curl -X GET "${BASE_URL}/honey-harvest/1"
curl -X POST "${BASE_URL}/honey-harvest" \
  -H "Content-Type: application/json" \
  -d '{"hive_id": 1, "harvest_date": "2023-06-01T10:00:00Z", "quantity": 5.5, "quality": "High"}'
curl -X PUT "${BASE_URL}/honey-harvest" \
  -H "Content-Type: application/json" \
  -d '{"hive_id": 1, "hive_id": 1, "harvest_date": "2023-06-01T10:00:00Z", "quantity": 6.0, "quality": "Premium"}'
curl -X DELETE "${BASE_URL}/honey-harvest/{harvestId}"

# Get All Honey Harvests
curl -X GET "${BASE_URL}/honey-harvests"
curl -X GET "${BASE_URL}/veterinary-passport/1"
curl -X POST "${BASE_URL}/veterinary-passport" -H "Content-Type: application/json" -d '{"bee_community_id": 1, "issue_date": "2023-01-01T10:00:00Z", "expiry_date": "2024-01-01T10:00:00Z"}'
curl -X PUT "${BASE_URL}/veterinary-passport" -H "Content-Type: application/json" -d '{"passport_id":1, "bee_community_id": 1, "issue_date": "2023-01-01T10:00:00Z", "expiry_date": "2024-06-01T10:00:00Z", "health_status": "not dead yet"}'
curl -X DELETE "${BASE_URL}/veterinary-passport/1"
curl -X GET "${BASE_URL}/veterinary-passports"

# Region Operations
curl -X GET "${BASE_URL}/region/1"
curl -X POST "${BASE_URL}/region" -H "Content-Type: application/json" -d '{"name": "Test Region"}'
curl -X PUT "${BASE_URL}/region" -H "Content-Type: application/json" -d '{"region_id": 1, "name": "Updated Region"}'
curl -X DELETE "${BASE_URL}/region/1"
curl -X GET "${BASE_URL}/regions"

# AllowedRegion Operations
curl -X GET "${BASE_URL}/allowed-region/1"
curl -X POST "${BASE_URL}/allowed-region" -H "Content-Type: application/json" -d '{"region_id": 1, "user_id": 1}'
curl -X PUT "${BASE_URL}/allowed-region" -H "Content-Type: application/json" -d '{"id": 1, "region_id": 1, "user_id": 2}'
curl -X DELETE "${BASE_URL}/allowed-region/1"
curl -X GET "${BASE_URL}/allowed-regions"

# RegionApiary Operations
curl -X GET "${BASE_URL}/region-apiary/1"
curl -X POST "${BASE_URL}/region-apiary" -H "Content-Type: application/json" -d '{"region_id": 1, "apiary_id": 1}'
curl -X PUT "${BASE_URL}/region-apiary" -H "Content-Type: application/json" -d '{"id": 1, "region_id": 1, "apiary_id": 2}'curl -X PUT "${BASE_URL}/region-apiary" -H "Content-Type: application/json" -d '{"id": 1, "region_id": 1, "apiary_id": 1}'
curl -X DELETE "${BASE_URL}/region-apiary/1"
curl -X GET "${BASE_URL}/region-apiaries"

# ObservationLog Operations
curl -X GET "${BASE_URL}/observation-log/1"
curl -X POST "${BASE_URL}/observation-log" -H "Content-Type: application/json" -d '{"hive_id": 1, "observer_id": 1, "observation_date": "2023-06-01T10:00:00Z", "notes": "Healthy hive"}'
curl -X PUT "${BASE_URL}/observation-log" -H "Content-Type: application/json" -d '{"log_id": 1, "hive_id": 1, "observer_id": 1, "observation_date": "2023-06-01T10:00:00Z", "notes": "Updated observation"}'
curl -X DELETE "${BASE_URL}/observation-log/1"
curl -X GET "${BASE_URL}/observation-logs"

# MaintenancePlan Operations
curl -X GET "${BASE_URL}/maintenance-plan/1"
curl -X POST "${BASE_URL}/maintenance-plan" -H "Content-Type: application/json" -d '{"apiary_id": 1, "work_type": "Construction", "planned_date": "2025-07-01T10:00:00Z"}'
curl -X PUT "${BASE_URL}/maintenance-plan" -H "Content-Type: application/json" -d '{"plan_id": 1, "apiary_id": 1, "planned_date": "2027-07-01T10:00:00Z", "description": "Updated maintenance plan"}'
curl -X DELETE "${BASE_URL}/maintenance-plan/1"
curl -X GET "${BASE_URL}/maintenance-plans"

# Incident Operations
curl -X GET "${BASE_URL}/incident/1"
curl -X POST "${BASE_URL}/incident" -H "Content-Type: application/json" -d '{"hive_id": 1, "incident_date": "2023-06-15T14:00:00Z", "description": "Pest infestation"}'
curl -X PUT "${BASE_URL}/incident" -H "Content-Type: application/json" -d '{"incident_id": 1, "hive_id": 1, "incident_date": "2023-06-15T14:00:00Z", "description": "Updated incident description"}'
curl -X DELETE "${BASE_URL}/incident/1"
curl -X GET "${BASE_URL}/incidents"

# VeterinaryRecord Operations
curl -X GET "${BASE_URL}/veterinary-record/1"
curl -X POST "${BASE_URL}/veterinary-record" -H "Content-Type: application/json" -d '{"passport_id": 1, "record_date": "2023-06-01T10:00:00Z", "description": "Regular check-up"}'
curl -X PUT "${BASE_URL}/veterinary-record" -H "Content-Type: application/json" -d '{"record_id": 1, "passport_id": 1, "record_date": "2023-06-01T10:00:00Z", "description": "Updated check-up results"}'
curl -X DELETE "${BASE_URL}/veterinary-record/1"
curl -X GET "${BASE_URL}/veterinary-records"

# User Operations
curl -X GET "${BASE_URL}/user/1"
curl -X POST "${BASE_URL}/user" -H "Content-Type: application/json" -d '{"username": "John Doe", "email": "john@example.com", "role": "ADMIN", "last_login": "1970-01-01T00:00:00Z"}'
curl -X PUT "${BASE_URL}/user" -H "Content-Type: application/json" -d '{"user_id": 1, "username": "John Doe", "email": "john.doe@example.com", "role": "WORKER","last_login": "1970-01-01T00:00:00Z"}'
curl -X DELETE "${BASE_URL}/user/1"
curl -X GET "${BASE_URL}/users"

# ProductionReport Operations
curl -X GET "${BASE_URL}/production-report/1"
curl -X POST "${BASE_URL}/production-report" -H "Content-Type: application/json" -d '{"apiary_id": 1, "report_date": "2023-06-30T10:00:00Z", "honey_production": 100.5, "wax_production": 10.2}'
curl -X PUT "${BASE_URL}/production-report" -H "Content-Type: application/json" -d '{"id": 1, "apiary_id": 1, "report_date": "2023-06-30T10:00:00Z", "honey_production": 110.5, "wax_production": 12.2}'
curl -X DELETE "${BASE_URL}/production-report/1"
curl -X GET "${BASE_URL}/production-reports"

# Sensor Operations
curl -X GET "${BASE_URL}/sensor/1"
curl -X POST "${BASE_URL}/sensor" -H "Content-Type: application/json" -d '{"hive_id": 1, "sensor_type": "temperature", "installation_date": "2023-05-01T10:00:00Z"}'
curl -X PUT "${BASE_URL}/sensor" -H "Content-Type: application/json" -d '{"sensor_id": 1, "hive_id": 1, "sensor_type": "humidity", "installation_date": "2023-05-01T10:00:00Z"}'
curl -X DELETE "${BASE_URL}/sensor/1"
curl -X GET "${BASE_URL}/sensors"

# SensorReading Operations
curl -X GET "${BASE_URL}/sensor-reading/1"
curl -X POST "${BASE_URL}/sensor-reading" -H "Content-Type: application/json" -d '{"sensor_id": 1, "reading_time": "2023-06-01T12:00:00Z", "value": 25.5}'
curl -X PUT "${BASE_URL}/sensor-reading" -H "Content-Type: application/json" -d '{"id": 1, "sensor_id": 1, "reading_time": "2023-06-01T12:00:00Z", "value": 26.0}'
curl -X DELETE "${BASE_URL}/sensor-reading/1"
curl -X GET "${BASE_URL}/sensor-readings"

# WeatherData Operations
curl -X GET "${BASE_URL}/weather-data/1"
curl -X POST "${BASE_URL}/weather-data" -H "Content-Type: application/json" -d '{"region_id": 1, "date": "2023-06-01T12:00:00Z", "temperature": 25.5, "humidity": 60, "wind_speed": 5.2, "precipitation": 0.5}'
curl -X PUT "${BASE_URL}/weather-data" -H "Content-Type: application/json" -d '{"id": 1, "apiary_id": 1, "record_time": "2023-06-01T12:00:00Z", "temperature": 26.0, "humidity": 62, "wind_speed": 5.5}'
curl -X DELETE "${BASE_URL}/weather-data/1"
curl -X GET "${BASE_URL}/weather-data-all"
