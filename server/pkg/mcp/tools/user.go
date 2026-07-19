package tools

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	"flec_blog/pkg/utils"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	userActionList   = "list"
	userActionGet    = "get"
	userActionCreate = "create"
	userActionUpdate = "update"
	userActionDelete = "delete"

	defaultUserListPageSize = 20
	maxUserListPageSize     = 100
)

// ============ MCP 类型定义============

// UserItem 用户列表项
type UserItem struct {
	ID           uint     `json:"id"`
	Email        string   `json:"email"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Badge        string   `json:"badge"`
	Website      string   `json:"website"`
	Role         string   `json:"role"`
	IsEnabled    bool     `json:"is_enabled"`
	HasPassword  bool     `json:"has_password"`
	LinkedOAuths []string `json:"linked_oauths"`
	LastLogin    *string  `json:"last_login"`
	CreatedAt    *string  `json:"created_at"`
	DeletedAt    *string  `json:"deleted_at,omitempty"`
}

// UserDetailItem 用户详情项
type UserDetailItem struct {
	ID             uint     `json:"id"`
	Email          string   `json:"email"`
	EmailHash      string   `json:"email_hash"`
	IsVirtualEmail bool     `json:"is_virtual_email"`
	Nickname       string   `json:"nickname"`
	Avatar         string   `json:"avatar"`
	Badge          string   `json:"badge"`
	Website        string   `json:"website"`
	Role           string   `json:"role"`
	HasPassword    bool     `json:"has_password"`
	LinkedOAuths   []string `json:"linked_oauths"`
	LastLogin      *string  `json:"last_login,omitempty"`
	CreatedAt      *string  `json:"created_at,omitempty"`
	IsEnabled      bool     `json:"is_enabled"`
}

// ============ 聚合 Tool 输入/输出类型============

// UserManageInput user_manage 聚合 tool 输入
type UserManageInput struct {
	Action  string            `json:"action"`
	Payload UserManagePayload `json:"payload"`
}

// UserManagePayload user_manage 载荷
type UserManagePayload struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Badge     string `json:"badge"`
	Website   string `json:"website"`
	Role      string `json:"role"`
	IsEnabled *bool  `json:"is_enabled"`
}

// UserManageOutput user_manage 聚合 tool 输出
type UserManageOutput struct {
	List          []UserItem      `json:"list,omitempty"`
	Total         int64           `json:"total,omitempty"`
	Page          int             `json:"page,omitempty"`
	PageSize      int             `json:"page_size,omitempty"`
	Item          *UserDetailItem `json:"item,omitempty"`
	DeleteSuccess *bool           `json:"delete_success,omitempty"`
	ID            *uint           `json:"id,omitempty"`
	Error         string          `json:"error,omitempty"`
}

// ============ 细粒度 Tool 输入/输出类型============

// UserListInput user_list 输入
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

// UserListOutput user_list 输出
type UserListOutput struct {
	Items    []UserItem `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// UserGetInput user_get 输入
type UserGetInput struct {
	ID uint `json:"id"`
}

// UserGetOutput user_get 输出
type UserGetOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserCreateInput user_create 输入
type UserCreateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
	Badge    string `json:"badge,omitempty"`
	Website  string `json:"website,omitempty"`
	Role     string `json:"role"`
}

// UserCreateOutput user_create 输出
type UserCreateOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserUpdateInput user_update 输入
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

// UserUpdateOutput user_update 输出
type UserUpdateOutput struct {
	Item UserDetailItem `json:"item"`
}

// UserDeleteInput user_delete 输入
type UserDeleteInput struct {
	ID uint `json:"id"`
}

// UserDeleteOutput user_delete 输出
type UserDeleteOutput struct {
	Deleted bool `json:"deleted"`
	ID      uint `json:"id"`
}

// ============ 服务包装器============

// UserWrapper 用户服务包装器
type UserWrapper struct {
	userService *service.UserService
}

// NewUserWrapper 创建用户服务包装器
func NewUserWrapper(userService *service.UserService) *UserWrapper {
	return &UserWrapper{userService: userService}
}

// ============ 聚合 Tool Handler============

