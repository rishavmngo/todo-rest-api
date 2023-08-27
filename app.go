package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rishavmngo/todo-http/todo"
	"github.com/rishavmngo/todo-http/user"
	"log"
	"net/http"
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

	respondWithJson(w, http.StatusOK, map[string]uint{"id": user.ID})
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

func (a *App) initilizeRouter() {
	//authentication
	a.Router.HandleFunc("/register", a.register).Methods("POST")
	a.Router.HandleFunc("/login", a.login).Methods("POST")
	a.Router.HandleFunc("/todo/add", a.addTodo).Methods("POST")
}
