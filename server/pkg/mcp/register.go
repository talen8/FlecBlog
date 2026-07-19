package mcp

import (
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ============ 工具注解辅助函数 ============

func readOnlyMCPToolAnnotations(title string) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:         title,
		ReadOnlyHint:  true,
		OpenWorldHint: &closedWorld,
	}
}

func mutatingMCPToolAnnotations(title string, destructive, idempotent bool) *sdkmcp.ToolAnnotations {
	closedWorld := false
	return &sdkmcp.ToolAnnotations{
		Title:           title,
		DestructiveHint: &destructive,
		IdempotentHint:  idempotent,
		OpenWorldHint:   &closedWorld,
	}
}

// ============ 注册入口 ============

func (s *publicServer) registerTools(server *sdkmcp.Server) {
	// 文章管理
	articleWrapper := tools.NewArticleWrapper(s.articleService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_list",
		Annotations: readOnlyMCPToolAnnotations("文章列表"),
		Description: "分页列出后台文章，可按发布状态筛选。page_size 默认 20，最大 100；不返回正文。",
	}, articleWrapper.ListArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_get",
		Annotations: readOnlyMCPToolAnnotations("文章详情"),
		Description: "按文章 ID 获取完整后台文章详情，包含正文、摘要、分类和标签。",
	}, articleWrapper.GetArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_search",
		Annotations: readOnlyMCPToolAnnotations("文章搜索"),
		Description: "按关键词搜索后台文章标题和正文，覆盖草稿与已发布文章；可按发布状态筛选。page_size 默认 20，最大 100。",
	}, articleWrapper.SearchArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_create_draft",
		Annotations: mutatingMCPToolAnnotations("创建文章草稿", false, false),
		Description: "创建文章草稿。该工具不会发布文章，也不接受发布状态参数。",
	}, articleWrapper.CreateArticleDraftTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_update_draft",
		Annotations: mutatingMCPToolAnnotations("更新文章草稿", true, false),
		Description: "更新现有草稿文章；拒绝修改已发布文章。字段省略时保留原值。",
	}, articleWrapper.UpdateArticleDraftTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_publish",
		Annotations: mutatingMCPToolAnnotations("发布文章", true, true),
		Description: "按文章 ID 显式发布文章。重复调用保持幂等；该工具不修改正文或其他编辑字段。",
	}, articleWrapper.PublishArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_update_published",
		Annotations: mutatingMCPToolAnnotations("更新已发布文章", true, false),
		Description: "更新现有已发布文章；拒绝修改草稿。字段省略时保留原值，不改变发布状态。",
	}, articleWrapper.UpdatePublishedArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_unpublish",
		Annotations: mutatingMCPToolAnnotations("取消发布文章", true, true),
		Description: "按文章 ID 显式取消发布。重复调用保持幂等；不修改正文或其他编辑字段。",
	}, articleWrapper.UnpublishArticleTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_delete",
		Annotations: mutatingMCPToolAnnotations("删除文章", true, false),
		Description: "按文章 ID 永久删除文章。不可恢复，仅在用户明确要求删除时调用。",
	}, articleWrapper.DeleteArticleTool)

	// 分类/标签管理
	taxonomyWrapper := tools.NewTaxonomyWrapper(s.categoryService, s.tagService, s.articleService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "taxonomy_manage",
		Annotations: mutatingMCPToolAnnotations("分类与标签管理", true, false),
		Description: "分类/标签管理。target：category/tag；action：list/create/update/delete/list_articles。",
		InputSchema: tools.TaxonomyManageInputSchema(),
	}, taxonomyWrapper.ManageTaxonomy)

	// 评论管理
	commentWrapper := tools.NewCommentWrapper(s.commentService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "comment_manage",
		Annotations: mutatingMCPToolAnnotations("评论管理", true, false),
		Description: "评论管理。action：list/get/toggle_status/delete。",
		InputSchema: tools.CommentManageInputSchema(),
	}, commentWrapper.ManageComment)

	// 友链管理
	friendWrapper := tools.NewFriendWrapper(s.friendService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "friend_manage",
		Annotations: mutatingMCPToolAnnotations("友链管理", true, false),
		Description: "友链管理。action：list/get/create/update/delete。",
		InputSchema: tools.FriendManageInputSchema(),
	}, friendWrapper.ManageFriend)

	// RSS 订阅管理
	rssFeedWrapper := tools.NewRssFeedWrapper(s.rssFeedService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "rssfeed_manage",
		Annotations: mutatingMCPToolAnnotations("RSS 阅读状态管理", false, true),
		Description: "RSS订阅管理。action：list/mark_read/mark_all_read。",
		InputSchema: tools.RssFeedManageInputSchema(),
	}, rssFeedWrapper.ManageRssFeed)

	// 站点统计查询
	statsWrapper := tools.NewStatsWrapper(s.statsService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "stats_query",
		Annotations: readOnlyMCPToolAnnotations("站点统计查询"),
		Description: "站点访问统计查询（只读）。action：dashboard/trend。",
		InputSchema: tools.StatsQueryInputSchema(),
	}, statsWrapper.QueryStats)

	// 动态管理
	momentWrapper := tools.NewMomentWrapper(s.momentService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "moment_manage",
		Annotations: mutatingMCPToolAnnotations("动态管理", true, false),
		Description: "动态管理。action：list/get/create/update/delete。",
		InputSchema: tools.MomentManageInputSchema(),
	}, momentWrapper.ManageMoment)

	// 用户管理
	userWrapper := tools.NewUserWrapper(s.userService)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_list",
		Annotations: readOnlyMCPToolAnnotations("用户列表"),
		Description: "分页查询用户，可按关键词、角色、启用状态、删除状态、登录方式和时间范围筛选。",
	}, userWrapper.ListUsersTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_get",
		Annotations: readOnlyMCPToolAnnotations("用户详情"),
		Description: "按用户 ID 获取用户详情。",
	}, userWrapper.GetUserTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_create",
		Annotations: mutatingMCPToolAnnotations("创建用户", false, false),
		Description: "创建用户并指定角色。",
	}, userWrapper.CreateUserTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_update",
		Annotations: mutatingMCPToolAnnotations("更新用户", true, false),
		Description: "更新用户资料、角色、启用状态或密码；字段省略时保留原值。",
	}, userWrapper.UpdateUserTool)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_delete",
		Annotations: mutatingMCPToolAnnotations("删除用户", true, false),
		Description: "按用户 ID 软删除用户。高风险操作。",
	}, userWrapper.DeleteUserTool)
}
