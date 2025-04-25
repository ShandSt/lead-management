# Lead Management API

## Overview

This application provides a RESTful API for lead management, allowing clients to manage leads, client data, and lead assignment processes efficiently. Built with Go and Gin framework, it offers a robust solution for lead distribution based on client priority, capacity, and working hours.

## Features

- Client management (CRUD operations)
- Intelligent lead assignment based on priority and capacity
- Working hours support for clients
- Thread-safe in-memory storage with optimized concurrency
- RESTful API design
- Docker support for easy deployment

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- Docker and Docker Compose (optional, for containerized deployment)

### Installation

```bash
git clone https://github.com/stasshander/lead-management.git
cd lead-management
go mod download
```

### Running the Application

#### Local Development

```bash
go run main.go
```

The server will start on port 8001 by default (http://localhost:8001).

#### Using Docker

Build and run the Docker container:

```bash
# Using docker-compose
docker-compose up -d

# Or using Make commands
make up
```

The application will be available at http://localhost:8001.

### Makefile Commands

The project includes a Makefile with helpful commands:

- `make build`: Build the Docker image
- `make up`: Run the containers in the background
- `make run`: Run the containers in the foreground
- `make down`: Stop and remove the containers
- `make logs`: Show logs from the containers
- `make test`: Run tests
- `make test-coverage`: Run tests with coverage
- `make dev`: Run local development server
- `make clean`: Remove all containers and images

## API Endpoints

### 1. Create Client
   `POST /clients`

Creates a new client in the system.
```json
{
  "id": "string",
  "name": "string",
  "workingHours": {
    "start": "2023-01-01T09:00:00Z",
    "end": "2023-01-01T17:00:00Z"
  },
  "priority": 10,
  "leadCount": 0,
  "capacity": 15
}
```

**Response**
- 201 Created: Successfully created client.
- 400 Bad Request: Invalid request data.

### 2. Get All Clients
   `GET /clients`

Retrieves a list of all clients.

**Response**
- 200 OK: Successfully retrieved all clients. Returns an array of clients.
- 404 Not Found: No clients found.

### 3. Get Client
   `GET /clients/:id`

Retrieves details of a specific client.

**Path Parameters**
- id (string): Unique identifier of the client.

**Response**
- 200 OK: Successfully retrieved client.
- 404 Not Found: Client not found.

### 4. Assign Lead
   `GET /clients/assignLead`

Assigns a lead to the most appropriate client based on priority, capacity, and working hours.

**Response**
- 200 OK: Successfully assigned a lead.
- 404 Not Found: No suitable client found.

## Lead Assignment Algorithm

The lead assignment algorithm selects clients based on the following criteria:

1. **Working Hours**: Only clients currently within their working hours are considered
2. **Capacity**: Clients must have available capacity (leadCount < capacity)
3. **Priority**: Higher priority clients are preferred
4. **Load Balancing**: When priorities are equal, the client with the lower load factor (leadCount/capacity ratio) is selected

## Development and Testing

This API is developed using Go with the Gin framework. For testing, use tools such as Postman or cURL to make requests to the local server.

### Running Tests

Execute the following command in the terminal to run the automated tests:

```bash
go test ./api_test

# Or using Make
make test
```

## Docker Deployment

The application includes Docker configuration for easy deployment:

1. **Dockerfile**: Multi-stage build for a compact production image
2. **docker-compose.yml**: Configuration for container orchestration

## Recent Optimizations

1. **Concurrency Improvements**:
   - Enhanced thread safety with read/write mutex
   - Defensive copying to prevent race conditions

2. **Bug Fixes**:
   - Fixed route conflicts between parameter-based and static routes
   - Resolved integer division issues in client selection algorithm
   - Improved time parsing to support multiple formats

3. **Error Handling**:
   - Standardized error types
   - Added more detailed error messages
   - Improved input validation

## License

MIT
