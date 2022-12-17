package repositories

import (
	"discusiin/models"

	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func NewGorm(db *gorm.DB) IDatabase {
	return &GormSql{
		DB: db,
	}
}

// User
func (db GormSql) SaveNewUser(user models.User) error {
	result := db.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (db GormSql) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) GetUserById(userId int) (models.User, error) {
	var user models.User
	err := db.DB.
		Where("id = ?", userId).
		First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) GetUsersAdminNotIncluded(page int) ([]models.User, error) {
	var users []models.User
	err := db.DB.Where("is_admin = 0").Order("username ASC").Offset((page - 1) * 20).Limit(20).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (db GormSql) GetProfile(id int) (models.User, error) {
	var user models.User
	err := db.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (db GormSql) UpdateProfile(user models.User) error {
	err := db.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (db GormSql) DeleteUser(userId int) error {
	err := db.DB.Delete(&models.User{}, userId).Error
	if err != nil {
		return err
	}

	return nil
}

// Topic -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) GetAllTopics() ([]models.Topic, error) {
	var topics []models.Topic

	result := db.DB.Find(&topics)

	if result.Error != nil {
		return nil, result.Error
	} else {
		if result.RowsAffected <= 0 {
			return nil, result.Error
		} else {
			return topics, nil
		}
	}
}

func (db GormSql) GetTopicByName(name string) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("name = ?", name).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db GormSql) GetTopicByID(id int) (models.Topic, error) {
	var topic models.Topic
	err := db.DB.Where("id = ?", id).First(&topic).Error

	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}

func (db GormSql) SaveNewTopic(topic models.Topic) error {
	result := db.DB.Create(&topic)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db GormSql) SaveTopic(topic models.Topic) error {
	err := db.DB.Where("id = ?", topic.ID).Save(&topic)
	if err != nil {
		return err.Error
	}
	return nil
}

func (db GormSql) RemoveTopic(id int) error {
	err := db.DB.Delete(&models.Topic{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

// Post -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveNewPost(post models.Post) error {
	err := db.DB.Create(&post).Error
	if err != nil {
		return err
	}

	return nil
}
func (db GormSql) GetRecentPost(page int, search string) ([]models.Post, error) {
	var result []models.Post
	err := db.DB.Where("title LIKE ?", "%"+search+"%").Order("created_at DESC").Preload("User").Preload("Topic").Offset((page - 1) * 20).Limit(20).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (db GormSql) GetAllPostByTopic(id int, page int, search string) ([]models.Post, error) {
	var posts []models.Post

	//find topic id
	err := db.DB.Where("topic_id = ?", id).
		Where("title LIKE ?", "%"+search+"%").
		Order("created_at DESC").
		Preload("User").
		Preload("Topic").
		Offset((page - 1) * 20).
		Limit(20).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (db GormSql) GetPostById(id int) (models.Post, error) {
	var post models.Post

	err := db.DB.Where("id = ?", id).
		Preload("User").
		Preload("Topic").
		First(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (db GormSql) GetPostByUserId(userId int, page int) ([]models.Post, error) {
	var posts []models.Post

	//find topic id
	err := db.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("User").
		Preload("Topic").
		Offset((page - 1) * 20).
		Limit(20).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (db GormSql) GetAllPostByLike(page int) ([]models.Post, error) {
	var posts []models.Post

	//find topic with like
	err := db.DB.Order("like_count DESC").Preload("User").Preload("Topic").Offset((page - 1) * 20).Limit(20).Find(&posts).Error
	if err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (db GormSql) SavePost(post models.Post) error {
	err := db.DB.Save(&post).Error

	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeletePost(id int) error {
	err := db.DB.Delete(&models.Post{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetPostByIdWithAll(id int) (models.Post, error) {
	var post models.Post
	err := db.DB.Model(&models.Post{}).Where("id = ?", id).Preload("Comments").First(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

// Comment -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveNewComment(comment models.Comment) error {
	err := db.DB.Create(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllCommentByPost(id int) ([]models.Comment, error) {
	var comments []models.Comment

	err := db.DB.Where("post_id = ?", id).Order("created_at DESC").Preload("User").Find(&comments).Error
	if err != nil {
		return []models.Comment{}, err
	}

	return comments, nil
}

func (db GormSql) GetCommentById(co int) (models.Comment, error) {
	var comment models.Comment

	err := db.DB.Where("id = ?", co).First(&comment).Error
	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (db GormSql) GetCommentByUserId(userId int, page int) ([]models.Comment, error) {
	var comments []models.Comment

	//find topic id
	err := db.DB.Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Post").
		Offset((page - 1) * 30).
		Limit(30).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (db GormSql) SaveComment(comment models.Comment) error {
	err := db.DB.Save(&comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeleteComment(co int) error {
	err := db.DB.Delete(&models.Comment{}, co).Error
	if err != nil {
		return err
	}

	return nil
}

// Reply -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveNewReply(reply models.Reply) error {
	err := db.DB.Create(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllReplyByComment(commentId int) ([]models.Reply, error) {
	var replys []models.Reply
	err := db.DB.Where("comment_id = ?", commentId).Preload("User").Find(&replys).Error
	if err != nil {
		return []models.Reply{}, err
	}

	return replys, nil
}

func (db GormSql) GetReplyById(id int) (models.Reply, error) {
	var reply models.Reply
	err := db.DB.Where("id = ?", id).First(&reply).Error
	if err != nil {
		return models.Reply{}, err
	}

	return reply, nil
}

func (db GormSql) SaveReply(reply models.Reply) error {
	err := db.DB.Save(&reply).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) DeleteReply(re int) error {
	err := db.DB.Delete(&models.Reply{}, re).Error
	if err != nil {
		return err
	}

	return nil
}

// Like -------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) GetLikeByUserAndPostId(userId int, postId int) (models.Like, error) {
	var like models.Like

	err := db.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like).Error
	if err != nil {
		return models.Like{}, err
	}

	return like, nil
}

func (db GormSql) SaveNewLike(like models.Like) error {
	err := db.DB.Create(&like).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) SaveLike(like models.Like) error {
	err := db.DB.Save(&like).Error
	if err != nil {
		return err
	}

	return nil
}

// Bookmark ------------------------------------------------------------------------------------------------------------------------------------------------
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

// FollowedPost ------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) SaveFollowedPost(followedPost models.FollowedPost) error {
	err := db.DB.Create(&followedPost).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetFollowedPost(userId int, postId int) (models.FollowedPost, error) {
	var followedPost models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Where("post_id = ?", postId).First(&followedPost).Error
	if err != nil {
		return models.FollowedPost{}, err
	}

	return followedPost, nil
}

func (db GormSql) DeleteFollowedPost(followedPostId int) error {
	err := db.DB.Unscoped().Delete(&models.FollowedPost{}, followedPostId).Error
	if err != nil {
		return err
	}

	return nil
}

func (db GormSql) GetAllFollowedPost(userId int) ([]models.FollowedPost, error) {
	var followedPosts []models.FollowedPost

	err := db.DB.Where("user_id = ?", userId).Order("created_at DESC").Preload("Post").Find(&followedPosts).Error
	if err != nil {
		return nil, err
	}

	return followedPosts, nil
}

// Count ------------------------------------------------------------------------------------------------------------------------------------------------
func (db GormSql) CountAllUserNotIncludeDeletedUser() (int, error) {
	var numberOfUser int64

	err := db.DB.Table("users").Where("deleted_at IS NULL").Count(&numberOfUser).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfUser), nil
}
func (db GormSql) CountAllTopic() (int, error) {
	var numberOfTopic int64

	err := db.DB.Table("topics").Count(&numberOfTopic).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfTopic), nil
}
func (db GormSql) CountPostLike(postID int) (int, error) {
	var postLike int64

	err := db.DB.Table("likes").Where("post_id = ? AND is_like = 1", postID).Count(&postLike).Error
	if err != nil {
		return 0, err
	}

	return int(postLike), nil
}
func (db GormSql) CountPostDislike(postID int) (int, error) {
	var postDislike int64
	err := db.DB.Table("likes").Where("post_id = ? AND is_dislike = 1", postID).Count(&postDislike).Error
	if err != nil {
		return 0, err
	}
	return int(postDislike), nil
}
func (db GormSql) CountPostComment(postID int) (int, error) {
	var commentCount int64
	err := db.DB.Table("comments").Where("post_id = ?", postID).Count(&commentCount).Error
	if err != nil {
		return 0, err
	}
	return int(commentCount), nil
}
func (db GormSql) CountAllPost() (int, error) {
	var numberOfPost int64

	err := db.DB.Table("posts").Count(&numberOfPost).Error
	if err != nil {
		return 0, err
	}

	return int(numberOfPost), nil
}
func (db GormSql) CountPostByTopicID(topicId int) (int, error) {
	var postCount int64

	err := db.DB.Table("posts").Where("topic_id = ?", topicId).Count(&postCount).Error
	if err != nil {
		return 0, err
	}

	return int(postCount), nil
}
func (db GormSql) CountPostByUserID(userId int) (int, error) {
	var postCount int64

	err := db.DB.Table("posts").Where("user_id = ?", userId).Count(&postCount).Error
	if err != nil {
		return 0, err
	}

	return int(postCount), nil
}

func (db GormSql) CountCommentByUserID(userId int) (int, error) {
	var commentCount int64

	err := db.DB.Table("comments").Where("user_id = ?", userId).Count(&commentCount).Error
	if err != nil {
		return 0, err
	}

	return int(commentCount), nil
}
func (db GormSql) CountAllUserNotAdminNotIncludeDeletedUser() (int, error) {
	var userCount int64

	err := db.DB.Table("users").Where("is_admin = 0").Where("deleted_at IS NULL").Count(&userCount).Error
	if err != nil {
		return 0, err
	}
	return int(userCount), nil
}
