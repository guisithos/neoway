#!/bin/bash

# PF
persons=(
  "075.210.019-07,Alqua Ayala"
  "980.672.300-72,Alqua Blevins"
  "584.587.199-91,Alqua Bryant"
  "095.710.001-91,Alqua Bryant"
  "055.036.799-38,Alqua Bryant"
  "874.777.349-91,Alqua Camacho"
  "071.701.649-81,Alqua Cannon"
  "356.468.828-50,Alqua Cannon"
  "053.916.489-55,Alqua Casey"
  "559.637.739-20,Alqua Cervantes"
)

# PJ
businesses=(
  "36.885.434/0001-63,Alphahive Camacho"
  "70.858.125/0001-32,Alphahive Herring"
  "27.422.812/0001-81,Alphahive Hobbs"
  "85.737.514/0001-87,Alphahive Mcconnell"
  "26.564.802/0001-18,Alphahive Preston"
  "16.626.228/0001-21,Alphahive Stephens"
  "77.765.140/0001-85,Alpharon Bates"
  "83.884.506/0001-38,Alpharon Blevins"
)

# post
send_post() {
  local document=$1
  local name=$2
  local type=$3

  curl -X POST http://localhost:8080/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name": "'"$name"'",
    "document": "'"$document"'",
    "type": "'"$type"'"
  }'
}

# send PF
for entry in "${persons[@]}"; do
  IFS=',' read -r document name <<< "$entry"
  send_post "$document" "$name" "PERSON"
done

# send PJ
for entry in "${businesses[@]}"; do
  IFS=',' read -r document name <<< "$entry"
  send_post "$document" "$name" "BUSINESS"
done
