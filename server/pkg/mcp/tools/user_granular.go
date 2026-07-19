package tools

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultUserListPageSize = 20
	maxUserListPageSize     = 100
)

// UserListInput user_list 输入。
type UserListInput struct {
	Page           int    `json:"page,omitempty"`
	PageSize       int    `json:"page_size,omitempty"`
	Keyword        string `json:"keyword,omitempty"`
	Role           string `json:"role,omitempty"`
	IsEnabled      *bool  `json:"is_enabled,omitempty"`
	IsDeleted      *bool  `json:"is_deleted,omitempty"`
	LoginMethod    string `json:"login_method,omitempty"`
	LastLoginStart string `json:"last_login_start,omitempty"`
	LastLoginEnd   string `json:"last_login_end,omitempty"`
	StartTime      string `json:"start_time,omitempty"`
	EndTime        string `json:"end_time,omitempty"`
}

// UserListOutput user_list 输出。
type UserListOutput struct {
	Items    []UserItem `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// UserGetInput user_get 输入。
type UserGetInput struct {
	ID uint `json:"id"`
}

// UserGetOutput user_get 输出。
type UserGetOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserCreateInput user_create 输入。
type UserCreateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
	Badge    string `json:"badge,omitempty"`
	Website  string `json:"website,omitempty"`
	Role     string `json:"role"`
}

// UserCreateOutput user_create 输出。
type UserCreateOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserUpdateInput user_update 输入。
type UserUpdateInput struct {
	ID        uint    `json:"id"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
	Nickname  *string `json:"nickname,omitempty"`
	Avatar    *string `json:"avatar,omitempty"`
	Badge     *string `json:"badge,omitempty"`
	Website   *string `json:"website,omitempty"`
	Role      *string `json:"role,omitempty"`
	IsEnabled *bool   `json:"is_enabled,omitempty"`
}

// UserUpdateOutput user_update 输出。
type UserUpdateOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserDeleteInput user_delete 输入。
type UserDeleteInput struct {
	ID uint `json:"id"`
}

// UserDeleteOutput user_delete 输出。
type UserDeleteOutput struct {
	Deleted bool `json:"deleted"`
	ID      uint `json:"id"`
}

// ListUsersTool 分页查询用户。
func (w *UserWrapper) ListUsersTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserListInput,
) (*sdkmcp.CallToolResult, UserListOutput, error) {
	if input.Role != "" {
		if _, err := strictUserRole(input.Role); err != nil {
			return nil, UserListOutput{}, err
		}
	}
	if input.LoginMethod != "" && !validUserLoginMethod(input.LoginMethod) {
		return nil, UserListOutput{}, fmt.Errorf("不支持的 login_method %q", input.LoginMethod)
	}

	page, pageSize := normalizeUserListPage(input.Page, input.PageSize)
	users, total, err := w.userService.List(&dto.ListUsersRequest{
		Page:           page,
		PageSize:       pageSize,
		Keyword:        input.Keyword,
		Role:           input.Role,
		IsEnabled:      input.IsEnabled,
		IsDeleted:      input.IsDeleted,
		LoginMethod:    input.LoginMethod,
		LastLoginStart: input.LastLoginStart,
		LastLoginEnd:   input.LastLoginEnd,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
	})
	if err != nil {
		return nil, UserListOutput{}, fmt.Errorf("获取用户列表失败: %w", err)
	}
	items := make([]UserItem, len(users))
	for i, user := range users {
		items[i] = convertToUserItem(user)
	}
	return nil, UserListOutput{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// GetUserTool 获取用户详情。
func (w *UserWrapper) GetUserTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserGetInput,
) (*sdkmcp.CallToolResult, UserGetOutput, error) {
	if input.ID == 0 {
		return nil, UserGetOutput{}, fmt.Errorf("用户 ID 不能为空")
	}
	item, err := w.loadAccurateUserDetail(input.ID)
	if err != nil {
		return nil, UserGetOutput{}, fmt.Errorf("获取用户失败: %w", err)
	}
	return nil, UserGetOutput{Item: item}, nil
}

// CreateUserTool 创建用户。
func (w *UserWrapper) CreateUserTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserCreateInput,
) (*sdkmcp.CallToolResult, UserCreateOutput, error) {
	email := strings.TrimSpace(input.Email)
	nickname := strings.TrimSpace(input.Nickname)
	if err := validateUserEmail(email); err != nil {
		return nil, UserCreateOutput{}, err
	}
	if len(input.Password) < 6 || len(input.Password) > 32 {
		return nil, UserCreateOutput{}, fmt.Errorf("密码长度必须为 6 到 32 个字符")
	}
	if nickname == "" {
		return nil, UserCreateOutput{}, fmt.Errorf("昵称不能为空")
	}
	role, err := strictUserRole(input.Role)
	if err != nil {
		return nil, UserCreateOutput{}, err
	}

	if err := w.userService.Create(mcpSuperAdminOperator(), &dto.AdminCreateUserRequest{
		Email:    email,
		Password: input.Password,
		Nickname: nickname,
		Avatar:   input.Avatar,
		Badge:    input.Badge,
		Website:  input.Website,
		Role:     role,
	}, ""); err != nil {
		return nil, UserCreateOutput{}, fmt.Errorf("创建用户失败: %w", err)
	}
	created, err := w.userService.GetByEmail(email)
	if err != nil {
		return nil, UserCreateOutput{}, fmt.Errorf("获取新建用户失败: %w", err)
	}
	item, err := w.loadAccurateUserDetail(created.ID)
	if err != nil {
		return nil, UserCreateOutput{}, fmt.Errorf("获取新建用户详情失败: %w", err)
	}
	return nil, UserCreateOutput{Item: item}, nil
}

