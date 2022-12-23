package followedPosts

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IFollowedPostDatabase {
	return &GormSql{
		DB: db,
	}
}

func (db GormSql) SaveFollowedPost(followedPost models.FollowedPost) error {
	err := db.DB.Create(&followedPost).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetFollowedPost(userId int, postId int) (models.FollowedPost, error) {
	var followedPost models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Where("post_id = ?", postId).First(&followedPost).Error
	if err != nil {
		return models.FollowedPost{}, err
	}

	return followedPost, nil
}

func (db GormSql) DeleteFollowedPost(followedPostId int) error {
	err := db.DB.Unscoped().Delete(&models.FollowedPost{}, followedPostId).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllFollowedPost(userId int) ([]models.FollowedPost, error) {
	var followedPosts []models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Order("created_at DESC").Preload("Post").Find(&followedPosts).Error
	if err != nil {
		return nil, err
	}

	return followedPosts, nil
}
