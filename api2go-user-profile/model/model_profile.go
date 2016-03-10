package model

import (
	"strconv"
	"time"
)

// User and profile has a one-to-many relationship
// Profile has reference to userid but userid does not store any column reference to profiles
type Profile struct {
	Pid       int64         `db:"profile_id" json:"profile_id"`
	User      int64         `db:"user_id" json:"user"`
	NumTrips  int64      	  `db:"lifetime_trips" json:"lifetime_trips"`
	Created   time.Time   `db:"created_at" json:"created_at"`
	Updated   time.Time   `db:"last_update" json:"updated_at"`
	exists    bool        `db:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (p Profile) GetID() string {
	// int to string
	return strconv.FormatInt(p.Pid, 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (p *Profile) SetID(id string) error {
	var err error
	// string to int
	p.Pid, err = strconv.ParseInt(id, 10, 64)
	return err
}
