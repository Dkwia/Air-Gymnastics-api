package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var Users = []User{
	{ID: "1", Username: "user1", Email: "user1@example.com"},
	{ID: "2", Username: "user2", Email: "user2@example.com"},
}
