#!/bin/bash

make init
echo "Applying migrations..."
make migrate-up
echo "Starting server..."
air
