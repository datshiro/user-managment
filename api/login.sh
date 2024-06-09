#!/usr/bin/env bash

. "./api/.env" 

login_with_username() {
  local account=$1
  local password=$2
  curl -X POST "$API_URL/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "{\"username\":\"$account\", \"password\": \"$password\"}"
}

login_with_email() {
  local account=$1
  local password=$2
  curl -v -X POST "$API_URL/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "{\"email\":\"$account\", \"password\": \"$password\"}"
}

login_with_phone() {
  local account=$1
  local password=$2
  curl -X POST "$API_URL/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "{\"phone_number\":\"$account\", \"password\": \"$password\"}"
}

login_fail() {
  local account=$1
  local password=$2
  curl -v -X POST "$API_URL/login" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "{\"phonenumber\":\"$account\", \"password\": \"$password\"}"
}

"$@"
