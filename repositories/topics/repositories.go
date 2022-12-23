package topics

import (
	"discusiin/dto"
	"discusiin/models"
)

type ITopicDatabase interface {
	GetAllTopics() ([]models.Topic, error)
	GetTopicByName(topicName string) (models.Topic, error)
	GetTopicByID(topicID int) (models.Topic, error)
	GetTopTopics() ([]dto.TopTopics, error)
	SaveNewTopic(models.Topic) error
	SaveTopic(models.Topic) error
	RemoveTopic(topicID int) error
	CountAllTopic() (int, error)
}
