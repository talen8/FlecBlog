package tools

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/utils"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	momentActionList   = "list"
	momentActionGet    = "get"
	momentActionCreate = "create"
	momentActionUpdate = "update"
	momentActionDelete = "delete"
)

// MomentItem 动态项
type MomentItem struct {
	ID          uint              `json:"id"`
	Content     dto.MomentContent `json:"content"`
	IsPublish   bool              `json:"is_publish"`
	PublishTime *string           `json:"publish_time"`
}

// MomentManageInput moment_manage 聚合 tool 输入
type MomentManageInput struct {
	Action  string              `json:"action"`
	Payload MomentManagePayload `json:"payload"`
}

// MomentManagePayload moment_manage 载荷
type MomentManagePayload struct {
	// 用于 list
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 get/update/delete
	ID uint `json:"id"`

	// 用于 create/update
	Content     *dto.MomentContent `json:"content"`
	IsPublish   *bool              `json:"is_publish"`
	PublishTime *string            `json:"publish_time"`
}

// MomentManageOutput moment_manage 聚合 tool 输出
type MomentManageOutput struct {
	// list 结果
	List     []MomentItem `json:"list,omitempty"`
	Total    int64        `json:"total,omitempty"`
	Page     int          `json:"page,omitempty"`
	PageSize int          `json:"page_size,omitempty"`

	// get/create/update 结果
	Item *MomentItem `json:"item,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
}

// MomentWrapper 动态服务包装器
type MomentWrapper struct {
	momentService *service.MomentService
}

// NewMomentWrapper 创建动态服务包装器
func NewMomentWrapper(momentService *service.MomentService) *MomentWrapper {
	return &MomentWrapper{momentService: momentService}
}

// ManageMoment 动态管理聚合入口
func (w *MomentWrapper) ManageMoment(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input MomentManageInput,
) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	switch input.Action {
	case momentActionList:
		return w.listMoments(input.Payload)
	case momentActionGet:
		return w.getMoment(input.Payload)
	case momentActionCreate:
		return w.createMoment(input.Payload)
	case momentActionUpdate:
		return w.updateMoment(input.Payload)
	case momentActionDelete:
		return w.deleteMoment(input.Payload)
	default:
		return nil, MomentManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

func (w *MomentWrapper) listMoments(payload MomentManagePayload) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	req := dto.ListMomentsRequest{
		Page:     page,
		PageSize: pageSize,
	}
	moments, total, err := w.momentService.List(context.Background(), req)
	if err != nil {
		return nil, MomentManageOutput{Error: fmt.Sprintf("获取动态列表失败: %v", err)}, nil
	}

	list := make([]MomentItem, len(moments))
	for i, moment := range moments {
		list[i] = convertToMomentItem(moment)
	}

	return nil, MomentManageOutput{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (w *MomentWrapper) getMoment(payload MomentManagePayload) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	if payload.ID == 0 {
		return nil, MomentManageOutput{Error: "动态 ID 不能为空"}, nil
	}

	moment, err := w.momentService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, MomentManageOutput{Error: fmt.Sprintf("获取动态失败: %v", err)}, nil
	}

	item := convertToMomentItem(*moment)
	return nil, MomentManageOutput{Item: &item}, nil
}

func (w *MomentWrapper) createMoment(payload MomentManagePayload) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	if payload.Content == nil {
		return nil, MomentManageOutput{Error: "动态内容不能为空"}, nil
	}
	if payload.IsPublish == nil {
		return nil, MomentManageOutput{Error: "是否发布不能为空"}, nil
	}

	req := &dto.CreateMomentRequest{
		Content:   *payload.Content,
		IsPublish: *payload.IsPublish,
	}
	if payload.PublishTime != nil && *payload.PublishTime != "" {
		publishTime, err := parseJSONTime(*payload.PublishTime)
		if err != nil {
			return nil, MomentManageOutput{Error: fmt.Sprintf("发布时间格式错误: %v", err)}, nil
		}
		req.PublishTime = publishTime
	}

	moment, err := w.momentService.Create(context.Background(), req)
	if err != nil {
		return nil, MomentManageOutput{Error: fmt.Sprintf("创建动态失败: %v", err)}, nil
	}

	item := convertToMomentItem(*moment)
	return nil, MomentManageOutput{Item: &item}, nil
}

func (w *MomentWrapper) updateMoment(payload MomentManagePayload) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	if payload.ID == 0 {
		return nil, MomentManageOutput{Error: "动态 ID 不能为空"}, nil
	}
	if payload.Content == nil {
		return nil, MomentManageOutput{Error: "动态内容不能为空"}, nil
	}
	if payload.IsPublish == nil {
		return nil, MomentManageOutput{Error: "是否发布不能为空"}, nil
	}

	req := &dto.UpdateMomentRequest{
		Content:   *payload.Content,
		IsPublish: *payload.IsPublish,
	}
	if payload.PublishTime != nil && *payload.PublishTime != "" {
		publishTime, err := parseJSONTime(*payload.PublishTime)
		if err != nil {
			return nil, MomentManageOutput{Error: fmt.Sprintf("发布时间格式错误: %v", err)}, nil
		}
		req.PublishTime = publishTime
	}

	moment, err := w.momentService.Update(context.Background(), payload.ID, req)
	if err != nil {
		return nil, MomentManageOutput{Error: fmt.Sprintf("更新动态失败: %v", err)}, nil
	}

	item := convertToMomentItem(*moment)
	return nil, MomentManageOutput{Item: &item}, nil
}

func (w *MomentWrapper) deleteMoment(payload MomentManagePayload) (*sdkmcp.CallToolResult, MomentManageOutput, error) {
	if payload.ID == 0 {
		return nil, MomentManageOutput{Error: "动态 ID 不能为空"}, nil
	}

	if err := w.momentService.Delete(context.Background(), payload.ID); err != nil {
		return nil, MomentManageOutput{Error: fmt.Sprintf("删除动态失败: %v", err)}, nil
	}

	success := true
	return nil, MomentManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

// MomentManageInputSchema 返回 moment_manage 的自定义输入 schema
func MomentManageInputSchema() *jsonschema.Schema {
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
	contentSchema := momentContentSchema()
	createPayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"content":      contentSchema,
			"is_publish":   {Type: "boolean", Description: "发布状态。没有明确需求请保持为false草稿状态，谨慎操作"},
			"publish_time": {Type: "string"},
		},
		"content",
		"is_publish",
	)
	updatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":           {Type: "integer"},
			"content":      contentSchema,
			"is_publish":   {Type: "boolean", Description: "发布状态。没有明确需求请保持为false草稿状态，谨慎操作"},
			"publish_time": {Type: "string"},
		},
		"id",
		"content",
		"is_publish",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					momentActionList,
					momentActionGet,
					momentActionCreate,
					momentActionUpdate,
					momentActionDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(momentActionList, "获取动态列表", listPayload),
			BuildActionSchema(momentActionGet, "获取动态详情", idPayload),
			BuildActionSchema(momentActionCreate, "创建动态", createPayload),
			BuildActionSchema(momentActionUpdate, "更新动态内容", updatePayload),
			BuildActionSchema(momentActionDelete, "删除动态。风险操作，谨慎使用，不可恢复", idPayload),
		},
	}
}

func momentContentSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"text":     {Type: "string"},
			"images":   {Type: "array", Items: &jsonschema.Schema{Type: "string"}},
			"location": {Type: "string"},
			"tags":     {Type: "string"},
			"link": {
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"url":     {Type: "string"},
					"title":   {Type: "string"},
					"favicon": {Type: "string"},
				},
				AdditionalProperties: FalseSchema(),
			},
			"music": {
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"server": {Type: "string"},
					"type":   {Type: "string"},
					"id":     {Type: "string"},
				},
				AdditionalProperties: FalseSchema(),
			},
			"video": {
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"url":      {Type: "string"},
					"platform": {Type: "string"},
					"video_id": {Type: "string"},
				},
				AdditionalProperties: FalseSchema(),
			},
			"book":  {Type: "object"},
			"movie": {Type: "object"},
		},
		AdditionalProperties: FalseSchema(),
	}
}

func convertToMomentItem(item dto.MomentListResponse) MomentItem {
	return MomentItem{
		ID:          item.ID,
		Content:     item.Content,
		IsPublish:   item.IsPublish,
		PublishTime: ToTimeStringPtr(item.PublishTime),
	}
}

func parseJSONTime(value string) (*utils.JSONTime, error) {
	var t utils.JSONTime
	data := []byte{'"'}
	data = append(data, value...)
	data = append(data, '"')
	if err := t.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &t, nil
}
