package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func GetDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "riad.db"
	}
	appDir := filepath.Join(configDir, "riad_app")
	return filepath.Join(appDir, "riad.db")
}

func InitDB(path string) error {
	var err error
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erreur creation dossier db: %v", err)
		}
	}
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("erreur ouverture db: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		nom TEXT,
		prenom TEXT,
		role TEXT,
		synced BOOLEAN DEFAULT 1
	);
	CREATE TABLE IF NOT EXISTS rooms (
		id TEXT PRIMARY KEY,
		numero INTEGER,
		type TEXT,
		prix REAL,
		description TEXT,
		equipements TEXT,
		statut TEXT,
		cleaning_status TEXT DEFAULT 'propre',
		synced BOOLEAN DEFAULT 1,
		created_at INTEGER,
		updated_at INTEGER
	);
	CREATE TABLE IF NOT EXISTS reservations (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		chambre_id TEXT,
		date_debut TEXT,
		date_fin TEXT,
		montant REAL,
		statut TEXT,
		synced BOOLEAN DEFAULT 0,
		created_at INTEGER,
		updated_at INTEGER
	);
	`
	_, err = DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("erreur schema: %v", err)
	}
	return nil
}

func GetRooms() ([]map[string]interface{}, error) {
	rows, err := DB.Query("SELECT * FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []map[string]interface{}
	for rows.Next() {
		var id, roomType, desc, equip, status, cleaningStatus string
		var num int
		var price float64
		var synced bool
		var createdAt, updatedAt int64
		if err := rows.Scan(&id, &num, &roomType, &price, &desc, &equip, &status, &cleaningStatus, &synced, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, map[string]interface{}{
			"id": id, "numero": num, "type": roomType, "prix": price, "description": desc, "equipements": equip, "statut": status, "cleaning_status": cleaningStatus, "synced": synced, "created_at": createdAt, "updated_at": updatedAt,
		})
	}
	return rooms, nil
}

func GetReservations() ([]map[string]interface{}, error) {
	rows, err := DB.Query("SELECT * FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			m[colName] = columns[i]
		}
		results = append(results, m)
	}
	return results, nil
}

func SaveRoom(id string, num int, roomType string, price float64, desc, equip, status, cleaningStatus string) error {
	now := time.Now().Unix()
	_, err := DB.Exec("INSERT OR REPLACE INTO rooms (id, numero, type, prix, description, equipements, statut, cleaning_status, synced, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", id, num, roomType, price, desc, equip, status, cleaningStatus, 0, now, now)
	return err
}

func SaveReservation(id, userId, roomId, start, end string, amount float64, status string) error {
	now := time.Now().Unix()
	_, err := DB.Exec("INSERT OR REPLACE INTO reservations (id, user_id, chambre_id, date_debut, date_fin, montant, statut, synced, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", id, userId, roomId, start, end, amount, status, 0, now, now)
	return err
}

func MarkSynced(table, id string) error {
	query := fmt.Sprintf("UPDATE %s SET synced = 1 WHERE id = ?", table)
	_, err := DB.Exec(query, id)
	return err
}

func GetUnsynced(table string) ([]map[string]interface{}, error) {
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE synced = 0", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			m[colName] = columns[i]
		}
		results = append(results, m)
	}
	return results, nil
}
