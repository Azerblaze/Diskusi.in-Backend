package configs

import (
	"discusiin/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	DB_User := Cfg.DB_USER
	DB_Password := Cfg.DB_PASSWORD
	DB_Host := Cfg.DB_HOST
	DB_Port := Cfg.DB_PORT
	DB_Name := Cfg.DB_NAME

	connectionString := fmt.
		Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			DB_User,
			DB_Password,
			DB_Host,
			DB_Port,
			DB_Name)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.Topic{},
		&models.Post{},
		&models.Comment{},
		&models.Reply{},
	)

	if err != nil {
		panic(err)
	}

	log.Print("Init DB Done")
}
