# Slyther Go Backend

Slyther Go Backend is a smart contract analysis service built in Go. It provides both static and dynamic analysis capabilities for smart contracts, integrating with Etherscan and Large Language Models (LLMs).

## Key Features

- Static smart contract analysis
- Dynamic smart contract analysis
- Etherscan integration for contract information
- LLM integration for analysis results
- RESTful API with Swagger documentation
- Rate limiting and CORS middleware
- Data persistence using GORM

## System Requirements

- Go 1.x
- Docker and Docker Compose
- PostgreSQL (or any GORM-compatible database)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Azzurriii/slythr-go-backend.git
cd slythr-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Start services with Docker Compose:

```bash
docker-compose up -d
```

## Project Structure

```
.
├── cmd/                  # Application entry points
├── config/              # Application configuration
├── docs/                # API documentation (Swagger)
├── internal/            # Internal application code
│   ├── application/     # Business logic
│   ├── domain/         # Domain models and interfaces
│   ├── infrastructure/ # External services implementations
│   └── interface/      # HTTP handlers and routes
├── pkg/                # Reusable packages
└── scripts/            # Utility scripts
```

## Development

1. Run the application in development mode:

```bash
make run
```

2. Run tests:

```bash
make test
```

3. Generate Swagger documentation:

```bash
make swagger
```

## API Documentation

API documentation is automatically generated using Swagger and can be accessed at `/swagger/index.html` after starting the server.

## Features in Detail

### Smart Contract Analysis

- Static Analysis: Code review and vulnerability detection
- Dynamic Analysis: Runtime behavior analysis
- Integration with popular blockchain networks

### API Endpoints

- Contract management
- Analysis execution
- Historical analysis results
- LLM-powered insights
