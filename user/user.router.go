package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "github.com/rishavmngo/todo-http/jsonUtils"
	jwtutil "github.com/rishavmngo/todo-http/jwtUtil"
)

func RouteHandler(db *sql.DB, r *mux.Router) {
	var user User

	r.HandleFunc("/register", controller(db, user.register)).Methods("POST")
	r.HandleFunc("/login", controller(db, user.login)).Methods("POST")
}

type controllerHandler func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func controller(db *sql.DB, handler controllerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}

}

func (user *User) register(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	// fmt.Println(user)
	if err := user.AddUser(db); err != nil {
		jsonutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonutils.RespondWithJson(w, http.StatusCreated, user)
}

func (user *User) login(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	if err := user.FindUserByUsernameAndPassword(db); err != nil {
		jsonutils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token := jwtutil.GenerateToken(user.ID)
	jsonutils.RespondWithJson(w, http.StatusOK, map[string]string{"token": token})
}
