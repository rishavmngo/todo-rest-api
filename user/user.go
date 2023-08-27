package user

import "database/sql"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) AddUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users(username,password) VALUES($1, $2) returning id", u.Username, u.Password).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) FindUserByUsernameAndPassword(db *sql.DB) error {
	err := db.QueryRow("SELECT id FROM users WHERE username=$1 and password=$2", u.Username, u.Password).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}
