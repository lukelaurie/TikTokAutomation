#!/bin/bash

# Define the PostgreSQL data directory
DATA_DIR="C:/Program Files/PostgreSQL/16/data"

# Check if PostgreSQL is running
echo "Starting PostgreSQL..."
pg_ctl -D "$DATA_DIR" start > /dev/null 2>&1

cd ..
go build ./cmd/app/main.go

if [ "$1" == "test" ]; then
    echo "Running in test mode..."
    ./main -test   
else
    echo "Running in normal mode..."
    ./main
fi