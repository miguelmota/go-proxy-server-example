#!/bin/bash

curl "https://localhost:8000/ping" --cacert ../certs/ca.crt

curl -X POST "https://localhost:8000/proxy" -d "@payload.json" --cacert ../certs/ca.crt
