#!/bin/bash

# Database connection details
DB_CONTAINER_NAME="monolith-db-1" # Use the correct container name
DB_USER="postgres"
DB_NAME="monolith"
SQL_FILE="./products.sql"
DEST_PATH="/products.sql"

# Copy the SQL file into the running container
docker cp $SQL_FILE $DB_CONTAINER_NAME:$DEST_PATH

# Execute the SQL file inside the container
docker exec -i $DB_CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -f $DEST_PATH

echo "Database seeded successfully."
