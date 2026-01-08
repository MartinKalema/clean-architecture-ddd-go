# Library System

A library management system built with **Go**, following **Clean Architecture** and **Domain-Driven Design (DDD)** principles.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Delivery Layer                           │
│                     (HTTP Handlers, Routes)                      │
├─────────────────────────────────────────────────────────────────┤
│                      Application Layer                           │
│                   (Commands, Queries, DTOs)                      │
├─────────────────────────────────────────────────────────────────┤
│                        Domain Layer                              │
│              (Entities, Value Objects, Interfaces)               │
├─────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                          │
│               (PostgreSQL, External Services)                    │
└─────────────────────────────────────────────────────────────────┘
```

### Key Principles

- **Dependency Rule**: Inner layers don't know about outer layers
- **Domain-Centric**: Business logic lives in the domain layer
- **CQRS**: Commands (writes) and Queries (reads) are separated
- **Testability**: Business logic is easily testable without infrastructure

## Project Structure

```
library-system/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Enterprise business rules
│   │   ├── catalog/
│   │   │   ├── book.go             # Book entity + value objects
│   │   │   ├── errors.go           # Domain errors
│   │   │   ├── events.go           # Domain events
│   │   │   └── repository.go       # Repository interface
│   │   ├── patron/                 # (TODO)
│   │   ├── lending/                # (TODO)
│   │   └── shared/
│   │       ├── errors.go           # Shared errors
│   │       └── events.go           # Domain event interface
│   ├── application/                # Application business rules
│   │   ├── commands/
│   │   │   ├── add_book.go
│   │   │   ├── borrow_book.go
│   │   │   └── return_book.go
│   │   └── queries/
│   │       ├── get_book.go
│   │       └── list_books.go
│   ├── infrastructure/             # External concerns
│   │   ├── external/
│   │   │   └── postgres.go         # Database connection
│   │   └── adapters/
│   │       └── catalog/
│   │           └── book_repository.go
│   └── delivery/                   # Interface adapters
│       └── http/
│           ├── handlers/
│           │   └── book_handler.go
│           ├── models/
│           │   └── book_models.go
│           └── routes/
│               └── routes.go
├── migrations/                     # Database migrations
│   ├── 000001_create_books_table.up.sql
│   └── 000001_create_books_table.down.sql
├── tests/
│   └── load/                       # Load tests (k6)
│       ├── smoke.js
│       ├── load.js
│       └── stress.js
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- golang-migrate (for database migrations)
- k6 (for load testing)

```bash
# Install golang-migrate
brew install golang-migrate

# Install k6
brew install k6
```

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd library-system

# Install dependencies
make setup

# Start PostgreSQL
make db-start

# Run migrations
make migrate-up

# Start the API server
make run
```

The API will be available at `http://localhost:8080`.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/books` | List all books |
| `POST` | `/api/v1/books` | Add a new book |
| `GET` | `/api/v1/books/:id` | Get book by ID |
| `POST` | `/api/v1/books/:id/borrow` | Borrow a book |
| `POST` | `/api/v1/books/:id/return` | Return a book |

### Examples

**Add a book:**
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Clean Code", "author": "Robert Martin"}'
```

**List books:**
```bash
curl http://localhost:8080/api/v1/books
```

**Borrow a book:**
```bash
curl -X POST http://localhost:8080/api/v1/books/{id}/borrow \
  -H "Content-Type: application/json" \
  -d '{"borrower_email": "user@example.com"}'
```

**Return a book:**
```bash
curl -X POST http://localhost:8080/api/v1/books/{id}/return
```

## Development

### Running Tests

```bash
# Run all unit tests
make test

# Run tests with coverage
make test-coverage
```

### Database Migrations

```bash
# Apply all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check current version
make migrate-version

# Create a new migration
make migrate-create name=add_patrons_table
```

### Load Testing

```bash
# Smoke test (1 user, 10s)
make load-smoke

# Load test (up to 1000 users, 12m)
make load-test

# Stress test (up to 10000 users, 17m)
make load-stress
```

Load tests use [k6](https://k6.io/) with a web dashboard at `http://localhost:5665`.

### Load Testing Results

**Test Environment:** PostgreSQL cluster (1 primary + 2 replicas) with synchronous replication.

| Test | VUs | Throughput | p95 Latency | Error Rate |
|------|-----|------------|-------------|------------|
| Smoke | 1 | ~10 req/s | < 50ms | 0% |
| Load | 1,000 | ~448 req/s | 3.08s | 0% |
| Stress | 10,000 | ~1,328 req/s | 12.89s | 0.76% |

**Note:** Load test p95 targets should be measured under expected production load (Load test), not extreme conditions (Stress test).

### Industry Standard p95 Latency Targets

| Application Type | p95 Target | p99 Target |
|------------------|------------|------------|
| Real-time trading | < 10ms | < 50ms |
| Search engines | < 200ms | < 500ms |
| E-commerce checkout | < 500ms | < 1s |
| Standard web APIs | < 500ms - 1s | < 2s |
| Internal tools | < 2s | < 5s |

**Latency Guidelines:**
- **< 500ms** - Excellent
- **< 1s** - Good
- **< 2s** - Acceptable
- **> 3s** - Needs optimization

### Database Replication

The system uses PostgreSQL synchronous replication for **read-your-writes consistency**:

```
Write Request
     │
     ▼
  Primary ────────▶ Replica (sync)
     │                   │
     │              Confirmed
     │◀──────────────────┘
     │
     ▼
  ACK to client
```

**Configuration:**
- `POSTGRESQL_SYNCHRONOUS_COMMIT_MODE: "on"` - Writes wait for replica confirmation
- `POSTGRESQL_NUM_SYNCHRONOUS_REPLICAS: 1` - At least 1 replica must confirm

**Trade-off:** Higher write latency in exchange for strong consistency (no stale reads).

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/library?sslmode=disable` |
| `PORT` | HTTP server port | `8080` |

## Architecture Details

### Domain Layer

The domain layer contains the core business logic:

- **Entities**: `Book` - aggregate root with business rules
- **Value Objects**: `BookID`, `Title`, `Author` - immutable, validated
- **Domain Events**: `BookBorrowed`, `BookReturned` - capture state changes
- **Repository Interfaces**: Define persistence contracts

### Application Layer

The application layer orchestrates use cases:

- **Commands**: `AddBook`, `BorrowBook`, `ReturnBook` - write operations
- **Queries**: `GetBook`, `ListBooks` - read operations
- **Handlers**: Execute commands/queries using domain entities

### Infrastructure Layer

The infrastructure layer implements external concerns:

- **PostgreSQL Repository**: Implements `BookRepository` interface
- **Database Connection**: Connection pool management

### Delivery Layer

The delivery layer handles HTTP concerns:

- **Handlers**: Convert HTTP requests to commands/queries
- **Routes**: Define API endpoints
- **Models**: Request/response DTOs with validation

## Testing Strategy

| Layer | Test Type | Tools |
|-------|-----------|-------|
| Domain | Unit tests | `go test` |
| Application | Unit tests with mocks | `go test` |
| Infrastructure | Integration tests | `go test` + testcontainers |
| Delivery | HTTP tests | `httptest` |
| System | Load tests | k6 |

## Roadmap

- [x] Catalog domain (Book entity)
- [x] CQRS command/query handlers
- [x] PostgreSQL repository
- [x] HTTP API
- [x] Unit tests
- [x] Load tests (k6)
- [x] Database migrations
- [ ] Patron domain
- [ ] Lending domain
- [ ] Circuit breaker pattern
- [ ] Elasticsearch integration
- [ ] Redis caching

## License

MIT
