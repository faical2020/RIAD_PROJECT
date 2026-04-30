# 🚀 Guide de Construction - API Server Riad

**Objectif**: Construire uniquement l'API REST Go (Gin) pour le projet Riad Manager.

---

## 📋 Prérequis

| Outil | Version | Vérification |
|-------|---------|--------------|
| Go | 1.22+ | `go version` |
| PostgreSQL | 15+ | `psql --version` |
| Docker (optionnel) | 20+ | `docker --version` |
| curl (test) | - | `curl --version` |

---

## 📁 Structure du Server

```
server/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── router.go       # Routes + CORS
│   │   ├── middleware/
│   │   │   └── auth.go   # JWT + Role auth
│   │   └── handlers/
│   │       ├── auth.go      # Register, Login, Me
│   │       ├── chambre.go  # CRUD Chambres
│   │       └── reservation.go # CRUD + Checkin/Checkout
│   ├── db/
│   │   └── postgres.go    # Connexion + Migrations
│   └── logic/
│       ├── models.go      # Structures (User, Chambre, etc.)
│       ├── auth.go        # JWT, rôles
│       ├── chambre.go     # Règles chambres
│       └── reservation.go # Règles réservations
├── .env                    # Variables (DATABASE_URL, JWT_SECRET)
├── go.mod
└── Dockerfile
```

---

## 🛠️ Étapes de Construction

### Task 1: Initialisation du Module Go

```bash
cd /home/faical/Projects/riad-projet/riad/server
go mod init riad-server
```

**Dépendances nécessaires:**
```bash
go get github.com/gin-gonic/gin           # Web framework
go get github.com/golang-jwt/jwt/v5      # JWT authentication
go get golang.org/x/crypto/bcrypt        # Password hashing
go get gorm.io/gorm                     # ORM
go get gorm.io/driver/postgres           # PostgreSQL driver
```

---

### Task 2: Création des Modèles (Logic)

**Fichier**: `internal/logic/models.go`

```go
package logic

import "encoding/json"

type User struct {
    ID        string `json:"id" gorm:"type:uuid;primaryKey"`
    Email     string `json:"email" gorm:"uniqueIndex"`
    Password  string `json:"password,omitempty" gorm:"type:varchar(255)"`
    Nom       string `json:"nom"`
    Prenom    string `json:"prenom"`
    Role      string `json:"role" gorm:"type:varchar(50);default:'client'"`
    Telephone string `json:"telephone"`
}

type Chambre struct {
    ID          string `json:"id" gorm:"type:uuid;primaryKey"`
    Numero      int    `json:"numero"`
    Type        string `json:"type" gorm:"type:varchar(50)"`
    Prix        float64 `json:"prix"`
    Statut      string `json:"statut" gorm:"type:varchar(50);default:'libre'"`
    Description string `json:"description"`
    Equipements string `json:"equipements"`
}

type Reservation struct {
    ID         string  `json:"id" gorm:"type:uuid;primaryKey"`
    UserID     string  `json:"user_id" gorm:"type:uuid"`
    ChambreID  string  `json:"chambre_id" gorm:"type:uuid"`
    DateDebut  string  `json:"date_debut"`
    DateFin    string  `json:"date_fin"`
    Statut     string  `json:"statut" gorm:"type:varchar(50);default:'confirmée'"`
    Montant    float64 `json:"montant"`
}

type Tache struct {
    ID          string `json:"id" gorm:"type:uuid;primaryKey"`
    EmployeID   string `json:"employe_id" gorm:"type:uuid"`
    Description string `json:"description"`
    Statut      string `json:"statut" gorm:"type:varchar(50);default:'à faire'"`
}

type Service struct {
    ID          string  `json:"id" gorm:"type:uuid;primaryKey"`
    Nom         string  `json:"nom"`
    Description string  `json:"description"`
    Prix        float64 `json:"prix"`
}

type Paiement struct {
    ID            string `json:"id" gorm:"type:uuid;primaryKey"`
    ReservationID string `json:"reservation_id" gorm:"type:uuid"`
    Montant       float64 `json:"montant"`
    ModePaiement  string `json:"mode_paiement" gorm:"type:varchar(50)"`
    Statut        string `json:"statut" gorm:"type:varchar(50);default:'en attente'"`
}
```

---

### Task 3: Logique Métier

**Fichier**: `internal/logic/auth.go`

```go
package logic

import "strings"

const (
    RoleClient         = "client"
    RoleEmploye        = "employe"
    RoleReceptionniste = "receptionniste"
    RoleManager        = "manager"
)

func (u User) HasPermission(requiredRoles ...string) bool {
    for _, role := range requiredRoles {
        if u.Role == strings.TrimSpace(role) {
            return true
        }
    }
    return false
}
```

**Fichier**: `internal/logic/chambre.go`

