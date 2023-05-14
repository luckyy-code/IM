package models

import (
	"IM/utils"
	"fmt"
	"gorm.io/gorm"
)

//人员关系

type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系
	TargetId uint //对应谁
	Type     int  //对应类型  1好友 2群组 3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.Db.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>", v)
		objIds = append(objIds, uint64(v.TargetId))
	}

	users := make([]UserBasic, 0)
	utils.Db.Where("id in ?", objIds).Find(&users)

	return users
}

func AddFriend(userId uint, targetName string) (int, string) {
	//user := UserBasic{}
	if targetName != "" {
		targetUser := FindUserByName(targetName)
		if targetUser.Salt != "" {
			if targetUser.ID == userId {
				return -1, "不能加自己"
			}

			contact0 := Contact{}
			utils.Db.Where("owner_id = ? and target_id = ? and type = 1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不能重复添加"
			}

			tx := utils.Db.Begin()
			//事务开始不论什么异常都回滚
			defer func() {
				r := recover()
				if r != nil {
					tx.Rollback()
				}
			}()

			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetUser.ID
			contact.Type = 1
			err := utils.Db.Create(&contact).Error
			if err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId, contact1.TargetId = targetUser.ID, userId
			contact1.Type = 1
			err1 := utils.Db.Create(&contact1).Error
			if err1 != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()

			return 0, "添加好友成功"
		}
		return -1, "没有找到此用户"
	}
	return -1, "好友ID不能为空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.Db.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
