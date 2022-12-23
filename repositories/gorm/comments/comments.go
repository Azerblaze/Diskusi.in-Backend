package comments

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.ICommentRepository {
	return &DBGorm{
		DB: db,
	}
}

func (db DBGorm) SaveNewComment(comment models.Comment) error {
	err := db.DB.Create(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetAllCommentByPost(id int) ([]models.Comment, error) {
	var comments []models.Comment

	err := db.DB.Where("post_id = ?", id).Order("created_at DESC").Preload("User").Find(&comments).Error
	if err != nil {
		return []models.Comment{}, err
	}

	return comments, nil
}

func (db DBGorm) GetCommentById(co int) (models.Comment, error) {
	var comment models.Comment

	err := db.DB.Where("id = ?", co).First(&comment).Error
	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (db DBGorm) GetCommentByUserId(userId int, page int) ([]models.Comment, error) {
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

func (db DBGorm) SaveComment(comment models.Comment) error {
	err := db.DB.Save(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) DeleteComment(co int) error {
	err := db.DB.Unscoped().Delete(&models.Comment{}, co).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) CountCommentByUserID(userId int) (int, error) {
	var commentCount int64

	err := db.DB.Table("comments").Where("user_id = ?", userId).Count(&commentCount).Error
	if err != nil {
		return 0, err
	}

	return int(commentCount), nil
}
