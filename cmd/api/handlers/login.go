package handlers

import (
	"net/http"
	"os"
	"testezhik/cmd/api/data"
	"testezhik/cmd/api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func LoginGET(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusOK, "login.html", data.LinkList)
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		secretKey := os.Getenv("secret")
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		c.HTML(http.StatusOK, "login.html", data.LinkList)
		return
	}
	c.Redirect(http.StatusPermanentRedirect, "user/profile")
}

func LoginPOST(c *gin.Context) {
	var u models.User
	u.Username = c.PostForm("username")
	pswd := c.PostForm("password")
	// Get user
	err := u.GetUserByUsername()
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(u.PswdHash), []byte(pswd))
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	token, err := u.NewToken()
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"menu": data.Menulist,
			"msg":  err.Error(),
		})
		return
	}
	c.SetCookie("token", token, 3600, "/", "", true, false)
	c.HTML(http.StatusOK, "logged.html", gin.H{
		"user": u,
		"menu": data.Menulist,
	})
}
