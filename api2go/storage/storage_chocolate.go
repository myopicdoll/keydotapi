package storage

import (
	"fmt"
	"strconv"
	"gopkg.in/gorp.v1"
	"api-sandbox/api2go/model"
)

// NewChocolateStorage initializes the storage
func NewChocolateStorage(db *gorp.DbMap) *ChocolateStorage {
	return &ChocolateStorage{db}
}

// ChocolateStorage stores all of the tasty chocolate, needs to be injected into
// User and Chocolate Resource. In the real world, you would use a database for that.
type ChocolateStorage struct {
	db *gorp.DbMap
}

// GetAll of the chocolate
func (s ChocolateStorage) GetAll() []model.Chocolate {
	var chocolates []model.Chocolate
	s.db.Select(&chocolates, "select * from chocolates order by cid")
	// s.db.Order("id").Find(&chocolates)
	return chocolates
}

// GetOne tasty chocolate
func (s ChocolateStorage) GetOne(id string) (model.Chocolate, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return model.Chocolate{}, fmt.Errorf("Chocolate id must be integer: %s", id)
	}
	var choc model.Chocolate
	error := s.db.SelectOne(&choc, "select * from chocolates where cid = $1 limit 1", intID)
	//s.db.First(&choc, intID)
	if error == nil {
		return choc, nil
	} else {
		return model.Chocolate{}, fmt.Errorf("Chocolate for id %s not found", id)
	}	
}

func newChoc(name string, taste string) model.Chocolate {
    return model.Chocolate{
        // Created: time.Now().UnixNano(),
        Name: name, 
        Taste: taste,
    }
}
// Insert a chocolate
func (s *ChocolateStorage) Insert(c model.Chocolate) (string, error) {
	// create a new user
    choc := newChoc("test_choc", "sweet")
    // insert rows - auto increment PKs will be set properly after the insert
    error := s.db.Insert(&choc)
	if error != nil {
		return "", error
	}
	return c.GetID(), nil
}

// Delete one 
func (s *ChocolateStorage) Delete(id string) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("Chocolate id must be integer: %s", id)
	}
	_, error := s.db.Exec("delete from chocolates where cid=?", intID)
	
	if error != nil {
		return fmt.Errorf("Chocolate with id %s does not exist", id)
	}
	return error
}

// Update updates an existing chocolate
func (s *ChocolateStorage) Update(c model.Chocolate) error {
	var choc model.Chocolate
	error := s.db.SelectOne(&choc, "select * from chocolates where cid = $1 limit 1", c.ID)
	//s.db.First(&choc, intID)
	choc.Name = "new_updated_name"
	if error != nil {
		return error
	}	
	choc.Taste = c.Taste
	_, err  := s.db.Update(&choc)
	return err
}
