package service

import (
	"ginchat/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex godoc
// @Tags         首页
// @Success      200  {string}  welcome
// @Router       /index [get]
func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login.html", nil)
}

func ToRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ToChat(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	c.HTML(http.StatusOK, "chat/index.html", gin.H{
		"user": user,
	})
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}