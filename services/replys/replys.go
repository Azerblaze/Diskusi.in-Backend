package replys

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
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
		return err
	}

	//input empty field in reply
	reply.CommentID = int(comment.ID)
	reply.UserID = int(token.ID)

	//create reply
	err = r.IDatabase.SaveNewReply(reply)
	if err != nil {
		return err
	}

	return nil
}

func (r *replyServices) GetAllReply(commentId int) ([]dto.PublicReply, error) {
	replys, err := r.IDatabase.GetAllReplyByComment(commentId)
	if err != nil {
		return []dto.PublicReply{}, err
	}

	var result []dto.PublicReply
	for _, reply := range replys {
		result = append(result, dto.PublicReply{
			Model:     reply.Model,
			UserID:    reply.UserID,
			CommentID: reply.CommentID,
			Body:      reply.Body,
			Username:  reply.User.Username,
		})
	}

	return result, nil
}

func (r *replyServices) UpdateReply(newReply models.Reply, replyId int, token dto.Token) error {
	//find reply
	reply, err := r.IDatabase.GetReplyById(replyId)
	if err != nil {
		return err
	}

	//check if user are correct
	if reply.UserID != int(token.ID) {
		return errors.New("user not eligible")
	}

	//update reply field
	reply.Body += " "
	reply.Body += newReply.Body

	//update reply
	err = r.IDatabase.SaveReply(reply)
	if err != nil {
		return err
	}

	return nil
}

func (r *replyServices) DeleteReply(replyId int, token dto.Token) error {
	//find reply
	reply, err := r.IDatabase.GetReplyById(replyId)
	if err != nil {
		return err
	}

	//check if user are correct
	if reply.UserID != int(token.ID) {
		return errors.New("user not eligible")
	}

	//delete reply
	err = r.IDatabase.DeleteReply(replyId)
	if err != nil {
		return err
	}

	return nil
}
