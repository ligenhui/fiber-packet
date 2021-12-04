package redigo

import (
	"errors"
	"fiber/core/config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
)

var Pool *redis.Pool

func InitRedis() {
	// 使用redigo 默认连接池
	Pool = defaultPool()
}

func Invoke(f func(conn redis.Conn) (reply interface{}, err error)) (reply interface{}, err error) {
	conn := Pool.Get()

	connErr := conn.Err()
	if connErr != nil {
		return nil, connErr
	}

	defer func(conn redis.Conn) {
		if e := recover(); e != nil {
			msg := utils.ToString(e)
			zap.L().Error(msg)
			err = errors.New(msg)
		}
		e := conn.Close()
		if e != nil {
			zap.L().Error("redis conn close err " + e.Error())
		}
	}(conn)

	reply, err = f(conn)
	if err != nil {
		zap.L().Error(err.Error())
	}

	return reply, err
}

func defaultPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", config.Viper.GetString("redis.host"), config.Viper.GetInt("redis.port")),
				redis.DialReadTimeout(config.Viper.GetDuration("redis.readTimeout")),
				redis.DialWriteTimeout(config.Viper.GetDuration("redis.writeTimeout")),
				redis.DialConnectTimeout(config.Viper.GetDuration("redis.connectTimeout")),
				redis.DialDatabase(config.Viper.GetInt("redis.db")),
				redis.DialPassword(config.Viper.GetString("redis.password")),
			)
			if err != nil {
				fmt.Println("redis conn err", err)
				zap.L().Error("redis conn err " + err.Error())
			}
			return conn, err
		},
		MaxIdle:     config.Viper.GetInt("redis.maxIdle"),          //最大空闲数
		MaxActive:   config.Viper.GetInt("redis.maxActive"),        //最大连接数
		IdleTimeout: config.Viper.GetDuration("redis.idleTimeout"), //多少时间后关闭空闲连接
		/*TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},*/
	}
}
