package replies

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.IReplyRepository {
	return &DBGorm{
		DB: db,
	}
}

func (db DBGorm) SaveNewReply(reply models.Reply) error {
	err := db.DB.Create(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetAllReplyByComment(commentId int) ([]models.Reply, error) {
	var replys []models.Reply
	err := db.DB.Where("comment_id = ?", commentId).Preload("User").Find(&replys).Error
	if err != nil {
		return []models.Reply{}, err
	}

	return replys, nil
}

func (db DBGorm) GetReplyById(id int) (models.Reply, error) {
	var reply models.Reply
	err := db.DB.Where("id = ?", id).First(&reply).Error
	if err != nil {
		return models.Reply{}, err
	}

	return reply, nil
}

func (db DBGorm) SaveReply(reply models.Reply) error {
	err := db.DB.Save(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) DeleteReply(re int) error {
	err := db.DB.Unscoped().Delete(&models.Reply{}, re).Error
	if err != nil {
		return err
	}

	return nil
}
