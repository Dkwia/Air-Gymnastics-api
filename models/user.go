package models

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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
	Role     string       `json:"role"`
	WhatsApp WhatsAppInfo `json:"whatsapp"`
}

type News struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Competition struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type Schedule struct {
	Event string `json:"event"`
	Time  string `json:"time"`
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
            role TEXT NOT NULL DEFAULT 'user',
            whatsapp_phone TEXT NOT NULL UNIQUE,
            whatsapp_opt_in BOOLEAN NOT NULL,
            whatsapp_last_msg_id TEXT
        )
    `)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS news (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL UNIQUE,
            content TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS competitions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            date TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS schedule (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            event TEXT NOT NULL UNIQUE,
            time TEXT NOT NULL
        )
    `)
	return err
}

func CreateUser(u User) error {
	u.ID = uuid.New().String()

	exists, err := isPhoneNumberUnique(u.WhatsApp.Phone)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("phone number %s is already registered", u.WhatsApp.Phone)
	}

	_, err = DB.Exec(
		`INSERT INTO users 
        (id, username, role, whatsapp_phone, whatsapp_opt_in, whatsapp_last_msg_id) 
        VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Role, u.WhatsApp.Phone, u.WhatsApp.OptIn, u.WhatsApp.LastMsgID,
	)
	return err
}

func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT id, username, role, whatsapp_phone, whatsapp_opt_in, whatsapp_last_msg_id FROM users")
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
			&u.Role,
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

func isPhoneNumberUnique(phone string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE whatsapp_phone = ?", phone).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func UpdateUser(userid string, updatedUser User) error {
	existingUser, err := GetUserByID(userid)
	if err != nil {
		return err
	}
	updatedUser.WhatsApp.Phone = existingUser.WhatsApp.Phone

	_, err = DB.Exec(`
        UPDATE users 
        SET username = ?, role = ?, whatsapp_opt_in = ?, whatsapp_last_msg_id = ?
        WHERE id = ?`,
		updatedUser.Username, updatedUser.Role, updatedUser.WhatsApp.OptIn, updatedUser.WhatsApp.LastMsgID, userid,
	)
	return err
}

func GetUserByID(userid string) (User, error) {
	var user User
	err := DB.QueryRow(`
        SELECT id, username, role, whatsapp_phone, whatsapp_opt_in, whatsapp_last_msg_id
        FROM users WHERE id = ?`, userid).
		Scan(&user.ID, &user.Username, &user.Role, &user.WhatsApp.Phone, &user.WhatsApp.OptIn, &user.WhatsApp.LastMsgID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}

func UpdateNews(news News) error {
	_, err := DB.Exec(`
        INSERT INTO news (title, content) VALUES (?, ?)
        ON CONFLICT(title) DO UPDATE SET content = ?`,
		news.Title, news.Content, news.Content)
	return err
}

func GetNews() (News, error) {
	var news News
	err := DB.QueryRow(`SELECT title, content FROM news ORDER BY id DESC LIMIT 1`).
		Scan(&news.Title, &news.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return news, fmt.Errorf("no news found")
		}
		return news, err
	}
	return news, nil
}

func UpdateCompetition(competition Competition) error {
	_, err := DB.Exec(`
        INSERT INTO competitions (name, date) VALUES (?, ?)
        ON CONFLICT(name) DO UPDATE SET date = ?`,
		competition.Name, competition.Date, competition.Date)
	return err
}

func GetCompetition() (Competition, error) {
	var competition Competition
	err := DB.QueryRow(`SELECT name, date FROM competitions ORDER BY id DESC LIMIT 1`).
		Scan(&competition.Name, &competition.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return competition, fmt.Errorf("no competition found")
		}
		return competition, err
	}
	return competition, nil
}

func UpdateSchedule(schedule Schedule) error {
	_, err := DB.Exec(`
        INSERT INTO schedule (event, time) VALUES (?, ?)
        ON CONFLICT(event) DO UPDATE SET time = ?`,
		schedule.Event, schedule.Time, schedule.Time)
	return err
}

func GetSchedule() (Schedule, error) {
	var schedule Schedule
	err := DB.QueryRow(`SELECT event, time FROM schedule ORDER BY id DESC LIMIT 1`).
		Scan(&schedule.Event, &schedule.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return schedule, fmt.Errorf("no schedule found")
		}
		return schedule, err
	}
	return schedule, nil
}
