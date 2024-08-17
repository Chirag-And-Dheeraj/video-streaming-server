#!/bin/bash

# Define directory names
dir1="segments"
dir2="data"

# Check and create 'segments' directory if it doesn't exist
[ ! -d "$dir1" ] && mkdir "$dir1" || { echo "Failed to create $dir1 directory"; exit 1; }

# Check and create 'data' directory if it doesn't exist
[ ! -d "$dir2" ] && mkdir "$dir2" || { echo "Failed to create $dir2 directory"; exit 1; }
