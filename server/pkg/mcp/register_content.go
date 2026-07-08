package mcp

import (
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *publicServer) registerContentTools(server *sdkmcp.Server) {
	// 创建分类/标签服务包装器
	taxonomyWrapper := tools.NewTaxonomyWrapper(s.categoryService, s.tagService, s.articleService)

	// taxonomy_manage - 分类/标签管理聚合 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "taxonomy_manage",
		Annotations: mutatingMCPToolAnnotations("分类与标签管理", true, false),
		Description: "分类/标签管理。target：category/tag；action：list/create/update/delete/list_articles。",
		InputSchema: tools.TaxonomyManageInputSchema(),
	}, mcpauth.ScopeManage, taxonomyWrapper.ManageTaxonomy)

	// 创建评论服务包装器
	commentWrapper := tools.NewCommentWrapper(s.commentService)

	// comment_manage - 评论管理聚合 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "comment_manage",
		Annotations: mutatingMCPToolAnnotations("评论管理", true, false),
		Description: "评论管理。action：list/get/toggle_status/delete。",
		InputSchema: tools.CommentManageInputSchema(),
	}, mcpauth.ScopeManage, commentWrapper.ManageComment)

	// 创建友链服务包装器
	friendWrapper := tools.NewFriendWrapper(s.friendService)

	// friend_manage - 友链管理聚合 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "friend_manage",
		Annotations: mutatingMCPToolAnnotations("友链管理", true, false),
		Description: "友链管理。action：list/get/create/update/delete。",
		InputSchema: tools.FriendManageInputSchema(),
	}, mcpauth.ScopeManage, friendWrapper.ManageFriend)

	// 创建 RSS 订阅服务包装器
	rssFeedWrapper := tools.NewRssFeedWrapper(s.rssFeedService)

	// rssfeed_manage - RSS订阅管理聚合 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "rssfeed_manage",
		Annotations: mutatingMCPToolAnnotations("RSS 阅读状态管理", false, true),
		Description: "RSS订阅管理。action：list/mark_read/mark_all_read。",
		InputSchema: tools.RssFeedManageInputSchema(),
	}, mcpauth.ScopeManage, rssFeedWrapper.ManageRssFeed)

	// 创建统计服务包装器
	statsWrapper := tools.NewStatsWrapper(s.statsService)

	// stats_query - 统计查询只读 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "stats_query",
		Annotations: readOnlyMCPToolAnnotations("站点统计查询"),
		Description: "站点访问统计查询（只读）。action：dashboard/trend。",
		InputSchema: tools.StatsQueryInputSchema(),
	}, mcpauth.ScopeRead, statsWrapper.QueryStats)

	// 创建动态服务包装器
	momentWrapper := tools.NewMomentWrapper(s.momentService)

	// moment_manage - 动态管理聚合 Tool
	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "moment_manage",
		Annotations: mutatingMCPToolAnnotations("动态管理", true, false),
		Description: "动态管理。action：list/get/create/update/delete。",
		InputSchema: tools.MomentManageInputSchema(),
	}, mcpauth.ScopeManage, momentWrapper.ManageMoment)
}
