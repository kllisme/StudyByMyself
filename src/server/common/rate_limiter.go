package common

func SetupRateLimiter() {

	// addr := viper.GetString("resource.redis.rate.redis.addr")
	// password := viper.GetString("resource.redis.rate.redis.password")
	// database := viper.GetString("resource.redis.rate.redis.database")
	// prefix := viper.GetString("resource.redis.rate.redis.prefix")
	// maxRetry := viper.GetInt("resource.redis.rate.redis.maxRetry")
	// _rate := viper.GetString("server.rate.value")

	// rate, err := limiter.NewRateFromFormatted(_rate)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// pool := redis.NewPool(func() (redis.Conn, error) {
	// 	c, err := redis.Dial("tcp", addr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if password != "" {
	// 		if _, err := c.Do("AUTH", password); err != nil {
	// 			c.Close()
	// 			return nil, err
	// 		}
	// 	}
	// 	if c.Do("SELECT", database); err != nil {
	// 		c.Close()
	// 		return nil, err
	// 	}
	// 	return c, err
	// }, 100)

	// store, err := limiter.NewRedisStoreWithOptions(
	// 	pool,
	// 	limiter.StoreOptions{Prefix: prefix, MaxRetry: maxRetry})

	// if err != nil {
	// 	panic(err.Error())
	// }

	// iris.UseFunc(middleware.NewRateLimiter(limiter.NewLimiter(store, rate)).MiddlewareFunc)
}
