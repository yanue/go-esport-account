/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountRpc.go
 Time    : 2018/9/11 15:20
 Author  : yanue
 Desc    : account微服务rpc接口

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/micro/go-micro/metadata"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/errcode"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/sms"
	"github.com/yanue/go-esport-common/util"
	"github.com/yanue/go-esport-common/validator"
	"golang.org/x/net/context"
)

/* 微服务接口 */
type AccountRpc struct {
	proto.AccountHandler
	account *AccountService
}

const (
	// 0尚未提交资料,1已提交资料,2审核通过,3审核失败
	IdentityStatusInitial  = iota // 初始状态
	IdentityStatusChecking        // 已提交资料,待审核
	IdentityStatusOk              // 审核已通过
	IdentityStatusFail            // 审核失败
)

/**
@note 验证jwt token
 */
func (this *AccountRpc) verifyToken(ctx context.Context) (token *proto.PJwtToken, err error) {
	// 获取token
	md, ok := metadata.FromContext(ctx)
	tokenStr, ok := md["authorization"]
	if !ok {
		err = errcode.GetError(errcode.ErrInvalidParam, "authorization")
		return
	}

	// 验证token
	jwt, payload1, err := util.JwtToken.Verify(tokenStr)
	if err != nil {
		common.Logs.Info("Verify JwtToken err:", err.Error())
		err = errcode.GetError(errcode.ErrAccountTokenVerify)
		return
	}

	// 验证redis内token信息
	payload2, err := this.account.cache.GetUserToken(int(jwt.Payload.Uid))
	if err != nil {
		common.Logs.Info("Verify JwtToken GetUserToken:", err.Error())
		err = errcode.GetError(errcode.ErrAccountTokenGet)
		return
	}

	if payload1 != payload2 {
		err = errcode.GetError(errcode.ErrAccountTokenNotEqual)
		return
	}

	// 返回token
	return jwt.PJwtToken, nil
}

