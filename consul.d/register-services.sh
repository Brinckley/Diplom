#!/bin/bash

curl -s -XPUT -d"{
  \"Name\": \"postgres\",
  \"ID\": \"postgres\",
  \"Tags\": [ \"postgres\" ],
  \"Address\": \"localhost\",
  \"Port\": 5432,
  \"Check\": {
    \"Name\": \"PostgreSQL TCP on port 5432\",
    \"ID\": \"postgres\",
    \"Interval\": \"10s\",
    \"TCP\": \"postgres:5432\",
    \"Timeout\": \"1s\",
    \"Status\": \"passing\"
  }
}" localhost:8500/v1/agent/service/register

curl -s -XPUT -d"{
  \"Name\": \"zookeeper\",
  \"ID\": \"zookeeper\",
  \"Tags\": [ \"zookeeper\" ],
  \"Address\": \"localhost\",
  \"Port\": 5181,
  \"Check\": {
    \"Name\": \"Zookeeper on port 5181\",
    \"ID\": \"zookeeper\",
    \"Interval\": \"10s\",
    \"TCP\": \"zookeeper:5181\",
    \"Timeout\": \"1s\",
    \"Status\": \"passing\"
  }
}" localhost:8500/v1/agent/service/register

curl -s -XPUT -d"{
  \"Name\": \"kafka\",
  \"ID\": \"kafka\",
  \"Tags\": [ \"kafka\" ],
  \"Address\": \"localhost\",
  \"Port\": 9092,
  \"Check\": {
    \"Name\": \"Kafka on port 9092\",
    \"ID\": \"kafka\",
    \"Interval\": \"10s\",
    \"TCP\": \"kafka:9092\",
    \"Timeout\": \"1s\",
    \"Status\": \"passing\"
  }
}" localhost:8500/v1/agent/service/register

curl -s -XPUT -d"{
  \"Name\": \"redis\",
  \"ID\": \"redis\",
  \"Tags\": [ \"redis\" ],
  \"Address\": \"localhost\",
  \"Port\": 6739,
  \"Check\": {
    \"Name\": \"Redis on port 6379\",
    \"ID\": \"redis\",
    \"Interval\": \"10s\",
    \"TCP\": \"redis:6379\",
    \"Timeout\": \"1s\",
    \"Status\": \"passing\"
  }
}" localhost:8500/v1/agent/service/register
