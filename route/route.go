package route

import (
	"fiber/core/db"
	"fiber/core/redigo"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
)

type User struct {
	Id     int
	Name   string
	Mobile string
}

func Init(app *fiber.App) {
	//app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		user := &User{}
		_ = db.Db.Debug().Last(user)
		reply, err := redigo.Invoke(func(conn redis.Conn) (reply interface{}, err error) {
			return conn.Do("GET", "key1")
		})
		fmt.Println(reply, "---", err)
		if err != nil {
			return c.Send([]byte(err.Error()))
		}
		//s, err := redis.String(redigo.Redigo.Do("GET", "foo"))
		//fmt.Println(s,err)
		return c.JSON(user)
	})
}
