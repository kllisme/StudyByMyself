package common

import (
	"time"

	iris "gopkg.in/kataras/iris.v5"

	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"
	"github.com/spf13/viper"
)

func SetupSession() {

	iris.Config.Sessions.Cookie = viper.GetString("server.session.cookie")
	iris.Config.Sessions.Expires = time.Duration(viper.GetInt("server.session.expires")) * time.Second
	iris.Config.Sessions.DisableSubdomainPersistence = true

	db := redis.New(service.Config{
		Network:       service.DefaultRedisNetwork,
		Addr:          viper.GetString("resource.redis.session.addr"),
		Password:      viper.GetString("resource.redis.session.password"),
		Database:      viper.GetString("resource.redis.session.database"),
		MaxIdle:       viper.GetInt("resource.redis.session.maxIdle"),
		MaxActive:     viper.GetInt("resource.redis.session.maxActive"),
		IdleTimeout:   time.Duration(viper.GetInt("resource.redis.session.idleTimeout")) * time.Second,
		Prefix:        viper.GetString("resource.redis.session.prefix"),
		MaxAgeSeconds: viper.GetInt("resource.redis.session.maxAgeSeconds"),
	})

	SESSION_DB = db
}

var (
	SESSION_DB *redis.Database
)
