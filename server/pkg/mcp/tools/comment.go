package tools

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	commentActionList         = "list"
	commentActionGet          = "get"
	commentActionToggleStatus = "toggle_status"
	commentActionDelete       = "delete"
)

// ============ MCP 类型定义============

// CommentItem 评论列表项
type CommentItem struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	Status    int           `json:"status"`
	ParentID  *uint         `json:"parent_id,omitempty"`
	CreatedAt *string       `json:"created_at"`
	DeletedAt *string       `json:"deleted_at,omitempty"`
	Target    CommentTarget `json:"target"`
	User      CommentActor  `json:"user"`
}

// CommentDetailItem 评论详情项
type CommentDetailItem struct {
	ID        uint           `json:"id"`
	Content   string         `json:"content"`
	Status    *int           `json:"status,omitempty"`
	ParentID  *uint          `json:"parent_id,omitempty"`
	CreatedAt *string        `json:"created_at"`
	DeletedAt *string        `json:"deleted_at,omitempty"`
	Target    *CommentTarget `json:"target,omitempty"`
	User      *CommentActor  `json:"user,omitempty"`
}

// CommentTarget 评论目标
type CommentTarget struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Title string `json:"title"`
}

// CommentActor 评论用户
type CommentActor struct {
	ID       uint   `json:"id"`
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Badge    string `json:"badge"`
}

// ============ 聚合 Tool 输入/输出类型============

// CommentManageInput comment_manage 聚合 tool 输入
type CommentManageInput struct {
	Action  string               `json:"action"`
	Payload CommentManagePayload `json:"payload"`
}

// CommentManagePayload comment_manage 载荷
type CommentManagePayload struct {
	// 用于 list
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Status   *int `json:"status"`

	// 用于 get/toggle_status/delete
	ID uint `json:"id"`
}

