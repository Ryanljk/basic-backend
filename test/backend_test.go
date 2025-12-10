package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ryanljk/basic-backend/model"
	"github.com/Ryanljk/basic-backend/service"
	"github.com/Ryanljk/basic-backend/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(users *[]model.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	svc := service.NewBackendService(users, "../data/test_data.json")
	controller := controller.NewBackendController(svc)

	r.GET("/api/get/:id", controller.GetUser)
	r.GET("/api/users", controller.GetAllUsers)
	r.POST("/api/add", controller.AddUser)
	r.DELETE("/api/delete/:id", controller.DeleteUser)

	return r
}

func TestAddUser(t *testing.T) {
	users := []model.User{}
	r := setupRouter(&users)

	//valid input
	payload := model.User{Email: "testadd@gmail.com", Password: "pass"}
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/add", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Len(t, users, 1)
	assert.Equal(t, "testadd@gmail.com", users[0].Email)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "pass", users[0].Password)

	//duplicate email
	req2, _ := http.NewRequest("POST", "/api/add", bytes.NewBuffer(data))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 500, w2.Code)

	//missing parameter
	payload2 := model.User{Password: "pass"}
	data2, _ := json.Marshal(payload2)
	req3, _ := http.NewRequest("POST", "/api/add", bytes.NewBuffer(data2))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 400, w3.Code)
}

func TestGetUser(t *testing.T) {
	users := []model.User{
		{ID: 1, Email: "testget@gmail.com", Password: "pass"},
	}
	r := setupRouter(&users)

	//valid input
	req, _ := http.NewRequest("GET", "/api/get/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var u model.User
	_ = json.Unmarshal(w.Body.Bytes(), &u)
	assert.Equal(t, "testget@gmail.com", u.Email)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "pass", u.Password)

	//user doesnt exist
	req2, _ := http.NewRequest("GET", "/api/get/1000", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)

	//invalid id
	req3, _ := http.NewRequest("GET", "/api/get/abc", nil)
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
	req, _ := http.NewRequest("DELETE", "/api/delete/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Len(t, users, 0)

	//missing user
	req2, _ := http.NewRequest("DELETE", "/api/delete/1000", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)

	//invalid id
	req3, _ := http.NewRequest("DELETE", "/api/delete/abc", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 400, w3.Code)
}
