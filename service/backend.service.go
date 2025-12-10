package service

import (
	"errors"
	"encoding/json"
	"os"
	"github.com/Ryanljk/basic-backend/model"
)

type BackendService struct {
	Users *[]model.User
	FilePath string
}

func NewBackendService(users *[]model.User, path string) *BackendService {
	return &BackendService{
		Users: users,
		FilePath: path,
	}
}

func (bs *BackendService) updateJSON() error {
	data, err := json.MarshalIndent(bs.Users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(bs.FilePath, data, 0644)
}

func (bs *BackendService) GetUser(id int) (model.User, error) {
	for _, u := range *bs.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (bs *BackendService) GetAllUsers() []model.User {
	return *bs.Users
}

func (bs *BackendService) AddUser(u model.User) error {
	for _, user := range *bs.Users {
		if user.Email == u.Email {
			return errors.New("email already exists")
		}
	}
	//get latest available ID
	id := len(*bs.Users)
	u.ID = id + 1

	//make shallow copy 
	originalUsers := make([]model.User, len(*bs.Users))
	copy(originalUsers, *bs.Users)

	*bs.Users = append(*bs.Users, u)
	if err := bs.updateJSON(); err != nil {
		//if JSON update fails, rollback
		*bs.Users = originalUsers
		return errors.New("failed to save to JSON: " + err.Error())
	}
	return nil
}

func (bs *BackendService) DeleteUser(id int) error {
	for i, u := range *bs.Users {
		if u.ID == id {
			//make shallow copy 
			originalUsers := make([]model.User, len(*bs.Users))
			copy(originalUsers, *bs.Users)

			*bs.Users = append((*bs.Users)[:i], (*bs.Users)[i+1:]...)
			if err := bs.updateJSON(); err != nil {
			//if JSON update fails, rollback
			*bs.Users = originalUsers
				return errors.New("failed to save to JSON: " + err.Error())
			}
			return nil
		}
	}
	return errors.New("user not found")
}
