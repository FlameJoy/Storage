package handlers

import (
	"net/http"
	"testezhik/cmd/api/data"
	"testezhik/cmd/api/models"

	"github.com/gin-gonic/gin"
)

func RegGET(c *gin.Context) {
	c.HTML(http.StatusOK, "regForm.html", data.LinkList)
}

func RegPOST(c *gin.Context) {
	var u models.User
	u.Username = c.PostForm("username")
	u.Email = c.PostForm("email")
	pswd1 := c.PostForm("password1")
	pswd2 := c.PostForm("password2")
	// Username validation
	err := u.ValidateUsername()
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Email validation
	err = u.ValidateEmail()
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Password validation
	err = u.ValidatePswd(pswd1, pswd2)
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Check if user already exist
	err = u.UserExist()
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Check if email already exist
	err = u.EmailExist()
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	// Create new user
	err = u.New(pswd1)
	if err != nil {
		c.HTML(http.StatusBadRequest, "regForm.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
			"user": u,
		})
		return
	}
	c.HTML(http.StatusOK, "reg-success.html", gin.H{
		"menu": data.Menulist,
	})
}
