/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : login.go
 Time    : 2018/9/25 12:19
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/errcode"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/sms"
	"github.com/yanue/go-esport-common/validator"
	"time"
)

/**
登陆授权后注册
-- 已经确认未注册
 */
func (this *AccountService) RegByAuth(auth *AuthLogin) (userAuth *TUserAuth, err error) {
	this.orm.db.Begin()
	// 最后执行
	defer func() {
		if err != nil {
			this.orm.db.Rollback()
		} else {
			this.orm.db.Commit()
		}
	}()

	user := &auth.User
	userAuth = &auth.Auth
	user.Created = time.Now().Unix()

	err = this.orm.db.Create(user).Error
	if err != nil {
		return
	}

	userAuth.UserId = user.Id
	userAuth.Created = time.Now().Unix()
	err = this.orm.db.Create(userAuth).Error

	return
}

/**
通过手机号注册(一键登录)
-- 已经确认未注册
 */
func (this *AccountService) RegByPhone(phone string) (user *TUser, err error) {
	user = new(TUser)
	user.Phone = phone
	user.Created = time.Now().Unix()

	err = this.orm.db.Create(user).Error

	return
}

/**
绑定账号名(账号密码登陆)
 */
func (this *AccountService) BindAccount(in *proto.PBindData, user *TUser) (err error) {
	u := new(TUser)

	if errno := validator.Verify.IsAccount(in.Account); errno > 0 {
		err = errcode.GetError(errno)
		return
	}

	// 查找用户
	res := this.orm.db.First(u, " account = ? and id != ?", in.Account, user.Id)
	// 数据已经存在
	if !res.RecordNotFound() && u.Id > 0 {
		err = errcode.GetError(errcode.ErrAccountExist)
	}

	// set phone
	err = this.orm.db.Model(user).Update("account", in.Account).Error

	return
}

/**
绑定手机号(手机号登陆)
-- 已绑定
 */
func (this *AccountService) BindPhone(in *proto.PBindData, user *TUser) (err error) {
	u := new(TUser)
	if errno := validator.Verify.IsPhone(in.Phone); errno > 0 {
		err = errcode.GetError(errno)
		return
	}

	// 验证手机验证码
	if !smsUtil.VerifyCode(in.Phone, in.VerifyCode, sms.SmsCodeTypeBind, true) {
		err = errcode.GetError(errcode.ErrSmsVerifyCodeCheck)
		return
	}

	// 查找用户
	res := this.orm.db.First(u, " phone = ? and id != ?", in.Phone, user.Id)
	// 数据已经存在
	if !res.RecordNotFound() && u.Id > 0 {
		err = errcode.GetError(errcode.ErrAccountBindExistPhone)
		return
	}

	// set phone
	user.Phone = in.Phone
	err = this.orm.db.Model(user).Update("phone", in.Phone).Error

	return
}

/**
绑定qq
 */
func (this *AccountService) BindQQ(in *proto.PBindData, user *TUser) (err error) {
	// 检查AccessToken
	if len(in.QqAccessToken) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "QqAccessToken")
		return
	}

	// 检查Openid
	if len(in.QqOpenid) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "QqOpenid")
		return
	}

	var auth *AuthLogin = new(AuthLogin)

	// 查找用户
	if !auth.AuthQQ(in.QqOpenid, in.QqAccessToken) {
		err = errcode.GetError(errcode.ErrAccountVerifyQQ)
		return
	}

	// 绑定授权信息
	return this.saveAuthBind(user.Id, auth)
}

/**
绑定微信
 */
func (this *AccountService) BindWeChat(in *proto.PBindData, user *TUser) (err error) {
	// 检查账户
	if len(in.WxCode) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "QqAccessToken")
		return
	}

	var auth *AuthLogin = new(AuthLogin)
	// 查找用户
	if !auth.AuthWechat(in.WxCode) {
		err = errcode.GetError(errcode.ErrAccountVerifyWechat)
		return
	}

	// 绑定授权信息
	return this.saveAuthBind(user.Id, auth)
}

func (this *AccountService) saveAuthBind(uid int, auth *AuthLogin) (err error) {
	// 查找用户
	u := new(TUserAuth)
	res := this.orm.db.First(u, " auth_site=? and (auth_openid=? or auth_union_id=?)", auth.Auth.AuthSite, auth.Auth.AuthOpenid, auth.Auth.AuthUnionID)
	// 已经绑定其他账号了
	if !res.RecordNotFound() && u.Id > 0 && u.UserId != uid {
		err = errcode.GetError(errcode.ErrAccountBindExist)
		return
	}

	// 以前绑定过
	if u.UserId > 0 {
		return
	}

	// 插入数据
	userAuth := &auth.Auth
	userAuth.UserId = uid
	userAuth.Created = time.Now().Unix()

	// 插入openid一条信息
	if e := this.orm.db.Save(userAuth).Error; e != nil {
		common.Logs.Info(e)
		err = errcode.GetError(errcode.ErrAccountBindWechat)
	}

	return
}
