package v1

import (
	// "discusiin/controllers/topics"

	"discusiin/configs"
	"discusiin/controllers/bookmarks"
	"discusiin/controllers/comments"
	"discusiin/controllers/dashboard"
	"discusiin/controllers/followedPosts"
	"discusiin/controllers/likes"
	"discusiin/controllers/posts"
	"discusiin/controllers/replies"
	"discusiin/controllers/topics"
	"discusiin/controllers/users"
	mid "discusiin/middleware"
	"discusiin/routes"
	"io"
	"net/http"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(payload *routes.Payload) (*echo.Echo, io.Closer) {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	mid.LogMiddleware(e)
	e.Use(middleware.Recover())
	cors := middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPatch,
				http.MethodPost,
				http.MethodDelete},
			AllowHeaders: []string{
				"Accept",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"Origin",
			},
		})
	e.Use(cors)

	trace := jaegertracing.New(e, nil)

	dHandler := dashboard.DashboardHandler{
		IDashboardServices: payload.GetDashboardServices(),
	}
	uHandler := users.UserHandler{
		IUserServices: payload.GetUserServices(),
	}
	tHandler := topics.TopicHandler{
		ITopicServices: payload.GetTopicServices(),
	}
	pHandler := posts.PostHandler{
		IPostServices: payload.GetPostServices(),
	}
	cHandler := comments.CommentHandler{
		ICommentServices: payload.GetCommentServices(),
	}
	rHandler := replies.ReplyHandler{
		IReplyServices: payload.GetReplyServices(),
	}
	lHandler := likes.LikeHandler{
		ILikeServices: payload.GetLikeServices(),
	}
	fHandler := followedPosts.FollowedPostHandler{
		IFollowedPostServices: payload.GetFollowedPostServices(),
	}
	bHandler := bookmarks.BookmarkHandler{
		IBookmarkServices: payload.GetBookmarkServices(),
	}

	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "https://app.swaggerhub.com/apis/MIFTAHERS_2/Diskusiin-API/Latest#/")
	})
	api := e.Group("/api")
	v1 := api.Group("/v1")

	dashboard := v1.Group("/dashboard")
	dashboard.GET("", dHandler.GetTotalCountOfUserAndTopicAndPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints auth
	auth := v1.Group("/auth")
	auth.POST("/register", uHandler.Register)
	auth.POST("/register/admin", uHandler.RegisterAdmin, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	auth.POST("/login", uHandler.Login)

	//endpoints users
	users := v1.Group("/users")
	users.GET("", uHandler.GetUsersAdminNotIncluded, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	users.GET("/:user-id/post", uHandler.GetPostByUserIdForAdmin, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	users.GET("/:user-id/comment", uHandler.GetCommentByUserIdForAdmin, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	users.GET("/post", uHandler.GetPostByUserIdAsUser, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	users.PUT("/:user-id/ban", uHandler.BanUser, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	users.DELETE("/:user-id", uHandler.DeleteUser, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	profile := users.Group("/profile")
	profile.PUT("/edit", uHandler.UpdateProfile, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	profile.GET("", uHandler.GetProfile, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints topics
	topics := v1.Group("/topics")
	topics.GET("", tHandler.GetAllTopics, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.GET("/:topic-id", tHandler.GetTopic, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.GET("/top", tHandler.GetTopTopics, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.GET("/:topic-name/count", tHandler.GetNumberOfPostOnATopicByTopicName, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.POST("", tHandler.CreateNewTopic, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.PUT("/:topic-id", tHandler.UpdateTopicDescription, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.DELETE("/:topic-id", tHandler.DeleteTopic, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints posts
	posts := v1.Group("/posts")
	posts.GET("/:post-id", pHandler.GetPostByPostID, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.POST("/:topic-name", pHandler.CreateNewPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.PUT("/:post-id/suspend", pHandler.SuspendPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.PUT("/:post-id", pHandler.EditPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.DELETE("/:post-id", pHandler.DeletePost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	postAll := posts.Group("/all")
	postAll.GET("/:topic-name", pHandler.GetAllPostByTopicName, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	postAll.GET("/:topic-name/top", pHandler.GetAllPostByTopicByLike, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	recentPost := posts.Group("/recents")
	recentPost.GET("", pHandler.GetAllRecentPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	recentPost.GET("/top", pHandler.GetAllPostSortByLike, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoint followedPost
	followedPosts := posts.Group("/followed-posts")
	followedPosts.GET("/all", fHandler.GetAllFollowedPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	followedPosts.POST("/:post-id", fHandler.AddFollowedPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	followedPosts.DELETE("/:post-id", fHandler.DeleteFollowedPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints comments
	comments := posts.Group("/comments")
	comments.GET("/:post-id", cHandler.GetAllCommentByPostID)
	comments.POST("/:post-id", cHandler.CreateComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	comments.PUT("/:comment-id", cHandler.UpdateComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	comments.DELETE("/:comment-id", cHandler.DeleteComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints reply
	replys := comments.Group("/replies")
	replys.GET("/:comment-id", rHandler.GetAllReplyByCommentID)
	replys.POST("/:comment-id", rHandler.CreateReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	replys.PUT("/:reply-id", rHandler.UpdateReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	replys.DELETE("/:reply-id", rHandler.DeleteReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoint Like
	posts.PUT("/like/:post-id", lHandler.LikePost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.PUT("/dislike/:post-id", lHandler.DislikePost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoint bookmark
	bookmarks := posts.Group("/bookmarks")
	bookmarks.GET("/all", bHandler.GetAllBookmark, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	bookmarks.POST("/:post-id", bHandler.AddBookmark, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	bookmarks.DELETE("/:bookmark-id", bHandler.DeleteBookmark, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	return e, trace
}
