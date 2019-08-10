package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq" // postgres
	_ "github.com/mattn/go-sqlite3"
)

const APIPort = ":8080"
const SqliteDBName = "./hemato.db" // sqlite
const DRIVER_NAME = "postgres"
const DB_USER = "postgres"
const DB_PASSWORD = "12345"
const POSTGRES_DB_NAME = "hemato" // postgres
const SSLMODE = "disable"
const HOST = "172.17.0.2"
const PORT = "5432"

func main() {

	a := App{}
	a.Initialize(POSTGRES_DB_NAME)
	a.Run(APIPort)

}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbName string) {
	// a.DB = ConnectSqliteDB(dbName)
	a.DB = ConnectPostgresDB(dbName)
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

func ConnectSqliteDB(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("não pingou")
	} else {
		log.Println("pingou")
	}
	return db
}

func ConnectPostgresDB(dbName string) *sql.DB {
	// var err error
	connString := fmt.Sprintf(fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		DB_USER, DB_PASSWORD, HOST, PORT, dbName, SSLMODE,
	))

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("não pingou")
	} else {
		log.Println("pingou")
	}
	return db
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Filter struct {
	Uf  string `schema:"uf"`
	CID string `schema:"cid"`
	Ano string `schema:"ano"`
}

func validateFilter(params map[string][]string) (string, error) {
	var filterString string
	filter := Filter{}
	if err := schema.NewDecoder().Decode(&filter, params); err != nil {
		return filterString, err
	}
	if filter.Uf == "Brasil" {
		filter.Uf = ""
	}
	if filter.Uf != "" {
		filterString += fmt.Sprintf(` AND "UF_ZI" = '%s'`, filter.Uf)
	}

	if filter.CID == "Todos" {
		filter.CID = ""
	}
	if filter.CID != "" {
		filterString += fmt.Sprintf(` AND "DESCR_CID" = '%s'`, filter.CID)
	}

	if filter.Ano == "Todos" {
		filter.Ano = ""
	}
	if filter.Ano != "" {
		filterString += fmt.Sprintf(` AND "ANO_CMPT" = '%s'`, filter.Ano)
	}

	filterString += ` GROUP BY 1 ORDER by 1`
	return filterString, nil
}

func GetResponse(db *sql.DB, params map[string][]string) (interface{}, error) {
	filter, err := validateFilter(params)
	if err != nil {
		return nil, err
	}

	linePlots, err := GetData(db, filter)
	if err != nil {
		return nil, err
	}

	ageHistData, err := GetAgeHistogram(db, filter)
	if err != nil {
		return nil, err
	}

	USSHistData, err := GetAIHHistogram(db, filter)
	if err != nil {
		return nil, err
	}

	options, err := GetOptions(db, filter)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		map[string]interface{}{
			"lineplots": linePlots,
			"age_hist":  ageHistData,
			"USS_hist":  USSHistData,
			"options":   options,
		},
	}, nil
}

type ResponseRow struct {
	Year       int     `json:"year"`
	Admissions int     `json:"admissions"`
	Mortality  float64 `json:"mortality"`
}

func GetData(db *sql.DB, filter string) (interface{}, error) {

	query := `SELECT "ANO_CMPT" as Ano, COUNT ("N_AIH") as Internações, round (AVG ("MORTE"), 2) as Mortalidade
				FROM admissions
				WHERE 1=1`
	query += filter

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	res := []ResponseRow{}
	for rows.Next() {
		r := ResponseRow{}
		rows.Scan(&r.Year, &r.Admissions, &r.Mortality)
		res = append(res, r)
	}
	return res, nil
}

type HistData struct {
	Bin   int `json:"bin"`
	Count int `json:"count"`
}

func GetAgeHistogram(db *sql.DB, filter string) (interface{}, error) {
	query := `SELECT cast("IDADE"/10 AS INT) * 10 as Idades, count(*) as Internações
				FROM admissions
				WHERE 1=1`
	query += filter

	return GetHistogram(db, query)
}

func GetAIHHistogram(db *sql.DB, filter string) (interface{}, error) {
	query := `SELECT cast("US_TOT"/500 AS INT) * 500 as US$, count("US_TOT") as Internações
				FROM admissions
				WHERE 1=1`
	query += filter
	return GetHistogram(db, query)
}

func GetHistogram(db *sql.DB, query string) (interface{}, error) {

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	res := []HistData{}
	for rows.Next() {
		a := HistData{}
		rows.Scan(&a.Bin, &a.Count)
		res = append(res, a)
	}
	return res, nil
}

func GetOptions(db *sql.DB, filter string) (interface{}, error) {
	queries := [][]string{
		{"ufs", `SELECT DISTINCT "UF_ZI" as ufs FROM admissions WHERE 1=1`},
		{"years", `SELECT DISTINCT "ANO_CMPT" as anos FROM admissions WHERE 1=1`},
		{"cids", `SELECT DISTINCT "DESCR_CID" as cids FROM admissions WHERE 1=1`},
	}

	options := make(map[string][]string)
	for _, query := range queries {
		query[1] += ` GROUP BY 1 ORDER by 1`
		rows, err := db.Query(query[1])
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var option string
			err = rows.Scan(&option)
			if err != nil {
				return nil, err
			}
			options[query[0]] = append(options[query[0]], option)
		}
	}
	return options, nil
}
