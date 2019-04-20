package model

import (
	"time"
	
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Message struct {
	Id        int64          `gorm:"primary_key:true"`
	Srv       string
	Event     string
	Data      postgres.Jsonb
	CreatedAt time.Time
	Active    bool
}

func (Message) TableName() string {
	return "message.message"
}


type PublishResult struct {
	Id    int64     `gorm:"primary_key:true"`
	MsgId int64
	Ts    time.Time
}

func (PublishResult) TableName() string {
	return "message.publish_result"
}