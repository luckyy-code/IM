package main

import (
	"IM/models"
	"IM/utils"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	utils.InitConfig()

	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数:", err)
	}

	//迁移 schema
	//db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Community{})
	//db.AutoMigrate(&models.Contact{})
	//db.AutoMigrate(&models.GroupBasic{})
	//
	////Creat
	//user := &models.UserBasic{}
	//user.Name = "sh==伸展"
	//db.Create(user)
	//
	//db.Model(user).Update("PassWord", "1234")
}
