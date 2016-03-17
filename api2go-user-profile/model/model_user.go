 package model

import (
	"time"
	"errors"
	"strconv"
	"strings"
	"fmt"
	"github.com/manyminds/api2go/jsonapi"
)

// User is a generic database user
type User struct {
	Uid       int64         `db:"user_id" json:"user_id"`
	Username  string      `db:"username" json:"username"`
	Password  string      `db:"passwordhash" json:"password"`
	Created   time.Time   `db:"created_at" json:"created_at"`
	Updated   time.Time   `db:"last_update" json:"updated_at"`
	ProfilesIDs  string    `db:"profile_ids" json:"-"`
	ProfilesIDList  []int64    `db:"-" json:"profile_ids"`
	Profiles  []Profile   `db:"-" json:"profiles"`
	exists    bool        `db:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u User) GetID() string {
	// int to string
	return strconv.FormatInt(u.Uid, 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *User) SetID(id string) error {
	var err error
	// string to int
	u.Uid, err = strconv.ParseInt(id, 10, 64)
	return err
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u User) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "profiles",
			Name: "profiles",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u User) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	// split id string into ids
	s := strings.Split(u.ProfilesIDs, ",")
    for i := range s {
    	fmt.Println(s[i])
    	result = append(result, jsonapi.ReferenceID{
			ID: string(s[i]),
			Type: "user-profile",
			Name: "profiles",
		})
	}
	return result
}
// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
// func (u User) GetReferencedStructs() []jsonapi.MarshalIdentifier {
// 	result := []jsonapi.MarshalIdentifier{}
// 	for key := range u.Profiles {
// 		result = append(result, u.Profiles[key])
// 	}
// 	fmt.Println(result)
// 	return result
// }

// SetToManyReferenceIDs sets the profile reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *User) SetToManyReferenceIDs(name string, IDs []int64) error {
	if name == "profiles" {
		u.ProfilesIDList = IDs
	}
	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds new profiles to a user
func (u *User) AddToManyIDs(name string, IDs []int64) error {
	if name == "profiles" {
		u.ProfilesIDList = append(u.ProfilesIDList, IDs...)
	}
	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes profiles from a user
func (u *User) DeleteToManyIDs(name string, IDs []int64) error {
	if name == "profiles" {
		for _, ID := range IDs {
			for pos, oldID := range u.ProfilesIDList {
				if ID == oldID {
					// match, this ID must be removed
					u.ProfilesIDList = append(u.ProfilesIDList[:pos], u.ProfilesIDList[pos+1:]...)
				}
			}
		}
	}
	return errors.New("There is no to-many relationship with the " + name)
}
