/*
Create a new user: 
	curl -X POST http://localhost:31415/v0/users -d '{"data" : {"type" : "users" , "attributes": {"username" : "flamingo", "password": "pass12"}}}'
Read (List) all users:
	curl -X GET http://localhost:31415/v0/users
** Update doesn't work yet due to "profile_ids" not being an actual database column but exists in the User model
** in production gorp will be forgo due to too many limitations and possibly gorm will be used, or just plain pq.
Update:
	curl -vX PATCH http://localhost:31415/v0/users/6 -d '{ "data" : {"type" : "users", "id": "6", "attributes": {"username" : "better marvin", "passwordhash": "pass6"}}}'
Delete a user:
	curl -vX DELETE http://localhost:31415/v0/users/4
*/
package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"api-sandbox/api2go-user-profile/model"
	"api-sandbox/api2go-user-profile/resource"
	"api-sandbox/api2go-user-profile/storage"
	_ "github.com/lib/pq"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"
)

func main() {
	r := gin.Default()
	api := api2go.NewAPIWithRouting(
		"v0",
		api2go.NewStaticResolver("/"),
		gingonic.New(r),
	)

	dbmap, err := storage.InitDb()
	if err != nil {
		panic(err)
	}
	defer dbmap.Db.Close()
	userStorage := storage.NewUserStorage(dbmap)
	profStorage := storage.NewProfileStorage(dbmap)
	api.AddResource(model.User{}, resource.UserResource{ProfStorage: profStorage, UserStorage: userStorage})
	api.AddResource(model.Profile{}, resource.ProfileResource{ProfStorage: profStorage, UserStorage: userStorage})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Run(":31415") // listen and serve on 0.0.0.0:31415
}
