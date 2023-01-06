package models

import "gorm.io/gorm"

// 人员关系
type Contact struct {
	gorm.Model
	OwbedId  uint   //谁的关系
	TargetId uint   //对应的谁
	Type     int    //关系类型
	Desc     string //描述

}

func (table *Contact) TableName() string {
	return "contact"
}
