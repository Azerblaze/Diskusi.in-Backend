package bookmarks

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check if bookmark exist
	_, err = b.IDatabase.GetBookmarkByUserIDAndPostID(int(token.ID), int(post.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//insert to empty bookmark field
		newBookmark.UserID = int(token.ID)
		newBookmark.PostID = int(post.ID)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Post has been bookmarked")
	}

	err = b.IDatabase.SaveBookmark(newBookmark)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (b *bookmarkServices) DeleteBookmark(token dto.Token, bookmarkID int) error {

	//check if bookmark exist
	_, err := b.IDatabase.GetBookmarkByBookmarkID(bookmarkID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Bookmark not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		err = b.IDatabase.DeleteBookmark(bookmarkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return nil
	}

}

func (b *bookmarkServices) GetAllBookmark(token dto.Token) ([]dto.PublicBookmark, error) {
	//get all bookmark
	bookmarks, err := b.IDatabase.GetAllBookmark(int(token.ID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicBookmark
	for _, bookmark := range bookmarks {
		post, _ := b.IDatabase.GetPostById(int(bookmark.PostID))
		result = append(result, dto.PublicBookmark{
			Model: bookmark.Model,
			User: dto.BookmarkUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Post: dto.BookmarkPost{
				PostID:    int(post.ID),
				PostTopic: post.Topic.Name,
				Title:     post.Title,
				Body:      post.Body,
			},
		})
	}

	return result, nil
}
