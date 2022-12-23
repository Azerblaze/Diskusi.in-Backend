package bookmarks

import (
	"discusiin/models"
	"discusiin/repositories"

	"gorm.io/gorm"
)

type DBGorm struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) repositories.IBookmarkRepository {
	return &DBGorm{
		DB: db,
	}
}
func (db DBGorm) SaveBookmark(bookmark models.Bookmark) error {
	err := db.DB.Create(&bookmark).Error
	if err != nil {
		return err
	}

	return nil
}

func (db DBGorm) GetBookmarkByUserIDAndPostID(userID, postID int) (models.Bookmark, error) {
	var bookmark models.Bookmark

	err := db.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&bookmark).Error
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}
func (db DBGorm) GetBookmarkByBookmarkID(bookmarkID int) (models.Bookmark, error) {
	var bookmark models.Bookmark

	err := db.DB.Where("id = ?", bookmarkID).First(&bookmark).Error
	if err != nil {
		return models.Bookmark{}, err
	}

	return bookmark, nil
}

func (db DBGorm) DeleteBookmark(bookmarkId int) error {
	err := db.DB.Unscoped().Delete(&models.Bookmark{}, bookmarkId).Error
	if err != nil {
		return err
	}

	return nil
}
func (db DBGorm) GetAllBookmark(userId int) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark

	err := db.DB.Where("user_id = ?", userId).Order("created_at DESC").Preload("Post").Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}
