package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string //发送者
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_bassic"
}
