package main

import (
	"fmt"
	"os"
	"database/sql"
	// "gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

func main() {
	dataUnnest :=  selectAsString()
	fmt.Println("SELECT as String:", dataUnnest)
}
//var dbmap = initDb()
func selectAsString() []string {
	db := initDb()
	results := make([]string, 0)
	rows, _ := db.Query("select array_to_string(array_agg(profile_id),',') AS profile_ids from profiles group by user_id")
	var scanString string
	for rows.Next() {
		rows.Scan(&scanString)
		results = append(results, scanString)
	}
	return results
	// return nil
	//select array_to_string(array_agg(profile_id), ',') AS profile_ids from profiles group by user_id
}

//InitDB connects to postgres database "test"
func initDb() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    // open connection to postgres
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
		panic(err)
	}
	return db
    // construct a gorp DbMap using PostgresDialect
    // dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    // return dbmap, nil
}