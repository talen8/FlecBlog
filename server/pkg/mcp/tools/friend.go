package tools

import (
	"context"
	"fmt"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	friendActionList   = "list"
	friendActionGet    = "get"
	friendActionCreate = "create"
	friendActionUpdate = "update"
	friendActionDelete = "delete"
)

// ============ MCP 类型定义============

// FriendItem 友链列表项
type FriendItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Avatar      string  `json:"avatar"`
	Screenshot  string  `json:"screenshot"`
	Sort        int     `json:"sort"`
	IsInvalid   bool    `json:"is_invalid"`
	IsPending   bool    `json:"is_pending"`
	TypeID      *uint   `json:"type_id,omitempty"`
	TypeName    string  `json:"type_name,omitempty"`
	RSSUrl      string  `json:"rss_url"`
	RSSLatime   *string `json:"rss_latime,omitempty"`
	Accessible  int     `json:"accessible"`
}

// FriendDetailItem 友链详情项
type FriendDetailItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Avatar      string  `json:"avatar"`
	Screenshot  string  `json:"screenshot"`
	Sort        int     `json:"sort"`
	IsInvalid   bool    `json:"is_invalid"`
	IsPending   bool    `json:"is_pending"`
	TypeID      *uint   `json:"type_id,omitempty"`
	TypeName    string  `json:"type_name,omitempty"`
	RSSUrl      string  `json:"rss_url"`
	RSSLatime   *string `json:"rss_latime,omitempty"`
	Accessible  int     `json:"accessible"`
}

// ============ 聚合 Tool 输入/输出类型============

// FriendManageInput friend_manage 聚合 tool 输入
type FriendManageInput struct {
	Action  string              `json:"action"` // list|get|create|update|delete
	Payload FriendManagePayload `json:"payload"`
}

// FriendManagePayload friend_manage 载荷
type FriendManagePayload struct {
	// 用于 list
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 get/update/delete
	ID uint `json:"id"`

	// 用于 create/update
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Screenshot  string `json:"screenshot"`
	Sort        int    `json:"sort"`
	IsInvalid   *bool  `json:"is_invalid"`
	IsPending   *bool  `json:"is_pending"`
	TypeID      *uint  `json:"type_id"`
	RSSUrl      string `json:"rss_url"`
	Accessible  *int   `json:"accessible"`
}

