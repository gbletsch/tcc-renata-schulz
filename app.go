package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbName string) {
	a.DB = ConnectDB(dbName)
	a.Router = mux.NewRouter().StrictSlash(true)
	a.InitializeRoutes()
}

func (a *App) Run(addr string) {
	fmt.Println("API running on port", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) InitializeRoutes() {
	a.Router.Handle("/", a.getRoot()).Methods(http.MethodGet, http.MethodOptions)
	a.Router.Handle("/index", a.getIndex()).Methods(http.MethodGet, http.MethodOptions)
}

func (a *App) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Redirect(w, r, "/index", http.StatusPermanentRedirect)
	}
}

func (a *App) getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := r.ParseForm()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}
		resp, err := GetResponse(a.DB, r.Form)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, resp)
	}
}

func respondWithError(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.WriteHeader(status)
	w.Write(response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
