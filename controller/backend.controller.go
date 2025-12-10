package controller

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/Ryanljk/basic-backend/service"
	"github.com/Ryanljk/basic-backend/model"
)

type BackendController struct {
	Service *service.BackendService
}

func NewBackendController(s *service.BackendService) *BackendController {
	return &BackendController{Service: s}
}

func (bc *BackendController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

	user, err := bc.Service.GetUser(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// displays all user data in JSON file (not in requirements, for debugging)
func (bc *BackendController) GetAllUsers(c *gin.Context) {
	users := bc.Service.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func (bc *BackendController) AddUser(c *gin.Context) {
	var req model.User


	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if (req.Email == "" || req.Password == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing email/password"})
		return
	}

	err := bc.Service.AddUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user added"})
}

func (bc *BackendController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
            return
        }
	
	if err := bc.Service.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
