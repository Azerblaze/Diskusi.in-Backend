package configs

import (
	"discusiin/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	DBUser := Cfg.DBUser
	DBPassword := Cfg.DBPassword
	DBHost := Cfg.DBHost
	DBPort := Cfg.DBPort
	DBName := Cfg.DBName

	connectionString := fmt.
		Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			DBUser,
			DBPassword,
			DBHost,
			DBPort,
			DBName)

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
		&models.Like{},
		&models.FollowedPost{},
		&models.Bookmark{},
	)

	if err != nil {
		panic(err)
	}

}
