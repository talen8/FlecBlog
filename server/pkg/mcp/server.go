package mcp

import (
	"net/http"
	"strings"

	"flec_blog/internal/service"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverInstructions = "本服务仅用于管理本站文章的标题与正文。先用 article_list 或 article_search 定位文章，再用 article_get 获取正文。用 article_create 新建文章，用 article_update 修改标题或正文，省略的字段保持原值。新建与更新只对草稿文章生效。article_delete 会永久删除草稿内容，仅在用户明确要求删除时调用。绝不猜测文章 ID。"

// ============ 工具注解辅助函数 ============

func readOnlyToolAnnotations(title string) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:         title,
		ReadOnlyHint:  true,
		OpenWorldHint: &closedWorld,
	}
}

func mutatingToolAnnotations(title string, destructive, idempotent bool) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:           title,
		DestructiveHint: &destructive,
		IdempotentHint:  idempotent,
		OpenWorldHint:   &closedWorld,
	}
}

// Server 文章管理 MCP Server
type Server struct {
	articleService *service.ArticleService
}

// NewHandler 创建文章管理 MCP Handler
func NewHandler(articleService *service.ArticleService) http.Handler {
	implVersion := strings.TrimSpace(service.AppVersion)
	if implVersion == "" {
		implVersion = "dev"
	}

	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "flecblog-article",
		Version: implVersion,
	}, &sdkmcp.ServerOptions{
		Instructions: serverInstructions,
	})

	s := &Server{
		articleService: articleService,
	}

	// 注册 tools
	s.registerTools(server)

	streamableHandler := sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, nil)
	return streamableHandler
}

// registerTools 注册文章管理工具
func (s *Server) registerTools(server *sdkmcp.Server) {
	wrapper := NewWrapper(s.articleService)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_list",
		Annotations: readOnlyToolAnnotations("文章列表"),
		Description: "分页列出文章，包含草稿与已发布文章，返回文章 ID、标题与发布状态，不含正文。草稿排在前面，已发布的按发布时间倒序。is_publish 可按发布状态筛选。page 从 1 开始，page_size 默认 20、最大 100。需要正文时用 article_get 获取。返回的 items 始终为数组。",
	}, wrapper.ListTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_get",
		Annotations: readOnlyToolAnnotations("文章详情"),
		Description: "按文章 ID 获取文章的标题与正文，草稿与已发布文章都可获取。修改文章前先调用它确认当前内容。",
	}, wrapper.GetTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_search",
		Annotations: readOnlyToolAnnotations("文章搜索"),
		Description: "按关键词搜索文章，匹配标题与正文，不区分大小写的部分匹配，草稿与已发布文章都会被搜到。keyword 必填。page 从 1 开始，page_size 默认 20、最大 100。返回文章 ID、标题与发布状态，不含正文；需要正文时用返回的 ID 调用 article_get。返回的 items 始终为数组。",
	}, wrapper.SearchTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_create",
		Annotations: mutatingToolAnnotations("创建文章", false, false),
		Description: "新建文章，只需提供标题与正文（Markdown）。返回新文章的 ID、标题与正文。",
	}, wrapper.CreateTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_update",
		Annotations: mutatingToolAnnotations("更新文章", true, false),
		Description: "按文章 ID 更新草稿文章的标题与正文，已发布文章会被拒绝。只改动传入的字段，省略的字段保持原值。返回更新后的标题与正文。",
	}, wrapper.UpdateTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_delete",
		Annotations: mutatingToolAnnotations("删除草稿文章", true, false),
		Description: "删除草稿文章，永久删除且不可恢复。已发布文章不能删除，需要在管理后台处理。id 来自 article_list 或 article_search 的返回结果，仅在用户明确要求删除某篇文章时调用。返回删除结果、文章 ID 与删除瞬间的标题正文副本，该副本仅供误删后参考重建，服务端不提供恢复功能。",
	}, wrapper.DeleteTool)
}
