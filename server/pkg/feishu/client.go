package feishu

import (
	"context"
	"fmt"
	"sync"
	"time"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/pkg/logger"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

// CardActionHandler 卡片动作处理器函数类型
type CardActionHandler func(ctx context.Context, action string, value map[string]interface{}) error

// CardActionHandlers 卡片动作处理器映射
var (
	cardActionHandlers = make(map[string]CardActionHandler)
	handlersMu         sync.RWMutex
)

// FriendApprover 友链审核接口
type FriendApprover interface {
	ApproveFriend(ctx context.Context, id uint) error
}

// CommentReplier 评论回复接口
type CommentReplier interface {
	ReplyCommentFromFeishu(ctx context.Context, commentID uint, content, openID string) error
}

// UserBinder 用户绑定接口
type UserBinder interface {
	BindFeishuByEmail(email, openID string) error
}

// StatsProvider 站点统计接口
type StatsProvider interface {
	GetDashboardStats() (*dto.DashboardStats, error)
}

// SystemProvider 系统状态接口
type SystemProvider interface {
	GetSystemStatus(ctx context.Context) (*SystemStatus, error)
}

type SystemStatus struct {
	CPUUsage      float64
	MemoryTotal   uint64
	MemoryUsed    uint64
	DiskTotal     uint64
	DiskUsed      uint64
	DBStatus      string
	StorageStatus string
	EmailStatus   string
	FeishuStatus  string
}

// RssArticleReader RSS 文章已读接口
type RssArticleReader interface {
	MarkAllReadFromFeishu(ctx context.Context) error
}

var userBinderInstance UserBinder
var statsProviderInstance StatsProvider
var systemProviderInstance SystemProvider

var (
	globalClient *Client
	globalMu     sync.RWMutex
)

// Client 飞书客户端
type Client struct {
	appID     string
	appSecret string
	chatID    string
	enable    bool
	client    *lark.Client
	wsClient  *larkws.Client
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// NewClient 创建飞书客户端
func NewClient(appID, appSecret, chatID string) *Client {
	enable := appID != "" && appSecret != "" && chatID != ""

	var client *lark.Client
	if enable {
		client = lark.NewClient(appID, appSecret)
	}

	return &Client{
		appID:     appID,
		appSecret: appSecret,
		chatID:    chatID,
		enable:    enable,
		client:    client,
	}
}

// IsEnabled 检查客户端是否启用
func (c *Client) IsEnabled() bool {
	return c.enable
}

// HealthCheck 检查飞书服务可用性
func (c *Client) HealthCheck() error {
	if !c.enable || c.client == nil {
		return fmt.Errorf("未配置飞书客户端")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetTenantAccessTokenBySelfBuiltApp(ctx, &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     c.appID,
		AppSecret: c.appSecret,
	})
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("飞书认证失败")
	}
	return nil
}

// RegisterCardActionHandler 注册卡片动作处理器
func RegisterCardActionHandler(action string, handler CardActionHandler) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	cardActionHandlers[action] = handler
}

// handleCardAction 处理卡片动作回调
func handleCardAction(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	if event == nil || event.Event == nil || event.Event.Action == nil {
		return nil, nil
	}

	actionValue := event.Event.Action.Value
	if actionValue == nil {
		return nil, nil
	}

	action, _ := actionValue["action"].(string)
	if action == "" {
		return nil, nil
	}

	handlersMu.RLock()
	handler, exists := cardActionHandlers[action]
	handlersMu.RUnlock()

	if !exists {
		return nil, nil
	}

	var openID string
	if event.Event.Operator != nil {
		openID = event.Event.Operator.OpenID
	}

	if event.Event.Action.Name != "" && event.Event.Action.InputValue != "" {
		actionValue[event.Event.Action.Name] = event.Event.Action.InputValue
	}
	actionValue["open_id"] = openID

	if err := handler(ctx, action, actionValue); err != nil {
		logger.Error("[Feishu] 卡片动作失败: %v", err)
		return &callback.CardActionTriggerResponse{
			Toast: &callback.Toast{Type: "error", Content: "操作失败: " + err.Error()},
		}, nil
	}

	return &callback.CardActionTriggerResponse{
		Toast: &callback.Toast{Type: "success", Content: "操作成功"},
	}, nil
}

