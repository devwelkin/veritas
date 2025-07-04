# Veritas URL Shortener

A **cloud-native, microservice-driven** URL shortener built with Go, React, and a modern DevOps tool-chain. Veritas is designed as a **portfolio-grade** project that demonstrates scalable architecture, Infrastructure-as-Code, and fully automated CI/CD on Azure Kubernetes Service (AKS).

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Microservices](#microservices)
3. [Data Flow](#data-flow)
4. [Key Features](#key-features)
5. [Tech Stack](#tech-stack)
6. [Getting Started](#getting-started)
   - [Local Development](#local-development-docker-compose)
   - [Environment Variables](#environment-variables)
7. [Production Deployment](#production-deployment-azure--terraform)
8. [CI/CD Pipeline](#cicd-pipeline)
9. [Infrastructure as Code](#infrastructure-as-code)
10. [Contributing](#contributing)
11. [License](#license)

---

## Architecture Overview

Veritas follows the **separation of concerns** principle with a clear split between frontend and backend services:

```
┌──────────┐        create       ┌──-───────────────┐
│ Frontend │ ───────────────▶    | Creator Service  │
└──────────┘                     └─-────────────────┘
    ▲                                   │
    │                                   ▼
┌────────────────────────────────────────────────┐
│              PostgreSQL Database              │
└────────────────────────────────────────────────┘

┌──────────┐   redirect ({code}) ┌─────────────────┐
│ Browser  │ ───────────────────▶│ Redirector Svc  │
└──────────┘                     └─────────────────┘
                                        │ cache
                                        ▼
                               ┌─────────────────┐
                               │      Redis      │
                               └─────────────────┘
                                        │ miss
                                        ▼
                               ┌─────────────────┐
                               │  PostgreSQL DB  │
                               └─────────────────┘
                                        │ publish
                                        ▼
                               ┌─────────────────┐
                               │    NATS Bus     │
                               └─────────────────┘
                                        │ consume
                                        ▼
                               ┌─────────────────┐
                               │ Analytics Svc   │
                               └─────────────────┘
```

All services are containerised with **Docker**, orchestrated by **Kubernetes** (AKS), and exposed via **Traefik Ingress**.

## Microservices

| Service | Language | Responsibility |
|---------|----------|----------------|
| **Frontend** | React + TypeScript | User interface for creating short URLs |
| **Creator Service** | Go | Handles `POST /api/create` and persists new short URLs |
| **Redirector Service** | Go | Resolves `{short_code}` requests, uses Redis for ultra-fast look-ups |
| **Analytics Service** | Go | Consumes redirect events from NATS and processes analytics |

## Data Flow

1. **Create** ‑ Frontend → Creator Service → PostgreSQL
2. **Redirect** ‑ Browser → Redirector Service → Redis → (cache miss) PostgreSQL → NATS → Analytics Service

## Key Features

- ⚡ **High-Performance Redirects** – In-memory caching with Redis to minimise DB hits.
- 📈 **Asynchronous Analytics** – NATS decouples redirect latency from analytics processing.
- ☁️ **Infrastructure as Code** – Entire Azure footprint defined in Terraform.
- 🔄 **Zero-Touch CI/CD** – GitHub Actions builds, tests, pushes, and deploys on each version tag.
- 🛠 **Modern Tool-Chain** – Go modules, sqlc, React, Vite, Tailwind CSS.

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Backend | **Go** | High-performance, concurrent microservices |
| Storage | **PostgreSQL** | Relational data store |
| Cache | **Redis** | Ultra-low-latency key-value cache |
| Messaging | **NATS** | Lightweight event streaming |
| Code Gen | **sqlc** | Type-safe Go from SQL |
| Frontend | **React + TypeScript** | Modern SPA |
| Build | **Vite** | Lightning-fast dev server & bundler |
| Styling | **Tailwind CSS** / **shadcn/ui** | Design system |
| Containerisation | **Docker** | Image packaging |
| Orchestration | **Kubernetes (AKS)** | Deployment & scaling |
| IaC | **Terraform** | Declarative Azure resources |
| CI/CD | **GitHub Actions** | Automated pipeline |
| Networking | **Traefik** | Ingress controller |

## Getting Started

### Prerequisites

- [Docker ≥ 20.10](https://docs.docker.com/get-docker/)
- [Docker Compose Plugin](https://docs.docker.com/compose/)
- Make `DATABASE_URL`, `REDIS_URL`, etc. available via a `.env` file (see _Environment Variables_).

### Local Development (docker-compose)

```bash
# 1. Clone the repository
git clone https://github.com/<your-username>/veritas.git
cd veritas

# 2. Create a .env file
touch .env  # populate with DB, Redis, & NATS URLs

# 3. Fire up the stack (build images on first run)
docker compose up --build
```

Endpoints:

- Frontend – <http://localhost:5173>
- API Gateway – <http://localhost:8080>

### Environment Variables

| Variable | Example | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://veritas:veritas@db:5432/veritas?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `redis://redis:6379` | Redis server URL |
| `NATS_URL` | `nats://nats:4222` | NATS server URL |
| `BASE_URL` | `http://localhost:8080` | Base URL for generated links |

> _Tip_: A ready-to-fill [.env.example](./.env.example) file is included.

## Production Deployment (Azure & Terraform)

The entire platform can be provisioned and deployed with **one** `terraform apply` and a **single** Git tag push.

### One-Time Setup

1. Install [Azure CLI](https://learn.microsoft.com/cli/azure/install-azure-cli) & [Terraform](https://developer.hashicorp.com/terraform/downloads).
2. Create an **Azure Service Principal** and add its credentials to repository secrets:
   - `AZURE_CREDENTIALS`
   - `ACR_LOGIN_SERVER`
3. **Apply Kubernetes Secrets** – Sensitive variables such as `DATABASE_URL` and `REDIS_URL` live in [`k8s/secret.yaml`](./k8s/secret.yaml). Replace the placeholder values and run:

```bash
kubectl apply -f k8s/secret.yaml  # only once, or when secrets change
```



### Deployment Workflow

```bash
# 1. Provision Azure resources
cd terraform
terraform init && terraform apply  # creates AKS, ACR, Static IP, Traefik, ...

# 2. Tag the commit you want to deploy
git tag v1.0.0

# 3. Push the tag to trigger GitHub Actions
git push origin v1.0.0
```

The pipeline will:

1. Run unit tests for all Go services and the React frontend _in parallel_.
2. Build and cache Docker images (GitHub Cache + ACR layer cache).
3. Push changed images to **Azure Container Registry**.
4. Apply Kubernetes manifests to **AKS**.

## CI/CD Pipeline

The pipeline lives in `.github/workflows/*` and is intentionally **tag-driven** to ensure controlled releases.

```
name: veritas-pipeline
on:
  push:
    tags: ['v*.*.*']  # e.g. v1.0.0
```

## Infrastructure as Code

All Azure resources are declared in [`terraform/`](./terraform). Key modules:

- `aks` – Managed Kubernetes Cluster
- `acr` – Container Registry
- `network` – Static IP & DNS
- `ingress` – Helm-managed Traefik install

Apply workflow is kept minimal:

```bash
terraform init
terraform apply
```

State is stored securely in **Azure Storage** (see backend config).


## License

Distributed under the MIT License. See [`LICENSE`](./LICENSE) for more information. 