/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : login.go
 Time    : 2018/9/25 12:19
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import "time"

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
绑定qq
 */
func (this *AccountService) BindQQ() {

}

/**
绑定微信
 */
func (this *AccountService) BindWeChat() {

}

/**
绑定账号名(账号密码登陆)
 */
func (this *AccountService) BindAccount() {

}

/**
绑定手机号(手机号登陆)
 */
func (this *AccountService) BindPhone() {

}
