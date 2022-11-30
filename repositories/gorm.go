package repositories

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IDatabase {
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
		Where("username = ?",
			username).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) Login(email, password string) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("email = ? AND password = ?",
			email, password).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// Topic -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) GetAllTopics() ([]models.Topic, error) {
	var topics []models.Topic

	result := db.DB.Find(&topics)

	if result.Error != nil {
		return nil, result.Error
	} else {
		if result.RowsAffected <= 0 {
			return nil, result.Error
		} else {
			return topics, nil
		}
	}
}

func (db GormSql) GetTopicByName(name string) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("name = ?", name).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db GormSql) GetTopicByID(id int) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("id = ?", id).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db GormSql) SaveNewTopic(topic models.Topic) error {
	result := db.DB.Create(&topic)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db GormSql) SaveTopic(topic models.Topic) error {
	err := db.DB.Where("id = ?", topic.ID).Save(&topic)
	if err != nil {
		return err.Error
	}
	return nil
}

func (db GormSql) RemoveTopic(id int) error {
	err := db.DB.Delete(&models.Topic{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

// Post -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveNewPost(post models.Post) error {
	err := db.DB.Create(&post).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllPostByTopic(id int) ([]models.Post, error) {
	var posts []models.Post

	//find topic id
	err := db.DB.Where("topic_id = ?", id).Preload("User").Preload("Topic").Find(&posts).Error
	if err != nil {
		return []models.Post{}, nil
	}

	return posts, nil
}

func (db GormSql) GetPostById(id int) (models.Post, error) {
	var post models.Post

	err := db.DB.Where("id = ?", id).Preload("User").Preload("Topic").First(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (db GormSql) SavePost(post models.Post) error {
	err := db.DB.Save(&post).Error

	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeletePost(id int) error {
	err := db.DB.Delete(&models.Post{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetPostByIdWithAll(id int) (models.Post, error) {
	var post models.Post
	err := db.DB.Model(&models.Post{}).Where("id = ?", id).Preload("Comments").Find(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

// Comment -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveNewComment(comment models.Comment) error {
	err := db.DB.Create(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllCommentByPost(id int) ([]models.Comment, error) {
	var comments []models.Comment

	err := db.DB.Where("post_id = ?", id).Preload("User").Find(&comments).Error
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

func (db GormSql) SaveComment(comment models.Comment) error {
	err := db.DB.Save(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeleteComment(co int) error {
	err := db.DB.Delete(&models.Comment{}, co).Error
	if err != nil {
		return err
	}

	return nil
}

// Reply -------------------------------------------------------------------------------------------------------------------------------------------------
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

func (db GormSql) GetReplyById(re int) (models.Reply, error) {
	var reply models.Reply
	err := db.DB.Where("id = ?", re).Find(&reply).Error
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
	err := db.DB.Delete(&models.Reply{}, re).Error
	if err != nil {
		return err
	}

	return nil
}
