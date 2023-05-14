package router

import (
	"IM/docs"
	"IM/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/register", service.ToRegister)
	r.GET("/chat", service.Chat)
	r.GET("/toChat", service.ToChat)
	r.POST("/searchFriends", service.SearchFriend)
	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/login", service.FindUserByNameAndPwd)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/find", service.FindByID)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/attach/upload", service.Upload)
	r.POST("/contact/addfriend", service.AddFriend)
	r.POST("/contact/createCommunity", service.CreateCommunity)
	r.POST("/contact/loadcommunity", service.LoadCommunity)
	r.POST("/contact/joinGroup", service.JoinGroups)

	r.POST("/user/redisMsg", service.RedisMsg)
	return r
}
