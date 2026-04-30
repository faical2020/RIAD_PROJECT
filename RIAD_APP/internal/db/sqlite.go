package db

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite", filepath)
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
		synced BOOLEAN DEFAULT 1
	);
	CREATE TABLE IF NOT EXISTS reservations (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		chambre_id TEXT,
		date_debut TEXT,
		date_fin TEXT,
		montant REAL,
		statut TEXT,
		synced BOOLEAN DEFAULT 0
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
		var id, roomType, desc, equip, status string
		var num int
		var price float64
		var synced bool
		if err := rows.Scan(&id, &num, &roomType, &price, &desc, &equip, &status, &synced); err != nil {
			return nil, err
		}
		rooms = append(rooms, map[string]interface{}{
			"id": id, "numero": num, "type": roomType, "prix": price, "description": desc, "equipements": equip, "statut": status, "synced": synced,
		})
	}
	return rooms, nil
}

func SaveRoom(id string, num int, roomType string, price float64, desc, equip, status string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO rooms (id, numero, type, prix, description, equipements, statut, synced) VALUES (?, ?, ?, ?, ?, ?, ?, 0)", id, num, roomType, price, desc, equip, status)
	return err
}

func SaveReservation(id, userId, roomId, start, end string, amount float64, status string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO reservations (id, user_id, chambre_id, date_debut, date_fin, montant, statut, synced) VALUES (?, ?, ?, ?, ?, ?, ?, 0)", id, userId, roomId, start, end, amount, status)
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
