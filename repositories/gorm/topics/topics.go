package topics

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.ITopicRepository {
	return &DBGorm{
		DB: db,
	}
}

func (db DBGorm) CountAllTopic() (int, error) {
	var numberOfTopic int64

	err := db.DB.Table("topics").Count(&numberOfTopic).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfTopic), nil
}

func (db DBGorm) CountNumberOfPostByTopicName(topicName string) (int, error) {
	var postCount int64

	err := db.DB.Table("posts").Where("topic_id = (SELECT id FROM topics WHERE name = ?)", topicName).Where("deleted_at IS NULL").Count(&postCount).Error
	if err != nil {
		return 0, err
	}

	return int(postCount), nil
}

func (db DBGorm) GetAllTopics() ([]models.Topic, error) {
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
func (db DBGorm) GetTopTopics() ([]dto.TopTopics, error) {

	var results []dto.TopTopics

	err := db.DB.Table("posts").Select("topic_id, COUNT(*) as post_count").Group("topic_id").Order("post_count DESC").Limit(3).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (db DBGorm) GetTopicByName(name string) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("name = ?", name).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db DBGorm) GetTopicByID(id int) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("id = ?", id).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db DBGorm) SaveNewTopic(topic models.Topic) error {
	result := db.DB.Create(&topic)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db DBGorm) SaveTopic(topic models.Topic) error {
	err := db.DB.Where("id = ?", topic.ID).Save(&topic)
	if err != nil {
		return err.Error
	}
	return nil
}

func (db DBGorm) RemoveTopic(id int) error {
	err := db.DB.Unscoped().Delete(&models.Topic{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
