#!/bin/bash
# Build the Next.js app locally first
npm install
npm run build

# Then build the Docker image with the pre-built files
docker build -t frontend .
