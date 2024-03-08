package handlers

import (
	"net/http"
	"testezhik/cmd/api/data"

	"github.com/gin-gonic/gin"
)

func Error404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", data.LinkList)
}

func ErrorNoMethod(c *gin.Context) {
	c.HTML(http.StatusMethodNotAllowed, "methodNotAllowed.html", data.LinkList)
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", data.LinkList)
}
