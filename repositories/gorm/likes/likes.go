package likes

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.ILikeRepository {
	return &DBGorm{
		DB: db,
	}
}

func (db DBGorm) GetLikeByUserAndPostId(userId int, postId int) (models.Like, error) {
	var like models.Like

	err := db.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like).Error
	if err != nil {
		return models.Like{}, err
	}

	return like, nil
}

func (db DBGorm) SaveNewLike(like models.Like) error {
	err := db.DB.Create(&like).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) SaveLike(like models.Like) error {
	err := db.DB.Save(&like).Error
	if err != nil {
		return err
	}

	return nil
}
