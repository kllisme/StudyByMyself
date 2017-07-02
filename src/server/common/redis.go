package common

import (
	"github.com/spf13/viper"
	"gopkg.in/redis.v5"
	"time"
)

func SetupRedis() {

	addr := viper.GetString("resource.redis.default.addr")
	password := viper.GetString("resource.redis.default.password")
	database := viper.GetInt("resource.redis.default.database")
	maxActive := viper.GetInt("resource.redis.default.maxActive")
	idleTimeout := time.Duration(viper.GetInt("resource.redis.default.idleTimeout")) * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          database,
		MaxRetries:  3,
		IdleTimeout: idleTimeout,
		PoolSize:    maxActive,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic("failed to connect redis:" + err.Error())
	}

	Redis = client

}

var (
	Redis *redis.Client
)
