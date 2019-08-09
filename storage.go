package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/schema"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(dbName string) *sql.DB {
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
		filterString += fmt.Sprintf(` AND UF_ZI = '%s'`, filter.Uf)
	}

	if filter.CID == "Todos" {
		filter.CID = ""
	}
	if filter.CID != "" {
		filterString += fmt.Sprintf(` AND DESCR_CID = '%s'`, filter.CID)
	}

	if filter.Ano == "Todos" {
		filter.Ano = ""
	}
	if filter.Ano != "" {
		filterString += fmt.Sprintf(` AND ANO_CMPT = '%s'`, filter.Ano)
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
			"options": options,
		},
	}, nil
}

type ResponseRow struct {
	Year       int     `json:"year"`
	Admissions int     `json:"admissions"`
	Mortality  float64 `json:"mortality"`
}

func GetData(db *sql.DB, filter string) (interface{}, error) {

	query := `SELECT ano_cmpt as Ano, COUNT (n_aih) as Internações, round (AVG (morte), 2) as Mortalidade
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
	query := `SELECT cast(idade/10 AS INT) * 10 as Idades, count(*) as Internações
				FROM admissions
				WHERE 1=1`
	query += filter

	return GetHistogram(db, query)
}

func GetAIHHistogram(db *sql.DB, filter string) (interface{}, error) {
	query := `SELECT cast(us_tot/500 AS INT) * 500 as US$, count(us_tot) as Internações
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
		{"ufs", `SELECT DISTINCT uf_zi as ufs FROM 'admissions' WHERE 1=1`},
		{"years", `SELECT DISTINCT ano_cmpt as anos FROM 'admissions' WHERE 1=1`},
		{"cids", `SELECT DISTINCT descr_cid as cids FROM 'admissions' WHERE 1=1`},
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
