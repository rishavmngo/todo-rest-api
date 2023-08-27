package todo

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID        uint
	Title     string
	Status    bool
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorID  uint      `json:"author_id"`
}

func (t *Todo) AddCreatedAt() {
	t.CreatedAt = time.Now()
}

func (t *Todo) AddModifyAt() {
	t.UpdatedAt = time.Now()
}

func (t *Todo) AddTodo(db *sql.DB) error {
	t.AddCreatedAt()
	t.AddModifyAt()
	err := db.QueryRow("INSERT INTO todos(title, status,created_at,updated_at,author_id) values($1, $2, $3, $4, $5) returning id", t.Title, t.Status, t.CreatedAt, t.UpdatedAt, t.AuthorID).Scan(&t.ID)

	if err != nil {
		return err
	}
	return nil
}
