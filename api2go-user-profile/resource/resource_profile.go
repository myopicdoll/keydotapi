package resource

import (
	"errors"
	"net/http"
	"api-sandbox/api2go-user-profile/model"
	"api-sandbox/api2go-user-profile/storage"
	"github.com/manyminds/api2go"
)

// ProfileResource for api2go routes
type ProfileResource struct {
	ProfStorage *storage.ProfileStorage
	UserStorage *storage.UserStorage
}

// FindAll profiles
func (p ProfileResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []model.Profile
	profs, _ := p.ProfStorage.GetAll()
	for _, prof := range profs {
		result = append(result, *prof)
	}
	return &Response{Res: result}, nil
}

// FindOne prof
func (p ProfileResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := p.ProfStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new prof
func (p ProfileResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	prof, ok := obj.(model.Profile)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id, _ := p.ProfStorage.Insert(prof)
	err := prof.SetID(id)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(errors.New("Non-integer ID given"), "Non-integer ID given", http.StatusInternalServerError)
	}
	return &Response{Res: prof, Code: http.StatusCreated}, nil
}

// Delete a prof 
func (p ProfileResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := p.ProfStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a prof
func (p ProfileResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	prof, ok := obj.(model.Profile)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := p.ProfStorage.Update(prof)
	return &Response{Res: prof, Code: http.StatusNoContent}, err
}