```go
package logic

import "errors"

func ValidateChambre(c Chambre) error {
    if c.Numero <= 0 {
        return errors.New("numéro de chambre invalide")
    }
    if c.Prix <= 0 {
        return errors.New("prix invalide")
    }
    return nil
}

func (c *Chambre) CanBook() bool {
    return c.Statut == "libre"
}
```

**Fichier**: `internal/logic/reservation.go`

```go
package logic

import "errors"

func ValidateReservation(r Reservation, c Chambre) error {
    if r.DateDebut >= r.DateFin {
        return errors.New("dates invalides")
    }
    if !c.CanBook() {
        return errors.New("chambre non disponible")
    }
    return nil
}

func (r *Reservation) Checkin(c *Chambre) error {
    if r.Statut != "confirmée" {
        return errors.New("réservation non confirmée")
    }
    r.Statut = "checkin"
    c.Statut = "occupee"
    return nil
}
```

---

### Task 4: Base de Données (PostgreSQL)

**Fichier**: `internal/db/postgres.go`

```go
package db

import (
    "fmt"
    "log"
    "riad-server/internal/logic"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres(databaseURL string) error {
    var err error
    DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("échec connexion PostgreSQL: %w", err)
    }

    log.Println("Migration des modèles...")
    err = DB.AutoMigrate(
        &logic.User{},
        &logic.Chambre{},
        &logic.Reservation{},
        &logic.Tache{},
        &logic.Service{},
        &logic.Paiement{},
    )
    if err != nil {
        return fmt.Errorf("échec migration: %w", err)
    }

    log.Println("Base PostgreSQL initialisée")
    return nil
}

func GetDB() *gorm.DB {
    return DB
}
```

---

### Task 5: Middleware (Auth + Rôles)

**Fichier**: `internal/api/middleware/auth.go`

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("change_this_in_production")

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "non autorisé"})
            return
        }

        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "accès interdit"})
    }
}
```

---

### Task 6: Handlers (API Endpoints)

**Fichier**: `internal/api/handlers/auth.go`

```go
package handlers

import (
    "net/http"
    "riad-server/internal/db"
    "riad-server/internal/logic"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var jwtSecret = []byte("change_this_in_production")

func Register(c *gin.Context) {
    var user logic.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPassword)

    if err := db.GetDB().Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur création"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "utilisateur créé"})
}

func Login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user logic.User
    if err := db.GetDB().Where("email = ?", credentials.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "identifiants invalides"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
    })
    tokenString, _ := token.SignedString(jwtSecret)

    c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role})
}
```

**Fichier**: `internal/api/handlers/chambre.go`

```go
package handlers

import (
    "net/http"
    "riad-server/internal/db"
    "riad-server/internal/logic"
    "github.com/gin-gonic/gin"
)

func GetChambres(c *gin.Context) {
    var chambres []logic.Chambre
    db.GetDB().Find(&chambres)
    c.JSON(http.StatusOK, chambres)
}

func CreateChambre(c *gin.Context) {
    var chambre logic.Chambre
    if err := c.ShouldBindJSON(&chambre); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.GetDB().Create(&chambre)
    c.JSON(http.StatusCreated, chambre)
}
```

**Fichier**: `internal/api/handlers/reservation.go`

```go
package handlers

import (
    "net/http"
    "riad-server/internal/db"
    "riad-server/internal/logic"
    "github.com/gin-gonic/gin"
)

func GetReservations(c *gin.Context) {
    var reservations []logic.Reservation
    db.GetDB().Find(&reservations)
    c.JSON(http.StatusOK, reservations)
}

func CreateReservation(c *gin.Context) {
    var res logic.Reservation
    if err := c.ShouldBindJSON(&res); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.GetDB().Create(&res)
    c.JSON(http.StatusCreated, res)
}

func Checkin(c *gin.Context) {
    id := c.Param("id")
    var res logic.Reservation
    db.GetDB().First(&res, "id = ?", id)

    var chambre logic.Chambre
    db.GetDB().First(&chambre, "id = ?", res.ChambreID)

    res.Checkin(&chambre)
    db.GetDB().Save(&res).Save(&chambre)

    c.JSON(http.StatusOK, res)
}
```

---

### Task 7: Router (Routes + CORS)

**Fichier**: `internal/api/router.go`

```go
package api

import (
    "github.com/gin-gonic/gin"
    "riad-server/internal/api/handlers"
    "riad-server/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Public routes
    public := r.Group("/api/v1")
    {
        public.POST("/auth/register", handlers.Register)
        public.POST("/auth/login", handlers.Login)
    }

    // Protected routes
    protected := r.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/chambres", handlers.GetChambres)
        protected.POST("/chambres", middleware.RoleMiddleware("manager"), handlers.CreateChambre)
        protected.GET("/reservations", middleware.RoleMiddleware("manager", "receptionniste"), handlers.GetReservations)
        protected.POST("/reservations", handlers.CreateReservation)
        protected.PATCH("/reservations/:id/checkin", middleware.RoleMiddleware("manager", "receptionniste"), handlers.Checkin)
        protected.GET("/auth/me", handlers.GetCurrentUser)
    }

    return r
}
```

---

### Task 8: Entry Point

**Fichier**: `cmd/main.go`

```go
package main

