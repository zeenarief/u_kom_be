#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
EMAIL="admin@example.com"
PASSWORD="password" # Reset to 'password'

# Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"login\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ]; then
  echo "Login failed. Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "Login successful. Token: ${TOKEN:0:10}..."

# Helper function for searching
search() {
  ENDPOINT=$1
  QUERY=$2
  echo -e "\nSearching $ENDPOINT with query '$QUERY'..."
  curl -s -X GET "$BASE_URL/$ENDPOINT?q=$QUERY" \
    -H "Authorization: Bearer $TOKEN" | jq '.data[] | {id, name: (.name // .full_name), email, username, nisn, nip, code}'
}

# Test Users
search "users" "admin"

# Test Students
# Assuming some data exists. If not, results will be empty but no error.
search "students" "a"

# Test Parents
search "parents" "a"

# Test Guardians
search "guardians" "a"

# Test Employees
search "employees" "a"

# Test Subjects
search "subjects" "a"
