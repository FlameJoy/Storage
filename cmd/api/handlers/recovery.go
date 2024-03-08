package handlers

import (
	"log"
	"net/http"
	"testezhik/cmd/api/data"
	"testezhik/cmd/api/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RecoveryGET(c *gin.Context) {
	c.HTML(http.StatusOK, "recovery.html", gin.H{
		"menu": data.Menulist,
		"msg":  "",
	})
}

func RecoveryPOST(c *gin.Context) {
	var u models.User
	u.Email = c.PostForm("email")
	err := u.GetUserByEmail()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = u.NewEmailVerPswd()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "checkEmail.html", gin.H{
		"menu": data.Menulist,
	})
}

func AccountRecoveryGET(c *gin.Context) {
	var u models.User
	u.Username = c.Param("username")
	verPswd := c.Param("verpswd")
	err := u.GetUserByUsername()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.VerHash), []byte(verPswd))
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	currentTime := time.Now()
	if currentTime.After(u.Timeout) {
		c.HTML(http.StatusBadRequest, "recovery.html", gin.H{
			"menu": data.Menulist,
			"msg":  "link expired",
		})
		return
	}
	c.HTML(http.StatusOK, "changePswd2.html", gin.H{
		"menu":    data.Menulist,
		"verpswd": verPswd,
		"user":    u,
	})
}

func AccountRecoveryPOST(c *gin.Context) {
	var u models.User
	u.Username = c.Param("username")
	verPswd := c.Param("verpswd")
	err := u.GetUserByUsername()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery-unsuccess.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.VerHash), []byte(verPswd))
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery-unsuccess.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	pswd1 := c.PostForm("password1")
	pswd2 := c.PostForm("password2")
	err = u.ValidatePswd(pswd1, pswd2)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery-unsuccess.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = u.ChangePswd(pswd1)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "recovery-unsuccess.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "recovery-success.html", gin.H{
		"menu": data.Menulist,
	})
}
