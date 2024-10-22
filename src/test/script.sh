# Base URL
BASE_URL="http://localhost:4040"

# Apiary Operations

# Create Apiary
curl -X POST "${BASE_URL}/apiary" \
  -H "Content-Type: application/json" \
  -d '{"location": "Test Location", "manager_id": 1, "establishment_date": "2023-01-01"}'

# Get Apiary (replace {apiaryId} with the actual ID)
curl -X GET "${BASE_URL}/apiary/1"

# Update Apiary (replace {apiaryId} with the actual ID)
curl -X PUT "${BASE_URL}/apiary" \
  -H "Content-Type: application/json" \
  -d '{"id": {apiaryId}, "location": "Updated Location", "manager_id": 1, "establishment_date": "2023-01-01"}'

# Get All Apiaries
curl -X GET "${BASE_URL}/apiaries"

# Delete Apiary (replace {apiaryId} with the actual ID)
curl -X DELETE "${BASE_URL}/apiary/1"

# User Operations

# Create User
curl -X POST "${BASE_URL}/user" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "email": "testuser@example.com", "role": "admin"}'

# Get User (replace {userId} with the actual ID)
curl -X GET "${BASE_URL}/user/{userId}"

# Update User (replace {userId} with the actual ID)
curl -X PUT "${BASE_URL}/user" \
  -H "Content-Type: application/json" \
  -d '{"id": {userId}, "name": "Updated User", "email": "testuser@example.com", "role": "admin"}'

# Get All Users
curl -X GET "${BASE_URL}/users"

# Delete User (replace {userId} with the actual ID)
curl -X DELETE "${BASE_URL}/user/{userId}"
