package mcp

import (
	"net/http"
	"strings"

	"flec_blog/internal/service"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const publicServerInstructions = "Use article_list or article_search to locate articles before article_get. Create new content with article_create_draft. Use article_update_draft only for drafts and article_update_published only for published articles; omit fields that should be preserved. Call article_publish, article_unpublish, or article_delete only when the user explicitly intends that consequential action. Use article_image_upload for PNG/JPEG/WebP and insert the returned URL into content or cover. Never guess article IDs."

type publicServer struct {
	articleService       *service.ArticleService
	categoryService      *service.CategoryService
	tagService           *service.TagService
	commentService       *service.CommentService
	friendService        *service.FriendService
	rssFeedService       *service.RssFeedService
	momentService        *service.MomentService
	userService          *service.UserService
	statsService         *service.StatsService
	articleImageUploader tools.ArticleImageUploader
}

// PublicHandlerOptions configures optional runtime providers without breaking the legacy constructor.
type PublicHandlerOptions struct {
	ArticleImageUploader tools.ArticleImageUploader
	MaxSessions          int
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
	return NewPublicHandlerWithOptions(
		articleService, categoryService, tagService, commentService, friendService, rssFeedService, momentService,
		userService, statsService, PublicHandlerOptions{},
	)
}

func NewPublicHandlerWithOptions(
	articleService *service.ArticleService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	commentService *service.CommentService,
	friendService *service.FriendService,
	rssFeedService *service.RssFeedService,
	momentService *service.MomentService,
	userService *service.UserService,
	statsService *service.StatsService,
	options PublicHandlerOptions,
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
		articleService:       articleService,
		categoryService:      categoryService,
		tagService:           tagService,
		commentService:       commentService,
		friendService:        friendService,
		rssFeedService:       rssFeedService,
		momentService:        momentService,
		userService:          userService,
		statsService:         statsService,
		articleImageUploader: options.ArticleImageUploader,
	}

	// 注册 tools
	s.registerTools(server)

	streamableHandler := sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, nil)
	return limitMCPRequestBody(
		newSessionAdmissionHandler(streamableHandler, options.MaxSessions),
		defaultMaxMCPRequestBodyBytes,
	)
}
