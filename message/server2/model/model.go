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


type Publish struct {
	Id int64  `gorm:"primary_key:true"`
	MessageId int64
	Ts        time.Time
}

func (Publish) TableName() string {
	return "message.publish"
}