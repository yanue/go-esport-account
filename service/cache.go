/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-redis操作

------------------------------- go ---------------------------------*/

package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/util"
)

type AccountCache struct {
	redis *redis.Client
	key   *CRedisKey
}

func (this *AccountCache) initRedis(redisAddr, redisPass string) {
	// 建立连接
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
		PoolSize: 300, // 连接池大小
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		panic("redis连接失败:" + err.Error())
	}

	common.Logs.Info("redis connected.")
	// 设置
	this.redis = client
	//common.Logs.Info("client.PoolStats: Hits,Misses,Timeouts,TotalConns,IdleConns,StaleConns = ", this.redis.PoolStats())
}

/**
 获取用户信息
 */
func (this *AccountCache) GetUserInfo(uid int) (user *TUser, err error) {
	user = new(TUser)
	key := this.key.HUserInfo(uid)
	fmt.Println("key", key)
	//util.Proto2Byte()
	res, err := this.redis.HGetAll(key).Result()
	if err != nil {
		return user, err
	}

	if len(res) == 0 {
		return user, errors.New("未找到数据")
	}

	data := make(map[string]interface{}, 0)
	for key, value := range res {
		data[key] = value
	}

	err = util.Struct.MapToStruct(data, user)
	// 解析到结构体
	if err != nil {
		return user, err
	}

	return user, nil
}

/**
 获取用户信息
 */
func (this *AccountCache) SetUserInfo(uid int) (user *TUser, err error) {
	user = new(TUser)
	key := this.key.HUserInfo(uid)

	res, err := this.redis.HGetAll(key).Result()
	if err != nil {
		return user, err
	}

	if len(res) == 0 {
		return user, errors.New("未找到数据")
	}

	// 解析到结构体
	err = mapstructure.Decode(res, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

/**
 设置用户token信息
 */
func (this *AccountCache) SetUserToken(uid int, payload string) (err error) {
	key := this.key.SUserToken(uid)

	// 0 不过期
	_, err = this.redis.Set(key, payload, 0).Result()
	return err
}

/**
 获取用户token信息
 */
func (this *AccountCache) GetUserToken(uid int) (token string, err error) {
	key := this.key.SUserToken(uid)

	// 0 不过期
	return this.redis.Get(key).Result()
}
