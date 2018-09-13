/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-redis操作

------------------------------- go ---------------------------------*/

package service

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/yanue/go-esport-common"
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
