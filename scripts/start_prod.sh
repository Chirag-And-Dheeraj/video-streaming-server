#!/bin/bash

make clean init build
echo "Applying migrations..."
make migrate-up
