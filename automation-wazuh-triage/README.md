# Wazuh Security Event Triage Automation System

A sophisticated Go-based REST API service for automating security event triage and management with Wazuh SIEM integration. This system provides automated event processing, rule analysis, and incident closure workflows.

## 🏗️ Architecture Overview

The project follows Clean Architecture principles with a well-structured layered design:

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/          # Business logic interfaces
│   ├── entity/          # Core business entities
│   ├── handler/         # HTTP request handlers (Controller layer)
│   ├── model/           # Request/Response models and DTOs
│   ├── repository/      # Data access layer implementations
│   ├── route/           # HTTP routing configuration
│   └── usecase/         # Business logic implementations
├── pkg/
│   ├── database/        # Database connection utilities
│   ├── logger/          # Structured logging
│   ├── middleware/      # HTTP middlewares
│   ├── opensearch/      # OpenSearch/Elasticsearch client
│   └── wazuh/           # Wazuh API client
└── docs/                # API documentation (OpenAPI spec)
```

## 🚀 Key Features

### Core Functionality
- **Security Event Fetching**: Retrieve security events from Wazuh/OpenSearch with flexible filtering
- **Automated Event Closure**: Bulk automatic closure of events based on configurable criteria
- **Manual Event Management**: Individual event closure with custom reasoning
- **Rule Analysis**: Integration with Wazuh rules for detailed security context
- **Event History**: Comprehensive tracking of closed events with full audit trail

### Advanced Features
- **Auto-Close Functionality**: Automatically close events matching specific criteria during fetch operations
- **Duplicate Prevention**: Intelligent detection and prevention of duplicate event closures
- **Rule Intelligence**: Fetch related rules from the same file for comprehensive analysis
- **Flexible Filtering**: Support for level-based filtering with range queries
- **Structured Logging**: Request ID tracking and comprehensive audit logging

## 🛠️ Technology Stack

- **Language**: Go 1.25.3
- **Web Framework**: Fiber v2 (High-performance HTTP framework)
- **Database**: SQLite (Local storage for closed events)
- **Search Engine**: OpenSearch/Elasticsearch (Wazuh event storage)
- **Documentation**: OpenAPI 3.0 with Swagger UI
- **Logging**: Logrus (Structured logging)
- **HTTP Client**: Resty (Wazuh API integration)

### Key Dependencies
```go
github.com/gofiber/fiber/v2      // High-performance web framework
github.com/olivere/elastic/v7    // Elasticsearch/OpenSearch client
github.com/mattn/go-sqlite3      // SQLite database driver
github.com/sirupsen/logrus       // Structured logging
github.com/go-resty/resty/v2     // HTTP client for Wazuh API
github.com/google/uuid           // UUID generation
github.com/gofiber/swagger       // API documentation
```

## 📊 Database Schema

### Closed Events Table
```sql
CREATE TABLE closed_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL UNIQUE,
    rule_id TEXT,
    raw_event TEXT,           -- Full JSON event data
    reason TEXT,              -- Closure reason
    status TEXT NOT NULL,     -- Event status (closed)
    close_at DATETIME NOT NULL
);
```

## 🔌 API Endpoints

### Health Check
- `GET /health` - Service health status

### Security Events
- `POST /v1/events` - Fetch events with optional auto-close
- `POST /v1/events/{event_id}/close` - Manually close specific event
- `GET /v1/events/close` - List all closed events
- `GET /v1/events/close/{id}` - Get detailed closed event with rule context
- `PATCH /v1/events/close/{id}/reason` - Update closure reason

### Wazuh Rules
- `GET /v1/rules/{id}` - Get specific rule details
- `GET /v1/rules/file/{filename}` - Get all rules from specific file

### API Documentation
- `GET /swagger/*` - Interactive Swagger UI
- `GET /docs/openapi.yaml` - OpenAPI specification

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.25.3 or higher
- Access to Wazuh/OpenSearch cluster
- SQLite support

### Environment Variables
```bash
# Server Configuration
SERVER_PORT=8080

# OpenSearch/Elasticsearch Configuration
OPENSEARCH_URL=https://your-opensearch-cluster
OPENSEARCH_USERNAME=admin
OPENSEARCH_PASSWORD=your-password

# Wazuh API Configuration (optional)
WAZUH_URL=https://your-wazuh-manager
WAZUH_USERNAME=wazuh
WAZUH_PASSWORD=your-wazuh-password
```

### Installation & Running

1. **Clone and install dependencies**:
```bash
git clone <repository-url>
cd automation-wazuh-triage
go mod download
```

2. **Build the application**:
```bash
go build -o bin/server cmd/server/main.go
```

3. **Run the service**:
```bash
./bin/server
```

4. **Access the API**:
- Service: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/
- Health Check: http://localhost:8080/health

## 📖 Usage Examples

### Fetch Events with Auto-Close
```bash
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "level_range": {"lte": 3},
    "limit": 100,
    "auto_add_to_close": true
  }'
```

### Manually Close Event
```bash
curl -X POST http://localhost:8080/v1/events/1760850699.19418/close \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "False positive - legitimate system activity"
  }'
```

### Get Closed Event Details
```bash
curl -X GET http://localhost:8080/v1/events/close/1 \
  -H "Content-Type: application/json"
```

## 🔍 Monitoring & Observability

### Logging
- Structured JSON logging with Logrus
- Request ID correlation across all operations
- Error tracking with context
- Performance metrics logging

### Health Checks
- Service availability monitoring
- Database connectivity verification
- External service dependency checks

## 🚧 Development

### Project Structure Principles
1. **Clean Architecture**: Clear separation of concerns
2. **Dependency Injection**: Testable and maintainable code
3. **Interface-Driven Design**: Mockable dependencies
4. **Error Handling**: Comprehensive error management
5. **Documentation**: Self-documenting code with OpenAPI

## 📋 API Response Format

All API responses follow a consistent structure:

```json
{
  "success": true,
  "message": "success",
  "data": { ... },
  "timestamp": "2025-10-19T12:14:24+07:00"
}
```

Error responses:
```json
{
  "success": false,
  "message": "Error description",
  "timestamp": "2025-10-19T12:14:24+07:00"
}
```

## 📄 License

This project is part of a senior software engineer assessment and is provided for evaluation purposes.

---

**Note**: This system is designed for security event triage automation in enterprise environments. Ensure proper access controls and monitoring when deploying to production systems.
