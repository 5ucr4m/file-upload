#!/bin/bash

cd frontend
echo "Installing dependencies if needed..."
if [ ! -d "node_modules" ]; then
  npm install
fi

echo "Starting frontend development server..."
npm run dev
