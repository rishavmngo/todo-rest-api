package main

import (
	"database/sql"
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

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}
func (a *App) initilizeRouter() {
	a.Router.Use(logMW)
	a.handleRequest("/user", user.RouteHandler)
	a.handleRequest("/todo", todo.RouteHandler)

}

type RequestHandlerFunction func(db *sql.DB, router *mux.Router)

func (a *App) handleRequest(path string, init RequestHandlerFunction) {
	route := a.Router.PathPrefix(path).Subrouter()
	init(a.DB, route)
}
