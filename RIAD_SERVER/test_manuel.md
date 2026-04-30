Guide de Test Manuel Complet - API Riad + Inspection Base de Données
📡 Tests Manuels de Toutes les API
⚙️ Configuration Initiale
# URL de base
BASE_URL="http://localhost:8081"
# Préparation : Inscription et récupération des tokens
# 1. Créer un client
curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"client@test.com","password":"123456","nom":"Client","prenom":"Test"}'
# 2. Créer un manager
curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@test.com","password":"123456","nom":"Manager","prenom":"Test","role":"manager"}'
# 3. Login client (récupérer TOKEN_CLIENT)
TOKEN_CLIENT=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"client@test.com","password":"123456"}' | jq -r '.token')
# 4. Login manager (récupérer TOKEN_MGR)
TOKEN_MGR=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manager@test.com","password":"123456"}' | jq -r '.token')
echo "TOKEN_CLIENT: $TOKEN_CLIENT"
echo "TOKEN_MGR: $TOKEN_MGR"
---
🔐 1. Endpoints d'Authentification
1.1 Inscription (Public)
curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nouveau@test.com",
    "password": "123456",
    "nom": "Nouveau",
    "prenom": "User",
    "telephone": "0612345678"
  }' | jq .
1.2 Connexion (Public)
curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"nouveau@test.com","password":"123456"}' | jq .
1.3 Utilisateur Actuel (Protégé - Tous rôles)
curl -s $BASE_URL/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN_CLIENT" | jq .
---
🏨 2. Gestion des Chambres
2.1 Lister toutes les chambres (Protégé - Tous rôles)
curl -s $BASE_URL/api/v1/chambres \
  -H "Authorization: Bearer $TOKEN_CLIENT" | jq .
2.2 Créer une chambre (Protégé - Manager seulement)
curl -s -X POST $BASE_URL/api/v1/chambres \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_MGR" \
  -d '{
    "numero": 102,
    "type": "suite",
    "prix": 2500,
    "description": "Suite présidentielle",
    "equipements": "TV,WiFi,Clim,Mini-bar,Jacuzzi"
  }' | jq .
2.3 Tenter création sans le rôle manager (Doit échouer avec 403)
curl -s -X POST $BASE_URL/api/v1/chambres \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_CLIENT" \
  -d '{"numero": 103,"type":"simple","prix":800}' | jq .
---
📅 3. Gestion des Réservations
3.1 Lister toutes les réservations (Protégé - Manager/Receptionniste)
curl -s $BASE_URL/api/v1/reservations \
  -H "Authorization: Bearer $TOKEN_MGR" | jq .
3.2 Créer une réservation (Protégé - Tous rôles)
# Remplacer CHAMBRE_ID par l'ID d'une chambre existante
CHAMBRE_ID=$(curl -s $BASE_URL/api/v1/chambres \
  -H "Authorization: Bearer $TOKEN_CLIENT" | jq -r '.[0].id')
curl -s -X POST $BASE_URL/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_CLIENT" \
  -d "{
    \"user_id\": \"267508c1-5a0a-4f1f-976a-4dfd4ff0fd0a\",
    \"chambre_id\": \"$CHAMBRE_ID\",
    \"date_debut\": \"2026-05-10\",
    \"date_fin\": \"2026-05-15\",
    \"montant\": 7500
  }" | jq .
3.3 Check-in (Protégé - Manager/Receptionniste)
# Récupérer l'ID de la réservation
RES_ID=$(curl -s $BASE_URL/api/v1/reservations \
  -H "Authorization: Bearer $TOKEN_MGR" | jq -r '.[0].id')
curl -s -X PATCH $BASE_URL/api/v1/reservations/$RES_ID/checkin \
  -H "Authorization: Bearer $TOKEN_MGR" | jq .
3.4 Check-out (Protégé - Manager/Receptionniste)
curl -s -X PATCH $BASE_URL/api/v1/reservations/$RES_ID/checkout \
  -H "Authorization: Bearer $TOKEN_MGR" | jq .
