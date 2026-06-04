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
  - Vue Router
  - Pinia for state management
- Tailwind CSS for styling
- Axios for HTTP requests
- Planned: Chart.js for analytics and Leaflet for maps (not yet integrated)

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

4. Start the development environment:
   ```bash
   make docker-compose-up
   ```

5. Run the services:
   ```bash
   # In one terminal
   make run

   # In another terminal
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
- ✅ Authentication & role-based access control (JWT, `pkg/auth`)
- ✅ Kubernetes manifests (`deployments/k8s/`)
- ✅ Unit tests for auth and HTTP handlers
- 🔄 Frontend ↔ backend integration (inventory wired; shipment views pending)
- 🔄 Analytics charts and map visualisations (planned)

## Authentication

Authentication uses signed JWTs (HMAC-SHA256). The token manager and the
role-based access-control middleware live in [`pkg/auth`](pkg/auth).

- The shipment service protects its `/api/v1` routes — requests must send an
  `Authorization: Bearer <token>` header. Deleting a shipment additionally
  requires the `admin` or `manager` role.
- The inventory service endpoints are currently unauthenticated to keep the
  Vue frontend integration simple.
- Configure the signing secret via the `auth.jwt_secret` config value or the
  `JWT_SECRET` environment variable.

Roles, from most to least privileged: `admin`, `manager`, `operator`, `viewer`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
