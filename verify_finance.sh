#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD="admin123" # Updated based on previous context

echo "========================================================"
echo "      VERIFYING FINANCE DONATION MODULE"
echo "========================================================"

# 1. Login as Admin
echo -n "Logging in as Admin... "
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"login\":\"$ADMIN_EMAIL\", \"password\":\"$ADMIN_PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ]; then
  echo "FAILED"
  echo "Login Response: $LOGIN_RESPONSE"
  exit 1
fi
echo "SUCCESS"
echo "Token: ${TOKEN:0:10}..."

# 1.5 Get Admin Role ID
echo -n "Fetching Admin Role ID... "
ROLES_RES=$(curl -s -X GET "$API_URL/roles" \
  -H "Authorization: Bearer $TOKEN")

ADMIN_ROLE_ID=$(echo $ROLES_RES | jq -r '.data[] | select(.name=="admin") | .id')

if [ -z "$ADMIN_ROLE_ID" ] || [ "$ADMIN_ROLE_ID" == "null" ]; then
  echo "FAILED to find admin role"
  # Fallback for some systems or if structure differs
  # Try simple list if data is not wrapped or wrapped differently
  # But assuming standard list response
  echo "Response: $ROLES_RES"
  exit 1
fi
echo "SUCCESS ($ADMIN_ROLE_ID)"

# 2. Create Finance Staff User
echo -n "Creating Finance Staff User... "
STAFF_EMAIL="staff_finance_$(date +%s)@example.com"
STAFF_PASS="Password123!" # Stronger password
STAFF_NAME="Finance Staff"

USER_RES=$(curl -s -X POST "$API_URL/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"$STAFF_NAME\",
    \"email\": \"$STAFF_EMAIL\",
    \"username\": \"finance_$(date +%s)\",
    \"password\": \"$STAFF_PASS\",
    \"role_ids\": [\"$ADMIN_ROLE_ID\"]
  }")
  # Note: role_ids array with UUID

USER_ID=$(echo $USER_RES | jq -r '.data.id')

if [ "$USER_ID" == "null" ]; then
  echo "FAILED"
  echo "Response: $USER_RES"
  exit 1
fi
echo "SUCCESS (ID: $USER_ID)"

# 3. Create Employee Record for Staff
echo -n "Creating Employee Profile for Staff... "
PHONE_NUM="08$(date +%s | cut -c 3-12)" # Unique phone number based on timestamp
NIK_NUM="$(date +%s)123456" # Unique NIK based on timestamp
# Need to check Employee Create Payload format
# Assuming standard multipart or json. Usually multipart in this project based on handlers.
# But let's check EmployeeHandler.CreateEmployee signature or assumes simple JSON for now if handler supports it?
# Use curl with form-data just in case.
EMPLOYEE_RES=$(curl -s -X POST "$API_URL/employees" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\": \"$STAFF_NAME\",
    \"nik\": \"$NIK_NUM\",
    \"email\": \"$STAFF_EMAIL\",
    \"user_id\": \"$USER_ID\",
    \"phone_number\": \"$PHONE_NUM\",
    \"address\": \"Jl. Keuangan No. 1\",
    \"date_of_birth\": \"1990-01-01\",
    \"gender\": \"male\",
    \"job_title\": \"Staff Keuangan\",
    \"employment_status\": \"active\",
    \"join_date\": \"2023-01-01\"
  }")

EMPLOYEE_ID=$(echo $EMPLOYEE_RES | jq -r '.data.id')

if [ "$EMPLOYEE_ID" == "null" ]; then
  echo "FAILED"
  echo "Response: $EMPLOYEE_RES"
  exit 1
fi
echo "SUCCESS (ID: $EMPLOYEE_ID)"

# 4. Login as Finance Staff
echo -n "Logging in as Finance Staff... "
STAFF_LOGIN_RES=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"login\":\"$STAFF_EMAIL\", \"password\":\"$STAFF_PASS\"}")

STAFF_TOKEN=$(echo $STAFF_LOGIN_RES | jq -r '.data.access_token')