// CommentManageOutput comment_manage 聚合 tool 输出
type CommentManageOutput struct {
	// list 结果
	List     []CommentItem `json:"list,omitempty"`
	Total    int64         `json:"total,omitempty"`
	Page     int           `json:"page,omitempty"`
	PageSize int           `json:"page_size,omitempty"`

	// get 结果
	Item *CommentDetailItem `json:"item,omitempty"`

	// toggle_status/delete 结果
	Success *bool `json:"success,omitempty"`
	ID      *uint `json:"id,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
}

// ============ 服务包装器============

// CommentWrapper 评论服务包装器
type CommentWrapper struct {
	commentService *service.CommentService
}

// NewCommentWrapper 创建评论服务包装器
func NewCommentWrapper(commentService *service.CommentService) *CommentWrapper {
	return &CommentWrapper{commentService: commentService}
}

// ============ 聚合 Tool Handler============

// ManageComment 评论管理聚合入口
func (w *CommentWrapper) ManageComment(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input CommentManageInput,
) (*sdkmcp.CallToolResult, CommentManageOutput, error) {
	switch input.Action {
	case commentActionList:
		return w.listComments(input.Payload)
	case commentActionGet:
		return w.getComment(input.Payload)
	case commentActionToggleStatus:
		return w.toggleStatus(input.Payload)
	case commentActionDelete:
		return w.deleteComment(input.Payload)
	default:
		return nil, CommentManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

// listComments 获取评论列表
func (w *CommentWrapper) listComments(payload CommentManagePayload) (*sdkmcp.CallToolResult, CommentManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	req := &dto.CommentQueryRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   payload.Status,
	}
	comments, total, err := w.commentService.List(context.Background(), req)
	if err != nil {
		return nil, CommentManageOutput{Error: fmt.Sprintf("获取评论列表失败: %v", err)}, nil
	}

	list := make([]CommentItem, len(comments))
	for i, comment := range comments {
		list[i] = convertToCommentItem(comment)
	}

	return nil, CommentManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// getComment 获取评论详情
func (w *CommentWrapper) getComment(payload CommentManagePayload) (*sdkmcp.CallToolResult, CommentManageOutput, error) {
	if payload.ID == 0 {
		return nil, CommentManageOutput{Error: "评论 ID 不能为空"}, nil
	}

	comment, err := w.commentService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, CommentManageOutput{Error: fmt.Sprintf("获取评论失败: %v", err)}, nil
	}

	item := convertToCommentDetailItem(*comment)
	return nil, CommentManageOutput{Item: &item}, nil
}

// toggleStatus 切换评论状态
func (w *CommentWrapper) toggleStatus(payload CommentManagePayload) (*sdkmcp.CallToolResult, CommentManageOutput, error) {
	if payload.ID == 0 {
		return nil, CommentManageOutput{Error: "评论 ID 不能为空"}, nil
	}

	err := w.commentService.ToggleStatus(context.Background(), payload.ID)
	if err != nil {
		return nil, CommentManageOutput{Error: fmt.Sprintf("切换评论状态失败: %v", err)}, nil
	}

	success := true
	return nil, CommentManageOutput{Success: &success, ID: &payload.ID}, nil
}

// deleteComment 删除评论
func (w *CommentWrapper) deleteComment(payload CommentManagePayload) (*sdkmcp.CallToolResult, CommentManageOutput, error) {
	if payload.ID == 0 {
		return nil, CommentManageOutput{Error: "评论 ID 不能为空"}, nil
	}

	err := w.commentService.Delete(context.Background(), payload.ID)
	if err != nil {
		return nil, CommentManageOutput{Error: fmt.Sprintf("删除评论失败: %v", err)}, nil
	}

	success := true
	return nil, CommentManageOutput{Success: &success, ID: &payload.ID}, nil
}

// CommentManageInputSchema 返回 comment_manage 的自定义输入 schema
func CommentManageInputSchema() *jsonschema.Schema {
	listPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{
		"page":      {Type: "integer"},
		"page_size": PageSizeSchema(),
		"status":    {Type: "integer"},
	})
	idPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id": {Type: "integer"},
		},
		"id",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					commentActionList,
					commentActionGet,
					commentActionToggleStatus,
					commentActionDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(commentActionList, "获取评论列表", listPayload),
			BuildActionSchema(commentActionGet, "获取评论详情", idPayload),
			BuildActionSchema(commentActionToggleStatus, "切换评论显示/隐藏状态", idPayload),
			BuildActionSchema(commentActionDelete, "删除评论（硬删除，不可恢复）。风险操作，谨慎使用", idPayload),
		},
	}
}

// ============ 转换函数============

func convertCommentTarget(target struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Title string `json:"title"`
}) CommentTarget {
	return CommentTarget{Type: target.Type, Key: target.Key, Title: target.Title}
}

func convertCommentListActor(user struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Badge    string `json:"badge"`
}) CommentActor {
	return CommentActor{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Badge:    user.Badge,
	}
}

func convertToCommentItem(item dto.CommentListResponse) CommentItem {
	return CommentItem{
		ID:        item.ID,
		Content:   item.Content,
		Status:    item.Status,
		ParentID:  item.ParentID,
		CreatedAt: ToTimeStringPtr(&item.CreatedAt),
		DeletedAt: ToTimeStringPtr(item.DeletedAt),
		Target:    convertCommentTarget(item.Target),
		User:      convertCommentListActor(item.User),
	}
}

func convertToCommentDetailItem(item dto.CommentListResponse) CommentDetailItem {
	status := item.Status
	target := convertCommentTarget(item.Target)
	user := convertCommentListActor(item.User)

	return CommentDetailItem{
		ID:        item.ID,
		Content:   item.Content,
		Status:    &status,
		ParentID:  item.ParentID,
		CreatedAt: ToTimeStringPtr(&item.CreatedAt),
		DeletedAt: ToTimeStringPtr(item.DeletedAt),
		Target:    &target,
		User:      &user,
	}
}
