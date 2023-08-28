package todo

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	jsonutils "github.com/rishavmngo/todo-http/jsonUtils"
	jwtutil "github.com/rishavmngo/todo-http/jwtUtil"
	"log"
	"net/http"
	"strconv"
)

func RouteHandler(db *sql.DB, r *mux.Router) {
	var todo Todo

	r.HandleFunc("/add", jwtutil.Authenticate(controller(db, todo.addTodo))).Methods("POST")
	r.HandleFunc("/delete/{author_id}/{todo_id}", jwtutil.Authenticate(controller(db, todo.removeTodo))).Methods("DELETE")
	r.HandleFunc("/update/{author_id}/{todo_id}", jwtutil.Authenticate(controller(db, todo.updateTodo))).Methods("PUT")
	r.HandleFunc("/getall/{author_id}", jwtutil.Authenticate(controller(db, todo.getAllTodosById))).Methods("GET")
}

type controllerHandler func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func controller(db *sql.DB, handler controllerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}

}

func (todo *Todo) addTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&todo)

	if err != nil {
		log.Fatal(err)
	}

	if err := todo.AddTodo(db); err != nil {
		jsonutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonutils.RespondWithJson(w, http.StatusCreated, map[string]uint{"id": todo.ID})
}

func (todo *Todo) removeTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])
	todo_id, _ := strconv.Atoi(mux.Vars(r)["todo_id"])

	todo.AuthorID = uint(author_id)
	todo.ID = uint(todo_id)

	if err := todo.RemoveTodo(db); err != nil {
		jsonutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	jsonutils.RespondWithJson(w, http.StatusOK, todo)
}
func (todo *Todo) updateTodo(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])
	todo_id, _ := strconv.Atoi(mux.Vars(r)["todo_id"])

	decoder := json.NewDecoder(r.Body)
	todo.AuthorID = uint(author_id)
	todo.ID = uint(todo_id)

	if err := decoder.Decode(&todo); err != nil {
		log.Fatal(err)
	}

	if err := todo.UpdateTodo(db); err != nil {
		jsonutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonutils.RespondWithJson(w, http.StatusOK, todo)
}
func (todo *Todo) getAllTodosById(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])

	todos, err := GetAllTodos(db, uint(author_id))
	if err != nil {
		jsonutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonutils.RespondWithJson(w, http.StatusOK, todos)
}
