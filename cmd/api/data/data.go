package data

import "github.com/gin-gonic/gin"

type menuList []struct {
	Title, Link string
}

var Menulist = menuList{
	{Title: "Main", Link: "/"},
	{Title: "Log in", Link: "/login"},
	{Title: "Registration", Link: "/registration"},
	{Title: "Profile", Link: "/user/profile"},
}

var LinkList = gin.H{
	"menu": Menulist,
}
