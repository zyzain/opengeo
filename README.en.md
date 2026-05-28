# OpenGEO Intelligent Publishing Platform

> AI-era GEO (Generative Engine Optimization) content optimization and multi-platform intelligent publishing system

English | **[中文](README.md)**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.7-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com)

---

## Features

### GEO Content Optimization
- AI Semantic Enhancement (4-dimensional scoring: Structure / Readability / Intent / Schema Markup)
- Multi-model Adaptation (DeepSeek / Kimi / Doubao / ChatGPT)
- Chinese Word Segmentation Keyword Extraction (dictionary longest match + bigram fallback)
- Compliance Detection (sensitive words / advertising law / AIGC labeling)
- Knowledge Graph Entity Management to enhance AI citation authority

### Multi-Platform Publishing
- WeChat / Weibo / Douyin / Xiaohongshu / Zhihu / Toutiao adaptation
- Deduplication Engine (synonym replacement / paragraph reshuffling / sentence pattern transformation)
- Staggered Scheduling (Worker Pool + staggered delay + heatmap)
- 3-tier fault tolerance: Retry / Degradation / Manual Review

### Anti-Ban Engine
- Proxy Pool Management (loading / health check / automatic failure removal)
- Browser Fingerprint Binding (fixed fingerprint per account, least-used assignment)
- Behavior Simulation (normal distribution delay / typing simulation / scroll simulation)
- Platform-level Rate Limiting (per hour / per day / minimum interval)

### Account & Permissions
- Multi-tenant Isolation
- RBAC Permission Control (integrated into all routes, admin auto-bypass)
- Resource Ownership Verification (IDOR protection)

### Monitoring & Analytics
- AI Citation Tracking (citation position / sentiment / model source)
- Source Scoring & Competitor Comparison
- ROI Attribution Analysis

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | React 19 + Vite + TypeScript + Ant Design 5 + Zustand + TanStack Query |
| **Gateway** | Hertz (CloudWeGo) + RBAC Middleware + RateLimiter |
| **Services** | Kitex RPC + Hexagonal Architecture + CQRS + EDA |
| **Storage** | MySQL 8.0 + Redis 7 |
| **Monitoring** | Prometheus + Grafana + Jaeger (OpenTelemetry) |
| **Deployment** | Docker Compose + Nginx Load Balancing |
| **Proto** | Buf + Protobuf + gRPC |

---

## Project Structure

