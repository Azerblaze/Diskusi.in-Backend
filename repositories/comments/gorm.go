package comments

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) ICommentDatabase {
	return &GormSql{
		DB: db,
	}
}

func (db GormSql) SaveNewComment(comment models.Comment) error {
	err := db.DB.Create(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllCommentByPost(id int) ([]models.Comment, error) {
	var comments []models.Comment

	err := db.DB.Where("post_id = ?", id).Order("created_at DESC").Preload("User").Find(&comments).Error
	if err != nil {
		return []models.Comment{}, err
	}

	return comments, nil
}

func (db GormSql) GetCommentById(co int) (models.Comment, error) {
	var comment models.Comment

	err := db.DB.Where("id = ?", co).First(&comment).Error
	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (db GormSql) GetCommentByUserId(userId int, page int) ([]models.Comment, error) {
	var comments []models.Comment

	//find topic id
	err := db.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Post").
		Offset((page - 1) * 30).
		Limit(30).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (db GormSql) SaveComment(comment models.Comment) error {
	err := db.DB.Save(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeleteComment(co int) error {
	err := db.DB.Unscoped().Delete(&models.Comment{}, co).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) CountCommentByUserID(userId int) (int, error) {
	var commentCount int64

	err := db.DB.Table("comments").Where("user_id = ?", userId).Count(&commentCount).Error
	if err != nil {
		return 0, err
	}

	return int(commentCount), nil
}