import (
    "log"
    "os"
    "riad-server/internal/api"
    "riad-server/internal/db"
)

func main() {
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        databaseURL = "postgres://postgres:postgres@localhost:5432/riad?sslmode=disable"
    }

    if err := db.InitPostgres(databaseURL); err != nil {
        log.Fatal("Erreur DB:", err)
    }

    router := api.SetupRouter()
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Serveur démarré sur :%s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Erreur serveur:", err)
    }
}
```

---

## 🚀 Construction et Lancement

### 1. Base de données (PostgreSQL via Docker)

```bash
cd /home/faical/Projects/riad-projet/riad
docker compose up postgres -d
# Vérifier: docker ps
```

**Ou manuellement:**
```bash
createdb riad
psql -U postgres -d riad -c "CREATE DATABASE riad;"
```

---

### 2. Build l'API Server

```bash
cd /home/faical/Projects/riad-projet/riad/server

# Installer dépendances
go mod tidy

# Compiler
go build -o riad-server cmd/main.go

# Vérifier
ls -la riad-server  # Doit exister
```

---

### 3. Lancer le Serveur

```bash
cd /home/faical/Projects/riad-projet/riad/server

# Créer .env (optionnel)
cat > .env << EOF
DATABASE_URL=postgres://postgres:postgres@localhost:5432/riad?sslmode=disable
PORT=8080
JWT_SECRET=change_this_in_production
EOF

# Lancer
./riad-server
# Doit afficher: "Serveur démarré sur :8080"
```

---

## 🧪 Tests de l'API

### 1. Test Registration
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456","nom":"Test"}'
# Attendu: {"message":"utilisateur créé"}
```

### 2. Test Login
```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456"}'
# Attendu: {"token":"eyJ...","role":"client"}
```

### 3. Test Protected Endpoint (Chambres)
```bash
# Sans token (doit échouer)
curl -s http://localhost:8080/api/v1/chambres
# Attendu: {"error":"token manquant"}

# Avec token (réussite)
TOKEN="votre_token_ici"
curl -s http://localhost:8080/api/v1/chambres \
  -H "Authorization: Bearer $TOKEN"
# Attendu: [] (tableau vide)
```

### 4. Test Création Chambre (Manager seulement)
```bash
# D'abord créer un manager
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@test.com","password":"123456","nom":"Manager","role":"manager"}'

# Login en tant que manager
TOKEN_MGR=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@test.com","password":"123456"}' | jq -r '.token')

# Créer une chambre
curl -s -X POST http://localhost:8080/api/v1/chambres \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_MGR" \
  -d '{"numero":101,"type":"double","prix":1500}'
```

---

## 📡 Endpoints Disponibles

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

## 🐳 Dockerisation (Optionnel)

**Fichier**: `Dockerfile`

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

**Build et run:**
```bash
docker build -t riad-server .
docker run -p 8080:8080 --env-file .env riad-server
```

---

## ✅ Checklist de Construction

- [ ] Task 1: Module Go initialisé (`go.mod`)
- [ ] Task 2: Modèles créés (`internal/logic/models.go`)
- [ ] Task 3: Logique métier (`auth.go`, `chambre.go`, `reservation.go`)
- [ ] Task 4: DB PostgreSQL (`internal/db/postgres.go`)
- [ ] Task 5: Middleware (`middleware/auth.go`)
- [ ] Task 6: Handlers (`handlers/auth.go`, `chambre.go`, `reservation.go`)
- [ ] Task 7: Router (`api/router.go`)
- [ ] Task 8: Entry point (`cmd/main.go`)
- [ ] Build: `go build -o riad-server cmd/main.go`
- [ ] DB: `docker compose up postgres -d`
- [ ] Test: Registration + Login + Protected endpoints

---

## 🐛 Débogage

**Erreur**: `panic: URL parameters can not be used when serving a static file`
- **Solution**: Ne pas utiliser `r.StaticFile("/app/*", ...)` - utiliser `r.NoRoute()` à la place.

**Erreur**: `failed to parse field: Equipements`
- **Solution**: Utiliser `string` au lieu de `[]string` pour le champ Equipements dans le modèle.

**Erreur**: `connection refused`
- **Solution**: Vérifier que PostgreSQL est démarré: `docker ps` ou `pg_isready`.

---

*Guide créé le 29 Avril 2026 - Riad Manager v1.0*
