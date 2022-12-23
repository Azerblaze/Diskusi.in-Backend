package likes

import (
	"discusiin/dto"
	"discusiin/models"
	"discusiin/repositories"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewLikeServices(postRepo repositories.IPostRepository, likeRepo repositories.ILikeRepository) ILikeServices {
	return &likeServices{IPostRepository: postRepo, ILikeRepository: likeRepo}
}

type ILikeServices interface {
	LikePost(token dto.Token, postId int) error
	DislikePost(token dto.Token, postId int) error
}

type likeServices struct {
	repositories.IPostRepository
	repositories.ILikeRepository
}

func (l *likeServices) LikePost(token dto.Token, postId int) error {
	var like models.Like
	//cek jika post ada
	post, err := l.IPostRepository.GetPostById(postId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var likenumber int //untuk membantu dalam perhitungan like

	//cek jika like ada
	oldLike, err := l.ILikeRepository.GetLikeByUserAndPostId(int(token.ID), postId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//jika tidak ada
		like.UserID = int(token.ID)
		like.PostID = postId
		like.IsLike = true

		//simpan data like baru
		errSaveLike := l.ILikeRepository.SaveNewLike(like)
		if errSaveLike != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errSaveLike.Error())
		}

		//update like count in post table
		post.LikeCount += 1
		errUpdateLikeCount := l.IPostRepository.SavePost(post)
		if errUpdateLikeCount != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errSaveLike.Error())
		}

		return nil

	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//jika like ada
	if oldLike.IsDislike { //jika sudah dislike
		//simpan data baru
		oldLike.IsDislike = false //in case like di klik lagi
		oldLike.IsLike = true
		likenumber = 1
	} else { //jika tidak di dislike
		likenumber = 0
		if oldLike.IsLike { //dan sudah di like
			oldLike.IsLike = false
		} else {
			oldLike.IsLike = true
		}
	}

	//simpan like baru
	log.Println(oldLike)
	errLike := l.ILikeRepository.SaveLike(oldLike)
	if errLike != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errLike.Error())
	}

	//update post count
	if oldLike.IsLike {
		post.LikeCount = (post.LikeCount + 1) + likenumber
	} else {
		post.LikeCount -= 1
	}
	errUpdateLikeCount := l.IPostRepository.SavePost(post)
	if errUpdateLikeCount != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errUpdateLikeCount.Error())
	}

	return nil
}

func (l *likeServices) DislikePost(token dto.Token, postId int) error {

	var like models.Like
	// cek jika post ada
	post, err := l.IPostRepository.GetPostById(postId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var dislikenumber int //untuk membantu dalam perhitungan like

	//cek jika like ada
	oldLike, err := l.ILikeRepository.GetLikeByUserAndPostId(int(token.ID), postId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//jika tidak ada
		like.UserID = int(token.ID)
		like.PostID = postId
		like.IsDislike = true

		//simpan data like baru
		errSaveLike := l.ILikeRepository.SaveNewLike(like)
		if errSaveLike != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errSaveLike.Error())
		}

		//update like count in post table
		post.LikeCount -= 1
		errUpdateLikeCount := l.IPostRepository.SavePost(post)
		if errUpdateLikeCount != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errSaveLike.Error())
		}
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//jika ada
	if oldLike.IsLike { //jika sudah di dislike
		//simpan data baru
		oldLike.IsDislike = true //in case like di klik lagi supaya netral
		oldLike.IsLike = false
		dislikenumber = 1
	} else { //jika tidak di like
		dislikenumber = 0
		if oldLike.IsDislike { //dan jika sudah di dislike
			oldLike.IsDislike = false
		} else {
			oldLike.IsDislike = true
		}
	}

	//simpan like baru
	errLike := l.ILikeRepository.SaveLike(oldLike)
	if errLike != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errLike.Error())
	}

	//update post count
	if oldLike.IsDislike {
		post.LikeCount = (post.LikeCount - 1) - dislikenumber
	} else {
		post.LikeCount += 1
	}
	errUpdateLikeCount := l.IPostRepository.SavePost(post)
	if errUpdateLikeCount != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errUpdateLikeCount.Error())
	}

	return nil
}
