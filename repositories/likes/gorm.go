package likes

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) ILikeDatabase {
	return &GormSql{
		DB: db,
	}
}

func (db GormSql) GetLikeByUserAndPostId(userId int, postId int) (models.Like, error) {
	var like models.Like

	err := db.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like).Error
	if err != nil {
		return models.Like{}, err
	}

	return like, nil
}

func (db GormSql) SaveNewLike(like models.Like) error {
	err := db.DB.Create(&like).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) SaveLike(like models.Like) error {
	err := db.DB.Save(&like).Error
	if err != nil {
		return err
	}

	return nil
}
