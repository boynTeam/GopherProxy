package pkg

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Code for database. It may by MySQL or Mongo.
// Author:Boyn
// Date:2020/8/31

var DefaultDB *gorm.DB
var once sync.Once

func InitDB() {
	var err error
	once.Do(func() {
		addr := viper.GetString("Mysql.Addr")
		port := viper.GetString("Mysql.Port")
		user := viper.GetString("Mysql.User")
		password := viper.GetString("Mysql.Password")
		database := viper.GetString("Mysql.Database")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, addr, port, database)
		DefaultDB, err = gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{})
		DefaultDB.Debug()
		DefaultDB.Logger = logger.Default.LogMode(logger.Info)
	})
	if err != nil {
		panic(err)
	}
}
