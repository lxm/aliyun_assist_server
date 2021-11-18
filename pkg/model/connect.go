package model

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

func ConnectDB() error {
	var err error
	once.Do(func() {
		err = connect()
	})
	return err
}

func connect() error {
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database)
	logrus.Debugf("dsn:%v", dsn)
	db, err = gorm.Open(mysql.New(
		mysql.Config{
			DSN: dsn,
		}),
		&gorm.Config{
			// Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		logrus.Errorf("ConnectDB error %v", err)
		return err
	}
	return nil

}

func GetDB() *gorm.DB {
	ConnectDB()
	return db
}

func Migrate() {
	db.AutoMigrate(&RegisterInfo{}, &Command{}, &Task{}, &ActivationCode{})
}
