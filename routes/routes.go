package routes

import (
	"database/sql"
	"discusiin/configs"
	"discusiin/repositories"
	bookmarkGormRepo "discusiin/repositories/gorm/bookmarks"
	commentGormRepo "discusiin/repositories/gorm/comments"
	followedPostGormRepo "discusiin/repositories/gorm/followedPosts"
	likeGormRepo "discusiin/repositories/gorm/likes"
	postGormRepo "discusiin/repositories/gorm/posts"
	replyGormRepo "discusiin/repositories/gorm/replies"
	topicGormRepo "discusiin/repositories/gorm/topics"
	userGormRepo "discusiin/repositories/gorm/users"
	bService "discusiin/services/bookmarks"
	cService "discusiin/services/comments"
	dService "discusiin/services/dashboard"
	fService "discusiin/services/followedPosts"
	lService "discusiin/services/likes"
	pService "discusiin/services/posts"
	rService "discusiin/services/replies"
	tService "discusiin/services/topics"
	uService "discusiin/services/users"

	"gorm.io/gorm"
)

type Payload struct {
	Config           *configs.Config
	DBGorm           *gorm.DB
	DBSql            *sql.DB
	userRepo         repositories.IUserRepository
	topicRepo        repositories.ITopicRepository
	postRepo         repositories.IPostRepository
	commentRepo      repositories.ICommentRepository
	replyRepo        repositories.IReplyRepository
	likeRepo         repositories.ILikeRepository
	bookmarkRepo     repositories.IBookmarkRepository
	followedPostRepo repositories.IFollowedPostRepository
	dService         dService.IDashboardServices
	userService      uService.IUserServices
	topicService     tService.ITopicServices
	pService         pService.IPostServices
	cService         cService.ICommentServices
	rService         rService.IReplyServices
	lService         lService.ILikeServices
	fService         fService.IFollowedPostServices
	bService         bService.IBookmarkServices
}

// Init Repo -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitRepo() {
	p.InitUserRepo()
	p.InitTopicRepo()
	p.InitPostRepo()
	p.InitCommentRepo()
	p.InitReplyRepo()
	p.InitLikeRepo()
	p.InitBookmarkRepo()
	p.InitFollowedPostRepo()
}

func (p *Payload) InitUserRepo() {
	p.userRepo = userGormRepo.NewGorm(p.DBGorm)
	p.commentRepo = commentGormRepo.NewGorm(p.DBGorm)
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitTopicRepo() {
	p.userRepo = userGormRepo.NewGorm(p.DBGorm)
	p.topicRepo = topicGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitPostRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.topicRepo = topicGormRepo.NewGorm(p.DBGorm)
	p.userRepo = userGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitCommentRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.userRepo = userGormRepo.NewGorm(p.DBGorm)
	p.commentRepo = commentGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitReplyRepo() {
	p.commentRepo = commentGormRepo.NewGorm(p.DBGorm)
	p.replyRepo = replyGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitLikeRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.likeRepo = likeGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitBookmarkRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.bookmarkRepo = bookmarkGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitFollowedPostRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.followedPostRepo = followedPostGormRepo.NewGorm(p.DBGorm)
}
func (p *Payload) InitDashboardRepo() {
	p.postRepo = postGormRepo.NewGorm(p.DBGorm)
	p.topicRepo = topicGormRepo.NewGorm(p.DBGorm)
	p.userRepo = userGormRepo.NewGorm(p.DBGorm)
}

// User -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetUserServices() uService.IUserServices {
	if p.userService == nil {
		p.InitUserService()
	}
	return p.userService
}
func (p *Payload) InitUserService() {
	if p.userRepo == nil || p.commentRepo == nil || p.postRepo == nil {
		p.InitUserRepo()
	}

	p.userService = uService.NewUserServices(p.userRepo, p.commentRepo, p.postRepo)
}

// Topic -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetTopicServices() tService.ITopicServices {
	if p.topicService == nil {
		p.InitTopicService()
	}

	return p.topicService
}
func (p *Payload) InitTopicService() {
	if p.topicRepo == nil || p.userRepo == nil {
		p.InitTopicRepo()
	}

	p.topicService = tService.NewTopicServices(p.topicRepo, p.userRepo)
}

// Post -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetPostServices() pService.IPostServices {
	if p.pService == nil {
		p.InitPostService()
	}

	return p.pService
}

func (p *Payload) InitPostService() {
	if p.postRepo == nil || p.topicRepo == nil || p.userRepo == nil {
		p.InitPostRepo()
	}

	p.pService = pService.NewPostServices(p.topicRepo, p.postRepo, p.userRepo)
}

// Comment -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetCommentServices() cService.ICommentServices {
	if p.cService == nil {
		p.InitCommentService()
	}

	return p.cService
}
func (p *Payload) InitCommentService() {
	if p.commentRepo == nil || p.postRepo == nil || p.userRepo == nil {
		p.InitCommentRepo()
	}

	p.cService = cService.NewCommentServices(p.commentRepo, p.postRepo, p.userRepo)
}

// Reply -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetReplyServices() rService.IReplyServices {
	if p.rService == nil {
		p.InitReplyService()
	}

	return p.rService
}
func (p *Payload) InitReplyService() {
	if p.replyRepo == nil || p.commentRepo == nil {
		p.InitReplyRepo()
	}

	p.rService = rService.NewReplyServices(p.commentRepo, p.replyRepo)
}

// Like -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetLikeServices() lService.ILikeServices {
	if p.lService == nil {
		p.InitLikeService()
	}

	return p.lService
}
func (p *Payload) InitLikeService() {
	if p.likeRepo == nil || p.postRepo == nil {
		p.InitLikeRepo()
	}

	p.lService = lService.NewLikeServices(p.postRepo, p.likeRepo)
}

// Bookmark -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetBookmarkServices() bService.IBookmarkServices {
	if p.bService == nil {
		p.InitBookmarkService()
	}

	return p.bService
}
func (p *Payload) InitBookmarkService() {
	if p.bookmarkRepo == nil || p.postRepo == nil {
		p.InitBookmarkRepo()
	}

	p.bService = bService.NewBookmarkServices(p.bookmarkRepo, p.postRepo)
}

// FollowedPost -----------------------------------------------------------------------------------------------------------------

func (p *Payload) GetFollowedPostServices() fService.IFollowedPostServices {
	if p.fService == nil {
		p.InitFollowedPostService()
	}

	return p.fService
}
func (p *Payload) InitFollowedPostService() {
	if p.followedPostRepo == nil || p.postRepo == nil {
		p.InitFollowedPostRepo()
	}

	p.fService = fService.NewFollowedPostServices(p.postRepo, p.followedPostRepo)
}

// Dashboard -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetDashboardServices() dService.IDashboardServices {
	if p.dService == nil {
		p.InitDashboardService()
	}

	return p.dService
}
func (p *Payload) InitDashboardService() {
	if p.postRepo == nil || p.topicRepo == nil || p.userRepo == nil {
		p.InitPostRepo()
	}

	p.dService = dService.NewDashboardServices(p.userRepo, p.topicRepo, p.postRepo)
}
