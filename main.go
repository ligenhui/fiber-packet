package main

import (
	"fiber/core/config"
	"fiber/core/db"
	"fiber/core/logger"
	"fiber/core/redigo"
	"fiber/route"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func init() {
	//配置文件初始化
	config.Init()
	logger.Init()
	db.Init()
	redigo.InitRedis()
}

func main() {
	fiberConfig := fiber.Config{
		AppName:      config.Viper.GetString("server.name"),
		ReadTimeout:  config.Viper.GetDuration("server.readTimeout"),
		WriteTimeout: config.Viper.GetDuration("server.writeTimeout"),
	}
	app := fiber.New(fiberConfig)

	//init route
	route.Init(app)

	err := app.Listen(config.Viper.GetString("server.address"))
	if err != nil {
		zap.L().Error("Exit outside the project,info:" + err.Error())
	}
}
