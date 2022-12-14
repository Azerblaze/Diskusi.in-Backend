package followedPosts

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewFollowedPostServices(db repositories.IDatabase) IFollowedPostServices {
	return &followedPostServices{IDatabase: db}
}

type IFollowedPostServices interface {
	AddFollowedPost(token dto.Token, postID int) error
	DeleteFollowedPost(token dto.Token, postID int) error
	GetAllFollowedPost(token dto.Token) ([]dto.PublicFollowedPost, error)
}

type followedPostServices struct {
	repositories.IDatabase
}

func (b *followedPostServices) AddFollowedPost(token dto.Token, postID int) error {
	var newFollowedPost models.FollowedPost

	//check post if exist
	post, err := b.IDatabase.GetPostById(postID)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if followedPost exist
	_, err = b.IDatabase.GetFollowedPost(int(token.ID), int(post.ID))
	if err != nil {
		if err.Error() == "record not found" {
			//insert to empty followedPost field
			newFollowedPost.UserID = int(token.ID)
			newFollowedPost.PostID = int(post.ID)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Post has been followed")
	}

	err = b.IDatabase.SaveFollowedPost(newFollowedPost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *followedPostServices) DeleteFollowedPost(token dto.Token, postID int) error {
	//check post if needed
	post, err := b.IDatabase.GetPostById(postID)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if followedPost exist
	followedPost, err := b.IDatabase.GetFollowedPost(int(token.ID), int(post.ID))
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "You are not following this post")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//delete followedPost
	err = b.IDatabase.DeleteFollowedPost(int(followedPost.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *followedPostServices) GetAllFollowedPost(token dto.Token) ([]dto.PublicFollowedPost, error) {
	//get all followedPost
	followedPosts, err := b.IDatabase.GetAllFollowedPost(int(token.ID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicFollowedPost
	for _, followedPost := range followedPosts {
		post, _ := b.IDatabase.GetPostById(int(followedPost.ID))
		result = append(result, dto.PublicFollowedPost{
			Model: followedPost.Model,
			User: dto.FollowedPostUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Post: dto.FollowedPost{
				PostID:    int(post.ID),
				PostTopic: post.Topic.Name,
				Title:     post.Title,
				Body:      post.Body,
			},
		})
	}

	return result, nil
}
