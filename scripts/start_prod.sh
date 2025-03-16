#!/bin/bash

make clean init install-prod build
echo "Applying migrations..."
make migrate-up
