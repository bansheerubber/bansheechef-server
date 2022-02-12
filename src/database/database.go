package database

import (
	"database/sql"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Open() {
	db, dbError := sql.Open("sqlite3", "/home/me/.config/bansheechef/storage/data.db")
	database = db

	if dbError != nil {
		log.Fatal(dbError)
	}

	dbError = database.Ping()
	if dbError != nil {
		log.Fatal(dbError)
	}

	log.Println("Opened SQLITE database")
}

func Close(args ...interface{}) {
	database.Close()

	log.Println("Closed SQLITE database")
}

func Exec(query string, args ...interface{}) {
	transaction, err := database.Begin()
	statement, err := transaction.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec(args)
	if err != nil {
		log.Fatal(err)
	}

	err = transaction.Commit()
	err = statement.Close()
}

func Query(query string, args []interface{}, resultType reflect.Type) []interface{} {
	statement, err := database.Prepare(query)
	defer statement.Close()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := statement.Query(args...)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var results []interface{}
	for rows.Next() {
		result := reflect.New(resultType) // create the type based on the result type supplied in the parameters

		columns := make([]interface{}, result.Elem().NumField()) // build a list of destinations to scan into using reflect
		for i := 0; i < result.Elem().NumField(); i++ {
			columns[i] = result.Elem().Field(i).Addr().Interface()
		}

		rows.Scan(columns...)

		results = append(results, result.Interface())
	}

	return results
}

func QueryOne(query string, args []interface{}, resultType reflect.Type) interface{} {
	statement, err := database.Prepare(query)
	defer statement.Close()
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(args...)
	if row == nil {
		return nil
	}

	result := reflect.New(resultType) // create the type based on the result type supplied in the parameters

	columns := make([]interface{}, result.Elem().NumField()) // build a list of destinations to scan into using reflect
	for i := 0; i < result.Elem().NumField(); i++ {
		columns[i] = result.Elem().Field(i).Addr().Interface()
	}

	row.Scan(columns...)
	return result.Interface()
}
