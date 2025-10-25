# KineticaFS

[![CI](https://github.com/argon-chat/KineticaFS/actions/workflows/ci.yml/badge.svg)](https://github.com/argon-chat/KineticaFS/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/argon-chat/KineticaFS)](https://goreportcard.com/report/github.com/argon-chat/KineticaFS)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

KineticaFS is a region-aware file lifecycle manager for S3-based storage. It tracks file usage across regions, migrates
hot files closer to clients, and safely removes unused ones using reference-counting and GC.

## ✨ Features

- 📦 **S3-compatible**: Works with any S3-compatible storage backend (AWS, MinIO, Wasabi, etc.)
- 🌍 **Region-aware**: Detects file access patterns and migrates "hot" files closer to clients
- 🧠 **Smart pointers**: Tracks file references in your system to prevent premature deletion
- ♻️ **Garbage collection**: Removes unreferenced or expired files safely and automatically

## 🚀 Quick Start

### Prerequisites

- Go 1.22+ 
- Make
- Docker (optional, for containerized deployment)

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/argon-chat/KineticaFS.git
   cd KineticaFS
   ```

2. Setup development environment:
   ```bash
   make dev-setup
   ```

3. Build and run:
   ```bash
   make build
   make run
   ```

The server will start on `http://localhost:3000` with Swagger documentation available at `/swagger/index.html`.

### Available Make Targets

- `make help` - Show all available targets
- `make all` - Run full build pipeline (format, lint, test, build, docs)
- `make build` - Build the application
- `make test` - Run tests with coverage
- `make lint` - Run code linters
- `make format` - Format code
- `make docs` - Generate API documentation
- `make run` - Run the server locally
- `make clean` - Clean build artifacts
- `make docker-build` - Build Docker image

### Running with Docker

```bash
make docker-build
make docker-run
```

Or use the provided Docker Compose:

```bash
docker compose up -d
```

## �️ Database Migrations

KineticaFS supports both PostgreSQL and Scylla/Cassandra databases. Database migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate).

### Creating New Migrations

To create a new migration file, use the `migrate` CLI tool:

```bash
# For PostgreSQL migrations
migrate create -ext sql -seq -dir migrations/postgres/ migration_name

# For Scylla migrations  
migrate create -ext sql -seq -dir migrations/scylla/ migration_name
```

This will create two files:
- `XXXXXX_migration_name.up.sql` - Contains the migration logic
- `XXXXXX_migration_name.down.sql` - Contains the rollback logic

### Migration Naming Convention

- Use descriptive names: `create_users_table`, `add_email_index`, `update_schema_v2`
- Keep names concise but clear about the change being made
- Use snake_case for consistency

### Database Support

- **PostgreSQL**: Full relational database support with foreign keys and constraints
- **Scylla/Cassandra**: NoSQL distributed database with application-level relationship management

## �📈 Roadmap

- [ ] File reference tracking API (`CreateRef`, `DeleteRef`, `ListRefs`) 🔥
- [ ] File upload 🔥
- [ ] Scylla Cassandra Support 🔥
- [ ] Migration logic
- [ ] Per-region heatmap tracking
- [ ] GC for unreferenced files 🔥
- [ ] Basic observability (logs, metrics)
- [ ] Public and expiring file links
- [ ] Optional TTL per reference
- [ ] Support for batch import/export
- [ ] Multi-tenant support
- [ ] NATS hook for event-driven GC
- [ ] Custom metadata indexing
- [ ] Integration with Prometheus / Grafana
- [ ] WASM hooks for file filters / pre-upload logic
- [ ] Admin panel / metrics endpoint
- [ ] K8s Operator
- [ ] PgSql Support
- [ ] Dotnet Client Library 🔥

## 📜 License

KineticaFS is licensed under the **GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later)**.  
This ensures that all improvements and deployments based on this code must remain open source.

See [`LICENSE`](./LICENSE) for the full text.
