package main

import (
	"net/http"
	"testezhik/cmd/api/data"

	"github.com/gin-gonic/gin"
)

func error404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", data.LinkList)
}

func errorNoMethod(c *gin.Context) {
	c.HTML(http.StatusMethodNotAllowed, "methodNotAllowed.html", data.LinkList)
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", data.LinkList)
}
