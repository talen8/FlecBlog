package tools

import (
	"context"
	"fmt"

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

// UserDetailItem 用户详情项（用于聚合输出）
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
	// 用于 list
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 get/update/delete
	ID uint `json:"id"`

	// 用于 create/update
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
	// list 结果
	List     []UserItem `json:"list,omitempty"`
	Total    int64      `json:"total,omitempty"`
	Page     int        `json:"page,omitempty"`
	PageSize int        `json:"page_size,omitempty"`

	// get/create/update 结果
	Item *UserDetailItem `json:"item,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
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

// listUsers 获取用户列表
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

	return nil, UserManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// getUser 获取用户详情
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

// createUser 创建用户
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

// updateUser 更新用户
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

// deleteUser 删除用户
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
			"password": {Type: "string", Description: "用户密码，创建时必填"},
			"nickname": {Type: "string"},
			"avatar":   {Type: "string"},
			"badge":    {Type: "string"},
			"website":  {Type: "string"},
			"role": {
				Type:        "string",
				Enum:        []any{"super_admin", "admin", "user", "guest"},
				Description: "用户角色：super_admin(超级管理员)、admin(管理员)、user(普通用户)、guest(访客)。请谨慎分配高权限角色",
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
			"password": {Type: "string", Description: "新密码，留空表示不修改"},
			"nickname": {Type: "string"},
			"avatar":   {Type: "string"},
			"badge":    {Type: "string"},
			"website":  {Type: "string"},
			"role": {
				Type:        "string",
				Enum:        []any{"super_admin", "admin", "user", "guest"},
				Description: "用户角色：super_admin(超级管理员)、admin(管理员)、user(普通用户)、guest(访客)。请谨慎分配高权限角色",
			},
			"is_enabled": {Type: "boolean", Description: "启用状态。设为false将禁止用户登录，请谨慎操作"},
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
			BuildActionSchema(userActionDelete, "删除用户。风险操作，谨慎使用，不可恢复", idPayload),
		},
	}
}

// ============ 转换函数============

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
