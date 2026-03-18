#!/bin/bash

BASE_URL="http://localhost:3000"
PASSWORD="123456"

# Generate unique emails
EMAIL1="user1_$(date +%s)@test.com"
EMAIL2="user2_$(date +%s)@test.com"

echo "Creating user 1: $EMAIL1"
REGISTER1=$(curl -s -X POST $BASE_URL/auth/register \
-H "Content-Type: application/json" \
-d "{
  \"username\":\"testuser1\",
  \"email\":\"$EMAIL1\",
  \"password\":\"$PASSWORD\",
  \"role\":\"user\"
}")
echo $REGISTER1 | jq
USER1_ID=$(echo $REGISTER1 | jq -r '.id')

echo ""
echo "Creating user 2: $EMAIL2"
REGISTER2=$(curl -s -X POST $BASE_URL/auth/register \
-H "Content-Type: application/json" \
-d "{
  \"username\":\"testuser2\",
  \"email\":\"$EMAIL2\",
  \"password\":\"$PASSWORD\",
  \"role\":\"user\"
}")
echo $REGISTER2 | jq
USER2_ID=$(echo $REGISTER2 | jq -r '.id')

# Validate IDs
if [ -z "$USER1_ID" ] || [ "$USER1_ID" = "null" ]; then
  echo " USER1_ID is invalid"
  exit 1
fi
if [ -z "$USER2_ID" ] || [ "$USER2_ID" = "null" ]; then
  echo " USER2_ID is invalid"
  exit 1
fi

# Login users to get tokens
TOKEN1=$(curl -s -X POST $BASE_URL/auth/login \
-H "Content-Type: application/json" \
-d "{
  \"email\":\"$EMAIL1\",
  \"password\":\"$PASSWORD\"
}" | jq -r '.access_token')

TOKEN2=$(curl -s -X POST $BASE_URL/auth/login \
-H "Content-Type: application/json" \
-d "{
  \"email\":\"$EMAIL2\",
  \"password\":\"$PASSWORD\"
}" | jq -r '.access_token')

echo ""
echo "Depositing 100 to User 1..."
curl -s -X POST $BASE_URL/wallet/deposit \
-H "Authorization: Bearer $TOKEN1" \
-H "Content-Type: application/json" \
-d '{
  "amount": 100
}' | jq

echo ""
echo "Withdrawing 50 from User 1..."
curl -s -X POST $BASE_URL/wallet/withdraw \
-H "Authorization: Bearer $TOKEN1" \
-H "Content-Type: application/json" \
-d '{
  "amount": 50
}' | jq

echo ""
echo "Transferring 25 from User 1 to User 2..."
curl -s -X POST $BASE_URL/wallet/transfer \
-H "Authorization: Bearer $TOKEN1" \
-H "Content-Type: application/json" \
-d "{
  \"to_user_id\": $USER2_ID,
  \"amount\": 25
}" | jq

echo ""
echo "Checking balances..."
echo "User 1:"
curl -s -X GET $BASE_URL/wallet/balance \
-H "Authorization: Bearer $TOKEN1" | jq

echo "User 2:"
curl -s -X GET $BASE_URL/wallet/balance \
-H "Authorization: Bearer $TOKEN2" | jq

echo ""
echo "Fetching transaction history for User 1..."
curl -s -X GET $BASE_URL/wallet/transactions \
-H "Authorization: Bearer $TOKEN1" | jq

echo ""
echo "Fetching transaction history for User 2..."
curl -s -X GET $BASE_URL/wallet/transactions \
-H "Authorization: Bearer $TOKEN2" | jq

echo ""
echo "All tests completed successfully."