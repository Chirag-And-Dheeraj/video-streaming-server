#!/bin/bash

echo "Applying migrations..."
make migrate-up
echo "Starting server..."
air