```
opengeo/
├── gateway/                          # HTTP Gateway (Entry Point)
│   ├── cmd/main.go                   # Startup entry
│   └── internal/
│       ├── auth/                     # Authentication (JWT + RBAC)
│       ├── handler/                  # HTTP Handlers (split by domain)
│       │   ├── handler.go            # Common utilities + Handler struct
│       │   ├── auth_handler.go       # Login / Register / Refresh
│       │   ├── user_handler.go       # User / Role / Tenant
│       │   ├── content_handler.go    # Content CRUD + Optimization
│       │   ├── account_handler.go    # Account Management
│       │   ├── knowledge_handler.go  # Knowledge Graph
│       │   ├── publish_handler.go    # Publishing + Channels + Platforms
│       │   ├── schedule_handler.go   # Scheduling
│       │   ├── monitor_handler.go    # Monitoring
│       │   └── system_handler.go     # System Config + Plugins + Webhooks
│       ├── client/                   # Downstream Clients (split by domain)
│       ├── middleware/               # Middleware (CORS / JWT / RBAC / RateLimiter)
│       ├── router/                   # Route Registration
│       └── dal/                      # Data Access Layer
│
├── service/                          # Microservice Layer
│   ├── account/                      # Account Service (Kitex layered)
│   ├── content/                      # Content Service (Hexagonal Architecture)
│   │   └── internal/
│   │       ├── domain/               # Domain Models + GEO Optimization Logic
│   │       ├── application/          # Application Services (Use Case Orchestration)
│   │       ├── port/                 # Port Interfaces (inbound / outbound)
│   │       └── adapter/              # Adapters (Database / AI / Events)
│   ├── publish/                      # Publish Service (Hexagonal + EDA)
│   │   └── internal/
│   │       ├── service/              # Core Services (Anti-ban / Dedup / Retry / Worker Pool)
│   │       ├── domain/               # Domain Models + Events
│   │       ├── adapter/              # Platform Adapters + Kafka
│   │       └── port/                 # Port Interfaces
│   ├── scheduler/                    # Scheduler Service (Heatmap / Staggered / Priority Queue)
│   ├── monitor/                      # Monitor Service (CQRS)
│   └── system/                       # System Service (Plugin SDK / Webhook / Audit)
│
├── pkg/                              # Shared Components
│   ├── ai/                           # AI Service Interface (unified type definitions)
│   ├── similarity/                   # Text Similarity (SimHash / Cosine / Jaccard)
│   ├── crypto/                       # Encryption Utilities (bcrypt)
│   ├── jwt/                          # JWT Utilities
│   ├── config/                       # Configuration Loading
│   ├── database/                     # Database Connection
│   ├── errcode/                      # Error Codes
│   └── eventbus/                     # Event Bus
│
├── proto/                            # Protobuf Interface Definitions
│   ├── opengeo/
│   │   ├── common/v1/                # Common Types (Pagination / Tenant / Tracing)
│   │   ├── internal/v1/              # Internal RPC (account / content / publish / ...)
│   │   ├── cloud/v1/                 # External API (Plans / Subscriptions / API Key)
│   │   ├── brand/v1/                 # Brand API
│   │   ├── publish/v1/               # Publish API
│   │   ├── monitor/v1/               # Monitor API
│   │   └── tenant/v1/                # Tenant API
│   ├── buf.yaml                      # Buf Configuration
│   └── buf.gen.yaml                  # Code Generation Configuration
│
├── web/                              # Frontend (React 19 + Vite)
│   ├── src/
│   │   ├── pages/                    # Pages (30+ pages covering all features)
│   │   ├── components/               # Components
│   │   ├── hooks/                    # React Query Hooks
│   │   ├── stores/                   # Zustand State Management
│   │   ├── lib/                      # API Client (80+ endpoints)
│   │   ├── types/                    # TypeScript Types
│   │   └── i18n/                     # Internationalization (zh-CN / en-US)
│   └── package.json
│
├── configs/                          # Configuration Files
│   ├── config.example.json           # Application Config Example
│   ├── prometheus.yml                # Prometheus Configuration
│   └── nginx.conf                    # Nginx Load Balancing Configuration
│
├── scripts/                          # Scripts
│   └── init.sql                      # Database Initialization
│
├── deployments/                      # Deployment Configs
│   └── Dockerfile.gateway            # Gateway Dockerfile
│
├── docker-compose.yml                # Development Environment Orchestration
├── docker-compose.prod.yml           # Production Environment Orchestration (multi-replica + rolling update)
├── Makefile                          # Build Commands
└── go.mod                            # Go Module Definition
```

---

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7

### 1. Start Infrastructure

```bash
docker compose up -d mysql redis consul jaeger prometheus grafana
```

### 2. Configure Environment Variables

```bash
export MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/opengeo?charset=utf8mb4&parseTime=True&loc=Local"
export ADMIN_PASSWORD="YourSecurePassword@123"  # Required for production
export JWT_SECRET_KEY="your-secret-key"          # Required for production
```

### 3. Start Backend

```bash
make dev-gateway
# Admin: admin / <ADMIN_PASSWORD>
# API: http://localhost:8080
```

### 4. Start Frontend

```bash
cd web && npm install && npm run dev
# Visit: http://localhost:3000
```

---

## Makefile Commands

```bash
# Infrastructure
make dev-up              # Start MySQL / Redis / Consul / Jaeger / Prometheus / Grafana
make dev-down            # Stop all containers

# Development
make dev-gateway         # Start Gateway
make dev-content         # Start Content Service
make dev-publish         # Start Publish Service
make dev-scheduler       # Start Scheduler Service
make dev-monitor         # Start Monitor Service
make dev-system          # Start System Service

# Build
make build               # Build all services
make build-gateway       # Build Gateway

# Test
make test                # Run all tests
make test-unit           # Run unit tests (22 test files, 14 packages passing)

# Proto
make proto-gen           # Generate Go / TS code
make proto-lint          # Proto lint check

# Code Quality
make lint                # golangci-lint
make fmt                 # gofmt
make vet                 # go vet

# Docker
make docker              # Docker full-stack start
make docker-down         # Stop Docker
```

