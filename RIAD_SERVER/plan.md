# RIAD Server - Docker-First Implementation Plan
*Replanned with Docker first, `go mod init .` for automatic module naming, `RIAD_SERVER` module path, and all original `server_api.md` tasks preserved + fixes for missing components.*
---
## Phase 1: Docker Foundation (Priority 1)
All Docker setup is completed first to provide a running PostgreSQL instance before code implementation.
1. **Verify Docker Prerequisites**
   - Check installed versions:
     ```bash
     docker --version
     docker compose --version
     ```
2. **Create `docker-compose.yml`**
   - Save to project root `/home/faical/Projects/RIAD_PROJECT/docker-compose.yml`:
     ```yaml
     services:
       postgres:
         image: postgres:15-alpine
         container_name: riad-postgres
         environment:
           POSTGRES_USER: postgres
           POSTGRES_PASSWORD: postgres
           POSTGRES_DB: riad
         volumes:
           - riad-postgres-data:/var/lib/postgresql/data
         ports:
           - "5432:5432"
         healthcheck:
           test: ["CMD-SHELL", "pg_isready -U postgres"]
           interval: 5s
           timeout: 5s
           retries: 5
       server:
         build:
           context: ./RIAD_SERVER
           dockerfile: Dockerfile
         container_name: riad-server
         depends_on:
           postgres:
             condition: service_healthy
         ports:
           - "8081:8080"
         env_file:
           - ./RIAD_SERVER/.env
         environment:
           DATABASE_URL: postgres://postgres:postgres@postgres:5432/riad?sslmode=disable
         restart: unless-stopped
     volumes:
       riad-postgres-data:
     ```
3. **Configure Environment Variables**
   - Create `/home/faical/Projects/RIAD_PROJECT/RIAD_SERVER/.env`:
     ```env
     DATABASE_URL=postgres://postgres:postgres@postgres:5432/riad?sslmode=disable
     JWT_SECRET=change_this_in_production
     PORT=8080
     ```
4. **Create Server Dockerfile**
   - Save to `/home/faical/Projects/RIAD_PROJECT/RIAD_SERVER/Dockerfile` (multi-stage build):
     ```dockerfile
     FROM golang:1.22-alpine AS builder
     WORKDIR /app
     COPY go.mod go.sum ./
     RUN go mod download
     COPY . .
     RUN CGO_ENABLED=0 GOOS=linux go build -o riad-server cmd/main.go
     FROM alpine:latest
     RUN apk --no-cache add ca-certificates
     WORKDIR /root/
     COPY --from=builder /app/riad-server .
     EXPOSE 8080
     CMD ["./riad-server"]
     ```
5. **Start PostgreSQL Container**
   - Launch only the database service first:
     ```bash
     cd /home/faical/Projects/RIAD_PROJECT
     docker compose up postgres -d
     ```
   - Verify DB readiness:
     ```bash
     docker ps
     docker exec -it riad-postgres pg_isready -U postgres
     ```
---
## Phase 2: Go Module Initialization (Using `go mod init .`)
Avoids hardcoded module names, automatically matches the `RIAD_SERVER` directory name.
6. **Initialize Go Module**
   - Navigate to server directory and run `go mod init .`:
     ```bash
     cd /home/faical/Projects/RIAD_PROJECT/RIAD_SERVER
     go mod init .
     ```
   - This automatically sets the module name to `RIAD_SERVER` (matches directory casing/underscore exactly).
7. **Install Dependencies**
   - Add all required packages:
     ```bash
     go get github.com/gin-gonic/gin
     go get github.com/golang-jwt/jwt/v5
     go get golang.org/x/crypto/bcrypt
     go get gorm.io/gorm
     go get gorm.io/driver/postgres
     ```
---
## Phase 3: Core Implementation (Original Tasks 2–8)
All import paths use the auto-generated `RIAD_SERVER/` prefix from `go mod init .`.
8. **Task 2: Data Models**
   - Create `internal/logic/models.go` with exact structs from `server_api.md` (User, Chambre, Reservation, Tache, Service, Paiement with GORM tags). No external imports required for this file.
9. **Task 3: Business Logic**
   - Implement:
     - `internal/logic/auth.go`: Role constants, `HasPermission` method
     - `internal/logic/chambre.go`: `ValidateChambre`, `CanBook` methods
     - `internal/logic/reservation.go`: `ValidateReservation`, `Checkin` methods
10. **Task 4: PostgreSQL Integration**
    - Create `internal/db/postgres.go` with GORM connection and `AutoMigrate` for all models. Updated imports:
      ```go
      import (
          "RIAD_SERVER/internal/logic"
          "gorm.io/driver/postgres"
          "gorm.io/gorm"
      )
      ```
