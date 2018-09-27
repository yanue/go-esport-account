/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : third_login.go
 Time    : 2018/9/27 15:49
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"encoding/json"
	"fmt"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/util"
	"strings"
)

// QQ认证
const (
	qq_appid       string = "1105704814"
	qq_auth_url    string = "https://graph.qq.com/user/get_simple_userinfo?oauth_consumer_key=%v&access_token=%s&openid=%v"
	qq_unionid_url string = "https://graph.qq.com/oauth2.0/me?access_token=%v&unionid=1"
)

// 微信认证
const (
	bb_weixin_appid  string = "wx8124a3c5700d9ef8"
	bb_weixin_secret string = "33a3be5c307b039738915a57fb2958e2"

	bb_weixin_appid_gzh  string = "wx84f53ebe25d191e6"
	bb_weixin_secret_gzh string = "1ccc8c67d2bb630e37a0c6a746606e27"

	wx_auth_url string = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	wx_user_url string = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
)

// 用户登录信息
type AuthLogin struct {
	User     TUser
	Auth     TUserAuth
	WxPlType int64 //0: 微信开发平台   1:微信公众号
}

func (this *AuthLogin) AuthQQ(openid, accessToken string) bool {
	this.Auth.AuthOpenid = openid

	// 链接地址
	strUrl := fmt.Sprintf(qq_auth_url, qq_appid, accessToken, openid)

	// 发送请求,获取用户基本信息
	data := util.Http.RemoteCall(strUrl)
	if data == nil {
		return false
	}

	// 解析基本信息
	var info map[string]interface{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		common.Logs.Info("json.Unmarshal fail:" + err.Error())
		return false
	}

	common.Logs.Debug("---------------------qq info : %v", info)

	// 获取昵称
	val, ok := info["nickname"]
	if !ok {
		common.Logs.Info("get userInfo fail:" + string(data))
		return false
	}
	this.User.Name = val.(string)

	// 获取性别
	val, ok = info["gender"]
	if ok {
		if val.(string) == "男" {
			this.User.Gender = 1
		}
		if val.(string) == "女" {
			this.User.Gender = 2
		}
	}

	// 获取头像
	avatar, ok := info["figureurl_qq_2"]
	if !ok {
		common.Logs.Info("not found figureurl_qq_2 img, try get figureurl_qq_1")
		avatar, ok = info["figureurl_qq_1"]
		if !ok {
			common.Logs.Info("not found figureurl_qq_1")
		}
	}

	if ok {
		this.User.Avatar = avatar.(string)
	}

	// get unionid
	strUrl = fmt.Sprintf(qq_unionid_url, accessToken)
	data = util.Http.RemoteCall(strUrl)
	if data != nil {
		src := string(data)
		b := strings.Index(src, "{")
		e := strings.Index(src, "}") + 1
		dest := src[b:e]
		//Info("dest=%v", dest)

		var info map[string]interface{}
		err := json.Unmarshal([]byte(dest), &info)
		if err == nil {
			//Info("qq unionid map = %v", info)

			// 获取昵称
			val, ok := info["unionid"]
			if !ok {
				common.Logs.Info("qq unionid map = %v", info)
				common.Logs.Info("get unionid fail:" + string(data))
				return false
			} else {
				this.Auth.AuthUnionID = val.(string)
			}
		} else {
			common.Logs.Info("json.Unmarshal fail(%v), src=%v", err.Error(), dest)
		}
	}

	return true
}

func (this *AuthLogin) AuthWechat(code string) bool {
	// 取access token
	strUrl := ""

	// 0-微信开放平台
	if this.WxPlType == 0 {
		strUrl = fmt.Sprintf(wx_auth_url, bb_weixin_appid, bb_weixin_secret, code)
	} else {
		// 1-微信公众号
		strUrl = fmt.Sprintf(wx_auth_url, bb_weixin_appid_gzh, bb_weixin_secret_gzh, code)
	}

	common.Logs.Debug("AuthWechat url:%v", strUrl)
	data := util.Http.RemoteCall(strUrl)
	if data == nil {
		return false
	}

	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		common.Logs.Info("json.Unmarshal fail:" + err.Error())
		return false
	}

	common.Logs.Debug("wechat info = %v", m)

	val, ok := m["access_token"]
	if !ok {
		// 没找到access_token，表示有错误发生
		common.Logs.Info("AuthWechat get access_token fail:" + string(data))
		return false
	}
	accessToken := val.(string)

	//expire_in := m["expire_in"].(float64)
	//refresh_token := m["refresh_token"].(string)
	openid := m["openid"].(string)
	//scope := m["scope"].(string)

	// 取微信用户信息
	strUrl = fmt.Sprintf(wx_user_url, accessToken, openid)
	data = util.Http.RemoteCall(strUrl)
	if data == nil {
		return false
	}

	var info map[string]interface{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		common.Logs.Info("json.Unmarshal fail:" + err.Error())
		return false
	}

	nickname, ok := info["nickname"]
	if !ok {
		common.Logs.Info("get userinfo fail:" + string(data))
		return false
	}
	this.User.Name = nickname.(string)

	unionid, ok := info["unionid"]
	if !ok {
		common.Logs.Info("get unionid fail:" + string(data))
		return false
	} else {
		this.Auth.AuthOpenid = unionid.(string)
		this.Auth.AuthUnionID = unionid.(string)
	}

	avatar, ok := info["headimgurl"]
	if ok {
		this.User.Avatar = avatar.(string)
	}

	// 微信中 1为男，2为女
	sex := info["sex"].(float64)
	if sex == 2 {
		sex = 0
	}
	this.User.Gender = int(sex)

	return true
}
