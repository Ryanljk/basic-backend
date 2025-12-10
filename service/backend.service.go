package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/Ryanljk/basic-backend/model"
	"golang.org/x/crypto/argon2"
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

//hash password to prevent storage as plaintext	
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	//generate password hash using salt
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	//encode salt$hash via base64
	encoded := base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
	return encoded, nil
}

//check email domain to see if it is legitimate
func isValidDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	mxRecords, err := net.LookupMX(domain)
	return err == nil && len(mxRecords) > 0
}

//check provided email to see if it is syntatically valid
func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return (re.MatchString(email) && isValidDomain(email))
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

	//check email validity
	if (!isValidEmail(u.Email)) {
		return errors.New("invalid email")
	}

	//get latest available ID
	id := len(*bs.Users)
	u.ID = id + 1

	//hash password
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}
	u.Password = hashed

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
