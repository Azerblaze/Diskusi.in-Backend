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
	UpdateReply(newReply models.Reply, re int, userId int) error
	DeleteReply(re int, userId int) error
}

type replyServices struct {
	repositories.IDatabase
}

func (p *replyServices) CreateReply(reply models.Reply, co int, token dto.Token) error {
	//get comment
	comment, err := p.IDatabase.GetCommentById(co)
	if err != nil {
		return err
	}

	//input empty field in reply
	reply.CommentID = int(comment.ID)
	reply.UserID = int(token.ID)

	//create reply
	err = p.IDatabase.SaveNewReply(reply)
	if err != nil {
		return err
	}

	return nil
}

func (p *replyServices) UpdateReply(newReply models.Reply, re int, userId int) error {
	//find reply
	reply, err := p.IDatabase.GetReplyById(re)
	if err != nil {
		return err
	}

	//check if user are correct
	if reply.UserID != userId {
		return errors.New("user not eligible")
	}

	//update reply field
	reply.Body += " "
	reply.Body += newReply.Body

	//update reply
	err = p.IDatabase.SaveReply(reply)
	if err != nil {
		return err
	}

	return nil
}

func (p *replyServices) DeleteReply(re int, userId int) error {
	//find reply
	reply, err := p.IDatabase.GetReplyById(re)
	if err != nil {
		return err
	}

	//check if user are correct
	if reply.UserID != userId {
		return errors.New("user not eligible")
	}

	//delete reply
	err = p.IDatabase.DeleteReply(re)
	if err != nil {
		return err
	}

	return nil
}
