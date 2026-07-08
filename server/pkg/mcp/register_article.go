package mcp

import (
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *publicServer) registerArticleTools(server *sdkmcp.Server) {
	// 创建文章服务包装器
	articleWrapper := tools.NewArticleWrapper(s.articleService)

	// article_list - 安全分页列出文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_list",
		Annotations: readOnlyMCPToolAnnotations("文章列表"),
		Description: "分页列出后台文章，可按发布状态筛选。page_size 默认 20，最大 100；不返回正文。",
	}, mcpauth.ScopeRead, articleWrapper.ListArticleTool)

	// article_get - 获取完整文章详情
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_get",
		Annotations: readOnlyMCPToolAnnotations("文章详情"),
		Description: "按文章 ID 获取完整后台文章详情，包含正文、摘要、分类和标签。",
	}, mcpauth.ScopeRead, articleWrapper.GetArticleTool)

	// article_search - 搜索后台文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_search",
		Annotations: readOnlyMCPToolAnnotations("文章搜索"),
		Description: "按关键词搜索后台文章标题和正文，覆盖草稿与已发布文章；可按发布状态筛选。page_size 默认 20，最大 100。",
	}, mcpauth.ScopeRead, articleWrapper.SearchArticleTool)

	// article_create_draft - 创建草稿文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_create_draft",
		Annotations: mutatingMCPToolAnnotations("创建文章草稿", false, false),
		Description: "创建文章草稿。该工具不会发布文章，也不接受发布状态参数。",
	}, mcpauth.ScopeDraft, articleWrapper.CreateArticleDraftTool)

	// article_update_draft - 更新草稿文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_update_draft",
		Annotations: mutatingMCPToolAnnotations("更新文章草稿", true, false),
		Description: "更新现有草稿文章；拒绝修改已发布文章。字段省略时保留原值。",
	}, mcpauth.ScopeDraft, articleWrapper.UpdateArticleDraftTool)

	// article_publish - 显式发布文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_publish",
		Annotations: mutatingMCPToolAnnotations("发布文章", true, true),
		Description: "按文章 ID 显式发布文章。重复调用保持幂等；该工具不修改正文或其他编辑字段。",
	}, mcpauth.ScopePublish, articleWrapper.PublishArticleTool)

	// article_update_published - 更新已发布文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_update_published",
		Annotations: mutatingMCPToolAnnotations("更新已发布文章", true, false),
		Description: "更新现有已发布文章；拒绝修改草稿。字段省略时保留原值，不改变发布状态。",
	}, mcpauth.ScopePublish, articleWrapper.UpdatePublishedArticleTool)

	// article_unpublish - 显式取消发布
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_unpublish",
		Annotations: mutatingMCPToolAnnotations("取消发布文章", true, true),
		Description: "按文章 ID 显式取消发布。重复调用保持幂等；不修改正文或其他编辑字段。",
	}, mcpauth.ScopePublish, articleWrapper.UnpublishArticleTool)

	// article_delete - 显式删除文章
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_delete",
		Annotations: mutatingMCPToolAnnotations("删除文章", true, false),
		Description: "按文章 ID 永久删除文章。不可恢复，仅在用户明确要求删除时调用。",
	}, mcpauth.ScopeAdmin, articleWrapper.DeleteArticleTool)

	// article_image_upload - 安全上传文章图片
	imageWrapper := tools.NewArticleImageWrapper(s.articleImageUploader)
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "article_image_upload",
		Annotations: mutatingMCPToolAnnotations("上传文章图片", false, false),
		Description: "上传文章图片。输入 image_base64 可为原始 base64 或 data URL；仅允许 PNG/JPEG/WebP，最大 5 MiB，并校验实际格式与像素尺寸。",
	}, mcpauth.ScopeDraft, imageWrapper.UploadArticleImageTool)
}
