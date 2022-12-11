package replies

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewReplyServices(db repositories.IDatabase) IReplyServices {
	return &replyServices{IDatabase: db}
}

type IReplyServices interface {
	CreateReply(reply models.Reply, co int, token dto.Token) error
	GetAllReply(commentId int) ([]dto.PublicReply, error)
	UpdateReply(newReply models.Reply, replyId int, token dto.Token) error
	DeleteReply(replyId int, token dto.Token) error
}

type replyServices struct {
	repositories.IDatabase
}

func (r *replyServices) CreateReply(reply models.Reply, co int, token dto.Token) error {
	//get comment
	comment, err := r.IDatabase.GetCommentById(co)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//input empty field in reply
	reply.CommentID = int(comment.ID)
	reply.UserID = int(token.ID)

	//create reply
	err = r.IDatabase.SaveNewReply(reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *replyServices) GetAllReply(commentId int) ([]dto.PublicReply, error) {
	replies, err := r.IDatabase.GetAllReplyByComment(commentId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var result []dto.PublicReply
	for _, reply := range replies {
		result = append(result, dto.PublicReply{
			Model:     reply.Model,
			CommentID: reply.CommentID,
			Body:      reply.Body,
			User: dto.ReplyUser{
				UserID:   reply.UserID,
				Username: reply.User.Username,
				Photo:    reply.User.Photo,
			},
		})
	}

	return result, nil
}

func (r *replyServices) UpdateReply(newReply models.Reply, replyId int, token dto.Token) error {
	//find reply
	reply, err := r.IDatabase.GetReplyById(replyId)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Reply not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if user are correct
	if reply.UserID != int(token.ID) {
		return echo.NewHTTPError(http.StatusForbidden, "You are not the reply owner")
	}

	//update reply field
	reply.Body += " "
	reply.Body += newReply.Body

	//update reply
	err = r.IDatabase.SaveReply(reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *replyServices) DeleteReply(replyId int, token dto.Token) error {
	//find reply
	reply, err := r.IDatabase.GetReplyById(replyId)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Reply not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if user are correct
	if reply.UserID != int(token.ID) {
		return echo.NewHTTPError(http.StatusForbidden, "You are not the reply owner")
	}

	//delete reply
	err = r.IDatabase.DeleteReply(replyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
