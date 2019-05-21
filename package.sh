#!/bin/bash

git_tmp=$(mktemp -d)
git clone https://github.com/lfkeitel/inca-frontend "$git_tmp"
pushd "$git_tmp"
yarn install
yarn build
popd

# Build application
make build

# Setup folder structure
rm -rf dist
mkdir -p dist/{archive,logs,latest}

# Copy application
mv $git_tmp/dist dist/frontend
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
rm -rf $git_tmp
