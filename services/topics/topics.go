package topics

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewTopicServices(topicRepo repositories.ITopicRepository, userRepo repositories.IUserRepository) ITopicServices {
	return &topicServices{ITopicRepository: topicRepo, IUserRepository: userRepo}
}

type ITopicServices interface {
	GetTopics() ([]models.Topic, error)
	CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error)
	GetTopic(id int) (models.Topic, error)
	GetTopTopics() ([]dto.TopTopics, error)
	GetNumberOfPostOnATopicByTopicName(topicName string) (int, error)
	UpdateTopicDescription(topic models.Topic, token dto.Token) (models.Topic, error)
	SaveTopic(topic models.Topic, token dto.Token) error
	RemoveTopic(token dto.Token, id int) error
}

type topicServices struct {
	repositories.ITopicRepository
	repositories.IUserRepository
}

func (t *topicServices) GetTopics() ([]models.Topic, error) {
	//get all topics
	topics, err := t.ITopicRepository.GetAllTopics()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
	} else if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return topics, nil
}
func (t *topicServices) GetTopTopics() ([]dto.TopTopics, error) {
	//get all topics
	topTopics, err := t.ITopicRepository.GetTopTopics()
	if err != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		return nil, err
	}
	for i := range topTopics {
		Topic, err := t.ITopicRepository.GetTopicByID(int(topTopics[i].TopicID))
		if err != nil {
			err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			return nil, err
		}
		topTopics[i].TopicName = Topic.Name
		topTopics[i].TopicDescription = Topic.Description
	}
	return topTopics, nil
}
func (t *topicServices) GetNumberOfPostOnATopicByTopicName(topicName string) (int, error) {
	//get all topics
	postCount, err := t.ITopicRepository.CountNumberOfPostByTopicName(topicName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		return 0, err
	} else if err != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		return 0, err
	}

	return postCount, nil
}
func (t *topicServices) CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error) {
	// isAdmin?
	user, errGetUser := t.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
		err := echo.NewHTTPError(http.StatusNotFound, "User not found")
		return models.Topic{}, err
	} else if errGetUser != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		return models.Topic{}, err
	}

	if !user.IsAdmin {
		return models.Topic{}, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	// isExist?
	_, errGetTopicByName := t.ITopicRepository.GetTopicByName(topic.Name)
	if errors.Is(errGetTopicByName, gorm.ErrRecordNotFound) {
		errSaveNewTopic := t.ITopicRepository.SaveNewTopic(topic)
		if errSaveNewTopic != nil {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errSaveNewTopic.Error())
		}
		Topic, err := t.ITopicRepository.GetTopicByName(topic.Name)
		if err != nil {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return Topic, nil
	} else if errGetTopicByName != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetTopicByName.Error())
	} else {
		return models.Topic{}, echo.NewHTTPError(http.StatusConflict, "Topic already exist")
	}

}

func (t *topicServices) GetTopic(id int) (models.Topic, error) {
	topic, err := t.ITopicRepository.GetTopicByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
	} else if err != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return topic, nil
}

func (t *topicServices) SaveTopic(topic models.Topic, token dto.Token) error {
	// isAdmin?
	user, errGetUser := t.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	} else if errGetUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "admin access only")
	}

	err := t.ITopicRepository.SaveTopic(topic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (t *topicServices) RemoveTopic(token dto.Token, id int) error {
	//check user
	user, errGetUser := t.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if errGetUser != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	topic, errGetTopic := t.ITopicRepository.GetTopicByID(id)
	if errors.Is(errGetTopic, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Topic not found")
	} else if errGetTopic != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetTopic.Error())
	}

	err := t.ITopicRepository.RemoveTopic(int(topic.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (t *topicServices) UpdateTopicDescription(newTopic models.Topic, token dto.Token) (models.Topic, error) {

	user, errGetUser := t.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
		return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if errGetUser != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
	}

	if !user.IsAdmin {
		return models.Topic{}, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	topic, errGetTopicByID := t.ITopicRepository.GetTopicByID(int(newTopic.ID))
	if errors.Is(errGetTopicByID, gorm.ErrRecordNotFound) {
		return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
	} else if errGetTopicByID != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetTopicByID.Error())
	}

	//update topic
	topic.Description = newTopic.Description

	err := t.ITopicRepository.SaveTopic(topic)
	if err != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return topic, nil
}
