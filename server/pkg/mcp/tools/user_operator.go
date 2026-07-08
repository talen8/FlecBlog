package tools

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	mcpauth "flec_blog/pkg/mcp/auth"
)

const staticOperatorUserIDEnv = "MCP_STATIC_OPERATOR_USER_ID"

type userOperatorResolver struct {
	userService                  *service.UserService
	staticOperatorUserIDProvider func() uint
}

func newUserOperatorResolver(
	userService *service.UserService,
	staticOperatorUserIDProvider func() uint,
) *userOperatorResolver {
	return &userOperatorResolver{
		userService:                  userService,
		staticOperatorUserIDProvider: staticOperatorUserIDProvider,
	}
}

func (r *userOperatorResolver) resolve(ctx context.Context) (*model.User, error) {
	if r == nil || r.userService == nil {
		return nil, fmt.Errorf("用户服务不可用")
	}
	principal, ok := mcpauth.PrincipalFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("MCP authorization principal is missing")
	}

	var operator *model.User
	var err error
	switch principal.Method {
	case "static":
		operator, err = r.resolveStaticOperator()
	case "oauth":
		operator, err = r.resolveOAuthOperator(principal)
	default:
		return nil, fmt.Errorf("不支持的 MCP principal method %q", principal.Method)
	}
	if err != nil {
		return nil, err
	}
	if operator.Role != model.RoleAdmin && operator.Role != model.RoleSuperAdmin {
		return nil, fmt.Errorf("MCP operator %d 不是本地管理员", operator.ID)
	}
	return operator, nil
}

func (r *userOperatorResolver) resolveStaticOperator() (*model.User, error) {
	rawID := strings.TrimSpace(os.Getenv(staticOperatorUserIDEnv))
	if rawID != "" {
		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil || id == 0 || id > math.MaxUint {
			return nil, fmt.Errorf("%s 必须是有效的正整数用户 ID", staticOperatorUserIDEnv)
		}
		return r.findActiveUserByID(uint(id))
	}

	if r.staticOperatorUserIDProvider != nil {
		if configuredUserID := r.staticOperatorUserIDProvider(); configuredUserID != 0 {
			return r.findActiveUserByID(configuredUserID)
		}
	}

	enabled := true
	deleted := false
	users, total, err := r.userService.List(&dto.ListUsersRequest{
		Page:      1,
		PageSize:  100,
		Role:      string(model.RoleSuperAdmin),
		IsEnabled: &enabled,
		IsDeleted: &deleted,
	})
	if err != nil {
		return nil, fmt.Errorf("解析静态 MCP operator 失败: %w", err)
	}
	if total == 0 {
		return nil, fmt.Errorf("未找到启用且未删除的真实超级管理员")
	}
	if total != 1 || len(users) != 1 {
		return nil, fmt.Errorf("存在多个真实超级管理员，请显式配置 %s", staticOperatorUserIDEnv)
	}
	return modelUserFromListItem(users[0]), nil
}

func (r *userOperatorResolver) resolveOAuthOperator(principal *mcpauth.Principal) (*model.User, error) {
	if principal == nil {
		return nil, fmt.Errorf("OAuth principal 为空")
	}
	userID, ok := principal.VerifiedLocalUserID()
	if !ok {
		return nil, fmt.Errorf("OAuth principal 缺少已验证的本地用户绑定，无法解析管理员 operator")
	}
	return r.findActiveUserByID(userID)
}

func (r *userOperatorResolver) findActiveUserByID(id uint) (*model.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("operator 用户 ID 不能为空")
	}
	user, err := r.userService.Get(id)
	if err != nil {
		return nil, fmt.Errorf("获取 operator 用户 %d 失败: %w", id, err)
	}
	return r.findActiveUser(user.Email, func(item dto.UserListResponse) bool { return item.ID == id })
}

func (r *userOperatorResolver) findActiveUser(keyword string, match func(dto.UserListResponse) bool) (*model.User, error) {
	enabled := true
	deleted := false
	users, _, err := r.userService.List(&dto.ListUsersRequest{
		Page:      1,
		PageSize:  100,
		Keyword:   keyword,
		IsEnabled: &enabled,
		IsDeleted: &deleted,
	})
	if err != nil {
		return nil, fmt.Errorf("查询真实 MCP operator 失败: %w", err)
	}
	var matched []dto.UserListResponse
	for _, item := range users {
		if match(item) {
			matched = append(matched, item)
		}
	}
	if len(matched) != 1 {
		return nil, fmt.Errorf("无法唯一解析启用且未删除的真实 MCP operator")
	}
	return modelUserFromListItem(matched[0]), nil
}

func modelUserFromListItem(item dto.UserListResponse) *model.User {
	operator := &model.User{
		Email:     item.Email,
		Nickname:  item.Nickname,
		Role:      item.Role,
		IsEnabled: item.IsEnabled,
	}
	operator.ID = item.ID
	return operator
}
