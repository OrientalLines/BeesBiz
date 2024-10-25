#!/bin/bash
BASE_URL="http://localhost:4040"
# Check if BASE_URL is set
if [ -z "$BASE_URL" ]; then
    echo "Error: BASE_URL environment variable is not set"
    exit 1
fi

# Loop 100 times
for i in {1..100}; do
    echo "Making request $i of 100..."
    # curl -X POST "${BASE_URL}/user" -H "Content-Type: application/json" -d "{\"username\": \"user$i\", \"email\": \"user$i@example.com\", \"role\": \"ADMIN\", \"last_login\": \"2023-01-01T00:00:00Z\"}"

    echo -e "\n"

    # Optional: Add a small delay to prevent overwhelming the server
    sleep 0.01
done

echo "Completed all 100 requests"