---

## Security Features

| Feature | Description |
|---------|-------------|
| RBAC Permissions | All routes verified by resource:action, admin auto-bypass |
| IDOR Protection | Content / Account / Entity / Task / Schedule verify user_id ownership |
| XSS Protection | PreviewPublish uses html.EscapeString for escaping |
| SSRF Protection | Webhook URLs block internal addresses (localhost / 169.254 / private IPs) |
| RateLimiter | Global rate limiting + login endpoint 5 requests/min per IP |
| Input Validation | page_size capped at 100, error messages use safeError sanitization |
| Password Security | bcrypt encryption, proxy password response sanitization, Admin password from env vars |
| Prompt Injection | AI calls use delimiters to wrap user content |

---

## Architecture Patterns

| Service | Pattern | Description |
|---------|---------|-------------|
| Gateway | Standard Layered | router / handler / middleware / client |
| Content | Hexagonal | domain / application / port / adapter |
| Publish | Hexagonal + EDA | Event-driven + Worker Pool + Anti-ban engine |
| Monitor | CQRS | Read-write separation (command / query) |
| Scheduler | Hexagonal | Heatmap + Staggered scheduling + Priority queue |
| System | Registry | Plugin SDK + Webhook + Audit log |

---

## API Overview

### Authentication
```
POST   /api/v1/auth/login          # Login
POST   /api/v1/auth/register       # Register
POST   /api/v1/auth/refresh        # Refresh Token
```

### Content Management
```
GET    /api/v1/contents             # List (paginated)
POST   /api/v1/contents             # Create
GET    /api/v1/contents/:id         # Get by ID
PUT    /api/v1/contents/:id         # Update
DELETE /api/v1/contents/:id         # Delete
POST   /api/v1/contents/:id/optimize  # AI Optimization
POST   /api/v1/contents/:id/publish   # Publish
```

### Publishing Management
```
GET    /api/v1/publish/tasks        # Task List
POST   /api/v1/publish/tasks        # Create Task
POST   /api/v1/publish/tasks/:id/cancel  # Cancel
POST   /api/v1/publish/preview      # Preview
POST   /api/v1/publish/dedup/check  # Dedup Check
```

### Knowledge Graph
```
GET    /api/v1/knowledge/entities    # Entity List
POST   /api/v1/knowledge/entities    # Create Entity
GET    /api/v1/knowledge/entities/search  # Search
```

### System Management
```
GET    /api/v1/system/configs       # System Configs
GET    /api/v1/system/plugins       # Plugin List
GET    /api/v1/system/webhooks      # Webhook List
```

For complete API documentation, see [API.md](API.md)

---

## Testing

```bash
# Run all tests
go test ./... -short -v

# Run specific package tests
go test ./service/publish/internal/service/... -v
go test ./service/content/internal/domain/service/... -v
go test ./pkg/similarity/... -v
```

Test Coverage:
- 22 test files
- 14 packages passing
- Includes unit tests + Benchmark performance baselines

---

## Deployment

### Development Environment

```bash
docker compose up -d          # Start infrastructure
make dev-gateway              # Start Gateway
cd web && npm run dev         # Start frontend
```

### Production Environment

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

Production Configuration Highlights:
- Gateway 3 replicas, Content / Publish 2 replicas
- Rolling updates (start-first, zero-downtime)
- Nginx load balancing + rate limiting
- Resource limits (CPU / Memory)

---

## Documentation

| Document | Description |
|----------|-------------|
| [README.md](README.md) | This document (Chinese) |
| [README.en.md](README.en.md) | This document (English) |
| [API.md](API.md) | RESTful API Reference |
| [proto/README.md](proto/README.md) | Protobuf Interface Definitions & Code Generation |

---

## License

MIT
