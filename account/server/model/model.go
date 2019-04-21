package model

import (
    "time"
)

type Role struct {
    Id   int32  `gorm:"primary_key:true"`
    Name string
}

func (Role) TableName() string {
    return "account.role"
}

type User struct {
    Id       int64     `gorm:"primary_key: true"`
    Name     string 
    Password string
    JoinAt   time.Time
    RoleId   int32
    Active   bool
}

func (User) TableName() string {
    return "account.user"
}

type AuthHistory struct {
    Id        int64 `gorm:"primary_key: true"`
    UserId    int64
    Device    string
    Ip        string
    Fp        string
    CreatedAt time.Time
}

func (AuthHistory) TableName() string {
    return "account.auth_history"
}

type AuthToken struct {
    Id        int64     `gorm:"primary_key: true"`
    AuthId    int64
    UserId    int64
    UserName  string
    Signature string
    CreatedAt time.Time
    ExpiredAt time.Time
    Active    bool
}

func (AuthToken) TableName() string {
    return "account.auth_token"
}