package posts

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewPostServices(topicRepo repositories.ITopicRepository, postRepo repositories.IPostRepository, userRepo repositories.IUserRepository) IPostServices {
	return &postServices{ITopicRepository: topicRepo, IPostRepository: postRepo, IUserRepository: userRepo}
}

type IPostServices interface {
	CreatePost(post models.Post, name string, token dto.Token) error
	GetPosts(name string, page int, search string) ([]dto.PublicPost, int, error)
	GetPostsByTopicByLike(name string, page int) ([]dto.PublicPost, int, error)
	GetPost(id int) (dto.PublicPost, error)
	UpdatePost(newPost models.Post, id int, token dto.Token) error
	DeletePost(id int, token dto.Token) error
	GetRecentPost(page int, search string) ([]dto.PublicPost, int, error)
	GetAllPostByLike(page int, search string) ([]dto.PublicPost, int, error)
	SuspendPost(token dto.Token, postId int) error
}

type postServices struct {
	repositories.ITopicRepository
	repositories.IPostRepository
	repositories.IUserRepository
}

func (p *postServices) CreatePost(post models.Post, name string, token dto.Token) error {
	//find topic
	topic, err := p.ITopicRepository.GetTopicByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "topic not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//owner
	post.UserID = int(token.ID)
	//insert topic id and is active
	post.TopicID = int(topic.ID)
	//epoch time
	post.CreatedAt = int(time.Now().UnixMilli())
	// isActiveDefault
	post.IsActive = true

	//save post
	err = p.IPostRepository.SaveNewPost(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (p *postServices) GetPosts(name string, page int, search string) ([]dto.PublicPost, int, error) {
	//find topic
	topic, err := p.ITopicRepository.GetTopicByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, echo.NewHTTPError(http.StatusNotFound, "topic not found")
	} else if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	posts, err := p.IPostRepository.GetAllPostByTopic(int(topic.ID), page, search)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := p.IPostRepository.CountPostLike(int(post.ID))
		commentCount, _ := p.IPostRepository.CountPostComment(int(post.ID))
		dislikeCount, _ := p.IPostRepository.CountPostDislike(int(post.ID))

		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := p.IPostRepository.CountPostByTopicID(int(topic.ID))
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut
	if numberOfPost%pageSize != 0 {
		numberOfPage++
	}
	return result, int(numberOfPage), nil
}

func (p *postServices) GetPostsByTopicByLike(name string, page int) ([]dto.PublicPost, int, error) {
	//find topic
	topic, err := p.ITopicRepository.GetTopicByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, echo.NewHTTPError(http.StatusNotFound, "topic not found")
	} else if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	posts, err := p.IPostRepository.GetAllPostByTopicByLike(int(topic.ID), page)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := p.IPostRepository.CountPostLike(int(post.ID))
		commentCount, _ := p.IPostRepository.CountPostComment(int(post.ID))
		dislikeCount, _ := p.IPostRepository.CountPostDislike(int(post.ID))

		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Photo:    post.User.Photo,
				Username: post.User.Username,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := p.IPostRepository.CountPostByTopicID(int(topic.ID))
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut
	if numberOfPost%pageSize != 0 {
		numberOfPage++
	}
	return result, int(numberOfPage), nil
}

