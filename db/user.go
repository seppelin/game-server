package db

import (
	"game-server/models"
	"log"
)

func initUser() {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		pwdhash TEXT NOT NULL
	);`
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("User table created successfully.")
}

func UpdateUser(user models.User) {
	// Prepare the SQL statement
	sql := `UPDATE users SET name = ?, email = ?, pwdhash = ? WHERE id = ?`
	_, err := db.Exec(sql, user.Name, user.Email, user.PwdHash, user.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertUser(name, email, pwdhash string) (models.UserID, error) {
	sql := `INSERT INTO users(name, email, pwdhash) VALUES (?, ?, ?)`
	result, err := db.Exec(sql, name, email, pwdhash)
	if err != nil {
		return 0, err
	}
	log.Println("Inserted user successfully.")
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Insert get id:", err)
	}
	return models.UserID(id), nil
}

func User(id models.UserID) (models.User, error) {
	sql := `SELECT id, name, email, pwdhash FROM users WHERE id = ?`
	var user models.User
	err := db.QueryRow(sql, id).Scan(&user.ID, &user.Name, &user.Email, &user.PwdHash)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserByName(name string) (models.User, error) {
	sql := `SELECT id, name, email, pwdhash FROM users WHERE name = ?`
	var user models.User
	err := db.QueryRow(sql, name).Scan(&user.ID, &user.Name, &user.Email, &user.PwdHash)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UserFromSession(sessionName string) (models.User, bool) {
	userID, err := UserIDByName(sessionName)
	if err != nil {
		return models.User{}, false
	}
	user, err := User(userID)
	if err != nil {
		return models.User{}, false
	}
	return user, true
}
