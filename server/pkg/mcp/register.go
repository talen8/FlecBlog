package mcp

import (
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerTools 注册 MCP Tools
func (s *publicServer) registerTools(server *sdkmcp.Server) {
	// 创建文章服务包装器
	articleWrapper := tools.NewArticleWrapper(s.articleService)

	// article_manage - 文章管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "article_manage",
		Description: "文章管理。action：list/get/create/update/delete。",
		InputSchema: tools.ArticleManageInputSchema(),
	}, articleWrapper.ManageArticle)

	// 创建分类/标签服务包装器
	taxonomyWrapper := tools.NewTaxonomyWrapper(s.categoryService, s.tagService, s.articleService)

	// taxonomy_manage - 分类/标签管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "taxonomy_manage",
		Description: "分类/标签管理。target：category/tag；action：list/create/update/delete/list_articles。",
		InputSchema: tools.TaxonomyManageInputSchema(),
	}, taxonomyWrapper.ManageTaxonomy)

	// 创建评论服务包装器
	commentWrapper := tools.NewCommentWrapper(s.commentService)

	// comment_manage - 评论管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "comment_manage",
		Description: "评论管理。action：list/get/toggle_status/delete。",
		InputSchema: tools.CommentManageInputSchema(),
	}, commentWrapper.ManageComment)

	// 创建友链服务包装器
	friendWrapper := tools.NewFriendWrapper(s.friendService)

	// friend_manage - 友链管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "friend_manage",
		Description: "友链管理。action：list/get/create/update/delete。",
		InputSchema: tools.FriendManageInputSchema(),
	}, friendWrapper.ManageFriend)

	// 创建 RSS 订阅服务包装器
	rssFeedWrapper := tools.NewRssFeedWrapper(s.rssFeedService)

	// rssfeed_manage - RSS订阅管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "rssfeed_manage",
		Description: "RSS订阅管理。action：list/mark_read/mark_all_read。",
		InputSchema: tools.RssFeedManageInputSchema(),
	}, rssFeedWrapper.ManageRssFeed)

	// 创建统计服务包装器
	statsWrapper := tools.NewStatsWrapper(s.statsService)

	// stats_query - 统计查询只读 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "stats_query",
		Description: "站点访问统计查询（只读）。action：dashboard/trend。",
		InputSchema: tools.StatsQueryInputSchema(),
	}, statsWrapper.QueryStats)

	// 创建动态服务包装器
	momentWrapper := tools.NewMomentWrapper(s.momentService)

	// moment_manage - 动态管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "moment_manage",
		Description: "动态管理。action：list/get/create/update/delete。",
		InputSchema: tools.MomentManageInputSchema(),
	}, momentWrapper.ManageMoment)

	// 创建用户服务包装器
	userWrapper := tools.NewUserWrapper(s.userService)

	// user_manage - 用户管理聚合 Tool
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "user_manage",
		Description: "用户管理。action：list/get/create/update/delete。",
		InputSchema: tools.UserManageInputSchema(),
	}, userWrapper.ManageUser)
}