// SendMessage 发送卡片消息
func (c *Client) SendMessage(ctx context.Context, cardJSON string) error {
	if !c.enable || c.client == nil {
		return fmt.Errorf("飞书客户端未启用")
	}

	resp, err := c.client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(c.chatID).
			MsgType("interactive").
			Content(cardJSON).
			Build()).
		Build())
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("%s", resp.Msg)
	}
	return nil
}

// start 启动长连接
func (c *Client) start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}
	if !c.enable {
		return
	}

	// #nosec G118 - cancel 函数保存在 c.cancel 中，在 stop/Reload 时调用
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.wsClient = larkws.NewClient(c.appID, c.appSecret, larkws.WithEventHandler(createEventHandler()))

	go func(ctx context.Context) {
		if err := c.wsClient.Start(ctx); err != nil {
			logger.Error("[Feishu] WebSocket 连接失败: %v", err)
		}
	}(ctx)
}

// createEventHandler 创建事件处理器
func createEventHandler() *dispatcher.EventDispatcher {
	return dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			if err := handleCommand(ctx, event); err != nil {
				logger.Error("[Feishu] 处理命令失败: %v", err)
			}
			return nil
		}).
		OnP2CardActionTrigger(handleCardAction)
}

// Initialize 初始化飞书客户端
func Initialize(conf *config.Config) *Client {
	c := NewClient(conf.Notification.FeishuAppID, conf.Notification.FeishuSecret, conf.Notification.FeishuChatID)
	c.start()

	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()

	return c
}

// Reload 重新加载配置
func Reload(appID, appSecret, chatID string) {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalClient == nil {
		return
	}

	globalClient.mu.Lock()
	defer globalClient.mu.Unlock()

	if globalClient.cancel != nil {
		globalClient.cancel()
		globalClient.cancel = nil
	}

	globalClient.appID = appID
	globalClient.appSecret = appSecret
	globalClient.chatID = chatID
	globalClient.enable = appID != "" && appSecret != "" && chatID != ""

	if globalClient.enable {
		globalClient.client = lark.NewClient(appID, appSecret)
		// #nosec G118 - cancel 函数保存在 globalClient.cancel 中，在 Reload 时调用
		ctx, cancel := context.WithCancel(context.Background())
		globalClient.cancel = cancel
		globalClient.wsClient = larkws.NewClient(appID, appSecret, larkws.WithEventHandler(createEventHandler()))

		go func(ctx context.Context) {
			if err := globalClient.wsClient.Start(ctx); err != nil {
				logger.Error("[Feishu] WebSocket 连接失败: %v", err)
			}
		}(ctx)
	}
}

// InitCardHandlers 初始化卡片动作处理器
func InitCardHandlers(friendService FriendApprover, commentService CommentReplier, userService UserBinder, statsService StatsProvider, systemService SystemProvider, rssService RssArticleReader) {
	InitCommands(userService, statsService, systemService)

	RegisterCardActionHandler("approve_friend", func(ctx context.Context, action string, value map[string]interface{}) error {
		friendID := uint(value["friend_id"].(float64))
		return friendService.ApproveFriend(ctx, friendID)
	})

	RegisterCardActionHandler("reply_comment", func(ctx context.Context, action string, value map[string]interface{}) error {
		commentID := uint(value["comment_id"].(float64))
		replyContent, ok := value["reply_content"].(string)
		if !ok || replyContent == "" {
			return fmt.Errorf("回复内容不能为空")
		}
		openID, ok := value["open_id"].(string)
		if !ok || openID == "" {
			return fmt.Errorf("无法获取用户身份")
		}
		return commentService.ReplyCommentFromFeishu(ctx, commentID, replyContent, openID)
	})

	RegisterCardActionHandler("rss_mark_all_read", func(ctx context.Context, action string, value map[string]interface{}) error {
		return rssService.MarkAllReadFromFeishu(ctx)
	})
}
