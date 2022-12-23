package bookmarks

import (
	"discusiin/models"
)

type IBookmarkDatabase interface {
	SaveBookmark(bookmark models.Bookmark) error
	GetBookmarkByUserIDAndPostID(userID, postID int) (models.Bookmark, error)
	GetBookmarkByBookmarkID(bookmarkID int) (models.Bookmark, error)
	DeleteBookmark(bookmarkId int) error
	GetAllBookmark(userId int) ([]models.Bookmark, error)
}
