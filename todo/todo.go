package todo

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Status    bool      `json:"status"`
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

func (t *Todo) RemoveTodo(db *sql.DB) error {
	err := db.QueryRow("DELETE FROM todos WHERE id=$1 and author_id=$2 returning title,status,created_at,updated_at", t.ID, t.AuthorID).Scan(&t.Title, &t.Status, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}
func GetAllTodos(db *sql.DB, AuthorId uint) ([]Todo, error) {
	rows, err := db.Query("Select id,title,status,created_at,updated_at,author_id from todos where author_id=$1", AuthorId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt, &todo.AuthorID); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (t *Todo) UpdateTodo(db *sql.DB) error {
	t.AddModifyAt()
	err := db.QueryRow("UPDATE todos SET title = $1, status = $2, updated_at = $5 where id = $3 and author_id = $4 returning created_at", t.Title, t.Status, t.ID, t.AuthorID, t.UpdatedAt).Scan(&t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
