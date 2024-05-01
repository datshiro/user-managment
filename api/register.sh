#!/usr/bin/env bash

. "./api/.env" 

# Function to generate a random alphanumeric string of length 8
random_string() {
  local chars='abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  local str=''
  for ((i = 0; i < 8; i++)); do
    local rand=$((RANDOM % ${#chars}))
    str+="${chars:$rand:1}"
  done
  echo "$str"
}

# Function to register using the API
register() {
  local username=$(random_string)
  local email=$(random_string)
  curl -X POST "$API_URL/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$username\", \"email\":\"$email\"}"
}
#
# Function to register using the API as form
register_form() {
  local username=$(random_string)
  local email=$(random_string)
  curl -X POST "$API_URL/register" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "{\"username\":\"$username\", \"email\":\"$email\"}"
}
#
# Call register function or other functions depending on script arguments
"$@"
