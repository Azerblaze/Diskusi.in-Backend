package topics

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewTopicServices(db repositories.IDatabase) ITopicServices {
	return &topicServices{IDatabase: db}
}

type ITopicServices interface {
	GetTopics() ([]models.Topic, error)
	CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error)
	GetTopic(id int) (models.Topic, error)
	UpdateTopicDescription(topic models.Topic, token dto.Token) error
	SaveTopic(topic models.Topic, token dto.Token) error
	RemoveTopic(id int) error
}

type topicServices struct {
	repositories.IDatabase
}

func (t *topicServices) GetTopics() ([]models.Topic, error) {
	//get all topics
	topics, err := t.IDatabase.GetAllTopics()
	if err != nil {
		if err.Error() == "record not found" {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return topics, nil
}

func (t *topicServices) CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error) {
	// isAdmin?
	user, errGetUser := t.IDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser.Error() == "record not found" {
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
	_, errGetTopicByName := t.IDatabase.GetTopicByName(topic.Name)
	if errGetTopicByName != nil {
		if errGetTopicByName.Error() == "record not found" {
			errSaveNewTopic := t.IDatabase.SaveNewTopic(topic)
			if errSaveNewTopic != nil {
				return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, errSaveNewTopic.Error())
			}
			Topic, err := t.IDatabase.GetTopicByName(topic.Name)
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
	topic, err := t.IDatabase.GetTopicByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return models.Topic{}, echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return models.Topic{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return topic, nil
}

func (t *topicServices) SaveTopic(topic models.Topic, token dto.Token) error {
	// isAdmin?
	user, errGetUser := t.IDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "admin access only")
	}

	err := t.IDatabase.SaveTopic(topic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (t *topicServices) RemoveTopic(id int) error {
	err := t.IDatabase.RemoveTopic(id)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (t *topicServices) UpdateTopicDescription(newTopic models.Topic, token dto.Token) error {

	user, errGetUser := t.IDatabase.GetUserByUsername(token.Username)
	if errGetUser != nil {
		if errGetUser.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetUser.Error())
		}
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	topic, errGetTopicByID := t.IDatabase.GetTopicByID(int(newTopic.ID))
	if errGetTopicByID != nil {
		if errGetTopicByID.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Topic not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, errGetTopicByID.Error())
		}
	}

	err := t.IDatabase.SaveTopic(topic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
