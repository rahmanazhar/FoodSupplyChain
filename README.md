# Food Supply Chain Tracker

A decentralized event-driven system for tracking goods across multiple stakeholders in a supply chain.

## Features

- Real-time tracking of goods across the supply chain
- Role-based access control for different stakeholders
- Event-driven architecture for real-time updates
- Inventory management with alerts
- Analytics dashboard
- Comprehensive audit trail
- Multi-tenant support

## Architecture

The system is built using a microservices architecture with the following components:

- Inventory Service: Manages products and inventory levels
- Shipment Service: Handles shipment tracking and status updates
- API Gateway: Routes requests and handles authentication
- Frontend: Vue.js dashboard for visualization and management

## Technology Stack

### Backend
- Go (Golang) 1.21+
  - Standard library for HTTP servers
  - gorilla/mux for routing
  - GORM for database operations
- PostgreSQL 15+ for data storage
- NATS for event streaming
- Docker & Kubernetes for deployment

### Frontend
- Vue.js 3
  - Composition API
  - Vue Router (with auth route guards)
  - Pinia for state management
- Tailwind CSS design system with light/dark themes
- Custom SVG charts (no chart dependency); Vitest unit tests
- Axios for HTTP requests
- Planned: Leaflet for maps (not yet integrated)

### Infrastructure
- Docker & Docker Compose for local development
- Kubernetes for production deployment
- Helm for package management
- PostgreSQL with multi-tenant support
- NATS with JetStream enabled

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- NATS Server
- Make
- Git

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/rahmanazhar/FoodSupplyChain.git
   cd FoodSupplyChain
   ```

2. Install backend dependencies:
   ```bash
   make dev-deps
   go mod download
   ```

3. Install frontend dependencies:
   ```bash
   cd frontend/foodsupplychain
   npm install
   cd ../..
   ```

4. Start the backend stack (Postgres + NATS + all services) with one command:
   ```bash
   ./scripts/run.sh           # add --seed to insert a demo product/location
   ```
   This brings up Docker (Postgres + NATS), builds and launches the inventory,
   shipment and gateway services, and tears everything down on Ctrl-C.
   Equivalent: `make run` (or `make run-seed`).

5. In another terminal, start the frontend:
   ```bash
   cd frontend/foodsupplychain
   npm run dev
   ```

### Development Commands

- `make build`: Build all services
- `make test`: Run tests
- `make lint`: Run linter
- `make docker-build`: Build Docker images
- `make docker-compose-up`: Start development environment
- `make docker-compose-down`: Stop development environment

### Frontend Commands
```bash
cd frontend/foodsupplychain
npm run dev        # Start development server
npm run build      # Build for production
npm run test:unit  # Run unit tests
npm run test:e2e   # Run end-to-end tests
npm run lint       # Run linter
```

## Project Structure

```
.
├── api/                    # API definitions and protocols
├── build/                  # Build configurations and Dockerfiles
├── cmd/                    # Service entry points
│   ├── gateway/           # API Gateway service
│   ├── inventory/         # Inventory service
│   └── shipment/          # Shipment service
├── configs/               # Configuration files
├── deployments/           # Deployment configurations
│   ├── docker/           # Docker Compose files
│   └── k8s/              # Kubernetes manifests
├── frontend/             # Frontend application
│   └── foodsupplychain/  # Vue.js application
└── internal/             # Private application code
    ├── common/          # Shared utilities
    ├── gateway/         # Gateway service implementation
    ├── inventory/       # Inventory service implementation
    └── shipment/        # Shipment service implementation
```

## Current Status

- ✅ Project setup and architecture
- ✅ Core service structure
- ✅ Frontend Vue.js setup with Tailwind CSS
- ✅ Inventory and shipment services (HTTP handlers wired to the GORM service layer)
- ✅ Database schema + auto-migration (GORM models)
- ✅ Event publishing via NATS JetStream
- ✅ Username/password auth, registration, and role-based access control (JWT)
- ✅ Full Vue dashboard (sidebar, light/dark) wired to the API end-to-end
- ✅ Server-side pagination + search on inventory and shipments
- ✅ Admin user management (list users, change roles)
- ✅ Analytics charts on the dashboard (custom SVG)
- ✅ Observability: structured `slog` logs, request IDs, Prometheus `/metrics`
- ✅ Hardening: rate-limited auth, security headers, panic recovery
- ✅ Kubernetes manifests with probes, resource limits, and scrape annotations
- ✅ Tests (Go + Vitest) and GitHub Actions CI
- 🔄 Map visualisations (Leaflet, planned)

## Production hardening

- **Observability** — every service emits structured JSON logs (`log/slog`) with
  a propagated `X-Request-ID`, and exposes Prometheus metrics at `GET /metrics`
  (request counts, in-flight gauge, latency). Shared middleware lives in
  [`pkg/httpx`](pkg/httpx) and [`pkg/metrics`](pkg/metrics).
- **Security** — `/auth/login` and `/auth/register` are rate-limited per client
  IP; responses carry hardening headers (`X-Content-Type-Options`,
  `X-Frame-Options`, `Referrer-Policy`); panics are recovered into clean 500s.
- **API** — list endpoints are paginated and searchable
  (`/inventory?limit=&offset=&search=`, `/shipments?...&status=`) returning
  `{ data, total, limit, offset }`. `POST /auth/refresh` re-issues tokens.
- **CI** — [`.github/workflows/ci.yml`](.github/workflows/ci.yml) builds, vets,
  gofmt-checks and tests the Go services, and lints, unit-tests and builds the
  frontend on every push and PR.

## Authentication

Users sign in with a username and password. The API gateway owns a
database-backed user store (bcrypt-hashed passwords) and issues signed JWTs
(HMAC-SHA256); the JWT helpers and role-based middleware live in
[`pkg/auth`](pkg/auth).

Endpoints (on the gateway):

- `POST /auth/register` — create an account (new users get the `viewer` role)
- `POST /auth/login` — exchange credentials for a JWT
- `GET /auth/me` — the current user (requires a bearer token)

Downstream, the shipment service validates the gateway-issued token on its
`/api/v1` routes; deleting a shipment additionally requires the `admin` or
`manager` role. The gateway and shipment service share the signing secret via
the `JWT_SECRET` environment variable (`scripts/run.sh` generates a fresh random
one per run).

Seeded demo accounts (created on first run): `admin/admin123`,
`manager/manager123`, `operator/operator123`, `viewer/viewer123`.

Roles, from most to least privileged: `admin`, `manager`, `operator`, `viewer`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
