package routes

import (
	"database/sql"
	"discusiin/configs"
	"discusiin/repositories"
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
	Config   *configs.Config
	DBGorm   *gorm.DB
	DBSql    *sql.DB
	repoSql  repositories.IDatabase
	dService dService.IDashboardServices
	uService uService.IUserServices
	tService tService.ITopicServices
	pService pService.IPostServices
	cService cService.ICommentServices
	rService rService.IReplyServices
	lService lService.ILikeServices
	fService fService.IFollowedPostServices
	bService bService.IBookmarkServices
}

func (p *Payload) InitRepoMysql() {
	p.repoSql = repositories.NewGorm(p.DBGorm)
}

// User -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetUserServices() uService.IUserServices {
	if p.uService == nil {
		p.InitUserService()
	}
	return p.uService
}
func (p *Payload) InitUserService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.uService = uService.NewUserServices(p.repoSql)
}

// Dashboard -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetDashboardServices() dService.IDashboardServices {
	if p.dService == nil {
		p.InitDashboardService()
	}

	return p.dService
}

func (p *Payload) InitDashboardService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.dService = dService.NewDashboardServices(p.repoSql)
}

// Topic -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetTopicServices() tService.ITopicServices {
	if p.tService == nil {
		p.InitTopicService()
	}

	return p.tService
}

func (p *Payload) InitTopicService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.tService = tService.NewTopicServices(p.repoSql)
}

// Post -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetPostServices() pService.IPostServices {
	if p.pService == nil {
		p.InitPostService()
	}

	return p.pService
}

func (p *Payload) InitPostService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.pService = pService.NewPostServices(p.repoSql)
}

// Comment -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetCommentServices() cService.ICommentServices {
	if p.cService == nil {
		p.InitCommentService()
	}

	return p.cService
}

func (p *Payload) InitCommentService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.cService = cService.NewCommentServices(p.repoSql)
}

// Reply -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetReplyServices() rService.IReplyServices {
	if p.rService == nil {
		p.InitReplyService()
	}

	return p.rService
}

func (p *Payload) InitReplyService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.rService = rService.NewReplyServices(p.repoSql)
}

// Like -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetLikeServices() lService.ILikeServices {
	if p.lService == nil {
		p.InitLikeService()
	}

	return p.lService
}

func (p *Payload) InitLikeService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.lService = lService.NewLikeServices(p.repoSql)
}

// Bookmark -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetBookmarkServices() bService.IBookmarkServices {
	if p.bService == nil {
		p.InitBookmarkService()
	}

	return p.bService
}

func (p *Payload) InitBookmarkService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.bService = bService.NewBookmarkServices(p.repoSql)
}

// FollowedPost -----------------------------------------------------------------------------------------------------------------
func (p *Payload) GetFollowedPostServices() fService.IFollowedPostServices {
	if p.fService == nil {
		p.InitFollowedPostService()
	}

	return p.fService
}

func (p *Payload) InitFollowedPostService() {
	if p.repoSql == nil {
		p.InitRepoMysql()
	}

	p.fService = fService.NewFollowedPostServices(p.repoSql)
}
