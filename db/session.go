package db

import (
	"game-server/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func initSession() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        user_id INTEGER NOT NULL,
		create_address TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}
	log.Println("Session table created")

	// Create a trigger to update accessed_at on row update
	createTriggerSQL := `
    CREATE TRIGGER IF NOT EXISTS update_accessed_at
    AFTER UPDATE ON sessions
    FOR EACH ROW
    BEGIN
        UPDATE sessions SET accessed_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;`
	if _, err := db.Exec(createTriggerSQL); err != nil {
		log.Fatal(err)
	}
}

func SetAccessedAt(id models.SessionID) error {
	sql := `UPDATE sessions SET accessed_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(sql, id)
	return err
}

func SetSession(userID models.UserID, createAddress string, expiresAt time.Time) (models.SessionID, string, error) {
	name := uuid.New().String()
	sql := `INSERT INTO sessions (name, user_id, create_address, expires_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(sql, name, userID, createAddress, expiresAt.UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return models.SessionID(id), name, err
}

func SessionByName(name string) (models.Session, error) {
	sql := `SELECT id, name, user_id, create_address, created_at, accessed_at, expires_at FROM sessions WHERE name = ? AND expires_at > ?`
	var session models.Session
	err := db.QueryRow(sql, name, time.Now()).Scan(&session.Id, &session.Name, &session.UserID, &session.CreateAddress, &session.CreatedAt, &session.AccessedAt, &session.ExpiresAt)
	return session, err
}

func Session(id models.SessionID) (models.Session, error) {
	sql := `SELECT id, name, user_id, create_address, created_at, accessed_at, expires_at FROM sessions WHERE id = ? AND expires_at > ?`
	var session models.Session
	err := db.QueryRow(sql, id, time.Now()).Scan(&session.Id, &session.Name, &session.UserID, &session.CreateAddress, &session.CreatedAt, &session.AccessedAt, &session.ExpiresAt)
	return session, err
}

func UserIDByName(name string) (models.UserID, error) {
	var userID models.UserID
	sql := `SELECT user_id FROM sessions WHERE name = ? AND expires_at > ?`
	err := db.QueryRow(sql, name, time.Now()).Scan(&userID)
	return userID, err
}

func ExpireSession(name string) error {
	sql := `UPDATE sessions SET expires_at = CURRENT_TIMESTAMP WHERE name = ?`
	_, err := db.Exec(sql, name)
	return err
}