/**登陆*/
func (this *AccountRpc) Login(ctx context.Context, in *proto.PLoginData, out *proto.PUserAndToken) error {
	// 获取ip
	md, ok := metadata.FromContext(ctx)
	ipStr := ""
	if ok {
		ip, _ := md[":authority"]
		ipStr = ip
	}

	// 处理登陆逻辑
	user, err := this.account.Login(in, ipStr)
	if err != nil {
		common.Logs.Debug("login err=", err.Error(), user)
		return err
	}

	// 不能直接赋值
	// 取地址
	*out = *user

	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetAccountInfo(ctx context.Context, in *proto.PNoParam, out *proto.PUser) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	user, err := this.account.GetUserInfo(int(token.Payload.Uid))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 取地址
	*out = proto.PUser{
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

	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetUserInfo(ctx context.Context, in *proto.PInt32, out *proto.PUser) error {
	user, err := this.account.GetUserInfo(int(in.Val))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 取地址
	*out = proto.PUser{
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

	return nil
}

/**
获取验证码
 */
func (this *AccountRpc) SendSmsVerifyCode(ctx context.Context, in *proto.PSmsData, out *proto.PNoResponse) error {
	if errCode := validator.Verify.IsPhone(in.Phone); errCode > 0 {
		return errcode.GetError(errCode)
	}

	if len(in.Imei) < 6 {
		return errcode.GetError(errcode.ErrInvalidParam, "imei")
	}

	codeType := sms.CodeType(in.CodeType.String())
	if errCode := smsUtil.SendCode(in.Phone, codeType, in.Imei); errCode > 0 {
		return errcode.GetError(errCode)
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
账号解绑(仅支持解绑第三方qq,微信),对应ELoginType
 */
func (this *AccountRpc) Unbind(ctx context.Context, in *proto.PString, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	var authSite string

	switch in.Val {
	// qq
	case proto.ELoginType_QQ.String():
		authSite = "qq"
		// wx
	case proto.ELoginType_WECHAT.String():
		authSite = "wx"
	default:
		return errcode.GetError(errcode.ErrInvalidParam, "PString")
	}

	if e := this.account.orm.db.Where("user_id = ? and auth_site=?", token.Payload.Uid, authSite).Delete(&TUserAuth{}).Error; e != nil {
		common.Logs.Info(e)
		return errcode.GetError(errcode.ErrCustomMsg, "账号解绑失败")
	}

	return err
}

// 绑定账号(手机号,qq,微信)
// -- 已经通过某种方式登陆
func (this *AccountRpc) Bind(ctx context.Context, in *proto.PBindData, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	user, err := this.account.GetUserInfo(int(token.Payload.Uid))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	switch in.BindType {
	// 手机验证码登陆方式
	case proto.ELoginType_PHONE:
		err = this.account.BindPhone(in, user)
		// 手机验证码登陆方式
	case proto.ELoginType_QQ:
		err = this.account.BindQQ(in, user)
		// 手机验证码登陆方式
	case proto.ELoginType_WECHAT:
		err = this.account.BindWeChat(in, user)
	default:
		err = errcode.GetError(errcode.ErrInvalidParam, "BindType")
	}

	return err
}

/**
设置账户名
 */
func (this *AccountRpc) SetAccountName(ctx context.Context, in *proto.PString, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	user, err := this.account.GetUserInfo(int(token.Payload.Uid))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 设置过账户名了
	if len(user.Account) > 0 {
		return errcode.GetError(errcode.ErrCustomMsg, "账号名已经设置过了,不能修改")
	}

	uid := token.Payload.Uid
	if errCode := validator.Verify.IsAccount(in.Val); errCode > 0 {
		return errcode.GetError(errCode)
	}

	u := new(TUser)
	// 查找用户
	res := this.account.orm.db.First(u, " account=? and id!=?", in.Val, uid)
	// 数据已经存在
	if !res.RecordNotFound() && u.Id > 0 {
		return errcode.GetError(errcode.ErrAccountExist)
	}

	err = this.account.orm.db.Model(user).Update("account", in.Val).Error
	if err != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "账户名设置失败")
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
设置昵称
 */
func (this *AccountRpc) ChangeNickname(ctx context.Context, in *proto.PString, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	user, err := this.account.GetUserInfo(int(token.Payload.Uid))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// todo 设置昵称次数限制

	uid := token.Payload.Uid
	if errCode := validator.Verify.IsNickname(in.Val); errCode > 0 {
		return errcode.GetError(errCode)
	}

	u := new(TUser)
	// 查找用户
	res := this.account.orm.db.First(u, " name=? and id!=?", in.Val, uid)
	// 数据已经存在
	if !res.RecordNotFound() && u.Id > 0 {
		return errcode.GetError(errcode.ErrAccountNicknameConflict)
	}

	err = this.account.orm.db.Model(user).Update("name", in.Val).Error
	if err != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "昵称设置失败")
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
设置密码
-- 参数PKeyValList,键值对list
  key -> val:
  code - 验证码
  password - 密码(rsa加密过)
 */
func (this *AccountRpc) SetPassword(ctx context.Context, in *proto.PKeyValList, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	uid := int(token.Payload.Uid)

	user, err := this.account.GetUserInfo(uid)
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	phone := user.Phone

	if len(phone) == 0 {
		return errcode.GetError(errcode.ErrAccountPhoneNotExist, err.Error())
	}

	// 验证码,密码
	var code, password string

	for _, kv := range in.List {
		// 取出验证码
		if kv.Key == "code" {
			code = kv.Val
		}

		// 取出密码
		if kv.Key == "password" {
			password = kv.Val
		}
	}

	if len(password) == 0 {
		return errcode.GetError(errcode.ErrInvalidParam, "password")
	}

	if len(code) == 0 {
		return errcode.GetError(errcode.ErrInvalidParam, "code")
	}

	// 密码解密-私钥解密(客户端公钥加密)
	pass, err := util.Rsa.RsaDecryptPrivate(password)
	if err != nil {
		return errcode.GetError(errcode.ErrInvalidPassword)
	}

	// 校验验证码
	if smsUtil.VerifyCode(user.Phone, code, sms.SmsCodeTypeResetPass, true) {
		return errcode.GetError(errcode.ErrSmsVerifyCodeCheck)
	}

	// 验证密码格式
	if errno := validator.Verify.IsPassword(pass); errno > 0 {
		return errcode.GetError(errno)
	}

	// 密码加密
	passEncrypt, _ := util.Password.Generate(pass)
	// 更新密码
	if e := this.account.orm.db.Model(user).Update("password", passEncrypt).Error; e != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "密码设置失败")
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
重置密码
-- 参数PKeyValList,键值对list
  key -> val:
  code - 验证码
  password - 密码(rsa加密过)
 */
func (this *AccountRpc) ResetPassword(ctx context.Context, in *proto.PKeyValList, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	uid := int(token.Payload.Uid)

	user, err := this.account.GetUserInfo(uid)
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	phone := user.Phone

	if len(phone) == 0 {
		return errcode.GetError(errcode.ErrAccountPhoneNotExist, err.Error())
	}

	// 验证码,密码
	var code, password string

	for _, kv := range in.List {
		// 取出验证码
		if kv.Key == "code" {
			code = kv.Val
		}

		// 取出密码
		if kv.Key == "password" {
			password = kv.Val
		}
	}

	if len(password) == 0 {
		return errcode.GetError(errcode.ErrInvalidParam, "password")
	}

	if len(code) == 0 {
		return errcode.GetError(errcode.ErrInvalidParam, "code")
	}

	// 密码解密-私钥解密(客户端公钥加密)
	pass, err := util.Rsa.RsaDecryptPrivate(password)
	if err != nil {
		return errcode.GetError(errcode.ErrInvalidPassword)
	}

	// 校验验证码
	if smsUtil.VerifyCode(user.Phone, code, sms.SmsCodeTypeResetPass, true) {
		return errcode.GetError(errcode.ErrSmsVerifyCodeCheck)
	}

	// 验证密码格式
	if errno := validator.Verify.IsPassword(pass); errno > 0 {
		return errcode.GetError(errno)
	}

	// 密码加密
	passEncrypt, _ := util.Password.Generate(pass)
	// 更新密码
	if e := this.account.orm.db.Model(user).Update("password", passEncrypt).Error; e != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "密码设置失败")
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
头像设置
 */
func (this *AccountRpc) SetAvatar(ctx context.Context, in *proto.PString, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	uid := int(token.Payload.Uid)

	user, err := this.account.GetUserInfo(uid)
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	if errno := validator.Verify.IsUrl(in.Val); errno > 0 {
		return errcode.GetError(errno)
	}

	if e := this.account.orm.db.Model(user).Update("avatar", in.Val).Error; e != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "更新头像失败")
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
保存资料
 */
func (this *AccountRpc) SaveProfile(ctx context.Context, in *proto.PKeyValList, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	uid := int(token.Payload.Uid)

	user, err := this.account.GetUserInfo(uid)
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	cols := make(map[string]interface{}, 0)

	for _, kv := range in.List {
		if kv.Key == "gender" {
			cols["gender"] = util.ToInt(kv.Val)
		}
		if kv.Key == "schoolId" {
			cols["schoolId"] = util.ToInt(kv.Val)
		}
		if kv.Key == "classId" {
			cols["classId"] = util.ToInt(kv.Val)
		}
		if kv.Key == "areaId" {
			cols["areaId"] = util.ToInt(kv.Val)
		}
	}

	if len(cols) > 0 {
		if e := this.account.orm.db.Model(user).Update(cols); e != nil {
			return errcode.GetError(errcode.ErrCustomMsg, "更新个人信息失败")
		}
	}

	*out = proto.PNoResponse{}

	return nil
}

/**
提交身份认证资料
-- 参数PKeyValList,键值对list
  key -> val:
  identity_no - 身份证号
  identity_name - 真实姓名
  identity_pic - 身份证照片
 */
func (this *AccountRpc) SubmitIdentity(ctx context.Context, in *proto.PKeyValList, out *proto.PNoResponse) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	uid := int(token.Payload.Uid)

	user, err := this.account.GetUserInfo(uid)
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 认证状态:0尚未提交资料,1已提交资料,2审核通过,3审核失败
	if user.IdentityStatus == IdentityStatusOk {
		return errcode.GetError(errcode.ErrAccountVerificationOk, err.Error())
	}

	cols := make(map[string]interface{}, 0)

	var identityNo, identityName, identityPic string

	for _, kv := range in.List {
		if kv.Key == "identity_no" {
			identityNo = kv.Val
		}
		if kv.Key == "identity_name" {
			identityName = kv.Val
		}
		if kv.Key == "identity_pic" {
			identityPic = kv.Val
		}
	}

	// 身份证号码
	if errno := validator.Verify.IsIdCard(identityNo); errno > 0 {
		return errcode.GetError(errno)
	}

	// 真实姓名
	if errno := validator.Verify.IsRealName(identityName); errno > 0 {
		return errcode.GetError(errno)
	}

	// 图片
	if errno := validator.Verify.IsUrl(identityPic); errno > 0 {
		return errcode.GetError(errno)
	}

	cols["identity_no"] = identityNo
	cols["identity_name"] = identityName
	cols["identity_pic"] = identityPic
	cols["identity_status"] = IdentityStatusChecking // 状态为已提交资料

	if e := this.account.orm.db.Model(user).Update(cols); e != nil {
		return errcode.GetError(errcode.ErrCustomMsg, "提交认证资料失败")
	}

	*out = proto.PNoResponse{}

	return nil
}
