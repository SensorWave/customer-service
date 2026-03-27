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
USER=$(curl -s -X POST $BASE_URL/companies/$COMPANY_ID/users-company \
  -H "Content-Type: application/json" \
  -d '{
        "id_auth_kc": "kc-acme-admin-001",
        "role": "admin"
      }')

echo "$USER"
USER_ID=$(echo "$USER" | jq -r '.id')

echo "User ID: $USER_ID"


echo "==> Listing companies..."
curl -s $BASE_URL/companies | jq .

echo "==> Listing users of company..."
curl -s $BASE_URL/companies/$COMPANY_ID/users-company | jq .

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
curl -s -X PUT $BASE_URL/users-company/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
        "id_auth_kc": "kc-acme-admin-002",
        "role": "manager"
      }' | jq .


echo "==> Deleting user..."
curl -s -X DELETE $BASE_URL/users-company/$USER_ID

echo "==> Deleting company..."
curl -s -X DELETE $BASE_URL/companies/$COMPANY_ID

echo "==> Done!"
