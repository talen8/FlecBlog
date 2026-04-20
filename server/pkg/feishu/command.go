package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"flec_blog/internal/dto"
	"flec_blog/pkg/logger"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// CommandHandler 命令处理器函数类型
type CommandHandler func(ctx context.Context, args []string, openID string) error

var (
	commandHandlers = make(map[string]CommandHandler)
	commandsMu      sync.RWMutex
)

// RegisterCommand 注册命令处理器
func RegisterCommand(command string, handler CommandHandler) {
	commandsMu.Lock()
	defer commandsMu.Unlock()
	commandHandlers[command] = handler
}

// handleCommand 处理命令消息
func handleCommand(_ context.Context, event *larkim.P2MessageReceiveV1) error {
	if event.Event == nil || event.Event.Message == nil || event.Event.Message.Content == nil {
		return nil
	}

	globalMu.RLock()
	chatID := globalClient.chatID
	globalMu.RUnlock()

	if event.Event.Message.ChatId == nil || *event.Event.Message.ChatId != chatID {
		return nil
	}

	var content struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*event.Event.Message.Content), &content); err != nil {
		return err
	}

	text := strings.TrimSpace(content.Text)
	cmdIndex := strings.Index(text, "/")
	if cmdIndex == -1 {
		return nil
	}

	cmdPart := strings.TrimSpace(text[cmdIndex:])
	parts := strings.Fields(cmdPart)
	if len(parts) == 0 {
		return nil
	}

	command := parts[0]
	args := parts[1:]

	var openID string
	if event.Event.Sender != nil && event.Event.Sender.SenderId != nil && event.Event.Sender.SenderId.OpenId != nil {
		openID = *event.Event.Sender.SenderId.OpenId
	}
	if openID == "" {
		return nil
	}

	commandsMu.RLock()
	handler, exists := commandHandlers[command]
	commandsMu.RUnlock()

	if !exists {
		return nil
	}

	return handler(context.Background(), args, openID)
}

// InitCommands 初始化命令处理器
func InitCommands(userService UserBinder, statsService StatsProvider, systemService SystemProvider) {
	userBinderInstance = userService
	statsProviderInstance = statsService
	systemProviderInstance = systemService

	RegisterCommand("/bind", func(ctx context.Context, args []string, openID string) error {
		if len(args) == 0 {
			return sendBindErrorCard(ctx, "请提供邮箱地址\n\n用法：`/bind 邮箱`")
		}

		email := args[0]
		if strings.HasPrefix(email, "[") && strings.Contains(email, "](mailto:") {
			if endBracket := strings.Index(email, "]"); endBracket > 1 {
				email = email[1:endBracket]
			}
		}

		if userBinderInstance != nil {
			if err := userBinderInstance.BindFeishuByEmail(email, openID); err != nil {
				logger.Error("[Feishu] 绑定失败: %v", err)
				return sendBindErrorCard(ctx, err.Error())
			}
			logger.Info("[Feishu] 绑定成功: %s", email)
			return sendBindSuccessCard(ctx, email)
		}
		return nil
	})

	RegisterCommand("/stats", func(ctx context.Context, args []string, openID string) error {
		if statsProviderInstance == nil {
			return sendStatsErrorCard(ctx, "统计服务未初始化")
		}

		stats, err := statsProviderInstance.GetDashboardStats()
		if err != nil {
			logger.Error("[Feishu] 获取统计失败: %v", err)
			return sendStatsErrorCard(ctx, "获取统计数据失败，请稍后重试")
		}

		globalMu.RLock()
		client := globalClient
		globalMu.RUnlock()

		if client != nil && client.IsEnabled() {
			card := buildStatsCard(stats)
			return client.SendMessage(ctx, card)
		}
		return nil
	})

	RegisterCommand("/system", func(ctx context.Context, args []string, openID string) error {
		if systemProviderInstance == nil {
			return sendSystemErrorCard(ctx, "系统服务未初始化")
		}

		status, err := systemProviderInstance.GetSystemStatus(ctx)
		if err != nil {
			logger.Error("[Feishu] 获取系统状态失败: %v", err)
			return sendSystemErrorCard(ctx, "获取系统状态失败，请稍后重试")
		}

		globalMu.RLock()
		client := globalClient
		globalMu.RUnlock()

		if client != nil && client.IsEnabled() {
			card := buildSystemCard(status)
			return client.SendMessage(ctx, card)
		}
		return nil
	})

	RegisterCommand("/help", func(ctx context.Context, args []string, openID string) error {
		helpText := "**可用命令：**\n\n" +
			"• `/bind 邮箱` - 绑定飞书账号\n" +
			"• `/stats` - 查看站点仪表盘统计\n" +
			"• `/system` - 查看系统状态\n" +
			"• `/help` - 显示帮助信息"

		globalMu.RLock()
		client := globalClient
		globalMu.RUnlock()

		if client != nil && client.IsEnabled() {
			card := buildHelpCard(helpText)
			return client.SendMessage(ctx, card)
		}
		return nil
	})
}