func (p *postServices) GetPost(id int) (dto.PublicPost, error) {
	post, err := p.IPostRepository.GetPostById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.PublicPost{}, echo.NewHTTPError(http.StatusNotFound, "post not found")
	} else if err != nil {
		return dto.PublicPost{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	likeCount, _ := p.IPostRepository.CountPostLike(int(post.ID))
	commentCount, _ := p.IPostRepository.CountPostComment(int(post.ID))
	dislikeCount, _ := p.IPostRepository.CountPostDislike(int(post.ID))
	result := dto.PublicPost{
		Model:     post.Model,
		Title:     post.Title,
		Photo:     post.Photo,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		IsActive:  post.IsActive,
		User: dto.PostUser{
			UserID:   post.UserID,
			Photo:    post.User.Photo,
			Username: post.User.Username,
		},
		Topic: dto.PostTopic{
			TopicID:   post.TopicID,
			TopicName: post.Topic.Name,
		},
		Count: dto.PostCount{
			LikeCount:    likeCount,
			CommentCount: commentCount,
			DislikeCount: dislikeCount,
		},
	}

	return result, nil
}

func (p *postServices) UpdatePost(newPost models.Post, postID int, token dto.Token) error {
	//get previous post
	post, err := p.IPostRepository.GetPostById(postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if int(token.ID) != post.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "You are not the post owner")
	}

	//check if post is active
	if !post.IsActive {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Post is suspended, All activity stopped")
	}

	//update post body
	post.Body += " "
	post.Body += newPost.Body

	err = p.IPostRepository.SavePost(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (p *postServices) DeletePost(id int, token dto.Token) error {
	//find post
	post, err := p.IPostRepository.GetPostById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//check user
	user, err := p.IUserRepository.GetUserByUsername(token.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !user.IsAdmin {
		if int(token.ID) != post.UserID {
			return echo.NewHTTPError(http.StatusForbidden, "You are not the post owner")
		}
	}

	err = p.IPostRepository.DeletePostByPostID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (p *postServices) GetRecentPost(page int, search string) ([]dto.PublicPost, int, error) {
	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	posts, err := p.IPostRepository.GetRecentPost(page, search)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := p.IPostRepository.CountPostLike(int(post.ID))
		commentCount, _ := p.IPostRepository.CountPostComment(int(post.ID))
		dislikeCount, _ := p.IPostRepository.CountPostDislike(int(post.ID))
		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Username: post.User.Username,
				Photo:    post.User.Photo,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := p.IPostRepository.CountAllPost()
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut
	if numberOfPost%pageSize != 0 {
		numberOfPage++
	}

	return result, int(numberOfPage), nil
}

func (p *postServices) GetAllPostByLike(page int, search string) ([]dto.PublicPost, int, error) {
	//cek jika page kosong
	if page < 1 {
		page = 1
	}

	posts, err := p.IPostRepository.GetAllPostByLike(page)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result []dto.PublicPost
	for _, post := range posts {
		likeCount, _ := p.IPostRepository.CountPostLike(int(post.ID))
		commentCount, _ := p.IPostRepository.CountPostComment(int(post.ID))
		dislikeCount, _ := p.IPostRepository.CountPostDislike(int(post.ID))
		result = append(result, dto.PublicPost{
			Model:     post.Model,
			Title:     post.Title,
			Photo:     post.Photo,
			Body:      post.Body,
			CreatedAt: post.CreatedAt,
			IsActive:  post.IsActive,
			User: dto.PostUser{
				UserID:   post.UserID,
				Username: post.User.Username,
				Photo:    post.User.Photo,
			},
			Topic: dto.PostTopic{
				TopicID:   post.TopicID,
				TopicName: post.Topic.Name,
			},
			Count: dto.PostCount{
				LikeCount:    likeCount,
				CommentCount: commentCount,
				DislikeCount: dislikeCount,
			},
		})
	}

	//count page number
	numberOfPost, errPage := p.IPostRepository.CountAllPost()
	if errPage != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Jumlah data per page
	pageSize := 20

	// Hitung jumlah page dengan pembagian sederhana
	numberOfPage := math.Ceil(float64(numberOfPost) / float64(pageSize))

	// Jika ada sisa, tambahkan 1 page untuk menampung sisa data tersebut
	if numberOfPost%pageSize != 0 {
		numberOfPage++
	}

	return result, int(numberOfPage), nil
}

func (p *postServices) SuspendPost(token dto.Token, postId int) error {
	//check user
	user, errGetUserByUsername := p.GetUserByUsername(token.Username)
	if errGetUserByUsername != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetUserByUsername.Error())
	}

	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Admin access only")
	}

	//find post
	post, errGetPostById := p.IPostRepository.GetPostById(postId)
	if errors.Is(errGetPostById, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if errGetPostById != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGetPostById.Error())
	}

	if post.IsActive {
		post.IsActive = false
	} else {
		post.IsActive = true
	}

	err := p.IPostRepository.SavePost(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
