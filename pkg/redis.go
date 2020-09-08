package pkg

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// Code for redis.
// Author:Boyn
// Date:2020/8/31

var RedisCli *redis.Client

func InitRedisCli() {
	addr := viper.GetString("Redis.Addr")
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}