11. **Task 5: Auth Middleware**
    - Create `internal/api/middleware/auth.go` with `AuthMiddleware` (JWT validation) and `RoleMiddleware` (role-based access). No module-specific imports required.
12. **Task 6: API Handlers**
    - Implement:
      - `internal/api/handlers/auth.go`: `Register`, `Login`, **add missing `GetCurrentUser`** (referenced in original router but not provided in `server_api.md`)
      - `internal/api/handlers/chambre.go`: `GetChambres`, `CreateChambre`
      - `internal/api/handlers/reservation.go`: `GetReservations`, `CreateReservation`, `Checkin`, **add missing `Checkout`** (referenced in endpoint table)
    - All handler imports use auto-generated module prefix:
      ```go
      import (
          "RIAD_SERVER/internal/db"
          "RIAD_SERVER/internal/logic"
          "github.com/gin-gonic/gin"
      )
      ```
13. **Task 7: Router Setup**
    - Create `internal/api/router.go` with CORS config, public/protected route groups. Fixes:
      - Correct `GetCurrentUser` reference
      - Add missing `PATCH /reservations/:id/checkout` route
    - Updated imports:
      ```go
      import (
          "RIAD_SERVER/internal/api/handlers"
          "RIAD_SERVER/internal/api/middleware"
          "github.com/gin-gonic/gin"
      )
      ```
14. **Task 8: Entry Point**
    - Create `cmd/main.go` to load env vars, initialize DB, start Gin server. Updated imports:
      ```go
      import (
          "RIAD_SERVER/internal/api"
          "RIAD_SERVER/internal/db"
      )
      ```
---
## Phase 4: Build & Deploy
15. **Sync Dependencies**
    - Run `go mod tidy` to resolve all dependencies.
16. **Build Server**
    - Local test build:
      ```bash
      cd /home/faical/Projects/RIAD_PROJECT/RIAD_SERVER
      go build -o riad-server cmd/main.go
      ```
    - Docker build (via compose):
      ```bash
      cd /home/faical/Projects/RIAD_PROJECT
      docker compose build server
      ```
17. **Launch Full Stack**
    - Start all services:
      ```bash
      docker compose up -d
      ```
    - Verify server logs:
      ```bash
      docker compose logs server
      ```
---
## Phase 5: API Testing
18. **Test Public Endpoints**
    - Register user:
      ```bash
      curl -s -X POST http://localhost:8080/api/v1/auth/register \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"123456","nom":"Test"}'
      ```
    - Login and capture token:
      ```bash
      curl -s -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"123456"}'
      ```
19. **Test Protected Endpoints**
    - List chambres with valid token:
      ```bash
      TOKEN="<your_jwt_token>"
      curl -s http://localhost:8080/api/v1/chambres \
        -H "Authorization: Bearer $TOKEN"
      ```
20. **Test Role-Restricted Endpoints**
    - Create manager user, test chambre creation and check-in/checkout with manager token.
---
## Phase 6: Optional Enhancements
21. **Complete Missing Endpoints**
    - Finalize `Checkout` handler logic for reservations
    - Add unit tests for handlers and business logic
22. **Production Hardening**
    - Move `JWT_SECRET` to Docker secrets
    - Add request logging and rate limiting middleware
    - Enable GORM production-safe config (disable debug mode)
---
## Final Endpoint Table
| Méthode | Endpoint | Auth | Rôle | Description |
|---------|----------|------|------|-------------|
| POST | `/api/v1/auth/register` | ❌ | Tous | Inscription |
| POST | `/api/v1/auth/login` | ❌ | Tous | Connexion |
| GET | `/api/v1/auth/me` | ✅ | Tous | Utilisateur actuel |
| GET | `/api/v1/chambres` | ✅ | Tous | Liste chambres |
| POST | `/api/v1/chambres` | ✅ | manager | Créer chambre |
| GET | `/api/v1/reservations` | ✅ | manager, receptionniste | Liste réservations |
| POST | `/api/v1/reservations` | ✅ | Tous | Créer réservation |
| PATCH | `/api/v1/reservations/:id/checkin` | ✅ | manager, receptionniste | Check-in |
| PATCH | `/api/v1/reservations/:id/checkout` | ✅ | manager, receptionniste | Check-out |
---
## Final Checklist
- [ ] Phase 1: Docker Foundation complete (PostgreSQL running)
- [ ] Phase 2: Go module initialized with `go mod init .` (module name = `RIAD_SERVER`)
- [ ] Phase 3: All core implementation tasks 2–8 complete
- [ ] Phase 4: Server built and full stack running
- [ ] Phase 5: All API tests pass
- [ ] Phase 6: Optional enhancements implemented
Paste this full content into your existing /home/faical/Projects/RIAD_PROJECT/plan.md manually (file writes are disabled in plan mode).