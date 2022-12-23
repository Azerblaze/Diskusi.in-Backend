package followedPosts

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewFollowedPostServices(postRepo repositories.IPostRepository, followedPostRepo repositories.IFollowedPostRepository) IFollowedPostServices {
	return &followedPostServices{IPostRepository: postRepo, IFollowedPostRepository: followedPostRepo}
}

type IFollowedPostServices interface {
	AddFollowedPost(token dto.Token, postID int) error
	DeleteFollowedPost(token dto.Token, postID int) error
	GetAllFollowedPost(token dto.Token) ([]dto.PublicFollowedPost, error)
}

type followedPostServices struct {
	repositories.IPostRepository
	repositories.IFollowedPostRepository
}

func (b *followedPostServices) AddFollowedPost(token dto.Token, postID int) error {
	var newFollowedPost models.FollowedPost

	//check post if exist
	post, err := b.IPostRepository.GetPostById(postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check if followedPost exist
	_, err = b.IFollowedPostRepository.GetFollowedPost(int(token.ID), int(post.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//insert to empty followedPost field
		newFollowedPost.UserID = int(token.ID)
		newFollowedPost.PostID = int(post.ID)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Post has been followed")
	}

	err = b.IFollowedPostRepository.SaveFollowedPost(newFollowedPost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *followedPostServices) DeleteFollowedPost(token dto.Token, postID int) error {
	//check post if needed
	post, err := b.IPostRepository.GetPostById(postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check if followedPost exist
	followedPost, err := b.IFollowedPostRepository.GetFollowedPost(int(token.ID), int(post.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "You are not following this post")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//delete followedPost
	err = b.IFollowedPostRepository.DeleteFollowedPost(int(followedPost.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *followedPostServices) GetAllFollowedPost(token dto.Token) ([]dto.PublicFollowedPost, error) {
	//get all followedPost
	followedPosts, err := b.IFollowedPostRepository.GetAllFollowedPost(int(token.ID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicFollowedPost
	for _, followedPost := range followedPosts {
		post, _ := b.IPostRepository.GetPostById(int(followedPost.PostID))
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