// FriendManageOutput friend_manage 聚合 tool 输出
type FriendManageOutput struct {
	// list 结果
	List     *[]FriendItem `json:"list,omitempty"`
	Total    *int64        `json:"total,omitempty"`
	Page     int          `json:"page,omitempty"`
	PageSize int          `json:"page_size,omitempty"`

	// get/create/update 结果
	Item *FriendDetailItem `json:"item,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
}

// ============ 服务包装器============

// FriendWrapper 友链服务包装器
type FriendWrapper struct {
	friendService *service.FriendService
}

// NewFriendWrapper 创建友链服务包装器
func NewFriendWrapper(friendService *service.FriendService) *FriendWrapper {
	return &FriendWrapper{friendService: friendService}
}

// ============ 聚合 Tool Handler============

// ManageFriend 友链管理聚合入口
func (w *FriendWrapper) ManageFriend(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input FriendManageInput,
) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	switch input.Action {
	case friendActionList:
		return w.listFriends(input.Payload)
	case friendActionGet:
		return w.getFriend(input.Payload)
	case friendActionCreate:
		return w.createFriend(input.Payload)
	case friendActionUpdate:
		return w.updateFriend(input.Payload)
	case friendActionDelete:
		return w.deleteFriend(input.Payload)
	default:
		return nil, FriendManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

// listFriends 获取友链列表
func (w *FriendWrapper) listFriends(payload FriendManagePayload) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	req := &dto.ListFriendRequest{Page: page, PageSize: pageSize}
	friends, total, err := w.friendService.List(context.Background(), req)
	if err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("获取友链列表失败: %v", err)}, nil
	}

	list := make([]FriendItem, len(friends))
	for i, friend := range friends {
		list[i] = convertToFriendItem(friend)
	}

	listField, totalField := ExplicitListFields(list, total)

	return nil, FriendManageOutput{
		List:     listField,
		Total:    totalField,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// getFriend 获取友链详情
func (w *FriendWrapper) getFriend(payload FriendManagePayload) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	if payload.ID == 0 {
		return nil, FriendManageOutput{Error: "友链 ID 不能为空"}, nil
	}

	friend, err := w.friendService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("获取友链失败: %v", err)}, nil
	}

	item := convertToFriendDetailItem(*friend)
	return nil, FriendManageOutput{Item: &item}, nil
}

// createFriend 创建友链
func (w *FriendWrapper) createFriend(payload FriendManagePayload) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	if payload.Name == "" {
		return nil, FriendManageOutput{Error: "友链名称不能为空"}, nil
	}
	if payload.URL == "" {
		return nil, FriendManageOutput{Error: "友链 URL 不能为空"}, nil
	}

	req := &dto.CreateFriendRequest{
		Name:        payload.Name,
		URL:         payload.URL,
		Description: payload.Description,
		Avatar:      payload.Avatar,
		Screenshot:  payload.Screenshot,
		Sort:        payload.Sort,
		TypeID:      payload.TypeID,
		RSSUrl:      payload.RSSUrl,
	}

	friend, err := w.friendService.Create(context.Background(), req)
	if err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("创建友链失败: %v", err)}, nil
	}

	item := convertToFriendDetailItem(*friend)
	return nil, FriendManageOutput{Item: &item}, nil
}

// updateFriend 更新友链
func (w *FriendWrapper) updateFriend(payload FriendManagePayload) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	if payload.ID == 0 {
		return nil, FriendManageOutput{Error: "友链 ID 不能为空"}, nil
	}

	req := &dto.UpdateFriendRequest{
		Name:        payload.Name,
		URL:         payload.URL,
		Description: payload.Description,
		Avatar:      payload.Avatar,
		Screenshot:  payload.Screenshot,
		Sort:        payload.Sort,
		IsInvalid:   payload.IsInvalid,
		IsPending:   payload.IsPending,
		TypeID:      payload.TypeID,
		RSSUrl:      payload.RSSUrl,
		Accessible:  payload.Accessible,
	}

	if err := w.friendService.Update(context.Background(), payload.ID, req); err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("更新友链失败: %v", err)}, nil
	}

	friend, err := w.friendService.Get(context.Background(), payload.ID)
	if err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("获取更新后友链失败: %v", err)}, nil
	}

	item := convertToFriendDetailItem(*friend)
	return nil, FriendManageOutput{Item: &item}, nil
}

// deleteFriend 删除友链
func (w *FriendWrapper) deleteFriend(payload FriendManagePayload) (*sdkmcp.CallToolResult, FriendManageOutput, error) {
	if payload.ID == 0 {
		return nil, FriendManageOutput{Error: "友链 ID 不能为空"}, nil
	}

	if err := w.friendService.Delete(context.Background(), payload.ID); err != nil {
		return nil, FriendManageOutput{Error: fmt.Sprintf("删除友链失败: %v", err)}, nil
	}

	success := true
	return nil, FriendManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

// FriendManageInputSchema 返回 friend_manage 的自定义输入 schema
func FriendManageInputSchema() *jsonschema.Schema {
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
			"name":        {Type: "string"},
			"url":         {Type: "string"},
			"description": {Type: "string"},
			"avatar":      {Type: "string"},
			"screenshot":  {Type: "string"},
			"sort":        {Type: "integer"},
			"type_id":     {Type: "integer"},
			"rss_url":     {Type: "string"},
		},
		"name",
		"url",
	)
	updatePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id":          {Type: "integer"},
			"name":        {Type: "string"},
			"url":         {Type: "string"},
			"description": {Type: "string"},
			"avatar":      {Type: "string"},
			"screenshot":  {Type: "string"},
			"sort":        {Type: "integer"},
			"is_invalid":  {Type: "boolean"},
			"is_pending":  {Type: "boolean"},
			"type_id":     {Type: "integer"},
			"rss_url":     {Type: "string"},
			"accessible":  {Type: "integer"},
		},
		"id",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					friendActionList,
					friendActionGet,
					friendActionCreate,
					friendActionUpdate,
					friendActionDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(friendActionList, "获取友链列表", listPayload),
			BuildActionSchema(friendActionGet, "获取友链详情", idPayload),
			BuildActionSchema(friendActionCreate, "创建友链", createPayload),
			BuildActionSchema(friendActionUpdate, "更新友链信息", updatePayload),
			BuildActionSchema(friendActionDelete, "删除友链。风险操作，谨慎使用，不可恢复", idPayload),
		},
	}
}

// ============ 转换函数============

func friendStringPtr(t interface{ String() string }) *string {
	return ToTimeStringPtr(t)
}

func convertToFriendItem(item dto.FriendListResponse) FriendItem {
	return FriendItem{
		ID:          item.ID,
		Name:        item.Name,
		URL:         item.URL,
		Description: item.Description,
		Avatar:      item.Avatar,
		Screenshot:  item.Screenshot,
		Sort:        item.Sort,
		IsInvalid:   item.IsInvalid,
		IsPending:   item.IsPending,
		TypeID:      item.TypeID,
		TypeName:    item.TypeName,
		RSSUrl:      item.RSSUrl,
		RSSLatime:   friendStringPtr(item.RSSLatime),
		Accessible:  item.Accessible,
	}
}

func convertToFriendDetailItem(item model.Friend) FriendDetailItem {
	return FriendDetailItem{
		ID:          item.ID,
		Name:        item.Name,
		URL:         item.URL,
		Description: item.Description,
		Avatar:      item.Avatar,
		Screenshot:  item.Screenshot,
		Sort:        item.Sort,
		IsInvalid:   item.IsInvalid,
		IsPending:   item.IsPending,
		TypeID:      item.TypeID,
		TypeName:    convertFriendTypeName(item.Type),
		RSSUrl:      item.RSSUrl,
		RSSLatime:   friendTimePtr(item.RSSLatime),
		Accessible:  item.Accessible,
	}
}

func convertFriendTypeName(friendType *model.FriendType) string {
	if friendType == nil {
		return ""
	}
	return friendType.Name
}

func friendTimePtr(t interface{ String() string }) *string {
	return ToTimeStringPtr(t)
}
