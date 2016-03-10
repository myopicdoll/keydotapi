package main

import (
	"fmt"
	"os"
	"database/sql"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

func main() {
	var _, err = initDb()
	fmt.Println(err)
}
//var dbmap = initDb()

//InitDB connects to postgres database "test"
func initDb() (*gorp.DbMap, error) {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    fmt.Println(dbinfo)
    // open connection to postgres
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
		return nil, err
	}
    // construct a gorp DbMap using PostgresDialect
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    return dbmap, nil
}