// ManageUser 用户管理聚合入口
func (w *UserWrapper) ManageUser(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input UserManageInput,
) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	switch input.Action {
	case userActionList:
		return w.listUsers(input.Payload)
	case userActionGet:
		return w.getUser(input.Payload)
	case userActionCreate:
		return w.createUser(input.Payload)
	case userActionUpdate:
		return w.updateUser(input.Payload)
	case userActionDelete:
		return w.deleteUser(input.Payload)
	default:
		return nil, UserManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

func (w *UserWrapper) listUsers(payload UserManagePayload) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)
	req := &dto.ListUsersRequest{Page: page, PageSize: pageSize}
	users, total, err := w.userService.List(req)
	if err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("获取用户列表失败: %v", err)}, nil
	}
	list := make([]UserItem, len(users))
	for i, user := range users {
		list[i] = convertToUserItem(user)
	}
	return nil, UserManageOutput{List: list, Total: total, Page: page, PageSize: pageSize}, nil
}

func (w *UserWrapper) getUser(payload UserManagePayload) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	if payload.ID == 0 {
		return nil, UserManageOutput{Error: "用户 ID 不能为空"}, nil
	}
	user, err := w.userService.Get(payload.ID)
	if err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("获取用户失败: %v", err)}, nil
	}
	item := convertToUserDetailItem(user)
	item.IsEnabled = user.IsEnabled
	return nil, UserManageOutput{Item: &item}, nil
}

func (w *UserWrapper) createUser(payload UserManagePayload) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	if payload.Email == "" {
		return nil, UserManageOutput{Error: "邮箱不能为空"}, nil
	}
	if payload.Password == "" {
		return nil, UserManageOutput{Error: "密码不能为空"}, nil
	}
	if payload.Nickname == "" {
		return nil, UserManageOutput{Error: "昵称不能为空"}, nil
	}
	req := &dto.AdminCreateUserRequest{
		Email:    payload.Email,
		Password: payload.Password,
		Nickname: payload.Nickname,
		Avatar:   payload.Avatar,
		Badge:    payload.Badge,
		Website:  payload.Website,
		Role:     parseUserRole(payload.Role),
	}
	if err := w.userService.Create(mcpSuperAdminOperator(), req, ""); err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("创建用户失败: %v", err)}, nil
	}
	createdUser, err := w.userService.GetByEmail(payload.Email)
	if err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("获取新建用户失败: %v", err)}, nil
	}
	item := convertToUserDetailItem(createdUser)
	item.IsEnabled = createdUser.IsEnabled
	return nil, UserManageOutput{Item: &item}, nil
}

func (w *UserWrapper) updateUser(payload UserManagePayload) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	if payload.ID == 0 {
		return nil, UserManageOutput{Error: "用户 ID 不能为空"}, nil
	}
	req := &dto.AdminUpdateUserRequest{
		Email:     payload.Email,
		Nickname:  payload.Nickname,
		Avatar:    payload.Avatar,
		Badge:     payload.Badge,
		Website:   payload.Website,
		Role:      parseUserRoleForUpdate(payload.Role),
		IsEnabled: payload.IsEnabled,
		Password:  payload.Password,
	}
	if err := w.userService.Update(mcpSuperAdminOperator(), payload.ID, req); err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("更新用户失败: %v", err)}, nil
	}
	user, err := w.userService.Get(payload.ID)
	if err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("获取更新后用户失败: %v", err)}, nil
	}
	item := convertToUserDetailItem(user)
	item.IsEnabled = user.IsEnabled
	return nil, UserManageOutput{Item: &item}, nil
}

