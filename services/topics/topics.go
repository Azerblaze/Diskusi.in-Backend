package topics

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories/posts"
	"discusiin/repositories/topics"
	"discusiin/repositories/users"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewTopicServices(dbTopic topics.ITopicDatabase, dbPost posts.IPostDatabase, dbUser users.IUserDatabase) ITopicServices {
	return &topicServices{ITopicDatabase: dbTopic, IPostDatabase: dbPost, IUserDatabase: dbUser}
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
	topics.ITopicDatabase
	posts.IPostDatabase
	users.IUserDatabase
}

func (t *topicServices) GetTopics() ([]models.Topic, error) {
	//get all topics
	topics, err := t.ITopicDatabase.GetAllTopics()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return topics, nil
}
func (t *topicServices) GetTopTopics() ([]dto.TopTopics, error) {
	//get all topics
	topTopics, err := t.ITopicDatabase.GetTopTopics()
	if err != nil {
		err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		if err == nil {
			panic("unexpected nil error")
		}
		return nil, err
	}
	for i := range topTopics {
		Topic, err := t.ITopicDatabase.GetTopicByID(int(topTopics[i].TopicID))
		if err != nil {
			err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			if err == nil {
				panic("unexpected nil error")
			}
			return nil, err
		}
		topTopics[i].TopicName = Topic.Name
		topTopics[i].TopicDescription = Topic.Description
	}
	return topTopics, nil
}
func (t *topicServices) GetNumberOfPostOnATopicByTopicName(topicName string) (int, error) {
	//get all topics
	postCount, err := t.IPostDatabase.CountNumberOfPostByTopicName(topicName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := echo.NewHTTPError(http.StatusNotFound, "Topic not found")
			if err == nil {
				panic("unexpected nil error")
			}
			return 0, err
		} else {
			err := echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			if err == nil {
				panic("unexpected nil error")
			}
			return 0, err
		}
	}

	return postCount, nil
}
func (t *topicServices) CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error) {
	// isAdmin?
	user, errGetUser := t.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser == gorm.ErrRecordNotFound {
			return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return models.Topic{}, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	var result models.Topic
	// isExist?
	_, errGetTopicByName := t.ITopicDatabase.GetTopicByName(topic.Name)
	if errGetTopicByName != nil {
		if errGetTopicByName == gorm.ErrRecordNotFound {
			errSaveNewTopic := t.ITopicDatabase.SaveNewTopic(topic)
			if errSaveNewTopic != nil {
				return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errSaveNewTopic.Error())
			}
			Topic, err := t.ITopicDatabase.GetTopicByName(topic.Name)
			if err != nil {
				return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			result = Topic
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetTopicByName.Error())
		}
	} else {
		return models.Topic{}, echo.NewHTTPError(http.StatusConflict, "Topic already exist")
	}
	return result, nil
}

func (t *topicServices) GetTopic(id int) (models.Topic, error) {
	topic, err := t.ITopicDatabase.GetTopicByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return topic, nil
}

func (t *topicServices) SaveTopic(topic models.Topic, token dto.Token) error {
	// isAdmin?
	user, errGetUser := t.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "admin access only")
	}

	err := t.ITopicDatabase.SaveTopic(topic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (t *topicServices) RemoveTopic(token dto.Token, id int) error {
	//check user
	user, errGetUser := t.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	topic, errGetTopic := t.ITopicDatabase.GetTopicByID(id)
	if errGetTopic != nil {
		if errGetTopic == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetTopic.Error())
		}
	}

	err := t.ITopicDatabase.RemoveTopic(int(topic.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (t *topicServices) UpdateTopicDescription(newTopic models.Topic, token dto.Token) (models.Topic, error) {

	user, errGetUser := t.IUserDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser == gorm.ErrRecordNotFound {
			return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return models.Topic{}, echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	topic, errGetTopicByID := t.ITopicDatabase.GetTopicByID(int(newTopic.ID))
	if errGetTopicByID != nil {
		if errGetTopicByID == gorm.ErrRecordNotFound {
			return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errGetTopicByID.Error())
		}
	}

	//update topic
	topic.Description = newTopic.Description

	err := t.ITopicDatabase.SaveTopic(topic)
	if err != nil {
		return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return topic, nil
}
