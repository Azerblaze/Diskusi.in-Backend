package posts

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.IPostRepository {
	return &DBGorm{
		DB: db,
	}
}

func (db DBGorm) SaveNewPost(post models.Post) error {
	err := db.DB.Create(&post).Error
	if err != nil {
		return err
	}

	return nil
}
func (db DBGorm) GetRecentPost(page int, search string) ([]models.Post, error) {
	var result []models.Post
	err := db.DB.Where("title LIKE ?", "%"+search+"%").Order("created_at DESC").Preload("User").Preload("Topic").Offset((page - 1) * 20).Limit(20).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db DBGorm) GetAllPostByTopic(id int, page int, search string) ([]models.Post, error) {
	var posts []models.Post

	//find topic id
	err := db.DB.Where("topic_id = ?", id).
		Where("title LIKE ?", "%"+search+"%").
		Order("created_at DESC").
		Preload("User").
		Preload("Topic").
		Offset((page - 1) * 20).
		Limit(20).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (db DBGorm) GetAllPostByTopicByLike(topicID int, page int) ([]models.Post, error) {
	var posts []models.Post

	//find topic with like
	err := db.DB.Where("topic_id = ?", topicID).Order("like_count DESC").Preload("User").Preload("Topic").Offset((page - 1) * 20).Limit(20).Find(&posts).Error
	if err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (db DBGorm) GetPostById(id int) (models.Post, error) {
	var post models.Post

	err := db.DB.Where("id = ?", id).
		Preload("User").
		Preload("Topic").
		First(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (db DBGorm) GetPostByUserId(userId int, page int) ([]models.Post, error) {
	var posts []models.Post

	//find topic id
	err := db.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("User").
		Preload("Topic").
		Offset((page - 1) * 20).
		Limit(20).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (db DBGorm) GetAllPostByLike(page int) ([]models.Post, error) {
	var posts []models.Post

	//find topic with like
	err := db.DB.Order("like_count DESC").Preload("User").Preload("Topic").Offset((page - 1) * 20).Limit(20).Find(&posts).Error
	if err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (db DBGorm) SavePost(post models.Post) error {
	err := db.DB.Save(&post).Error

	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) DeletePostByPostID(id int) error {
	err := db.DB.Unscoped().Where("id = ?", id).Delete(&models.Post{}).Error
	if err != nil {
		return err
	}

	return nil
}
func (db DBGorm) DeletePostByUserID(userID int) error {
	err := db.DB.Unscoped().Where("user_id = ?", userID).Delete(&models.Post{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetPostByIdWithAll(id int) (models.Post, error) {
	var post models.Post
	err := db.DB.Model(&models.Post{}).Where("id = ?", id).Preload("Comments").First(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}
func (db DBGorm) CountPostLike(postID int) (int, error) {
	var postLike int64

	err := db.DB.Table("likes").Where("post_id = ? AND is_like = 1", postID).Count(&postLike).Error
	if err != nil {
		return 0, err
	}

	return int(postLike), nil
}
func (db DBGorm) CountPostDislike(postID int) (int, error) {
	var postDislike int64
	err := db.DB.Table("likes").Where("post_id = ? AND is_dislike = 1", postID).Count(&postDislike).Error
	if err != nil {
		return 0, err
	}
	return int(postDislike), nil
}
func (db DBGorm) CountPostComment(postID int) (int, error) {
	var commentCount int64
	err := db.DB.Table("comments").Where("post_id = ?", postID).Count(&commentCount).Error
	if err != nil {
		return 0, err
	}
	return int(commentCount), nil
}
func (db DBGorm) CountAllPost() (int, error) {
	var numberOfPost int64

	err := db.DB.Table("posts").Count(&numberOfPost).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfPost), nil
}
func (db DBGorm) CountPostByTopicID(topicId int) (int, error) {
	var postCount int64

	err := db.DB.Table("posts").Where("topic_id = ?", topicId).Count(&postCount).Error
	if err != nil {
		return 0, err
	}

	return int(postCount), nil
}
func (db DBGorm) CountPostByUserID(userId int) (int, error) {
	var postCount int64

	err := db.DB.Table("posts").Where("user_id = ?", userId).Count(&postCount).Error
	if err != nil {
		return 0, err
	}

	return int(postCount), nil
}
