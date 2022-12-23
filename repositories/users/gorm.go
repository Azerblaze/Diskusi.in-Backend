package users

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IUserDatabase {
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
		Where("username = ?", username).
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
func (db GormSql) GetUserById(userId int) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("id = ?", userId).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) GetUsersAdminNotIncluded(page int) ([]models.User, error) {
	var users []models.User
	err := db.DB.Where("is_admin = 0").Order("username ASC").Offset((page - 1) * 20).Limit(20).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (db GormSql) GetProfile(id int) (models.User, error) {
	var user models.User
	err := db.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) UpdateProfile(user models.User) error {
	err := db.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (db GormSql) DeleteUser(userId int) error {
	err := db.DB.Unscoped().Delete(&models.User{}, userId).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) CountAllUserNotIncludeDeletedUser() (int, error) {
	var numberOfUser int64

	err := db.DB.Table("users").Where("deleted_at IS NULL").Count(&numberOfUser).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfUser), nil
}

func (db GormSql) CountAllUserNotAdminNotIncludeDeletedUser() (int, error) {
	var userCount int64

	err := db.DB.Table("users").Where("is_admin = 0").Where("deleted_at IS NULL").Count(&userCount).Error
	if err != nil {
		return 0, err
	}
	return int(userCount), nil
}
