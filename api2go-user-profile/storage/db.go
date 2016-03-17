package storage

import (
	"fmt"
	"os"
	"database/sql"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
    "api-sandbox/api2go-user-profile/model"
)

// InitDB connects to postgres database "test"
func InitDb() (*gorp.DbMap, error) {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    // open connection to postgres
    db, err := sql.Open("postgres", dbinfo)
    
    if err != nil {
		return nil, err
	} 
    // construct a gorp DbMap using PostgresDialect
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    dbmap.AddTableWithName(model.User{}, "users").SetKeys(true, "user_id")
    // dbmap.AddTableWithName(model.User{}, "users").SetKeys(true, "user_id")
    return dbmap, nil
}