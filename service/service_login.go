/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : login.go
 Time    : 2018/9/25 12:19
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/errcode"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/sms"
	"github.com/yanue/go-esport-common/util"
	"github.com/yanue/go-esport-common/validator"
	"time"
)

/**
账号登陆
 */
func (this *AccountService) Login(in *proto.PLoginData, ip string) (out *proto.PUserAndToken, err error) {
	user := new(TUser)
	out = new(proto.PUserAndToken)

	defer func() {
		// 记录每次登陆信息
		uLogin := new(TUserLogin)

		// login type
		uLogin.LoginType = in.LoginType.String()
		// 设备信息
		uLogin.Device = util.Struct.ToJsonString(in.Device)
		// ip
		uLogin.Ip = ip
		// 登陆时间
		uLogin.Created = time.Now().Unix()
		// login data
		in.Device = nil // set nil
		uLogin.LoginData = util.Struct.ToJsonString(in)
		// 登陆成功,用户id
		uLogin.UserId = user.Id
		// 错误信息
		if err != nil {
			uLogin.ErrMsg = err.Error()
		}

		// 插入数据
		if e := this.orm.db.Create(uLogin).Error; e != nil {
			common.Logs.Info(e)
		}
	}()

	// 设备信息
	if in.Device == nil {
		err = errcode.GetError(errcode.ErrInvalidParam, "Device")
		return
	}

	// os参数
	if len(in.Device.Os.String()) == 0 {
		err = errcode.GetError(errcode.ErrInvalidParam, "os")
		return
	}

	// android ios 需要设备唯一码
	if in.Device.Os != common.OS_WEB && len(in.Device.Imei) == 0 {
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
		user, err = this.loginByQQ(in)
		// 手机验证码登陆方式
	case proto.ELoginType_WECHAT:
		user, err = this.loginByWeChat(in)
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
			fmt.Println("payloadInfo", payloadInfo, payloadInfo.Device)
			deviceId = payloadInfo.Device.Imei
		}
	}

	// token生成
	token, payload, err := util.JwtToken.Generate(user.Id, in.LoginType, in.Device)
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
	if deviceId != in.Device.Imei {
		// todo 发送通知到旧手机,踢出下线
		common.Logs.Info("send offline", deviceId)
	}

	return out, nil
}

/**
账号登陆方式
 */
func (this *AccountService) loginByAccount(in *proto.PLoginData) (user *TUser, err error) {
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
func (this *AccountService) loginByPhoneCode(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

	// 检查账户
	if code := validator.Verify.IsPhoneWithoutCode(in.Phone); code > 0 {
		err = errcode.GetError(code)
		return
	}

	// 验证手机验证码
	if !smsUtil.VerifyCode(in.Phone, in.VerifyCode, sms.SmsCodeTypeQuickLogin, true) {
		err = errcode.GetError(errcode.ErrVerifyCodeCheck)
		return
	}

	// 查找用户
	if err = this.orm.db.First(user, " phone = ? ", in.Phone).Error; err != nil {
		// 未找到数据
		if err == gorm.ErrRecordNotFound {
			// 通过手机号注册
			return this.RegByPhone(in.Phone)
		} else {
			err = errcode.GetError(errcode.ErrAccountNotExist)
		}

		return
	}

	// 设置session
	return user, nil
}

/**
qq登陆
 */
func (this *AccountService) loginByQQ(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

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

	common.Logs.Info("login:", auth)

	// 设置绑定信息
	return this.setBindInfo(auth)
}

/**
微信登陆
 */
func (this *AccountService) loginByWeChat(in *proto.PLoginData) (user *TUser, err error) {
	user = new(TUser)

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

	common.Logs.Info("login:", auth)

	// 设置绑定信息
	return this.setBindInfo(auth)
}

/**
授权成功后,设置绑定信息
1. 检查是否授权过(需要注册)
2. 检查是否其他方式授权过(统一qq不同appkey同一union_id)
 */
func (this *AccountService) setBindInfo(auth *AuthLogin) (user *TUser, err error) {
	user = new(TUser)
	userAuth := new(TUserAuth)

	defer func() {
		if err == nil && userAuth.UserId > 0 {
			user, err = this.GetUserInfo(userAuth.UserId)
			return
		}
	}()

	// 查询注册情况(unionid)
	if e := this.orm.db.First(userAuth, "auth_union_id=?", auth.Auth.AuthUnionID).Error; e != nil && e != gorm.ErrRecordNotFound {
		common.Logs.Info(e)
		err = errcode.GetError(errcode.ErrAccountBindGet)
		return
	}

	// 数据未找到
	if userAuth.Id == 0 {
		// 通过授权信息注册用户
		userAuth, err = this.RegByAuth(auth)
		return
	}

	// 重置数据
	uid := userAuth.UserId
	// id重置为0
	userAuth.Id = 0

	// 通过openid查询
	if e := this.orm.db.First(userAuth, "auth_openid=?", auth.Auth.AuthOpenid).Error; e != nil && e != gorm.ErrRecordNotFound {
		common.Logs.Info(e)
		err = errcode.GetError(errcode.ErrAccountBindGet)
		return
	}

	// openid数据未找到
	if userAuth.Id == 0 {
		// 插入数据
		userAuth = &auth.Auth
		userAuth.UserId = uid
		userAuth.Created = time.Now().Unix()

		// 插入openid一条信息
		if e := this.orm.db.Create(userAuth).Error; e != nil && e != gorm.ErrRecordNotFound {
			common.Logs.Info(e)
			err = errcode.GetError(errcode.ErrAccountBindGet)
		}
	}

	return
}
