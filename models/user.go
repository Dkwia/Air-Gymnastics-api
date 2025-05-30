package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type WhatsAppInfo struct {
	Phone     string `json:"phone"`
	OptIn     bool   `json:"opt_in"`
	LastMsgID string `json:"last_msg_id,omitempty"`
}

type User struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	WhatsApp WhatsAppInfo `json:"whatsapp"`
}

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "file:data/users.db?_fk=true")
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			whatsapp_phone TEXT NOT NULL,
			whatsapp_opt_in BOOLEAN NOT NULL,
			whatsapp_last_msg_id TEXT
		)
	`)
	return err
}

func CreateUser(u User) error {
	_, err := DB.Exec(
		`INSERT INTO users 
		(id, username, whatsapp_phone, whatsapp_opt_in, whatsapp_last_msg_id) 
		VALUES (?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.WhatsApp.Phone, u.WhatsApp.OptIn, u.WhatsApp.LastMsgID,
	)
	return err
}

func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.WhatsApp.Phone,
			&u.WhatsApp.OptIn,
			&u.WhatsApp.LastMsgID,
		)
		if err != nil {
			log.Println("DB scan error:", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}
