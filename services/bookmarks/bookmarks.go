package bookmarks

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewBookmarkServices(db repositories.IDatabase) IBookmarkServices {
	return &bookmarkServices{IDatabase: db}
}

type IBookmarkServices interface {
	AddBookmark(token dto.Token, postID int) error
	DeleteBookmark(token dto.Token, postID int) error
	GetAllBookmark(token dto.Token) ([]dto.PublicBookmark, error)
}

type bookmarkServices struct {
	repositories.IDatabase
}

func (b *bookmarkServices) AddBookmark(token dto.Token, postID int) error {
	var newBookmark models.Bookmark

	//check post if exist
	post, err := b.IDatabase.GetPostById(postID)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if bookmark exist
	_, err = b.IDatabase.GetBookmark(int(token.ID), int(post.ID))
	if err != nil {
		if err.Error() == "record not found" {
			//insert to empty bookmark field
			newBookmark.UserID = int(token.ID)
			newBookmark.PostID = int(post.ID)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Post has been bookmarked")
	}

	err = b.IDatabase.SaveBookmark(newBookmark)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *bookmarkServices) DeleteBookmark(token dto.Token, postID int) error {
	//check post if needed
	post, err := b.IDatabase.GetPostById(postID)
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//check if bookmark exist
	bookmark, err := b.IDatabase.GetBookmark(int(token.ID), int(post.ID))
	if err != nil {
		if err.Error() == "record not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	//delete bookmark
	err = b.IDatabase.DeleteBookmark(int(bookmark.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *bookmarkServices) GetAllBookmark(token dto.Token) ([]dto.PublicBookmark, error) {
	//get all bookmark
	bookmarks, err := b.IDatabase.GetAllBookmark(int(token.ID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicBookmark
	for _, bookmark := range bookmarks {
		post, _ := b.IDatabase.GetPostById(int(bookmark.ID))
		result = append(result, dto.PublicBookmark{
			Model: bookmark.Model,
			User: dto.BookmarkUser{
				UserID:   bookmark.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Post: dto.BookmarkPost{
				PostID:    bookmark.PostID,
				PostTopic: post.Topic.Name,
				Title:     bookmark.Post.Title,
				Body:      bookmark.Post.Body,
			},
		})
	}

	return result, nil
}
