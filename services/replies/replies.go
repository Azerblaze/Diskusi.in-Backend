package replies

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories/comments"
	"discusiin/repositories/replies"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewReplyServices(dbReply replies.IReplyDatabase, dbComment comments.ICommentDatabase) IReplyServices {
	return &replyServices{IReplyDatabase: dbReply, ICommentDatabase: dbComment}
}

type IReplyServices interface {
	CreateReply(reply models.Reply, co int, token dto.Token) error
	GetAllReply(commentId int) ([]dto.PublicReply, error)
	UpdateReply(newReply models.Reply, replyId int, token dto.Token) error
	DeleteReply(replyId int, token dto.Token) error
}

type replyServices struct {
	replies.IReplyDatabase
	comments.ICommentDatabase
}

func (r *replyServices) CreateReply(reply models.Reply, co int, token dto.Token) error {
	//get comment
	comment, err := r.ICommentDatabase.GetCommentById(co)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//input empty field in reply
	reply.CommentID = int(comment.ID)
	reply.UserID = int(token.ID)

	//create reply
	err = r.IReplyDatabase.SaveNewReply(reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *replyServices) GetAllReply(commentId int) ([]dto.PublicReply, error) {
	replies, err := r.IReplyDatabase.GetAllReplyByComment(commentId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	reply, err := r.IReplyDatabase.GetReplyById(replyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	err = r.IReplyDatabase.SaveReply(reply)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *replyServices) DeleteReply(replyId int, token dto.Token) error {
	//find reply
	reply, err := r.IReplyDatabase.GetReplyById(replyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	err = r.IReplyDatabase.DeleteReply(replyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
