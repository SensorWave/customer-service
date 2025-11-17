#!/bin/bash

BASE_URL="http://localhost:8080"

echo "==> Creating company..."
COMPANY=$(curl -s -X POST $BASE_URL/companies \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Acme Corp",
        "email": "contact@acme.com",
        "phone": "0601020304",
        "address": "42 rue de Paris"
      }')

echo "$COMPANY"
COMPANY_ID=$(echo "$COMPANY" | jq -r '.id')

echo "Company ID: $COMPANY_ID"

echo "==> Creating user in company..."
USER=$(curl -s -X POST $BASE_URL/companies/$COMPANY_ID/users \
  -H "Content-Type: application/json" \
  -d '{
        "first_name": "John",
        "last_name": "Doe",
        "email": "john@acme.com",
        "phone": "0611223344"
      }')

echo "$USER"
USER_ID=$(echo "$USER" | jq -r '.id')

echo "User ID: $USER_ID"


echo "==> Listing companies..."
curl -s $BASE_URL/companies | jq .

echo "==> Listing users of company..."
curl -s $BASE_URL/companies/$COMPANY_ID/users | jq .

echo "==> Updating company..."
curl -s -X PUT $BASE_URL/companies/$COMPANY_ID \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Updated Corp",
        "email": "support@acme.com",
        "phone": "0707070707",
        "address": "75 avenue de Toulouse"
      }' | jq .

echo "==> Updating user..."
curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
        "first_name": "Jane",
        "last_name": "Doe",
        "email": "jane@acme.com",
        "phone": "0699887766"
      }' | jq .


echo "==> Deleting user..."
curl -s -X DELETE $BASE_URL/users/$USER_ID

echo "==> Deleting company..."
curl -s -X DELETE $BASE_URL/companies/$COMPANY_ID

echo "==> Done!"