if [ "$STAFF_TOKEN" == "null" ]; then
  echo "FAILED"
  echo "Response: $STAFF_LOGIN_RES"
  exit 1
fi
echo "SUCCESS"

# 5. Create Donation (Money)
echo -n "Creating Money Donation... "
DONATION_RES=$(curl -s -X POST "$API_URL/finance/donations" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -F "donor_name=Hamba Allah" \
  -F "donor_phone=08111111111" \
  -F "type=MONEY" \
  -F "payment_method=CASH" \
  -F "total_amount=500000" \
  -F "description=Shodaqah Jumat")

DONATION_ID=$(echo $DONATION_RES | jq -r '.data.id')
DONOR_ID=$(echo $DONATION_RES | jq -r '.data.donor.id')

echo "DEBUG: Donation Response: $DONATION_RES"

if [ "$DONATION_ID" == "null" ] || [ -z "$DONATION_ID" ]; then
  echo "FAILED"
  echo "Response: $DONATION_RES"
  exit 1
fi
echo "SUCCESS (ID: $DONATION_ID, DonorID: $DONOR_ID)"

# 6. Verify Donation List
echo -n "Verifying Donation List... "
LIST_RES=$(curl -s -X GET "$API_URL/finance/donations?donor_id=$DONOR_ID" \
  -H "Authorization: Bearer $STAFF_TOKEN")

COUNT=$(echo $LIST_RES | jq -r '.data.items | length')
if [ "$COUNT" -ge 1 ]; then
  echo "SUCCESS (Found $COUNT records)"
else
  echo "FAILED (Expected >= 1, got $COUNT)"
  echo "Response: $LIST_RES"
fi

# 7. Create Donation (Goods)
echo -n "Creating Goods Donation... "
GOODS_ITEMS='[{"item_name":"Beras","quantity":10,"unit":"kg","notes":"Premium"},{"item_name":"Minyak","quantity":5,"unit":"liter","notes":"Fortune"}]'

# Need dummy file for proof
touch proof.jpg

GOODS_RES=$(curl -s -X POST "$API_URL/finance/donations" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -F "donor_name=Hamba Allah" \
  -F "donor_phone=08111111111" \
  -F "type=GOODS" \
  -F "payment_method=GOODS" \
  -F "items_json=$GOODS_ITEMS" \
  -F "proof_file=@proof.jpg")

GOODS_ID=$(echo $GOODS_RES | jq -r '.data.id')

if [ "$GOODS_ID" == "null" ]; then
  echo "FAILED"
  echo "Response: $GOODS_RES"
  rm proof.jpg
  exit 1
fi
echo "SUCCESS (ID: $GOODS_ID)"
rm proof.jpg

# 8. Verify Donor List
echo -n "Verifying Donor List... "
DONOR_RES=$(curl -s -X GET "$API_URL/finance/donors?name=Hamba" \
  -H "Authorization: Bearer $STAFF_TOKEN")

DONOR_COUNT=$(echo $DONOR_RES | jq -r '.data.items | length')
if [ "$DONOR_COUNT" -ge 1 ]; then
  echo "SUCCESS (Found $DONOR_COUNT donors)"
  # Check for address field presence (even if null or empty string, jq provides it if key exists)
  ADDRESS_CHECK=$(echo $DONOR_RES | jq -r '.data.items[0] | has("address")')
  if [ "$ADDRESS_CHECK" == "true" ]; then
      echo "SUCCESS (Address field present)"
  else
      echo "FAILED (Address field missing)"
      echo "Response: $DONOR_RES"
      exit 1
  fi
else
  echo "FAILED"
  echo "Response: $DONOR_RES"
fi

# 9. Test Duplicate Phone Different Name
echo -n "Testing Shared Phone Scenario... "
SHARED_PHONE="08111111111"
# Previous donor was "Hamba Allah" with this phone.
# Now create "Istri Hamba" with same phone.
SHARED_RES=$(curl -s -X POST "$API_URL/finance/donations" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -F "donor_name=Istri Hamba" \
  -F "donor_phone=$SHARED_PHONE" \
  -F "type=MONEY" \
  -F "payment_method=CASH" \
  -F "total_amount=100000")

