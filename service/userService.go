package service

import (
	"IM/models"
	"IM/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GetUserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1 失败
		"message": "获取成功",
		"data":    data,
	})
}

func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" || user.Name == "" || password == "" {

		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1 失败
			"message": "用户名已注册",
			"data":    data,
		})
		return
	}

	if password != repassword || password == "" || repassword == "" {

		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1 失败
			"message": "两次密码不一致",
			"data":    data,
		})
		return
	}

	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	//	user.PassWord = password
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1 失败
		"message": "新增用户成功",
		"data":    data,
	})

}

func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	//password := c.Query("password")

	name := c.Request.FormValue("mobile")
	password := c.Request.FormValue("passwd")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {

		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1 失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {

		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1 失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1 失败
		"message": "登录成功",
		"data":    data,
	})
}

func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Request.FormValue("id"))
	user.ID = uint(id)
	models.DeleteUser(user)

	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1 失败
		"message": "删除用户成功",
	})
}

func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.Avatar = c.PostForm("icon")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println("修改格式不对")
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1 失败
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)

		c.JSON(200, gin.H{
			"code":    0, // 0成功 -1 失败
			"message": "修改用户成功",
		})
	}

}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)

	MsgHandler(ws, c)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("MsgHandler 发送失败！", err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriend(c *gin.Context) {
	atoi, _ := strconv.Atoi(c.PostForm("userId"))

	u := uint(atoi)
	users := models.SearchFriend(u)

	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"message": "查询好友列表成功",
	//	"data":    users,
	//})
	utils.RespOKList(c.Writer, users, len(users))

}
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetName := c.PostForm("targetName")

	code, msg := models.AddFriend(uint(userId), targetName)

	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"message": "查询好友列表成功",
	//	"data":    users,
	//})
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}

}

//新建群

func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	Name := c.PostForm("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = Name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	models.JoinGroup(community.OwnerId, community.Name)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加载群列表

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))

	data, msg := models.LoadCommunity(uint(ownerId))

	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加入群 userId uint, comId uint

func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
