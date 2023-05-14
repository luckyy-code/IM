package models

import (
	"IM/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Avatar        string //头像
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}

	utils.Db.Where("name = ? and pass_word = ?", name, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)

	utils.Db.Model(&user).Where("id = ?", user.ID).Update("Identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}

	utils.Db.Where("name = ?", name).First(&user)
	return user
}
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.Db.Where("phone = ?", phone).First(&user)
}
func FindUserById(id uint) UserBasic {
	user := UserBasic{}

	utils.Db.Where("id = ?", id).First(&user)
	return user
}
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.Db.Where("email = ?", email).First(&user)
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.Db.Find(&data)

	for _, v := range data {
		fmt.Println(v)
	}

	return data
}
func CreateUser(user UserBasic) *gorm.DB {
	return utils.Db.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.Db.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.Db.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}

//查找某个用户

func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.Db.Where("id = ?", id).First(&user)
	return user
}
