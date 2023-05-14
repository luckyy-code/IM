package service

import (
	"IM/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {

	c.HTML(200, "index.html", "views/chat/head.html")

	//c.JSON(200, gin.H{
	//	"message": "welcome !! ",
	//})
}

func ToRegister(c *gin.Context) {

	c.HTML(200, "register.html", nil)

	//c.JSON(200, gin.H{
	//	"message": "welcome !! ",
	//})
}

func ToChat(c *gin.Context) {

	ind, _ := template.ParseFiles(
		"views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/foot.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
	)

	user := models.UserBasic{}
	atoi, _ := strconv.Atoi(c.Query("userId"))
	user.ID = uint(atoi)
	user.Identity = c.PostForm("token")

	ind.Execute(c.Writer, user)

}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
