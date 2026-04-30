# RIAD Server - API Endpoints Reference
## 📡 Complete API Endpoints Table
| # | Method | Endpoint | Auth Required | Role(s) | Description | Example cURL Command |
|---|--------|----------|---------------|---------|-------------|---------------------|
| 1 | `POST` | `/api/v1/auth/register` | ❌ No | All | Register a new user | `curl -s -X POST http://localhost:8081/api/v1/auth/register -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"123456","nom":"Nom","prenom":"Prenom"}'` |
| 2 | `POST` | `/api/v1/auth/login` | ❌ No | All | Login and get JWT token | `curl -s -X POST http://localhost:8081/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"123456"}'` |
| 3 | `GET` | `/api/v1/auth/me` | ✅ Yes | All | Get current user info | `curl -s http://localhost:8081/api/v1/auth/me -H "Authorization: Bearer YOUR_TOKEN"` |
| 4 | `GET` | `/api/v1/chambres` | ✅ Yes | All | List all rooms | `curl -s http://localhost:8081/api/v1/chambres -H "Authorization: Bearer YOUR_TOKEN"` |
| 5 | `POST` | `/api/v1/chambres` | ✅ Yes | **manager** | Create a new room | `curl -s -X POST http://localhost:8081/api/v1/chambres -H "Content-Type: application/json" -H "Authorization: Bearer MANAGER_TOKEN" -d '{"numero":105,"type":"double","prix":1500}'` |
| 6 | `GET` | `/api/v1/reservations` | ✅ Yes | **manager, receptionniste** | List all reservations | `curl -s http://localhost:8081/api/v1/reservations -H "Authorization: Bearer MANAGER_TOKEN"` |
| 7 | `POST` | `/api/v1/reservations` | ✅ Yes | All | Create a new reservation | `curl -s -X POST http://localhost:8081/api/v1/reservations -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_TOKEN" -d '{"user_id":"UUID","chambre_id":"UUID","date_debut":"2026-05-01","date_fin":"2026-05-05","montant":6000}'` |
| 8 | `PATCH` | `/api/v1/reservations/:id/checkin` | ✅ Yes | **manager, receptionniste** | Check-in a reservation | `curl -s -X PATCH http://localhost:8081/api/v1/reservations/RESERVATION_ID/checkin -H "Authorization: Bearer MANAGER_TOKEN"` |
| 9 | `PATCH` | `/api/v1/reservations/:id/checkout` | ✅ Yes | **manager, receptionniste** | Check-out a reservation | `curl -s -X PATCH http://localhost:8081/api/v1/reservations/RESERVATION_ID/checkout -H "Authorization: Bearer MANAGER_TOKEN"` |
---
## 🔐 Available Roles
| Role | Description | Access Level |
|------|-------------|---------------|
| `client` | Default role on registration | Can create reservations, view rooms |
| `employe` | Hotel employee | Basic employee access |
| `receptionniste` | Reception staff | Can manage reservations, check-in/out |
| `manager` | Manager | Full access (create rooms, manage all) |
---
## 📦 Request/Response Structures
### User Registration (`POST /auth/register`)
**Request:**
```json
{
  "email": "user@test.com",
  "password": "123456",
  "nom": "Nom",
  "prenom": "Prenom",
  "telephone": "0612345678",
  "role": "client"  // Optional, defaults to "client"
}
Response (201):
{
  "message": "utilisateur créé"
}
---
User Login (POST /auth/login)
Request:
{
  "email": "user@test.com",
  "password": "123456"
}
Response (200):
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "client"
}
---
Create Room (POST /chambres) - Manager Only
Request:
{
  "numero": 101,
  "type": "double",
  "prix": 1500,
  "description": "Chambre luxe",
  "equipements": "TV,WiFi,Clim"
}
Response (201):
{
  "id": "uuid-here",
  "numero": 101,
  "type": "double",
  "prix": 1500,
  "statut": "libre",
  "description": "Chambre luxe",
  "equipements": "TV,WiFi,Clim"
}
---
Create Reservation (POST /reservations)
Request:
{
  "user_id": "user-uuid",
  "chambre_id": "room-uuid",
  "date_debut": "2026-05-01",
  "date_fin": "2026-05-05",
  "montant": 6000
}
Response (201):
{
  "id": "reservation-uuid",
  "user_id": "user-uuid",
  "chambre_id": "room-uuid",
  "date_debut": "2026-05-01",
  "date_fin": "2026-05-05",
  "statut": "confirmée",
  "montant": 6000
}
---
⚡ HTTP Status Codes
Code	Meaning
200	OK
201	Created
400	Bad Request
401	Unauthorized
403	Forbidden
404	Not Found
---
## 🧪 Recommended Test Order
1. `POST /auth/register` → Register a client
2. `POST /auth/register` → Register a manager (add `"role":"manager"`)
3. `POST /auth/login` → Get client token
4. `POST /auth/login` → Get manager token
5. `GET /auth/me` → Verify token works
6. `POST /chambres` → Create room (use manager token)
7. `GET /chambres` → List rooms
8. `POST /reservations` → Create reservation (use client token)
9. `GET /reservations` → List reservations (use manager token)
10. `PATCH /reservations/:id/checkin` → Check-in (manager)
11. `PATCH /reservations/:id/checkout` → Check-out (manager)
---
🔍 Database Inspection Commands
# View all users
docker exec riad-postgres psql -U postgres -d riad -c "SELECT id, email, nom, role FROM users;"
# View all rooms
docker exec riad-postgres psql -U postgres -d riad -c "SELECT numero, type, statut, prix FROM chambres;"
# View all reservations with status
docker exec riad-postgres psql -U postgres -d riad -c "SELECT id, user_id, chambre_id, statut, date_debut, date_fin FROM reservations;"
# Check room status in real-time
docker exec riad-postgres psql -U postgres -d riad -c "SELECT numero, statut FROM chambres;"
---
🚨 Error Response Format
All errors return this format:
{
  "error": "error message description"
}
Common errors:
- "token manquant" → No Authorization header
- "token invalide" → Invalid/expired JWT
- "accès interdit" → Insufficient role
- "identifiants invalides" → Wrong email/password
- "erreur création" → Database error (check server logs)