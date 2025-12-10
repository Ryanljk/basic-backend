package main

import (
	"encoding/json"
	"os"

	"github.com/Ryanljk/basic-backend/controller"
	"github.com/Ryanljk/basic-backend/model"
	"github.com/Ryanljk/basic-backend/service"
	"github.com/gin-gonic/gin"
)

var users []model.User

//function to load JSON file
func loadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	port := os.Getenv("PORT")
	filepath := os.Getenv("FILEPATH")

	//load JSON
	err := loadJSON(filepath)
	if err != nil {
		panic("Failed to load data.json: " + err.Error())
	}

	service := service.NewBackendService(&users, filepath)
	controller := controller.NewBackendController(service)

	//define endpoints
	api := server.Group("/api")
	{
		api.GET("/", controller.GetAllUsers) //displays all users in JSON
		api.GET("/get/:id", controller.GetUser) //displays only 1 user, search by ID
		api.POST("/add", controller.AddUser) //adds 1 user to json
		api.DELETE("/delete/:id", controller.DeleteUser) //
	}

	server.Run(":" + port)
}