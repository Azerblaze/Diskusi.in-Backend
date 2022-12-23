package posts

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories/comments"
	"discusiin/repositories/posts"
	"discusiin/repositories/users"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewCommentServices(dbComment comments.ICommentDatabase, dbPost posts.IPostDatabase, dbUser users.IUserDatabase) ICommentServices {
	return &commentServices{ICommentDatabase: dbComment, IPostDatabase: dbPost, IUserDatabase: dbUser}
}

type ICommentServices interface {
	CreateComment(comment models.Comment, postID int, token dto.Token) error
	GetAllComments(id int) ([]dto.PublicComment, error)
	UpdateComment(newComment models.Comment, token dto.Token) error
	DeleteComment(commentID int, token dto.Token) error
}

type commentServices struct {
	comments.ICommentDatabase
	posts.IPostDatabase
	users.IUserDatabase
}

func (c *commentServices) CreateComment(comment models.Comment, postID int, token dto.Token) error {
	//get post
	post, err := c.IPostDatabase.GetPostById(postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if post is active
	if !post.IsActive {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Post is suspended, All activity stopped")
	}

	//fill empty comment field
	comment.UserID = int(token.ID)
	comment.PostID = int(post.ID)
	comment.IsFollowed = true

	err = c.ICommentDatabase.SaveNewComment(comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (c *commentServices) GetAllComments(id int) ([]dto.PublicComment, error) {
	comments, err := c.ICommentDatabase.GetAllCommentByPost(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
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
	comment, err := c.ICommentDatabase.GetCommentById(int(newComment.ID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//get post
	post, err := c.IPostDatabase.GetPostById(comment.PostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
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
	err = c.ICommentDatabase.SaveComment(comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (c *commentServices) DeleteComment(commentID int, token dto.Token) error {
	//get comment
	user, err := c.IUserDatabase.GetUserByUsername(token.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	comment, err := c.ICommentDatabase.GetCommentById(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check user
	if !user.IsAdmin {
		if comment.UserID != int(token.ID) {
			return echo.NewHTTPError(http.StatusForbidden, "You are not the comment owner")
		}
	}

	err = c.ICommentDatabase.DeleteComment(commentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
