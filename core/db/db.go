package db

import (
	"fiber/core/config"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s&writeTimeout=%s&readTimeout=%s",
		config.Viper.GetString("db.user"),
		config.Viper.GetString("db.password"),
		config.Viper.GetString("db.host"),
		config.Viper.GetString("db.port"),
		config.Viper.GetString("db.database"),
		config.Viper.GetDuration("db.connTimeOut"),
		config.Viper.GetDuration("db.writeTimeout"),
		config.Viper.GetDuration("db.readTimeout"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Viper.GetString("db.prefix"),
			SingularTable: true, //表复数禁用
		},
		SkipDefaultTransaction: true, //关闭默认事务
		PrepareStmt:            true, // 开启缓存预编译，可以提高后续的调用速度
		QueryFields:            true, //自动将struct结构体的字段设置为查询字段
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             2 * time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Silent,
		}),
	})

	if err != nil {
		fmt.Printf("DB connect Error %s \n", err)
		zap.L().Fatal("DB connect Error " + err.Error())
	}

	//一个坑，不设置这个参数，gorm会把表名转义后加个s，导致找不到数据库的表
	SqlDB, _ := db.DB()
	// 设置连接池中最大的闲置连接数
	SqlDB.SetMaxIdleConns(config.Viper.GetInt("db.maxIdleConns"))
	// 设置数据库的最大连接数量
	SqlDB.SetMaxOpenConns(config.Viper.GetInt("db.maxOpenConns"))
	// 这是连接的最大可复用时间
	SqlDB.SetConnMaxLifetime(config.Viper.GetDuration("db.connMaxLifetime"))
	Db = db
}
