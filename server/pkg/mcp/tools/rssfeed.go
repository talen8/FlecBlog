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
	rssFeedActionList        = "list"
	rssFeedActionMarkRead    = "mark_read"
	rssFeedActionMarkAllRead = "mark_all_read"
)

// RssFeedItem RSS订阅文章项
type RssFeedItem struct {
	ID          uint    `json:"id"`
	FriendID    uint    `json:"friend_id"`
	FriendName  string  `json:"friend_name"`
	FriendURL   string  `json:"friend_url"`
	Title       string  `json:"title"`
	Link        string  `json:"link"`
	PublishedAt *string `json:"published_at,omitempty"`
	IsRead      bool    `json:"is_read"`
	CreatedAt   *string `json:"created_at,omitempty"`
}

// RssFeedManageInput rssfeed_manage 聚合 tool 输入
type RssFeedManageInput struct {
	Action  string               `json:"action"`
	Payload RssFeedManagePayload `json:"payload"`
}

// RssFeedManagePayload rssfeed_manage 载荷
type RssFeedManagePayload struct {
	// 用于 list
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// 用于 mark_read
	ID uint `json:"id"`
}

// RssFeedManageOutput rssfeed_manage 聚合 tool 输出
type RssFeedManageOutput struct {
	// list 结果
	List        []RssFeedItem `json:"list,omitempty"`
	Total       int64         `json:"total,omitempty"`
	Page        int           `json:"page,omitempty"`
	PageSize    int           `json:"page_size,omitempty"`
	UnreadCount int64         `json:"unread_count,omitempty"`

	// mark_read / mark_all_read 结果
	Success  *bool  `json:"success,omitempty"`
	ID       *uint  `json:"id,omitempty"`
	Affected *int64 `json:"affected,omitempty"`

	// 错误信息
	Error string `json:"error,omitempty"`
}

// RssFeedWrapper RSS订阅服务包装器
type RssFeedWrapper struct {
	rssFeedService *service.RssFeedService
}

// NewRssFeedWrapper 创建 RSS 订阅服务包装器
func NewRssFeedWrapper(rssFeedService *service.RssFeedService) *RssFeedWrapper {
	return &RssFeedWrapper{rssFeedService: rssFeedService}
}

// ManageRssFeed RSS订阅管理聚合入口
func (w *RssFeedWrapper) ManageRssFeed(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input RssFeedManageInput,
) (*sdkmcp.CallToolResult, RssFeedManageOutput, error) {
	switch input.Action {
	case rssFeedActionList:
		return w.listRssFeeds(input.Payload)
	case rssFeedActionMarkRead:
		return w.markRead(input.Payload)
	case rssFeedActionMarkAllRead:
		return w.markAllRead()
	default:
		return nil, RssFeedManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

func (w *RssFeedWrapper) listRssFeeds(payload RssFeedManagePayload) (*sdkmcp.CallToolResult, RssFeedManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	result, err := w.rssFeedService.List(context.Background(), &dto.ListRssArticleRequest{Page: page, PageSize: pageSize})
	if err != nil {
		return nil, RssFeedManageOutput{Error: fmt.Sprintf("获取 RSS 订阅列表失败: %v", err)}, nil
	}

	list := make([]RssFeedItem, len(result.List))
	for i, item := range result.List {
		list[i] = convertToRssFeedItem(item)
	}

	return nil, RssFeedManageOutput{
		List:        list,
		Total:       result.Total,
		Page:        result.Page,
		PageSize:    result.PageSize,
		UnreadCount: result.UnreadCount,
	}, nil
}

func (w *RssFeedWrapper) markRead(payload RssFeedManagePayload) (*sdkmcp.CallToolResult, RssFeedManageOutput, error) {
	if payload.ID == 0 {
		return nil, RssFeedManageOutput{Error: "RSS 文章 ID 不能为空"}, nil
	}

	if err := w.rssFeedService.MarkRead(context.Background(), payload.ID); err != nil {
		return nil, RssFeedManageOutput{Error: fmt.Sprintf("标记 RSS 文章已读失败: %v", err)}, nil
	}

	success := true
	return nil, RssFeedManageOutput{Success: &success, ID: &payload.ID}, nil
}

func (w *RssFeedWrapper) markAllRead() (*sdkmcp.CallToolResult, RssFeedManageOutput, error) {
	affected, err := w.rssFeedService.MarkAllRead(context.Background())
	if err != nil {
		return nil, RssFeedManageOutput{Error: fmt.Sprintf("全部标记 RSS 文章已读失败: %v", err)}, nil
	}

	success := true
	return nil, RssFeedManageOutput{Success: &success, Affected: &affected}, nil
}

// RssFeedManageInputSchema 返回 rssfeed_manage 的自定义输入 schema
func RssFeedManageInputSchema() *jsonschema.Schema {
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
	emptyPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{})

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					rssFeedActionList,
					rssFeedActionMarkRead,
					rssFeedActionMarkAllRead,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(rssFeedActionList, "获取RSS订阅文章列表", listPayload),
			BuildActionSchema(rssFeedActionMarkRead, "标记单篇RSS文章为已读", idPayload),
			BuildActionSchema(rssFeedActionMarkAllRead, "标记所有RSS文章为已读", emptyPayload),
		},
	}
}

func convertToRssFeedItem(item dto.RssArticleResponse) RssFeedItem {
	return RssFeedItem{
		ID:          item.ID,
		FriendID:    item.FriendID,
		FriendName:  item.FriendName,
		FriendURL:   item.FriendURL,
		Title:       item.Title,
		Link:        item.Link,
		PublishedAt: ToTimeStringPtr(item.PublishedAt),
		IsRead:      item.IsRead,
		CreatedAt:   ToTimeStringPtr(item.CreatedAt),
	}
}