func (w *UserWrapper) deleteUser(payload UserManagePayload) (*sdkmcp.CallToolResult, UserManageOutput, error) {
	if payload.ID == 0 {
		return nil, UserManageOutput{Error: "用户 ID 不能为空"}, nil
	}
	if err := w.userService.Delete(mcpSuperAdminOperator(), payload.ID); err != nil {
		return nil, UserManageOutput{Error: fmt.Sprintf("删除用户失败: %v", err)}, nil
	}
	success := true
	return nil, UserManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

// ============ 细粒度 Tool Handler============

// ListUsersTool 分页查询用户
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

// GetUserTool 获取用户详情
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

// CreateUserTool 创建用户
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

// UpdateUserTool 更新用户
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

// DeleteUserTool 软删除用户
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

// ============ 辅助函数============

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

// UserManageInputSchema 返回 user_manage 的自定义输入 schema
func UserManageInputSchema() *jsonschema.Schema {
	listPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{
		"page":      {Type: "integer"},
		"page_size": PageSizeSchema(),
	})
	idPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id": {Type: "integer"},
		},
		"id",
	)
	createPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"email":    {Type: "string"},
			"password": {Type: "string"},
			"nickname": {Type: "string"},
			"avatar":   {Type: "string"},
			"badge":    {Type: "string"},
			"website":  {Type: "string"},
			"role": {
				Type: "string",
				Enum: []any{"super_admin", "admin", "user", "guest"},
			},
		},
		"email",
		"password",
		"nickname",
		"role",
	)
	updatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":       {Type: "integer"},
			"email":    {Type: "string"},
			"password": {Type: "string"},
			"nickname": {Type: "string"},
			"avatar":   {Type: "string"},
			"badge":    {Type: "string"},
			"website":  {Type: "string"},
			"role": {
				Type: "string",
				Enum: []any{"super_admin", "admin", "user", "guest"},
			},
			"is_enabled": {Type: "boolean"},
		},
		"id",
	)
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					userActionList,
					userActionGet,
					userActionCreate,
					userActionUpdate,
					userActionDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(userActionList, "获取用户列表", listPayload),
			BuildActionSchema(userActionGet, "获取用户详情", idPayload),
			BuildActionSchema(userActionCreate, "创建用户", createPayload),
			BuildActionSchema(userActionUpdate, "更新用户信息", updatePayload),
			BuildActionSchema(userActionDelete, "删除用户", idPayload),
		},
	}
}

func convertToUserItem(user dto.UserListResponse) UserItem {
	createdAt := userTimePtrFromJSONTime(user.CreatedAt)
	return UserItem{
		ID:           user.ID,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Badge:        user.Badge,
		Website:      user.Website,
		Role:         string(user.Role),
		IsEnabled:    user.IsEnabled,
		HasPassword:  user.HasPassword,
		LinkedOAuths: extractLinkedOAuthsFromList(user),
		LastLogin:    ToTimeStringPtr(user.LastLogin),
		CreatedAt:    createdAt,
		DeletedAt:    ToTimeStringPtr(user.DeletedAt),
	}
}

func convertToUserDetailItem(user *dto.UserResponse) UserDetailItem {
	createdAt := userTimePtrFromJSONTime(user.CreatedAt)
	return UserDetailItem{
		ID:             user.ID,
		Email:          user.Email,
		EmailHash:      user.EmailHash,
		IsVirtualEmail: user.IsVirtualEmail,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Badge:          user.Badge,
		Website:        user.Website,
		Role:           string(user.Role),
		HasPassword:    user.HasPassword,
		LinkedOAuths:   user.LinkedOAuths,
		LastLogin:      ToTimeStringPtr(user.LastLogin),
		CreatedAt:      createdAt,
		IsEnabled:      user.IsEnabled,
	}
}

func extractLinkedOAuthsFromList(user dto.UserListResponse) []string {
	linked := make([]string, 0, 5)
	if user.GithubID != "" {
		linked = append(linked, "github")
	}
	if user.GoogleID != "" {
		linked = append(linked, "google")
	}
	if user.QQID != "" {
		linked = append(linked, "qq")
	}
	if user.MicrosoftID != "" {
		linked = append(linked, "microsoft")
	}
	if user.FeishuOpenID != "" {
		linked = append(linked, "feishu")
	}
	return linked
}

func parseUserRole(role string) model.UserRole {
	switch role {
	case "super_admin":
		return model.RoleSuperAdmin
	case "admin":
		return model.RoleAdmin
	case "guest":
		return model.RoleGuest
	default:
		return model.RoleUser
	}
}

func parseUserRoleForUpdate(role string) model.UserRole {
	if role == "" {
		return ""
	}
	return parseUserRole(role)
}

func mcpSuperAdminOperator() *model.User {
	return &model.User{Role: model.RoleSuperAdmin}
}

func userTimePtrFromJSONTime(t utils.JSONTime) *string {
	return ToTimeStringPtr(&t)
}
