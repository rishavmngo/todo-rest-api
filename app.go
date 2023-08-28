package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rishavmngo/todo-http/jwtUtil"
	"github.com/rishavmngo/todo-http/todo"
	"github.com/rishavmngo/todo-http/user"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initilize(dbname, port, password, user string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", user, password, dbname, port)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initilizeRouter()
}

func (a *App) Run(addr string) {

	log.Fatal(http.ListenAndServe(addr, a.Router))

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}
func (a *App) register(w http.ResponseWriter, r *http.Request) {

	var user user.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	// fmt.Println(user)
	if err := user.AddUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, user)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var user user.User

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	if err := user.FindUserByUsernameAndPassword(a.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token := jwtutil.GenerateToken(user.ID)
	respondWithJson(w, http.StatusOK, map[string]string{"token": token})
}

func (a *App) addTodo(w http.ResponseWriter, r *http.Request) {
	var todo todo.Todo

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&todo)

	if err != nil {
		log.Fatal(err)
	}

	if err := todo.AddTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, map[string]uint{"id": todo.ID})
}

func (a *App) removeTodo(w http.ResponseWriter, r *http.Request) {

	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])
	todo_id, _ := strconv.Atoi(mux.Vars(r)["todo_id"])

	var todo todo.Todo

	todo.AuthorID = uint(author_id)
	todo.ID = uint(todo_id)

	if err := todo.RemoveTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusOK, todo)
}
func (a *App) updateTodo(w http.ResponseWriter, r *http.Request) {

	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])
	todo_id, _ := strconv.Atoi(mux.Vars(r)["todo_id"])

	decoder := json.NewDecoder(r.Body)
	var todo todo.Todo
	todo.AuthorID = uint(author_id)
	todo.ID = uint(todo_id)

	if err := decoder.Decode(&todo); err != nil {
		log.Fatal(err)
	}

	if err := todo.UpdateTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, todo)
}
func (a *App) getAllTodosById(w http.ResponseWriter, r *http.Request) {
	author_id, _ := strconv.Atoi(mux.Vars(r)["author_id"])

	todos, err := todo.GetAllTodos(a.DB, uint(author_id))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, todos)
}
func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}
func (a *App) initilizeRouter() {
	//authentication
	a.Router.Use(logMW)
	userRoute := a.Router.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("/register", a.register).Methods("POST")
	userRoute.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/todo/add", jwtutil.Authenticate(a.addTodo)).Methods("POST")
	a.Router.HandleFunc("/todo/delete/{author_id}/{todo_id}", jwtutil.Authenticate(a.removeTodo)).Methods("DELETE")
	a.Router.HandleFunc("/todo/update/{author_id}/{todo_id}", jwtutil.Authenticate(a.updateTodo)).Methods("PUT")
	a.Router.HandleFunc("/todo/getall/{author_id}", jwtutil.Authenticate(a.getAllTodosById)).Methods("GET")
}
