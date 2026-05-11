# Monorepo Structure: Platforma SP (Enterprise Blueprint)

```text
/platforma-sp-monorepo
├── /apps
│   ├── /landing-service       # Frontend: Landing, SEO, Marketing
│   ├── /dashboard-app         # Frontend: Main Application (React/Next.js)
│   └── /admin-panel           # Frontend: Moderator/Admin view
├── /services
│   ├── /auth-service          # Node.js: Authentication, JWT, OAuth2
│   ├── /matching-service      # Python/FastAPI: AI Logic, Embeddings, Scorpio
│   ├── /proposal-engine       # Node.js: Task lifecycle management
│   ├── /parser-service        # Python/Scrapy/Playwright: Website scraping
│   └── /notification-service  # Node.js: Email, Push, WebSockets
├── /packages
│   ├── /shared-ui             # React: UI-кит (Base components, design tokens)
│   ├── /shared-types          # TS/JSON Schema: DTOs, API models
│   ├── /shared-utils          # Common logic: logger, error helper
│   └── /database-client       # Prisma/Drizzle: Shared DB access layer
├── /infra
│   ├── /docker                # Dockerfiles, docker-compose.yml
│   ├── /k8s                   # Kubernetes manifests / Helm charts
│   ├── /terraform             # IaC: Cloud resources (AWS/GCP/Azure)
│   └── /monitoring            # Prometheus/Grafana configs
├── /scripts                   # Deployment & CI/CD scripts
├── /tests                     # End-to-end and integration tests
├── .gitignore
├── turbo.json                 # Monorepo build orchestrator (Turborepo)
├── package.json
└── README.md
```

## Details on Scaling
*   **Horizontal Scaling:** Each service in `/services` can be scaled independently using K8s.
*   **Matching Service:** Requires GPU-accelerated instances for embedding generation if load is high.
*   **Shared Libraries:** Help maintain consistency between frontend and backend.
