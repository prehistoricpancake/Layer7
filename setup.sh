#!/bin/bash

# Exit on any error
set -e

echo "Creating chat platform project structure..."


# Create Go chat server structure
echo "Setting up chat server..."
mkdir -p chat-server/handlers chat-server/models
touch chat-server/main.go
touch chat-server/handlers/websocket.go
touch chat-server/handlers/rest.go
touch chat-server/models/message.go
touch chat-server/Dockerfile
touch chat-server/go.mod

# Create Python moderation service structure
echo "Setting up moderation service..."
mkdir -p moderation-service
touch moderation-service/app.py
touch moderation-service/requirements.txt
touch moderation-service/Dockerfile

# Create frontend structure
echo "Setting up frontend..."
mkdir -p frontend/src
touch frontend/src/index.html
touch frontend/src/app.js
touch frontend/src/styles.css
touch frontend/Dockerfile

# Create K8s operator
echo "Setting up Kubernetes operator..."
mkdir -p operator
touch operator/main.go
touch operator/Dockerfile
touch operator/go.mod

# Create deployment files
echo "Setting up deployment configurations..."
mkdir -p deploy/kubernetes
touch deploy/docker-compose.yml
touch deploy/kubernetes/manifests.yaml

# Create README
echo "Creating README..."
touch README.md

echo "Project structure created successfully!"

# Print tree structure for verification
if command -v tree >/dev/null 2>&1; then
    echo -e "\nProject structure:"
    tree
else
    echo -e "\nInstall 'tree' command to view the directory structure"
    ls -R
fi