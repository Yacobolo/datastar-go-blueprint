<h1 align="center">Datastar Go Blueprint</h1>

<p align="center">
  <strong>A modern, opinionated blueprint for building reactive web applications.</strong>
  <br>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
  </a>
</p>

<p align="center">
  <img src="mascot.png" alt="Datastar Go Blueprint mascot" width="600">
</p>

---

## The Stack

- **Backend:** [Go](https://go.dev/), [SQLite](https://www.sqlite.org/) + [sqlc](https://sqlc.dev/), [Datastar](https://github.com/starfederation/datastar), [gomponents](https://www.gomponents.com/), [gomponents-datastar](https://github.com/maragudk/gomponents-datastar), [Chi](https://github.com/go-chi/chi)
- **Frontend:** [Lit](https://lit.dev/), [Tailwind CSS](https://tailwindcss.com/), [DaisyUI](https://daisyui.com/), [esbuild](https://esbuild.github.io/)
- **Tools:** [Air](https://github.com/air-verse/air), [Hivemind](https://github.com/DarthSim/hivemind), [UPX](https://upx.github.io/)

---

## Quick Start

1. **Clone & Install**

   Shell

   ```
   git clone https://github.com/yacobolo/datastar-go-blueprint.git
   cd datastar-go-blueprint
   go mod tidy
   pnpm --dir web install
   ```

2. **Environment Setup**

   Shell

   ```
   task tools:install  # Installs air, sqlc, golangci-lint, hivemind
   task generate:all   # Generates SQL code
   ```

3. **Run Development Server**

   Shell

   ```
   task dev            # Starts Air, esbuild, and Tailwind via Hivemind
   ```

   Visit the URL printed by `task dev`. Each worktree gets its own app port, NATS port, SQLite DB, and session cookie name.

---

## Project Structure

```
.
├── cmd/server/          # App entry point & build scripts
├── data/                # Local SQLite database files
├── internal/
│   ├── app/             # Application lifecycle & initialization
│   ├── domain/          # Core business logic and entities
│   ├── platform/        # Shared infra (Router, PubSub)
│   ├── services/        # Business logic services and UI state models
│   └── store/           # Database layer (Migrations, SQLC, Repositories)
├── web/
│   ├── resources/       # Static assets & embedded Go files
│   └── ui/              # Frontend source (Lit components, Tailwind entry, TS)
├── internal/ui/
│   ├── core/            # Shared layouts and reusable UI components
│   └── todo/            # Todo page, components, handlers, and routes
├── Procfile.dev         # Dev process management (Hivemind)
├── Taskfile.yml         # Project automation tasks
└── sqlc.yaml            # SQL compiler configuration
```

---

## Key Features

- **Flat UI + Services Architecture:** Shared UI lives in `internal/ui/core`, feature UI in `internal/ui/<feature>`, and business logic in `internal/services`.
- **Type-Safe Everything:** \* [sqlc](https://sqlc.dev/) for database queries.
  - [gomponents-datastar](https://github.com/maragudk/gomponents-datastar) for Datastar attributes.
- **Tailwind + DaisyUI:** Styling stays close to framework defaults with no handwritten app CSS.
- **Single Binary:** Static assets embedded using `go:embed`.

---

## Essential Commands

| **Command**         | **Description**                             |
| ------------------- | ------------------------------------------- |
| `task dev`          | Start hot-reload dev env with worktree-specific ports/state |
| `task generate:all` | Run all code generators (SQL)               |
| `task build`        | Create a production-ready compressed binary |
| `task check`        | Run tests and linters                       |
| `task docker:build` | Build optimized Docker image                |

---

## Deployment

**Docker Quick Start:**

Shell

```
docker build -t datastar-app .
docker run -p 8080:9001 datastar-app
```

**Binary Build:**

Executables are optimized with UPX compression for minimal footprint. Use `task build`.

---

**License:** MIT