// buildHelpCard 构建帮助卡片
func buildHelpCard(helpText string) string {
	elements := []interface{}{
		&MarkdownElement{Tag: "markdown", Content: helpText},
	}
	card, _ := buildCard("📖 命令帮助", "blue", elements)
	return card
}

// buildStatsCard 构建统计卡片
func buildStatsCard(stats *dto.DashboardStats) string {
	content := fmt.Sprintf(
		"**总览**\n"+
			"文章：`%d`  友链：`%d`  动态：`%d`\n"+
			"评论：`%d`  用户：`%d`\n"+
			"总 PV：`%d`  总 UV：`%d`\n\n"+
			"**今日**\n"+
			"PV：`%d`  UV：`%d`\n"+
			"评论：`%d`  新用户：`%d`\n\n"+
			"**较昨日**\n"+
			"PV：`%s`  UV：`%s`\n"+
			"评论：`%s`  用户：`%s`",
		stats.TotalArticles,
		stats.TotalFriends,
		stats.TotalMoments,
		stats.TotalComments,
		stats.TotalUsers,
		stats.TotalViews,
		stats.TotalVisitors,
		stats.TodayViews,
		stats.TodayVisitors,
		stats.TodayComments,
		stats.TodayUsers,
		formatGrowthRate(stats.ViewsGrowth),
		formatGrowthRate(stats.VisitorsGrowth),
		formatGrowthRate(stats.CommentsGrowth),
		formatGrowthRate(stats.UsersGrowth),
	)

	elements := []interface{}{
		&MarkdownElement{Tag: "markdown", Content: content},
	}
	card, _ := buildCard("📊 站点统计", "blue", elements)
	return card
}

// buildSystemCard 构建系统状态卡片
func buildSystemCard(status *SystemStatus) string {
	content := fmt.Sprintf(
		"**服务状态**\n"+
			"数据库：`%s`\n"+
			"存储：`%s`\n"+
			"邮件：`%s`\n"+
			"飞书：`%s`\n\n"+
			"**资源使用**\n"+
			"CPU：`%.2f%%`\n"+
			"内存：`%s / %s`\n"+
			"磁盘：`%s / %s`",
		status.DBStatus,
		status.StorageStatus,
		status.EmailStatus,
		status.FeishuStatus,
		status.CPUUsage,
		formatBytes(status.MemoryUsed),
		formatBytes(status.MemoryTotal),
		formatBytes(status.DiskUsed),
		formatBytes(status.DiskTotal),
	)

	elements := []interface{}{
		&MarkdownElement{Tag: "markdown", Content: content},
	}
	card, _ := buildCard("🔧 系统状态", "blue", elements)
	return card
}

func formatGrowthRate(rate float64) string {
	return fmt.Sprintf("%+.2f%%", rate)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// sendBindSuccessCard 发送绑定成功卡片
func sendBindSuccessCard(ctx context.Context, email string) error {
	globalMu.RLock()
	client := globalClient
	globalMu.RUnlock()

	if client != nil && client.IsEnabled() {
		content := fmt.Sprintf("绑定成功\n\n邮箱：`%s`", email)
		elements := []interface{}{
			&MarkdownElement{Tag: "markdown", Content: content},
		}
		card, _ := buildCard("✅ 绑定成功", "green", elements)
		return client.SendMessage(ctx, card)
	}
	return nil
}

// sendBindErrorCard 发送绑定失败卡片
func sendBindErrorCard(ctx context.Context, errMsg string) error {
	globalMu.RLock()
	client := globalClient
	globalMu.RUnlock()

	if client != nil && client.IsEnabled() {
		content := fmt.Sprintf("绑定失败\n\n%s", errMsg)
		elements := []interface{}{
			&MarkdownElement{Tag: "markdown", Content: content},
		}
		card, _ := buildCard("❌ 绑定失败", "red", elements)
		return client.SendMessage(ctx, card)
	}
	return nil
}

func sendStatsErrorCard(ctx context.Context, errMsg string) error {
	globalMu.RLock()
	client := globalClient
	globalMu.RUnlock()

	if client != nil && client.IsEnabled() {
		elements := []interface{}{
			&MarkdownElement{Tag: "markdown", Content: errMsg},
		}
		card, _ := buildCard("❌ 统计失败", "red", elements)
		return client.SendMessage(ctx, card)
	}
	return nil
}

func sendSystemErrorCard(ctx context.Context, errMsg string) error {
	globalMu.RLock()
	client := globalClient
	globalMu.RUnlock()

	if client != nil && client.IsEnabled() {
		elements := []interface{}{
			&MarkdownElement{Tag: "markdown", Content: errMsg},
		}
		card, _ := buildCard("❌ 系统状态失败", "red", elements)
		return client.SendMessage(ctx, card)
	}
	return nil
}
