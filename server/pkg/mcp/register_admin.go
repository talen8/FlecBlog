package mcp

import (
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *publicServer) registerAdminTools(server *sdkmcp.Server) {
	// 用户管理：保留原 user_manage 的完整能力，但拆为细粒度工具并绑定本地用户身份。
	userWrapper := tools.NewUserWrapperWithStaticOperatorUserIDProvider(
		s.userService,
		s.staticOperatorUserIDProvider,
	)

	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "user_list",
		Annotations: readOnlyMCPToolAnnotations("用户列表"),
		Description: "分页查询用户，可按关键词、角色、启用状态、删除状态、登录方式和时间范围筛选。",
	}, mcpauth.ScopeAdmin, requireMCPAdminToolsEnabled(s.adminToolsEnabledProvider, userWrapper.ListUsersTool))

	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "user_get",
		Annotations: readOnlyMCPToolAnnotations("用户详情"),
		Description: "按用户 ID 获取用户详情。",
	}, mcpauth.ScopeAdmin, requireMCPAdminToolsEnabled(s.adminToolsEnabledProvider, userWrapper.GetUserTool))

	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "user_create",
		Annotations: mutatingMCPToolAnnotations("创建用户", false, false),
		Description: "创建用户并指定角色。",
	}, mcpauth.ScopeAdmin, requireMCPAdminToolsEnabled(s.adminToolsEnabledProvider, userWrapper.CreateUserTool))

	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "user_update",
		Annotations: mutatingMCPToolAnnotations("更新用户", true, false),
		Description: "更新用户资料、角色、启用状态或密码；字段省略时保留原值。",
	}, mcpauth.ScopeAdmin, requireMCPAdminToolsEnabled(s.adminToolsEnabledProvider, userWrapper.UpdateUserTool))

	addScopedMCPTool(server, &sdkmcp.Tool{
		Name:        "user_delete",
		Annotations: mutatingMCPToolAnnotations("删除用户", true, false),
		Description: "按用户 ID 软删除用户。高风险操作。",
	}, mcpauth.ScopeAdmin, requireMCPAdminToolsEnabled(s.adminToolsEnabledProvider, userWrapper.DeleteUserTool))
}
