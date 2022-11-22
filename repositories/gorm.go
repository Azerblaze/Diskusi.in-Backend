package repositories

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IDatabase {
	return &GormSql{
		DB: db,
	}
}

// User
func (db GormSql) SaveNewUser(user models.User) error {
	result := db.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (db GormSql) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("username = ?",
			username).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
