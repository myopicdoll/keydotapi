package storage

import (
	"fmt"
	"strconv"
	"gopkg.in/gorp.v1"
	"api-sandbox/api2go-user-profile/model"
)

// NewProfileStorage initializes the storage
func NewProfileStorage(db *gorp.DbMap) *ProfileStorage {
	return &ProfileStorage{db}
}

type ProfileStorage struct {
	db *gorp.DbMap
}

func (s ProfileStorage) GetAll() []model.Profile {
	var profiles []model.Profile
	s.db.Select(&profiles, "select * from profiles order by profile_id")
	return profiles
}

// get one profile by profile id
func (s ProfileStorage) GetOne(id string) (model.Profile, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return model.Profile{}, fmt.Errorf("Profile id must be an integer: %s", id)
	}
	var prof model.Profile
	error := s.db.SelectOne(&prof, "select * from profiles where profile_id = $1 limit 1", intID)
	if error == nil {
		return prof, nil
	} else {
		return model.Profile{}, fmt.Errorf("Profile with id %s not found", id)
	}	
}

func newProf(user int64, total_trips int64) model.Profile {
    return model.Profile{
        User: user,
		NumTrips: total_trips,
    }
}
// Insert a profile
func (s *ProfileStorage) Insert(p model.Profile) (string, error) {
	// create a new profile
    prof := newProf(3, 3)
    // insert rows - auto increment PKs will be set properly after the insert
    error := s.db.Insert(&prof)
	if error != nil {
		return "", error
	}
	return p.GetID(), nil
}

// Delete one 
func (s *ProfileStorage) Delete(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Profile id must be an integer: %s", id)
	}
	_, error := s.db.Exec("delete from profiles where profile_id=?", intID)
	
	if error != nil {
		return fmt.Errorf("Profile with id %s does not exist", id)
	}
	return error
}

// Update updates an existing chocolate
func (s *ProfileStorage) Update(p model.Profile) error {
	var prof model.Profile
	error := s.db.SelectOne(&prof, "select * from profiles where profile_id = $1 limit 1", p.Pid)
	prof.NumTrips = p.NumTrips
	if error != nil {
		return error
	}	
	prof.User = p.User
	_, err  := s.db.Update(&prof)
	return err
}