// UpdateUserTool 更新用户。
func (w *UserWrapper) UpdateUserTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserUpdateInput,
) (*sdkmcp.CallToolResult, UserUpdateOutput, error) {
	if input.ID == 0 {
		return nil, UserUpdateOutput{}, fmt.Errorf("用户 ID 不能为空")
	}

	req := &dto.AdminUpdateUserRequest{IsEnabled: input.IsEnabled}
	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		if email != "" {
			if err := validateUserEmail(email); err != nil {
				return nil, UserUpdateOutput{}, err
			}
			req.Email = email
		}
	}
	if input.Password != nil {
		if *input.Password != "" && (len(*input.Password) < 6 || len(*input.Password) > 32) {
			return nil, UserUpdateOutput{}, fmt.Errorf("密码长度必须为 6 到 32 个字符")
		}
		req.Password = *input.Password
	}
	if input.Nickname != nil {
		req.Nickname = strings.TrimSpace(*input.Nickname)
	}
	if input.Avatar != nil {
		req.Avatar = *input.Avatar
	}
	if input.Badge != nil {
		req.Badge = *input.Badge
	}
	if input.Website != nil {
		req.Website = *input.Website
	}
	if input.Role != nil && strings.TrimSpace(*input.Role) != "" {
		role, err := strictUserRole(*input.Role)
		if err != nil {
			return nil, UserUpdateOutput{}, err
		}
		req.Role = role
	}

	if err := w.userService.Update(mcpSuperAdminOperator(), input.ID, req); err != nil {
		return nil, UserUpdateOutput{}, fmt.Errorf("更新用户失败: %w", err)
	}
	item, err := w.loadAccurateUserDetail(input.ID)
	if err != nil {
		return nil, UserUpdateOutput{}, fmt.Errorf("获取更新后用户失败: %w", err)
	}
	return nil, UserUpdateOutput{Item: item}, nil
}

// DeleteUserTool 软删除用户。
func (w *UserWrapper) DeleteUserTool(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserDeleteInput,
) (*sdkmcp.CallToolResult, UserDeleteOutput, error) {
	if input.ID == 0 {
		return nil, UserDeleteOutput{}, fmt.Errorf("用户 ID 不能为空")
	}
	if err := w.userService.Delete(mcpSuperAdminOperator(), input.ID); err != nil {
		return nil, UserDeleteOutput{}, fmt.Errorf("删除用户失败: %w", err)
	}
	return nil, UserDeleteOutput{Deleted: true, ID: input.ID}, nil
}

func (w *UserWrapper) loadAccurateUserDetail(id uint) (UserDetailItem, error) {
	user, err := w.userService.Get(id)
	if err != nil {
		return UserDetailItem{}, err
	}
	listItem, err := w.findCurrentUserListItem(id, user.Email)
	if err != nil {
		return UserDetailItem{}, err
	}
	item := convertToUserDetailItem(user)
	item.IsEnabled = listItem.IsEnabled
	return item, nil
}

func (w *UserWrapper) findCurrentUserListItem(id uint, email string) (dto.UserListResponse, error) {
	deleted := false
	for page := 1; ; page++ {
		users, total, err := w.userService.List(&dto.ListUsersRequest{
			Page:      page,
			PageSize:  100,
			Keyword:   email,
			IsDeleted: &deleted,
		})
		if err != nil {
			return dto.UserListResponse{}, err
		}
		for _, user := range users {
			if user.ID == id {
				return user, nil
			}
		}
		if int64(page*100) >= total {
			break
		}
	}
	return dto.UserListResponse{}, fmt.Errorf("用户 %d 不存在或已删除", id)
}

func strictUserRole(raw string) (model.UserRole, error) {
	switch strings.TrimSpace(raw) {
	case string(model.RoleSuperAdmin):
		return model.RoleSuperAdmin, nil
	case string(model.RoleAdmin):
		return model.RoleAdmin, nil
	case string(model.RoleUser):
		return model.RoleUser, nil
	case string(model.RoleGuest):
		return model.RoleGuest, nil
	default:
		return "", fmt.Errorf("不支持的用户角色 %q", raw)
	}
}

func validUserLoginMethod(value string) bool {
	switch value {
	case "password", "github", "google", "qq", "microsoft", "oidc":
		return true
	default:
		return false
	}
}

func validateUserEmail(value string) error {
	if value == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	address, err := mail.ParseAddress(value)
	if err != nil || !strings.EqualFold(address.Address, value) {
		return fmt.Errorf("邮箱格式无效")
	}
	return nil
}

func normalizeUserListPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultUserListPageSize
	}
	if pageSize > maxUserListPageSize {
		pageSize = maxUserListPageSize
	}
	return page, pageSize
}
