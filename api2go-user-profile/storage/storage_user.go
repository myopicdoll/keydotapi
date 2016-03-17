package storage

import (
	"errors"
	"strings"
	"fmt"
	// "log"
	// "database/sql"
	"time"
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
	// GET_ALL = "select users.user_id, username, passwordhash, created_at, last_update, user_profile.profile_ids FROM users left join (select array_to_string(array_agg(profile_id),',') AS profile_ids, user_id from profiles group by user_id) as user_profile on users.user_id = user_profile.user_id"
	GET_ALL = "select users.user_id, username, passwordhash, created_at, last_update, coalesce(user_profile.profile_ids, '0') as profile_ids FROM users left join (select array_to_string(array_agg(profile_id),',') AS profile_ids, user_id from profiles group by user_id) as user_profile on users.user_id = user_profile.user_id"
)

// GetAll returns the user map (because we need the ID as key too)
func (s UserStorage) GetAll() (map[int64]*model.User, error) {
	var users []model.User
	// select all users and their associated profile ids
	_, err := s.db.Select(&users, GET_ALL)
	fmt.Println(users)
	fmt.Println(err)
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

func NewUser(username string, password string) model.User {
    return model.User{
        Username:   username,
        Password: password,
        Created: time.Now(), 
    	Updated: time.Now(),
    }
}
// Insert a user
func (s *UserStorage) Insert(u model.User) error {
	// create a new user
    //user := NewUser(u.Username, u.Password)
	//user.Profiles = make([]model.Profile, len(u.ProfilesIDList))
    // insert rows - auto increment PKs will be set properly after the insert
    // error := s.db.Insert(&user)
	// NOTE: user raw SQL here because otherwise db will complain that the column
	// profile_ids does not exist
	_, err := s.db.Exec(`insert into users (username, passwordhash) values ($1, $2) returning user_id`, u.Username, u.Password)
	if err != nil {
		return err
	} 
	return nil
}

// Delete one 
func (s *UserStorage) Delete(id string) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("User id must be integer: %s", id)
	}
	_, error := s.db.Exec("delete from users where user_id = $1", intID)
	
	if error != nil {
		return fmt.Errorf("User with id %s does not exist", id)
	}
	return error
}

// Update a user
func (s *UserStorage) Update(u model.User) error {
    _, err := s.db.Update(&u)
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