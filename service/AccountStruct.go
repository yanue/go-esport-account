/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-数据库操作

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 用户表
type TUser struct {
	Id             int
	Account        string `gorm:"not null;size:50;unique;comment:'登陆账号,唯一';"`
	Name           string `gorm:"not null;index;size:100;"`
	Gender         int    `gorm:"not null;comment:'0未设置,1男,2女';"`
	Phone          string `gorm:"not null;type:char(11);comment:'手机号';"`
	Email          string `gorm:"not null;comment:'邮箱';"`
	SchoolId       int    `gorm:"not null;index"`
	ClassId        int    `gorm:"not null;index"`
	AreaId         int    `gorm:"not null;index"`
	IdentityNo     string `gorm:"not null;type:char(18);comment:'身份证号';"`
	IdentityName   string `gorm:"not null;size:60;comment:'身份姓名(真实姓名)';"`
	IdentityStatus int    `gorm:"not null;index;comment:'身份认证状态';"`
	Created        int    `gorm:"not null;"`
	Modified       int    `gorm:"not null;"`
}

// 第三方登陆信息
type TUserAuth struct {
	Id         int
	UserId     int    `gorm:"not null;"`
	AuthSite   string `gorm:"type:ENUM('wx', 'qq', 'wb');default:'wx';"`
	AuthOpenid int    `gorm:"not null;comment:'微信/qq等openid';"`
	AuthToken  int    `gorm:"not null;comment:'授权信息';"`
	AuthExpire int    `gorm:"not null;comment:'授权过期时间';"`
	Created    int    `gorm:"not null;"`
	Modified   int    `gorm:"not null;"`
}

// 省份
type TAreaProvince struct {
	Id   int
	Name string `gorm:"not null;index"`
}

// 城市
type TAreaCity struct {
	Id         int
	ProvinceId int    `gorm:"not null;index"`
	Name       string `gorm:"not null;index"`
}

// 学校表
type TSchool struct {
	Id      int
	Name    string `gorm:"not null;size:255"`
	AreaId  string `gorm:"not null;"`
	Created int    `gorm:"not null;"`
}

// 班级表
type TClass struct {
	Id      int
	Name    string `gorm:"not null;size:255"`
	AreaId  string `gorm:"not null;"`
	Created int    `gorm:"not null;"`
}

type CRedisKey struct{}

/**
用户信息
 */
func (this *CRedisKey) HUserInfo(uid int) (string) {
	return RedisPrefix + fmt.Sprintf("user:%d", uid)
}
