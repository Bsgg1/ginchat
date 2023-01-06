package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromId   uint   //发送者
	TargetId uint   //接收者
	Type     string //消息类型 群聊私聊广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string //图片
	Url      string
	Desc     string //描述
	Amount   int    //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}
