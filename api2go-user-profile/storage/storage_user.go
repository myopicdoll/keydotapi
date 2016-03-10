package storage

import (
	"errors"
	"strings"
	"fmt"
	// "log"
	// "database/sql"
	"net/http"
	"strconv"
	"gopkg.in/gorp.v1"
	"api-sandbox/api2go-user-profile/model"
	"github.com/manyminds/api2go"
)

// NewUserStorage initializes the storage
func NewUserStorage(db *gorp.DbMap) *UserStorage {
	return &UserStorage{db}
}

// UserStorage stores all users
type UserStorage struct {
	db *gorp.DbMap
}

// Store long SQL queries
const (
	GET_ALL = "select users.user_id, username, passwordhash, created_at, last_update, user_profile.profile_ids FROM users left join (select array_to_string(array_agg(profile_id),',') AS profile_ids, user_id from profiles group by user_id) as user_profile on users.user_id = user_profile.user_id"
)

// GetAll returns the user map (because we need the ID as key too)
func (s UserStorage) GetAll() (map[int64]*model.User, error) {
	var users []model.User
	// select all users and their associated profile ids
	_, err := s.db.Select(&users, GET_ALL)
	if err == nil {
		userMap := make(map[int64]*model.User)
		for i := range users {
			u := &users[i]
			fmt.Println("storing user in map with id %s", u.Uid)
			userMap[u.Uid] = u
		}
		return userMap, nil
	} else {
		// return the error
		return nil, err
	}
}

// Get one particular user given user id
func (s UserStorage) GetOne(id string) (model.User, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("User id must be an integer: %s", id)
		return model.User{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusBadRequest)
	}
	var user model.User
	str := []string{GET_ALL, "where users.user_id = $1 limit 1"}
	GET_ONE := strings.Join(str, " ")
	error := s.db.SelectOne(&user, GET_ONE, intID)
	if error == nil {
		return user, nil
	} else {
		return model.User{}, fmt.Errorf("User for id %s not found", id)
	}	
	// return s.getOneWithAssociations(intID)
}

// func (s UserStorage) getOneWithAssociations(id int64) (model.User, error) {
// 	var user model.User
// 	s.db.First(&user, id)
// 	s.db.Model(&user).Related(&user.Chocolates, "Chocolates")
// 	if err := s.db.Error; err == gorm.RecordNotFound {
// 		errMessage := fmt.Sprintf("User for id %s not found", id)
// 		return model.User{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
// 	} else if err != nil {
// 		return model.User{}, err
// 	}
// 	user.ChocolatesIDs = make([]string, len(user.Chocolates))
// 	for i, choc := range user.Chocolates {
// 		user.ChocolatesIDs[i] = choc.GetID()
// 	}
// 	return user, nil
// }

func newUser(username string, password string) model.User {
    return model.User{
        Username:   username,
        Password: password,
    }
}
// Insert a user
func (s *UserStorage) Insert(c model.User) (string, error) {
	// c.Chocolates = make([]model.Chocolate, len(c.ChocolatesIDs))
	// err := s.updateChocolatesByChocolatesIDs(&c)
	// if err != nil {
	// 	return "", err
	// }

	// create a new user
    u := newUser("test_user", "test_password")
    // insert rows - auto increment PKs will be set properly after the insert
    error := s.db.Insert(&u)
    
	if error != nil {
		return "", error
	}
	return c.GetID(), nil
}

// Delete one 
func (s *UserStorage) Delete(id string) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("User id must be integer: %s", id)
	}
	_, error := s.db.Exec("delete from user where uid=?", intID)
	
	if error != nil {
		return fmt.Errorf("User with id %s does not exist", id)
	}
	return error
}

// Update a user
func (s *UserStorage) Update(c model.User) error {
	c.Username = "test_update"
    _, err := s.db.Update(&c)
    return err
}

func indexOf(s string, items []string) int {
	for i, item := range items {
		if s == item {
			return i
		}
	}
	return -1
}

func (s *UserStorage) updateChocolatesByChocolatesIDs(u *model.User) error {
	// u.Chocolates = make([]model.Chocolate, len(u.ChocolatesIDs))
	// for i, chocID := range u.ChocolatesIDs {
	// 	intID, err := strconv.ParseInt(chocID, 10, 64)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	s.db.SelectOne(&u.Chocolates[i], "select * from chocolates where cid = $1 limit 1", intID)
	// }
	return nil
}
