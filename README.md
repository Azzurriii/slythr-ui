# Slythr - Solidity Contract Analyze

This repository contains the backend service for Slythr, a comprehensive smart contract analysis platform. Built with Go, it provides a robust API for performing static and dynamic security analysis on Solidity source code. The system is designed with a Clean Architecture approach to ensure maintainability, scalability, and separation of concerns.

## Features

- **Static Analysis**: Leverages the powerful Slither static analysis tool within a dedicated Docker container to detect a wide range of vulnerabilities and code quality issues.
- **Dynamic AI-Powered Analysis**: Utilizes Google's Gemini AI model to perform a deep, contextual security audit, providing a human-like assessment of vulnerabilities, logical flaws, and best practice deviations.
- **Etherscan Integration**: Fetches verified smart contract source code directly from Etherscan for various EVM-compatible networks.
- **Multi-Layer Caching**: Implements a two-layer caching strategy (Redis for L1, PostgreSQL for L2) to optimize performance and reduce redundant analysis and API calls.
- **Clean Architecture**: Organized following Clean Architecture principles, promoting a clear separation between domain logic, application use cases, and infrastructure details.
- **Containerized Environment**: Fully containerized using Docker and Docker Compose for consistent development, testing, and deployment environments.
- **API Documentation**: Provides comprehensive API documentation through Swagger (OpenAPI).

## Architecture

The project is structured following the principles of Clean Architecture, which organizes the codebase into a series of concentric layers. This design enforces the dependency rule: source code dependencies can only point inwards.

- **Domain**: The innermost layer, containing enterprise-wide business rules and entities. It is the most stable and independent layer. It includes entities, value objects, and repository interfaces.
- **Application**: This layer contains application-specific business rules and use cases. It orchestrates the flow of data to and from the Domain layer. It includes services, DTOs (Data Transfer Objects), and handlers.
- **Interface**: The outermost layer, responsible for communication with the outside world. This includes the Gin web framework, routers, handlers, and middleware.
- **Infrastructure**: This layer provides implementations for the interfaces defined in the layers above. It contains database repositories (GORM), caching clients (Redis), and clients for external services (Etherscan, Gemini, Slither).

## Technology Stack

- **Backend**: Go 1.24
- **Framework**: Gin Gonic
- **Database**: PostgreSQL
- **ORM**: GORM
- **Caching**: Redis
- **Containerization**: Docker, Docker Compose
- **Static Analysis**: Slither
- **AI Analysis**: Google Gemini

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.24 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/Azzurriii/slythr-go-backend.git
    cd slythr-go-backend
    ```

2.  **Set up environment variables:**
    Copy the example environment file and populate it with your configuration.

    ```sh
    cp .env.example .env
    ```

    You will need to provide your own API keys for `ETHERSCAN_API_KEY` and `GEMINI_API_KEY` in the `.env` file.

3.  **Build and run the services:**
    The easiest way to get all services running is by using Docker Compose. The `Makefile` provides a convenient command for this.

    ```sh
    make docker-run
    ```

    This command will build the Go application image, the Slither image, and start all necessary containers (app, postgres, redis, slither) in detached mode.

4.  **Verify the setup:**
    Check the logs to ensure the application has started correctly.
    ```sh
    docker-compose logs -f app
    ```
    You should see a message indicating that the server has started on the configured port (default: 8080).

## Usage

The application can be run directly with Go for development or through Docker for a production-like environment.

## Project Structure

The project follows a logical structure to support Clean Architecture principles.

```
└── azzurriii-slythr-go-backend/
    ├── cmd/api/            # Application entry point (main.go)
    ├── config/             # Configuration loading (Viper)
    ├── docs/               # Swagger API documentation files
    ├── internal/           # Core application code
    │   ├── application/    # Application layer
    │   ├── domain/         # Domain layer
    │   ├── infrastructure/ # Infrastructure layer
    │   └── interface/      # Interface layer (HTTP routes, server)
    ├── pkg/                # Shared, reusable packages (logger, utils)
    ├── scripts/            # Standalone scripts (e.g., migrations)
    ├── .env.example        # Example environment file
    ├── Dockerfile          # Dockerfile for the main Go application
    ├── Dockerfile.slither  # Dockerfile for the Slither analysis environment
    ├── docker-compose.yml  # Docker Compose configuration for all services
    └── Makefile            # Makefile for common development tasks
```
