# Linux ACL Management System

A system for managing Linux Access Control Lists (ACLs) through a web interface with gRPC backend communication.

Progress Docs: [https://pythonhacker24.github.io/linux-acl-management/](https://pythonhacker24.github.io/linux-acl-management/)

## Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- Node.js (for frontend)

## Installation

1. Install Protocol Buffers compiler (protoc):
   ```bash
   # macOS
   brew install protobuf

   # Linux
   apt-get install protobuf-compiler
   ```

2. Install Go protocol buffer plugins:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. Install project dependencies:
   ```bash
   cd backend-server
   go mod tidy
   ```

## Build Instructions

1. Generate Protocol Buffer code:
   ```bash
   cd backend-server
   protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/daemon.proto
   ```

2. Build the backend server:
   ```bash
   cd backend-server
   go build
   ```

3. Run the server:
   ```bash
   ./backend-server
   ```

## Configuration

The backend server configuration is stored in `backend.yaml`. Example configuration:

```yaml
servers:
  - name: "daemon1"
    address: "localhost:50051"
  - name: "daemon2"
    address: "localhost:50052"
```

## API Documentation

The backend server provides the following endpoints:

- `GET /health` - Health check endpoint
- Additional endpoints to be documented

## License

This project is licensed under the MIT License - see the LICENSE file for details.
