#!/bin/sh

# Run migrations before starting the application
echo "Running migrations..."
/app/migrate

# Start the application
echo "Starting the Go app..."
/app/main
