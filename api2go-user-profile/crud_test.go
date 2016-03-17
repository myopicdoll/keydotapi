package main_test

import (
	"os"
	"fmt"
	"testing"
	"net/http"
	//"reflect"
 	"net/http/httptest"
 	// "strings"

 	"api-sandbox/api2go-user-profile/model"
 	"api-sandbox/api2go-user-profile/storage"
 	"api-sandbox/api2go-user-profile/resource"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
 	"github.com/manyminds/api2go"
 	// "github.com/manyminds/api2go-adapter/gingonic"
 	// "github.com/gin-gonic/gin"
)

var rec *httptest.ResponseRecorder
var dbmap *gorp.DbMap
var api  *api2go.API

func Crasher(details string) {
    fmt.Println("Going down in flames!", details)
    os.Exit(1)
}

func setup() {
	//r := gin.Default()
	// api = api2go.NewAPIWithRouting(
	// 	"v0",
	// 	api2go.NewStaticResolver("/"),
	// 	gingonic.New(r),
	// )
	api = api2go.NewAPIWithBaseURL("v0", "http://localhost:31415")
	var err error
	dbmap, err = storage.InitDb()
	if err != nil {
		Crasher("Database initialization error!")
	}
	userStorage := storage.NewUserStorage(dbmap)
	profStorage := storage.NewProfileStorage(dbmap)
	api.AddResource(model.User{}, resource.UserResource{ProfStorage: profStorage, UserStorage: userStorage})
	api.AddResource(model.Profile{}, resource.ProfileResource{ProfStorage: profStorage, UserStorage: userStorage})
	rec = httptest.NewRecorder()
}

func teardown() {
	fmt.Println("tearing DOWN!")
	defer dbmap.Db.Close()
}
// ~...~...~...~...~...~...~...~...~ //
// CREATE (POST)
func TestCreateUser(t *testing.T) {
	t.Log("testing user creation")
	var createUser = func() {
		rec = httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v0/users", strings.NewReader(`
		{
			"data": {
				"type": "users",
				"attributes": {
					"username": "Holygarian",
					"password": "pass6"
				}
			}
		}
		`))
		// Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusCreated))
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"meta": {
				"author": "The api2go examples crew",
				"license": "wtfpl",
				"license-url": "http://www.wtfpl.net"
			},
			"data": {
				"id": "1",
				"type": "users",
				"attributes": {
					"user-name": "marvin"
				},
				"relationships": {
					"sweets": {
						"data": [],
						"links": {
							"related": "http://localhost:31415/v0/users/1/sweets",
							"self": "http://localhost:31415/v0/users/1/relationships/sweets"
						}
					}
				}
			}
		}
		`))
	}
	
}

func TestCreateProfile(t *testing.T) {
	t.Log("testing profile creation")
}
// READ (GET)
func TestGetUsers(t *testing.T) {
    t.Log("testing the GET method on all users")
    rec = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v0/users", nil)
	api.Handler().ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	if err != nil {
		t.Error("error getting users.")
	}
}

func TestGetOneUser(t *testing.T) {
	t.Log("testing the GET method on one users")
	rec = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v0/users/1", nil)
	api.Handler().ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	if err != nil {
		t.Error("error getting user with userid 1.")
	}
}
// UPDATE (PATCH)

// DELETE

func TestMain(m *testing.M) {
	setup()
    eval := m.Run()
    teardown()
    os.Exit(eval)
}