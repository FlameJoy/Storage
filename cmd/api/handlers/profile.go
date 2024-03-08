package handlers

import (
	"net/http"
	"strconv"
	"testezhik/cmd/api/data"
	"testezhik/cmd/api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckAuth(c *gin.Context) {
	token, err := models.GetToken(c)
	if err != nil || !token.Valid {
		c.Redirect(http.StatusPermanentRedirect, "/login")
		c.Abort()
		return
	}
	c.Next()
}

func Profile(c *gin.Context) {
	var u models.User
	token, err := models.GetToken(c)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, ok1 := claims["userID"].(string)
	userMongoID, ok2 := claims["userMongoID"].(string)
	if !ok1 && !ok2 {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  "Cannot retrive a user id from token",
		})
		return
	}
	parsedID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	u.ID = uint(parsedID)
	u.MongoID = userMongoID
	err = u.GetUserByID()
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"menu": data.Menulist,
		"user": u,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, false)
	c.HTML(http.StatusOK, "login.html", gin.H{
		"menu": data.Menulist,
		"msg":  "Logged out",
	})
}

func ChangePswdGET(c *gin.Context) {
	c.HTML(http.StatusOK, "changePswd.html", gin.H{
		"menu": data.Menulist,
	})
}

func ChangePswdPOST(c *gin.Context) {
	var u models.User
	oldPswd := c.PostForm("old_password")
	pswd1 := c.PostForm("password1")
	pswd2 := c.PostForm("password2")
	token, err := models.GetToken(c)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, ok1 := claims["userID"].(string)
	userMongoID, ok2 := claims["userMongoID"].(string)
	if !ok1 && !ok2 {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  "Cannot retrive a userID from token",
		})
		return
	}
	parsedID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	u.ID = uint(parsedID)
	u.MongoID = userMongoID
	err = u.GetUserByID()
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PswdHash), []byte(oldPswd))
	if err != nil {
		c.HTML(http.StatusBadRequest, "changePswd.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = u.ValidatePswd(pswd1, pswd2)
	if err != nil {
		c.HTML(http.StatusBadRequest, "changePswd.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	err = u.ChangePswd(pswd1)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "changePswd.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "changePswd.html", gin.H{
		"menu": data.Menulist,
		"msg":  "password successfuly changed",
	})
}
