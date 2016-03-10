package main

import (
    "fmt"
    "time"
	"database/sql"
    // "strconv"
    // "log"
    // "net/http"
    "gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
    "github.com/gin-gonic/gin"
)

const (
	DB_USER = "ion"
	DB_PASSWORD = "pass"
	DB_NAME = "test"
)

type User struct {
    Id string  `db:"uid" json:"id"`
    Username string `db:"username" json:"username"`
    Department string `db:"departname" json:"department"`
    Created  time.Time `db:"created" json:"created"`
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        DB_USER, DB_PASSWORD, DB_NAME)
    // open connection to postgres
    db, err := sql.Open("postgres", dbinfo)
    checkErr(err)
    // construct a gorp DbMap using PostgresDialect
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    return dbmap
}

func main() {
    port := "8080"
    // default gin router
    r := gin.Default()
    //r.Use(Cors())
    v1 := r.Group("api/v1")
    {
        v1.GET("/users", GetUsers)
        v1.GET("/users/:uid", GetUser)
        v1.POST("/users/", PostUser)
    }
    r.Run(":" + port)
}

func GetUser(c *gin.Context) {
    uid := c.Params.ByName("uid")
    var user User
    err := dbmap.SelectOne(&user, "select * from user where uid = $1 limit 1", uid)
    if err == nil {
        info := &User{
            Id: user.Id,
            Username: user.Username,
            Department: user.Department,
            Created: user.Created,
        }
        c.JSON(200, info)
    } else {
        c.JSON(404, gin.H{"error": "user not found"})
    }
}

func GetUsers(c *gin.Context) {
    var users []User
    _, err := dbmap.Select(&users, "select * from user")
    fmt.Println(users)
    if err == nil {
        c.JSON(200, users)
    } else {
        c.JSON(404, gin.H{"error": "no users exist in database"})
    }
}

func PostUser(c *gin.Context) {
    var user User
    t := time.Now()
    c.Bind(&user)
    if user.Username != "" && user.Department != "" {
        if insert, _ := dbmap.Exec(`insert into user (username, departname, created) values ($1, $2, $3)`, user.Username, user.Department, t); insert != nil {
         info := &User{
            Username: user.Username,
            Department: user.Department,
            Created: t,
            }
            c.JSON(201, info)
        } else { 
            fmt.Println("insert error")
        }
     } else {
        c.JSON(400, gin.H{"error": "Fields are empty"})
     }
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}