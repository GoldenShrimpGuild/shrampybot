#!/bin/bash

ROOT_DIR=$(git rev-parse --show-toplevel)
BACKEND_DIR=$ROOT_DIR/function
CURRENT_DIR=$PWD

echo "Using tygo to generate typescript from go structs in $BACKEND_DIR"
cd $BACKEND_DIR
tygo generate
echo "Returning to $CURRENT_DIR"
cd $CURRENT_DIR
