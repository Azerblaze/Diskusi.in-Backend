package routes

import (
	"database/sql"
	"discusiin/configs"
	"discusiin/repositories/bookmarks"
	"discusiin/repositories/comments"
	"discusiin/repositories/followedPosts"
	"discusiin/repositories/likes"
	"discusiin/repositories/posts"
	"discusiin/repositories/replies"
	"discusiin/repositories/topics"
	"discusiin/repositories/users"
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
	Config              *configs.Config
	DBGorm              *gorm.DB
	DBSql               *sql.DB
	repoUserSql         users.IUserDatabase
	repoTopicSql        topics.ITopicDatabase
	repoPostSql         posts.IPostDatabase
	repoCommentSql      comments.ICommentDatabase
	repoReplySql        replies.IReplyDatabase
	repoLikeSql         likes.ILikeDatabase
	repoBookmarkSql     bookmarks.IBookmarkDatabase
	repoFollowedPostSql followedPosts.IFollowedPostDatabase
	dService            dService.IDashboardServices
	uService            uService.IUserServices
	tService            tService.ITopicServices
	pService            pService.IPostServices
	cService            cService.ICommentServices
	rService            rService.IReplyServices
	lService            lService.ILikeServices
	fService            fService.IFollowedPostServices
	bService            bService.IBookmarkServices
}

// User -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitUserRepoMysql() {
	p.repoUserSql = users.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
	p.repoCommentSql = comments.NewGorm(p.DBGorm)
}

func (p *Payload) GetUserServices() uService.IUserServices {
	if p.uService == nil {
		p.InitUserService()
	}
	return p.uService
}
func (p *Payload) InitUserService() {
	if p.repoUserSql == nil || p.repoPostSql == nil || p.repoCommentSql == nil {
		p.InitUserRepoMysql()
	}

	p.uService = uService.NewUserServices(p.repoUserSql, p.repoPostSql, p.repoCommentSql)
}

// Dashboard -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitDashboardRepoMysql() {
	p.repoUserSql = users.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
	p.repoTopicSql = topics.NewGorm(p.DBGorm)
}

func (p *Payload) GetDashboardServices() dService.IDashboardServices {
	if p.dService == nil {
		p.InitDashboardService()
	}

	return p.dService
}

func (p *Payload) InitDashboardService() {
	if p.repoUserSql == nil || p.repoPostSql == nil || p.repoTopicSql == nil {
		p.InitDashboardRepoMysql()
	}

	p.dService = dService.NewDashboardServices(p.repoUserSql, p.repoPostSql, p.repoTopicSql)
}

// Topic -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitTopicRepoMysql() {
	p.repoUserSql = users.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
	p.repoTopicSql = topics.NewGorm(p.DBGorm)
}

func (p *Payload) GetTopicServices() tService.ITopicServices {
	if p.tService == nil {
		p.InitTopicService()
	}

	return p.tService
}

func (p *Payload) InitTopicService() {
	if p.repoUserSql == nil || p.repoPostSql == nil || p.repoTopicSql == nil {
		p.InitTopicRepoMysql()
	}

	p.tService = tService.NewTopicServices(p.repoTopicSql, p.repoPostSql, p.repoUserSql)
}

// Post -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitPostRepoMysql() {
	p.repoUserSql = users.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
	p.repoTopicSql = topics.NewGorm(p.DBGorm)
}

func (p *Payload) GetPostServices() pService.IPostServices {
	if p.pService == nil {
		p.InitPostService()
	}

	return p.pService
}

func (p *Payload) InitPostService() {
	if p.repoUserSql == nil || p.repoPostSql == nil || p.repoTopicSql == nil {
		p.InitPostRepoMysql()
	}

	p.pService = pService.NewPostServices(p.repoPostSql, p.repoUserSql, p.repoTopicSql)
}

// Comment -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitCommentRepoMysql() {
	p.repoUserSql = users.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
	p.repoCommentSql = comments.NewGorm(p.DBGorm)
}

func (p *Payload) GetCommentServices() cService.ICommentServices {
	if p.cService == nil {
		p.InitCommentService()
	}

	return p.cService
}

func (p *Payload) InitCommentService() {
	if p.repoUserSql == nil || p.repoPostSql == nil || p.repoCommentSql == nil {
		p.InitUserRepoMysql()
	}

	p.cService = cService.NewCommentServices(p.repoCommentSql, p.repoPostSql, p.repoUserSql)
}

// Reply -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitReplyRepoMysql() {
	p.repoCommentSql = comments.NewGorm(p.DBGorm)
	p.repoReplySql = replies.NewGorm(p.DBGorm)
}

func (p *Payload) GetReplyServices() rService.IReplyServices {
	if p.rService == nil {
		p.InitReplyService()
	}

	return p.rService
}

func (p *Payload) InitReplyService() {
	if p.repoCommentSql == nil || p.repoReplySql == nil {
		p.InitReplyRepoMysql()
	}

	p.rService = rService.NewReplyServices(p.repoReplySql, p.repoCommentSql)
}

// Like -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitLikeRepoMysql() {
	p.repoLikeSql = likes.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
}

func (p *Payload) GetLikeServices() lService.ILikeServices {
	if p.lService == nil {
		p.InitLikeService()
	}

	return p.lService
}

func (p *Payload) InitLikeService() {
	if p.repoLikeSql == nil || p.repoPostSql == nil {
		p.InitLikeRepoMysql()
	}

	p.lService = lService.NewLikeServices(p.repoLikeSql, p.repoPostSql)
}

// Bookmark -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitBookmarkRepoMysql() {
	p.repoBookmarkSql = bookmarks.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
}

func (p *Payload) GetBookmarkServices() bService.IBookmarkServices {
	if p.bService == nil {
		p.InitBookmarkService()
	}

	return p.bService
}

func (p *Payload) InitBookmarkService() {
	if p.repoBookmarkSql == nil || p.repoPostSql == nil {
		p.InitBookmarkRepoMysql()
	}

	p.bService = bService.NewBookmarkServices(p.repoBookmarkSql, p.repoPostSql)
}

// FollowedPost -----------------------------------------------------------------------------------------------------------------
func (p *Payload) InitFollowedPostRepoMysql() {
	p.repoFollowedPostSql = followedPosts.NewGorm(p.DBGorm)
	p.repoPostSql = posts.NewGorm(p.DBGorm)
}

func (p *Payload) GetFollowedPostServices() fService.IFollowedPostServices {
	if p.fService == nil {
		p.InitFollowedPostService()
	}

	return p.fService
}

func (p *Payload) InitFollowedPostService() {
	if p.repoFollowedPostSql == nil || p.repoPostSql == nil {
		p.InitFollowedPostRepoMysql()
	}

	p.fService = fService.NewFollowedPostServices(p.repoFollowedPostSql, p.repoPostSql)
}
