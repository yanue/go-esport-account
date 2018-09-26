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
	"github.com/yanue/go-esport-common/util"
	"github.com/yanue/go-esport-common/validator"
)

/**
账号登陆
 */
func (this AccountService) Login(in *proto.PLoginData) (out *proto.PUserAndToken, err error) {
	user := new(TUser)
	out = new(proto.PUserAndToken)

	// os参数
	if len(in.Os.String()) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "os")
		return
	}

	// android ios 需要设备唯一码
	if in.Os != common.OS_WEB && len(in.DeviceId) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "DeviceId")
		return
	}

	switch in.LoginType {
	// 账号密码登陆
	case proto.ELoginType_ACCOUNT:
		user, err = this.loginByAccount(in)
		// 手机验证码登陆方式
	case proto.ELoginType_PHONE:
		user, err = this.loginByPhoneCode(in)
		// 手机验证码登陆方式
	case proto.ELoginType_QQ:
		user, err = this.loginByPhoneCode(in)
		// 手机验证码登陆方式
	case proto.ELoginType_WECHAT:
		user, err = this.loginByPhoneCode(in)
	default:
		err = errcode.GetError(errcode.ErrInvalidParam, "LoginType")
	}

	// 出现错误
	if err != nil {
		return
	}

	// 获取原来登陆设备信息(如果已登陆)
	var deviceId string
	if payloadOldStr, err := this.cache.GetUserToken(user.Id); err == nil {
		if payloadInfo, err := util.JwtToken.ParsePayload(payloadOldStr); err == nil {
			deviceId = payloadInfo.DeviceId
		}
	}

	// token生成
	token, payload, err := util.JwtToken.Generate(user.Id, in.Os, in.LoginType, in.DeviceId)
	if err != nil {
		err = errcode.GetError(errcode.ErrAccountTokenGenerate)
		return
	}

	out.Token = token
	out.User = &proto.PUser{
		Id:             int32(user.Id),
		Name:           user.Name,
		Phone:          user.Phone,
		Email:          user.Email,
		Gender:         int32(user.Gender),
		SchoolId:       int32(user.SchoolId),
		ClassId:        int32(user.ClassId),
		AreaId:         int32(user.AreaId),
		IdentityStatus: int32(user.IdentityStatus),
		Created:        int32(user.Created),
	}

	// 保存token
	err = this.cache.SetUserToken(user.Id, payload)
	if err != nil {
		common.Logs.Info("SetUserToken err", err.Error())
		err = errcode.GetError(errcode.ErrCustomMsg, "保存用户token失败!")
		return
	}

	// 不同设备登陆,踢原下线
	if deviceId != in.DeviceId {
		// todo 发送通知到旧手机,踢出下线
		common.Logs.Info("send offline", deviceId)
	}

	return out, nil
}

/**
账号登陆方式
 */
func (this AccountService) loginByAccount(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

	// 检查账户
	if code := validator.Verify.IsAccount(in.Account); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 密码解密-私钥解密(客户端公钥加密)
	pass, err := util.Rsa.RsaDecryptPrivate(in.Password)
	if err != nil {
		return
	}

	// 检查密码格式
	if code := validator.Verify.IsPassword(pass); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 查找用户
	if err = this.orm.db.First(user, " account = ? ", in.Account).Error; err != nil {
		err = errcode.GetError(errcode.ErrAccountNotExist)
		return
	}

	// 密码长度60
	if len(user.Password) != 60 {
		err = errcode.GetError(errcode.ErrAccountPassNotSet)
		return
	}

	// 验证密码
	if !util.Password.Verify(pass, user.Password) {
		err = errcode.GetError(errcode.ErrAccountPassIncorrect)
		return
	}

	// 设置session
	return user, nil
}

/**
手机验证码方式登陆
 */
func (this AccountService) loginByPhoneCode(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

	// 检查账户
	if code := validator.Verify.IsPhoneWithoutCode(in.Phone); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 查找用户
	if err = this.orm.db.First(user, " phone = ? ", in.Phone).Error; err != nil {
		err = errcode.GetError(errcode.ErrAccountNotExist)
		return
	}

	// 验证手机验证码
	if !smsUtil.VerifyCode(in.Phone, in.VerifyCode, sms.SmsCodeTypeQuickLogin, true) {
		err = errcode.GetError(errcode.ErrVerifyCodeCheck)
		return
	}

	// 设置session
	return user, nil
}

/**
qq登陆
 */
func (this AccountService) loginByQQ(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

	// 检查账户
	if code := validator.Verify.IsPhoneWithoutCode(in.Phone); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 查找用户
	if err = this.orm.db.First(user, " phone = ? ", in.Phone).Error; err != nil {
		err = errcode.GetError(errcode.ErrAccountNotExist)
		return
	}

	// 设置session
	return user, nil
}

/**
微信登陆
 */
func (this AccountService) loginByWeChat(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

	// 检查账户
	if code := validator.Verify.IsPhoneWithoutCode(in.Phone); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 查找用户
	if err = this.orm.db.First(user, " phone = ? ", in.Phone).Error; err != nil {
		err = errcode.GetError(errcode.ErrAccountNotExist)
		return
	}

	// 验证手机验证码 todo
	sms.PrefixBindVercode

	// 设置session
	return user, nil
}
