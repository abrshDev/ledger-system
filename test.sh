#!/bin/bash

BASE_URL="http://localhost:3000"

echo "Registering user..."

curl -s -X POST $BASE_URL/auth/register \
-H "Content-Type: application/json" \
-d '{
  "username":"john",
  "email":"john@test.com",
  "password":"123456",
  "role":"user"
}' | jq

echo ""
echo "Logging in..."

LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email":"john@test.com",
  "password":"123456"
}')

echo $LOGIN_RESPONSE | jq

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')

echo ""
echo "Access Token:"
echo $TOKEN

echo ""
echo "Getting wallet balance..."

curl -s -X GET $BASE_URL/wallet/balance \
-H "Authorization: Bearer $TOKEN" | jq