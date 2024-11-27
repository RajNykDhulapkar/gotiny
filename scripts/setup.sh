#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up local development environment...${NC}"

# Check if Docker is installed
if ! command -v docker &>/dev/null; then
  echo -e "${RED}Docker is not installed. Please install Docker first.${NC}"
  exit 1
fi

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo -e "${RED}Go is not installed. Please install Go first.${NC}"
  exit 1
fi

# Create docker network if it doesn't exist
echo "Creating Docker network..."
docker network create gotiny-network 2>/dev/null || true

# Stop existing containers if they exist
echo "Stopping existing containers..."
docker stop mongodb redis postgres range-allocator 2>/dev/null || true
docker rm mongodb redis postgres range-allocator 2>/dev/null || true

# Create necessary directories
mkdir -p ./data/mongodb
mkdir -p ./data/redis
mkdir -p ./data/postgres

# Start PostgreSQL
echo "Starting PostgreSQL..."
docker run --name postgres \
  --network gotiny-network \
  -e POSTGRES_USER=username \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=alloc \
  -p 5432:5432 \
  -v "$(pwd)/data/postgres:/var/lib/postgresql/data" \
  -d postgres:latest

# Start MongoDB
echo "Starting MongoDB..."
docker run -d \
  --name mongodb \
  --network gotiny-network \
  -p 27017:27017 \
  -e MONGODB_INITDB_ROOT_USERNAME= \
  -e MONGODB_INITDB_ROOT_PASSWORD= \
  mongo:latest

# Start Redis
echo "Starting Redis..."
docker run -d \
  --name redis \
  --network gotiny-network \
  -p 6379:6379 \
  redis:latest

echo "Waiting for PostgreSQL to be ready..."
until docker exec postgres pg_isready -U username -d alloc; do
  sleep 2
done

# Start Range Allocator Service
echo "Starting Range Allocator service..."
docker run -d \
  --name range-allocator \
  --network gotiny-network \
  -p 50051:50051 \
  -e RANGE_ALLOCATOR_DATABASE_URL="postgres://username:password@postgres:5432/alloc?sslmode=disable" \
  -e RANGE_ALLOCATOR_GRPC_PORT=50051 \
  -e RANGE_ALLOCATOR_RANGE_DEFAULT_SIZE=1000 \
  -e RANGE_ALLOCATOR_RANGE_MIN_SIZE=100 \
  -e RANGE_ALLOCATOR_RANGE_MAX_SIZE=10000 \
  gotiny-range-allocator

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Export required environment variables
export RANGE_ALLOCATOR_GRPC_PORT=50051

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Verify services are running
echo "Verifying services..."

# Check MongoDB
if docker exec mongodb mongosh --eval "db.stats()" &>/dev/null; then
  echo -e "${GREEN}MongoDB is running${NC}"
else
  echo -e "${RED}MongoDB failed to start${NC}"
fi

# Check Redis
if docker exec redis redis-cli ping &>/dev/null; then
  echo -e "${GREEN}Redis is running${NC}"
else
  echo -e "${RED}Redis failed to start${NC}"
fi

# Check PostgreSQL
if docker exec postgres pg_isready -U username -d alloc &>/dev/null; then
  echo -e "${GREEN}PostgreSQL is running${NC}"
else
  echo -e "${RED}PostgreSQL failed to start${NC}"
fi

# Check Range Allocator
if docker exec range-allocator grpcurl -plaintext localhost:${RANGE_ALLOCATOR_GRPC_PORT} rangeallocator.v1.RangeAllocator/GetHealth &>/dev/null; then
  echo -e "${GREEN}Range Allocator is running${NC}"
else
  echo -e "${RED}Range Allocator failed to start${NC}"
fi

# Print setup information
echo -e "\n${GREEN}Local development environment is ready!${NC}"
echo -e "\nService endpoints:"
echo "MongoDB: localhost:27017"
echo "Redis: localhost:6379"
echo "PostgreSQL: localhost:5432"
echo "Range Allocator gRPC: localhost:50051"
echo "API Server will run on: localhost:8080"

echo -e "\nNetwork Information:"
echo "All services are running on the 'gotiny-network' Docker network"
echo "Service hostnames:"
echo "  - MongoDB: mongodb"
echo "  - Redis: redis"
echo "  - PostgreSQL: postgres"
echo "  - Range Allocator: range-allocator"

echo -e "\nDatabase credentials:"
echo "PostgreSQL:"
echo "  Username: username"
echo "  Password: password"
echo "  Database: alloc"

echo -e "\nUseful commands:"
echo "Connect to PostgreSQL:"
echo "  docker exec -it postgres psql -U username -d alloc"
echo "Check Range Allocator health:"
echo "  docker exec range-allocator grpcurl -plaintext localhost:50051 rangeallocator.v1.RangeAllocator/GetHealth"
echo "Check containers on network:"
echo "  docker network inspect gotiny-network"
echo "View all container logs:"
echo "  docker-compose logs -f"

echo -e "\nTo start the application:"
echo "go run main.go"

echo -e "\nTo stop all services:"
echo "docker stop mongodb redis postgres range-allocator"
