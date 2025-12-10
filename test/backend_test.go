package test

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ryanljk/basic-backend/controller"
	"github.com/Ryanljk/basic-backend/model"
	"github.com/Ryanljk/basic-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
)

func setupRouter(users *[]model.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	svc := service.NewBackendService(users, "../data/test_data.json")
	controller := controller.NewBackendController(svc)

	api := router.Group("/api")
	{
		api.GET("/", controller.GetAllUsers) //displays all users in JSON
		api.GET("/:id", controller.GetUser) //displays only 1 user, search by ID
		api.POST("/", controller.AddUser) //adds 1 user to json
		api.DELETE("/:id", controller.DeleteUser) //delete 1 user, by ID
	}

	return router
}

//helper function to verify that the hashed password matches the plaintext password 
func verifyPassword(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false
	}
	//decode salt & hash from base64
	salt, err1 := base64.RawStdEncoding.DecodeString(parts[0])
	hash, err2 := base64.RawStdEncoding.DecodeString(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

func TestAddUser(t *testing.T) {
	users := []model.User{}
	r := setupRouter(&users)

	//valid input
	payload := model.User{Email: "testadd@gmail.com", Password: "pass"}
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Len(t, users, 1)
	assert.Equal(t, "testadd@gmail.com", users[0].Email)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, true, verifyPassword(users[0].Password, "pass")) //ensure that plaintext password matches hashed

	//duplicate email
	req2, _ := http.NewRequest("POST", "/api/", bytes.NewBuffer(data))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)	

	//missing parameter
	payload2 := model.User{Password: "pass"}
	data2, _ := json.Marshal(payload2)
	req3, _ := http.NewRequest("POST", "/api/", bytes.NewBuffer(data2))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 400, w3.Code)

	//invalid email
	payload3 := model.User{Email: "hello"}
	data3, _ := json.Marshal(payload3)
	req4, _ := http.NewRequest("POST", "/api/", bytes.NewBuffer(data3))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	assert.Equal(t, 400, w2.Code)
}

func TestGetUser(t *testing.T) {
	users := []model.User{
		{ID: 1, Email: "testget@gmail.com", Password: "pass"}, //not using hashed password as that was already tested in beforehand test case
	}
	r := setupRouter(&users)

	//valid input
	req, _ := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var u model.User
	_ = json.Unmarshal(w.Body.Bytes(), &u)
	assert.Equal(t, "testget@gmail.com", u.Email)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "pass", u.Password)

	//user doesnt exist
	req2, _ := http.NewRequest("GET", "/api/1000", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)

	//invalid id
	req3, _ := http.NewRequest("GET", "/api/abc", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 400, w3.Code)
}

func TestDeleteUser(t *testing.T) {
	users := []model.User{
		{ID: 1, Email: "testdelete@gmail.com", Password: "pass"},
	}
	r := setupRouter(&users)

	//valid input
	req, _ := http.NewRequest("DELETE", "/api/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Len(t, users, 0)

	//missing user
	req2, _ := http.NewRequest("DELETE", "/api/1000", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)

	//invalid id
	req3, _ := http.NewRequest("DELETE", "/api/abc", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 400, w3.Code)
}
