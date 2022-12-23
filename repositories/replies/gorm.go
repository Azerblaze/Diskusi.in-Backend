package replies

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IReplyDatabase {
	return &GormSql{
		DB: db,
	}
}

func (db GormSql) SaveNewReply(reply models.Reply) error {
	err := db.DB.Create(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllReplyByComment(commentId int) ([]models.Reply, error) {
	var replys []models.Reply
	err := db.DB.Where("comment_id = ?", commentId).Preload("User").Find(&replys).Error
	if err != nil {
		return []models.Reply{}, err
	}

	return replys, nil
}

func (db GormSql) GetReplyById(id int) (models.Reply, error) {
	var reply models.Reply
	err := db.DB.Where("id = ?", id).First(&reply).Error
	if err != nil {
		return models.Reply{}, err
	}

	return reply, nil
}

func (db GormSql) SaveReply(reply models.Reply) error {
	err := db.DB.Save(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeleteReply(re int) error {
	err := db.DB.Unscoped().Delete(&models.Reply{}, re).Error
	if err != nil {
		return err
	}

	return nil
}
