package mcp

import (
	"net/http"
	"strings"

	"flec_blog/internal/service"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const publicServerInstructions = "Use article_list or article_search to locate articles before article_get. Create new content with article_create_draft. Use article_update_draft only for drafts and article_update_published only for published articles; omit fields that should be preserved. Call article_publish, article_unpublish, or article_delete only when the user explicitly intends that consequential action. Never guess article IDs."

type publicServer struct {
	articleService  *service.ArticleService
	categoryService *service.CategoryService
	tagService      *service.TagService
	commentService  *service.CommentService
	friendService   *service.FriendService
	rssFeedService  *service.RssFeedService
	momentService   *service.MomentService
	userService     *service.UserService
	statsService    *service.StatsService
}

func NewPublicHandler(
	articleService *service.ArticleService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	commentService *service.CommentService,
	friendService *service.FriendService,
	rssFeedService *service.RssFeedService,
	momentService *service.MomentService,
	userService *service.UserService,
	statsService *service.StatsService,
) http.Handler {
	implVersion := strings.TrimSpace(service.AppVersion)
	if implVersion == "" {
		implVersion = "dev"
	}

	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "flecblog-public",
		Version: implVersion,
	}, &sdkmcp.ServerOptions{
		Instructions: publicServerInstructions,
	})

	s := &publicServer{
		articleService:  articleService,
		categoryService: categoryService,
		tagService:      tagService,
		commentService:  commentService,
		friendService:   friendService,
		rssFeedService:  rssFeedService,
		momentService:   momentService,
		userService:     userService,
		statsService:    statsService,
	}

	// 注册 tools
	s.registerTools(server)

	streamableHandler := sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, nil)
	return streamableHandler
}
