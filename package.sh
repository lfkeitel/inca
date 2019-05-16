#!/bin/bash

# Build application
yarn run build
make build

# Setup folder structure
rm -rf dist
mkdir -p dist/{archive,logs,latest,frontend}

# Copy application
cp -R frontend/dist dist/frontend/dist
cp -R config dist/
cp -R scripts dist/

cp bin/inca dist/
cp LICENSE dist/
cp README.md dist/

# Create archive
cd dist
tar -czf ../inca-dist.tar.gz  ./*
cd ..

# Cleanup
rm -rf dist
