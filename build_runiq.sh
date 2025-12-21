#!/bin/bash
set -e
echo "🔨 Building Runiq SDK..."
rm -f go.mod go.sum
go mod init runiq
go mod tidy
go build -o runiq cmd/runiq/main.go
sudo mv runiq /usr/local/bin/runiq
sudo chmod +x /usr/local/bin/runiq
echo "✅ RUNIQ INSTALLED: Type 'runiq' to start."