SHARED_DONOR_ID=$(echo $SHARED_RES | jq -r '.data.donor.id')
SHARED_DONOR_NAME=$(echo $SHARED_RES | jq -r '.data.donor.name')

# Verify it is different from first donor
FIRST_DONOR_ID=$DONOR_ID # From step 5

if [ "$SHARED_DONOR_ID" != "$FIRST_DONOR_ID" ] && [ "$SHARED_DONOR_NAME" == "Istri Hamba" ]; then
  echo "SUCCESS (Created valid separate donor: $SHARED_DONOR_ID)"
else
  echo "FAILED (Donor reused or name mismatch)"
  echo "Old ID: $FIRST_DONOR_ID, New ID: $SHARED_DONOR_ID"
  echo "Response: $SHARED_RES"
  exit 1
fi



# 10. Test Update Donor
echo -n "Testing Update Donor... "
# Use SHARED_DONOR_ID from previous step
UPDATE_RES=$(curl -s -X PUT "$API_URL/finance/donors/$SHARED_DONOR_ID" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "Jl. Baru No. 99",
    "email": "istri.baru@example.com"
  }')

UPDATED_ADDRESS=$(echo $UPDATE_RES | jq -r '.data.address')
if [ "$UPDATED_ADDRESS" == "Jl. Baru No. 99" ]; then
  echo "SUCCESS (Address updated)"
else
  echo "FAILED"
  echo "Response: $UPDATE_RES"
  exit 1
fi

echo "========================================================"
echo "FULL VERIFICATION COMPLETED"

# 11. Test Update Donation
echo "Testing Update Donation..."
# Use GOODS_ID from step 7
echo "Updating Donation ID: $GOODS_ID"
curl -v -X PUT "$API_URL/finance/donations/$GOODS_ID" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -F "description=Updated Description" \
  -F "items_json=[{\"item_name\":\"Beras Premium\",\"quantity\":20,\"unit\":\"kg\"}]" > response.json

cat response.json
UPDATE_DONATION_RES=$(cat response.json)

UPDATED_DESC=$(echo $UPDATE_DONATION_RES | jq -r '.data.description')
UPDATED_QUANTITY=$(echo $UPDATE_DONATION_RES | jq -r '.data.items[0].quantity')

if [ "$UPDATED_DESC" == "Updated Description" ] && [ "$UPDATED_QUANTITY" == "20" ]; then
  echo "SUCCESS (Description and Items updated)"
else
  echo "FAILED"
  echo "Response: $UPDATE_DONATION_RES"
  exit 1
fi

# 12. Test Get Donation Detail
echo -n "Testing Get Donation Detail... "
DETAIL_RES=$(curl -s -X GET "$API_URL/finance/donations/$GOODS_ID" \
  -H "Authorization: Bearer $STAFF_TOKEN")
DETAIL_ID=$(echo $DETAIL_RES | jq -r '.data.id')
if [ "$DETAIL_ID" == "$GOODS_ID" ]; then
  echo "SUCCESS"
else
  echo "FAILED"
  echo "Response: $DETAIL_RES"
  exit 1
fi

# 13. Test Get Donor Detail
echo -n "Testing Get Donor Detail... "
# Use DONOR_ID from step 2 (captured in CREATE_RES)
DONOR_ID=$(echo $CREATE_RES | jq -r '.data.donor.id')
DONOR_DETAIL=$(curl -s -X GET "$API_URL/finance/donors/$DONOR_ID" \
  -H "Authorization: Bearer $STAFF_TOKEN")
RETRIEVED_ID=$(echo $DONOR_DETAIL | jq -r '.data.id')
if [ "$RETRIEVED_ID" == "$DONOR_ID" ]; then
  echo "SUCCESS"
else
  echo "FAILED"
  echo "Response: $DONOR_DETAIL"
  exit 1
fi


