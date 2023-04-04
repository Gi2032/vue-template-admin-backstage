package main

import (
	"ggg/handle"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	si, err := handle.InitInterface("root:guochang123@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/login", si.UserLogin)
	router.POST("/register", si.UserRegister)
	router.Run(":8080")
}
