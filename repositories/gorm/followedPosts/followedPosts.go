package followedposts

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.IFollowedPostRepository {
	return &DBGorm{
		DB: db,
	}
}
func (db DBGorm) SaveFollowedPost(followedPost models.FollowedPost) error {
	err := db.DB.Create(&followedPost).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetFollowedPost(userId int, postId int) (models.FollowedPost, error) {
	var followedPost models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Where("post_id = ?", postId).First(&followedPost).Error
	if err != nil {
		return models.FollowedPost{}, err
	}

	return followedPost, nil
}

func (db DBGorm) DeleteFollowedPost(followedPostId int) error {
	err := db.DB.Unscoped().Delete(&models.FollowedPost{}, followedPostId).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetAllFollowedPost(userId int) ([]models.FollowedPost, error) {
	var followedPosts []models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Order("created_at DESC").Preload("Post").Find(&followedPosts).Error
	if err != nil {
		return nil, err
	}

	return followedPosts, nil
}
