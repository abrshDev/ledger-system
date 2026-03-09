#!/bin/bash

BASE_URL="http://localhost:3000"

EMAIL="user$(date +%s)@test.com"
PASSWORD="123456"

echo "Creating user with email: $EMAIL"

REGISTER=$(curl -s -X POST $BASE_URL/auth/register \
-H "Content-Type: application/json" \
-d "{
  \"username\":\"testuser\",
  \"email\":\"$EMAIL\",
  \"password\":\"$PASSWORD\",
  \"role\":\"user\"
}")

echo "Register response:"
echo $REGISTER | jq

echo ""
echo "Logging in..."

LOGIN=$(curl -s -X POST $BASE_URL/auth/login \
-H "Content-Type: application/json" \
-d "{
  \"email\":\"$EMAIL\",
  \"password\":\"$PASSWORD\"
}")

echo $LOGIN | jq

TOKEN=$(echo $LOGIN | jq -r '.access_token')

echo ""
echo "Access Token:"
echo $TOKEN

echo ""
echo "Depositing 100..."

DEPOSIT=$(curl -s -X POST $BASE_URL/wallet/deposit \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "amount": 100
}')

echo $DEPOSIT | jq

echo ""
echo "Checking balance..."

BALANCE=$(curl -s -X GET $BASE_URL/wallet/balance \
-H "Authorization: Bearer $TOKEN")

echo $BALANCE | jq

echo ""
echo "Test completed."