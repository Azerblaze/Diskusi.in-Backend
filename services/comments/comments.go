package posts

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewCommentServices(commentRepo repositories.ICommentRepository, postRepo repositories.IPostRepository, userRepo repositories.IUserRepository) ICommentServices {
	return &commentServices{ICommentRepository: commentRepo, IPostRepository: postRepo, IUserRepository: userRepo}
}

type ICommentServices interface {
	CreateComment(comment models.Comment, post_id int, token dto.Token) error
	GetAllComments(id int) ([]dto.PublicComment, error)
	UpdateComment(newComment models.Comment, token dto.Token) error
	DeleteComment(commentID int, token dto.Token) error
}

type commentServices struct {
	repositories.ICommentRepository
	repositories.IPostRepository
	repositories.IUserRepository
}

func (c *commentServices) CreateComment(comment models.Comment, postID int, token dto.Token) error {
	//get post
	post, err := c.IPostRepository.GetPostById(postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check if post is active
	if !post.IsActive {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Post is suspended, All activity stopped")
	}

	//fill empty comment field
	comment.UserID = int(token.ID)
	comment.PostID = int(post.ID)
	comment.IsFollowed = true

	err = c.ICommentRepository.SaveNewComment(comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil

}

func (c *commentServices) GetAllComments(id int) ([]dto.PublicComment, error) {
	comments, err := c.ICommentRepository.GetAllCommentByPost(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicComment
	for _, comment := range comments {
		result = append(result, dto.PublicComment{
			Model:  comment.Model,
			PostID: comment.PostID,
			Body:   comment.Body,
			User: dto.CommentUser{
				UserID:   (comment.UserID),
				Username: comment.User.Username,
				Photo:    comment.User.Photo,
			},
		})
	}

	return result, nil
}

func (c *commentServices) UpdateComment(newComment models.Comment, token dto.Token) error {
	//get comment
	comment, err := c.ICommentRepository.GetCommentById(int(newComment.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//get post
	post, err := c.IPostRepository.GetPostById(comment.PostID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check if post is active
	if !post.IsActive {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Post is suspended, All activity stopped")
	}

	//check user
	if comment.UserID != int(token.ID) {
		return echo.NewHTTPError(http.StatusForbidden, "You are not the comment owner")
	}

	//update comment field
	comment.Body += " "
	comment.Body += newComment.Body

	//save comment
	err = c.ICommentRepository.SaveComment(comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (c *commentServices) DeleteComment(commentID int, token dto.Token) error {
	//get comment
	user, err := c.IUserRepository.GetUserByUsername(token.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	comment, err := c.ICommentRepository.GetCommentById(commentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check user
	if !user.IsAdmin {
		if comment.UserID != int(token.ID) {
			return echo.NewHTTPError(http.StatusForbidden, "You are not the comment owner")
		}
	}

	err = c.ICommentRepository.DeleteComment(commentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
