package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	AuthContrliier "github.com/worapong778/GoJwtApi_project/controller/auth"
	UserController "github.com/worapong778/GoJwtApi_project/controller/user"
	"github.com/worapong778/GoJwtApi_project/middleware"
	"github.com/worapong778/GoJwtApi_project/orm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// run Db
	orm.InitDB()

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", AuthContrliier.Register)
	r.POST("/login", AuthContrliier.Login)
	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadUsersAll)
	authorized.GET("/profile", UserController.Profile)

	// run API
	r.Run("localhost:8080")
}