---
🔍 Inspection de la Base de Données (PostgreSQL dans Docker)
Méthode 1: Connexion directe via psql (Depuis l'hôte)
# Se connecter à PostgreSQL via le port mappé (5433)
psql -h localhost -p 5433 -U postgres -d riad
# Une fois connecté, vous pouvez exécuter:
\dt                          # Lister toutes les tables
SELECT * FROM users;         # Voir tous les utilisateurs
SELECT * FROM chambres;      # Voir toutes les chambres
SELECT * FROM reservations;   # Voir toutes les réservations
SELECT * FROM taches;
SELECT * FROM services;
SELECT * FROM paiements;
# Pour voir la structure d'une table:
\d users
\d chambres
\d reservations
Méthode 2: Connexion via Docker (Depuis le conteneur)
# Se connecter au conteneur PostgreSQL
docker exec -it riad-postgres psql -U postgres -d riad
# Ou ouvrir un shell dans le conteneur:
docker exec -it riad-postgres sh
# Puis: psql -U postgres -d riad
Méthode 3: Requêtes rapides sans connexion interactive
# Voir le nombre d'utilisateurs
docker exec riad-postgres psql -U postgres -d riad -c "SELECT COUNT(*) FROM users;"
# Voir toutes les chambres
docker exec riad-postgres psql -U postgres -d riad -c "SELECT * FROM chambres;"
# Voir toutes les réservations avec leurs statuts
docker exec riad-postgres psql -U postgres -d riad -c "SELECT id, user_id, chambre_id, statut, date_debut, date_fin FROM reservations;"
# Voir les chambres et leur statut en temps réel
docker exec riad-postgres psql -U postgres -d riad -c "SELECT numero, type, statut, prix FROM chambres;"
Méthode 4: Surveillance en temps réel (Logs et activité)
# Voir les logs du serveur API en temps réel
docker compose logs server -f
# Voir les logs PostgreSQL
docker compose logs postgres -f
# Voir les conteneurs actifs
docker ps
# Statistiques d'utilisation des conteneurs
docker stats
---
📊 Script de Test Automatisé Complet (Optionnel)
Crée un fichier test_all.sh:
#!/bin/bash
BASE_URL="http://localhost:8081"
echo "=== TEST COMPLET DE L'API RIAD ==="
echo -e "\n1. Inscription..."
curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"testauto@test.com","password":"123456","nom":"Auto","prenom":"Test"}' | jq -r '.message'
echo -e "\n2. Connexion..."
TOKEN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testauto@test.com","password":"123456"}' | jq -r '.token')
echo "Token obtenu: ${TOKEN:0:20}..."
echo -e "\n3. Utilisateur actuel..."
curl -s $BASE_URL/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq '.nom'
echo -e "\n4. Test accès sans token (doit échouer)..."
curl -s $BASE_URL/api/v1/chambres | jq '.error'
echo -e "\n=== TESTS TERMINÉS ==="
---
🎯 Points de Vérification Importants
Test
Inscription
Connexion
Création chambre
Réservation
Check-in
Check-out
---
🚨 Dépannage Rapide
# Si l'API ne répond pas:
docker compose ps                    # Vérifier l'état des conteneurs
docker compose logs server --tail=50  # Voir les erreurs récentes
# Si la DB ne répond pas:
docker exec riad-postgres pg_isready -U postgres
# Pour tout redémarrer proprement:
docker compose down && docker compose up -d
Ce guide te permet de tester manuellement chaque endpoint et de voir exactement ce qui se passe dans la base de données PostgreSQL via Docker.

Le problème est que tu tapes les commandes sur plusieurs lignes sans continuation. Bash interprète -H et -d comme des commandes séparées.
🔧 Correction - 2 Méthodes
Méthode 1: Commande sur une seule ligne (Recommandé pour débuter)
curl -s -X POST http://localhost:8081/api/v1/auth/register -H "Content-Type: application/json" -d '{"email":"client5@test.com","password":"123456","nom":"sali","prenom":"faical"}'
Méthode 2: Avec continuation de ligne (backslash \ à la fin)
curl -s -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"client6@test.com","password":"123456","nom":"ali","prenom":"faal"}'
---
📋 Commandes Complètes à Executer (Copie-Colle)
# 1. Inscription client5
curl -s -X POST http://localhost:8081/api/v1/auth/register -H "Content-Type: application/json" -d '{"email":"client5@test.com","password":"123456","nom":"sali","prenom":"faical"}'
# 2. Login manager (sur une seule ligne)
curl -s -X POST http://localhost:8081/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"manager@test.com","password":"123456"}'
# 3. Récupérer le token manager (sur une seule ligne)
TOKEN_MGR=$(curl -s -X POST http://localhost:8081/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"manager@test.com","password":"123456"}' | jq -r '.token')
# 4. Créer une chambre (sur une seule ligne)
curl -s -X POST http://localhost:8081/api/v1/chambres -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN_MGR" -d '{"numero":103,"type":"simple","prix":1200}'
# 5. Voir toutes les chambres
curl -s http://localhost:8081/api/v1/chambres -H "Authorization: Bearer $TOKEN_MGR" | jq .
---
🔍 Vérifier la Base de Données
# Voir les utilisateurs
docker exec riad-postgres psql -U postgres -d riad -c "SELECT id, email, nom, role FROM users;"
# Voir les chambres
docker exec riad-postgres psql -U postgres -d riad -c "SELECT numero, type, statut, prix FROM chambres;"