package redisclient

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func GetClient() *redis.Client {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	database := viper.GetString("redis.database")

	databaseInt, err := strconv.Atoi(database)
	if err != nil {
		databaseInt = 0
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       databaseInt,
	})

	return rdb

}
