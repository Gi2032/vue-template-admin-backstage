package main

import (
	"ggg/handle"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	si, err := handle.InitInterface("root:guochang123@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}
	router.Use(cors.New(config))
	router.POST("/login", si.UserLogin)
	router.POST("/register", si.UserRegister)
	router.GET("/user", si.GetUserList)
	router.DELETE("/deleteUser", si.DeleteUser)
	router.POST("/addUser", si.AddUser)
	router.POST("/updateUser", si.UpdateUser)
	router.GET("/getMenuList", si.GetMenuList)
	router.POST("/addMenu", si.AddMenu)
	router.DELETE("/deleteMenu", si.DeleteMenu)
	router.POST("/updateMenu", si.UpdateMenu)
	router.Run(":8080")
}
