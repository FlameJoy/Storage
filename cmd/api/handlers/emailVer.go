package handlers

import (
	"log"
	"net/http"
	"testezhik/cmd/api/data"
	"testezhik/cmd/api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EmailVer(c *gin.Context) {
	var u models.User
	u.Username = c.Param("user")
	linkHash := c.Param("emailVerHash")
	// Check user exist
	err := u.GetUserByUsername()
	if err != nil {
		log.Println("Error: ", err)
		c.HTML(http.StatusBadRequest, "email-ver.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Compare hash
	err = bcrypt.CompareHashAndPassword([]byte(u.VerHash), []byte(linkHash))
	if err == nil {
		// Activate account
		err := u.VerifyAccount()
		if err != nil {
			log.Println("Error: ", err)
			c.HTML(http.StatusBadRequest, "email-ver.html", gin.H{
				"menu": data.Menulist,
				"msg":  err.Error(),
				"user": u,
			})
			return
		}
		c.HTML(http.StatusOK, "email-ver.html", gin.H{
			"menu": data.Menulist,
			"msg":  "Email verified",
		})
		return
	}
	log.Println("Error: ", err)
	c.HTML(http.StatusBadRequest, "email-ver.html", gin.H{
		"menu": data.Menulist,
		"msg":  "Please try to verify your email again: " + err.Error(),
	})
}
