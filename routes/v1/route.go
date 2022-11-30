package v1

import (
	// "discusiin/controllers/topics"

	"discusiin/configs"
	"discusiin/controllers/comments"
	"discusiin/controllers/posts"
	"discusiin/controllers/replys"
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
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{"Content-Type", "X-CSRF-Token"},
	})
	e.Pre(cors)

	trace := jaegertracing.New(e, nil)

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

	rHandler := replys.ReplyHandler{
		IReplyServices: payload.GetReplyServices(),
	}

	api := e.Group("/api")
	v1 := api.Group("/v1")

	//endpoints users
	users := v1.Group("/users")
	users.POST("/register", uHandler.Register)
	users.POST("/login", uHandler.Login)

	//endpoints topics
	topics := v1.Group("/topics")
	topics.POST("/create", tHandler.CreateNewTopic, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.PUT("/edit_description/:topic_id", tHandler.UpdateTopicDescription, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	topics.GET("", tHandler.GetAllTopics)
	topics.GET("/:topic_id", tHandler.GetTopic)
	topics.DELETE("/delete/:topic_id", tHandler.DeleteTopic, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints posts
	posts := v1.Group("/posts")
	posts.POST("/create/:topic_name", pHandler.CreateNewPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.GET("/all/:topic_name", pHandler.GetAllPost)
	posts.GET("/:post_id", pHandler.GetPost)
	posts.PUT("/edit/:post_id", pHandler.EditPost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	posts.DELETE("/delete/:post_id", pHandler.DeletePost, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints comments
	comments := posts.Group("/comments")
	comments.GET("/:post_id", cHandler.GetAllComment)
	comments.POST("/create/:post_id", cHandler.CreateComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	comments.PUT("/edit/:comment_id", cHandler.UpdateComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	comments.DELETE("/delete/:comment_id", cHandler.DeleteComment, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	//endpoints reply
	replys := comments.Group("/replys")
	replys.GET("/:comment_id", rHandler.GetAllReply)
	replys.POST("/create/:comment_id", rHandler.CreateReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	replys.PUT("/edit/:reply_id", rHandler.UpdateReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))
	replys.DELETE("/delete/:reply_id", rHandler.DeleteReply, middleware.JWT([]byte(configs.Cfg.TokenSecret)))

	return e, trace
}
