package main

import (
	"log"
	"os"
	"testezhik/cmd/api/handlers"
	"testezhik/cmd/api/initialization"
	"testezhik/cmd/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Server settings
	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())
	router.LoadHTMLGlob("./www/templates/*/*.html")
	// Endpoints
	router.NoRoute(handlers.Error404)
	router.NoMethod(handlers.ErrorNoMethod)
	// router.Static()
	router.GET("/", handlers.IndexHandler)
	// Registration
	router.GET("/registration", handlers.RegGET)
	router.POST("/registration", handlers.RegPOST)
	// Login
	router.GET("/login", handlers.LoginGET)
	router.POST("/login", handlers.LoginPOST)
	// // Email verification
	router.GET("verify-email/:user/:emailVerHash", handlers.EmailVer)
	// Profile
	authGroup := router.Group("/user", handlers.CheckAuth)
	authGroup.GET("/profile", handlers.Profile)
	authGroup.GET("/logout", handlers.Logout)
	authGroup.GET("/changePswd", handlers.ChangePswdGET)
	authGroup.POST("/changePswd", handlers.ChangePswdPOST)
	// Account recovery
	router.GET("/recovery", handlers.RecoveryGET)
	router.POST("/recovery", handlers.RecoveryPOST)
	router.GET("/account-recovery/:username/:verpswd", handlers.AccountRecoveryGET)
	router.POST("/account-recovery/:username/:verpswd", handlers.AccountRecoveryPOST)

	// ENV var init
	initialization.LoadEnv("../../.env")
	storage.ConnToDB()
	//Server run
	port := os.Getenv("port")
	log.Printf("Server starts in port %s", port)
	log.Fatalln(router.Run(port))
}
