package bookmarks

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IBookmarkDatabase {
	return &GormSql{
		DB: db,
	}
}

func (db GormSql) SaveBookmark(bookmark models.Bookmark) error {
	err := db.DB.Create(&bookmark).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetBookmarkByUserIDAndPostID(userID, postID int) (models.Bookmark, error) {
	var bookmark models.Bookmark

	err := db.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&bookmark).Error
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}
func (db GormSql) GetBookmarkByBookmarkID(bookmarkID int) (models.Bookmark, error) {
	var bookmark models.Bookmark

	err := db.DB.Where("id = ?", bookmarkID).First(&bookmark).Error
	if err != nil {
		return models.Bookmark{}, err
	}

	return bookmark, nil
}

func (db GormSql) DeleteBookmark(bookmarkId int) error {
	err := db.DB.Unscoped().Delete(&models.Bookmark{}, bookmarkId).Error
	if err != nil {
		return err
	}

	return nil
}
func (db GormSql) GetAllBookmark(userId int) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark

	err := db.DB.Where("user_id = ?", userId).Order("created_at DESC").Preload("Post").Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}
