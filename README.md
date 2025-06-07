# Go HTTP Server

A simple HTTP server implementation in Go to understand HTTP protocol fundamentals.

## Description

This project is a basic HTTP server built from scratch in Go, designed to help understand how HTTP requests and responses work at a low level. It handles HTTP requests, parses headers, and demonstrates concurrent request processing.

## Project Structure

```
.
├── go.mod              # Go module file
├── go.sum              # Go dependencies
├── trim_test.go        # Test file for string parsing
├── note.md             # Development notes
└── pkg/
    ├── main.go         # Main server entry point
    ├── models.go       # HTTP request/response models
    └── server.go       # Server implementation
```

## Features

- Basic HTTP request parsing
- Header extraction and processing
- User-Agent parsing
- Concurrent connection handling
- TCP server implementation

## Getting Started

### Prerequisites

- Go 1.24.0 or later

### Running the Server

1. Clone the repository:
```bash
git clone https://github.com/alirezapourchali/go-http-server.git
cd go-http-server
```

2. Run the server:
```bash
go run pkg/*.go
```

The server will start listening on `localhost:4221`.

### Testing

Run the test file to see string parsing functionality:
```bash
go run trim_test.go
```

## Usage

Send HTTP requests to the server:
```bash
curl http://localhost:4221/
curl -H "User-Agent: myapp/1.0" http://localhost:4221/user-agent
```

## Learning Goals

- Understanding HTTP protocol structure
- Learning TCP server implementation in Go
- Practicing request/response parsing
- Exploring concurrent programming patterns

## License

This project is for educational purposes. 
It's inspired by CodeCrafters