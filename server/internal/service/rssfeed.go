package service

import (
	"context"
	"errors"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/feishu"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/utils"

	"github.com/mmcdole/gofeed"
)

// RssFeedService RSS订阅服务
type RssFeedService struct {
	repo            *repository.RssFeedRepository
	parser          *gofeed.Parser
	notificationSvc *NotificationService
}

// NewRssFeedService 创建RSS订阅服务实例
func NewRssFeedService(repo *repository.RssFeedRepository, notificationSvc *NotificationService) *RssFeedService {
	return &RssFeedService{
		repo:            repo,
		parser:          gofeed.NewParser(),
		notificationSvc: notificationSvc,
	}
}

// List 获取RSS文章列表
func (s *RssFeedService) List(ctx context.Context, req *dto.ListRssArticleRequest) (*dto.RssArticleListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	articles, total, err := s.repo.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.repo.CountUnread(ctx)
	if err != nil {
		unreadCount = 0
	}

	list := make([]dto.RssArticleResponse, 0, len(articles))
	for _, article := range articles {
		item := dto.RssArticleResponse{
			ID:          article.ID,
			FriendID:    article.FriendID,
			Title:       article.Title,
			Link:        article.Link,
			IsRead:      article.IsRead,
			PublishedAt: utils.ToJSONTime(article.PublishedAt),
			CreatedAt:   utils.ToJSONTime(&article.CreatedAt),
		}

		if article.Friend != nil {
			item.FriendName = article.Friend.Name
			item.FriendURL = article.Friend.URL
		}

		list = append(list, item)
	}

	return &dto.RssArticleListResponse{
		List:        list,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		UnreadCount: unreadCount,
	}, nil
}

// MarkRead 标记文章已读
func (s *RssFeedService) MarkRead(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("文章不存在")
	}
	return s.repo.MarkRead(ctx, id)
}

// MarkAllRead 全部标记已读
func (s *RssFeedService) MarkAllRead(ctx context.Context) (int64, error) {
	return s.repo.MarkAllRead(ctx)
}

// MarkAllReadFromFeishu 从飞书调用全部标记已读
func (s *RssFeedService) MarkAllReadFromFeishu(ctx context.Context) error {
	_, err := s.repo.MarkAllRead(ctx)
	return err
}

// RefreshAllFeeds 刷新所有RSS订阅源
func (s *RssFeedService) RefreshAllFeeds() error {
	ctx := context.Background()
	friends, err := s.repo.GetFriendsWithRSS(ctx)
	if err != nil {
		return err
	}

	for _, friend := range friends {
		_ = s.refreshFriendFeed(ctx, &friend)
	}

	return nil
}

// refreshFriendFeed 刷新单个友链的RSS订阅
func (s *RssFeedService) refreshFriendFeed(ctx context.Context, friend *model.Friend) error {
	feed, err := s.parser.ParseURL(friend.RSSUrl)
	if err != nil {
		return err
	}

	isFirstSubscribe := friend.RSSLatime == nil

	var articlesToCreate []model.RssArticle

	for _, item := range feed.Items {
		if item.Link == "" {
			continue
		}

		exists, err := s.repo.ExistsByLink(ctx, item.Link)
		if err != nil || exists {
			continue
		}

		var publishedAt *time.Time
		switch {
		case item.PublishedParsed != nil:
			publishedAt = item.PublishedParsed
		case item.UpdatedParsed != nil:
			publishedAt = item.UpdatedParsed
		case item.Published != "":
			// Fallback: try to parse non-standard pubDate formats
			if t, err := parseRSSDate(item.Published); err == nil {
				publishedAt = &t
			}
		}

		article := model.RssArticle{
			FriendID:    friend.ID,
			Title:       item.Title,
			Link:        item.Link,
			PublishedAt: publishedAt,
			IsRead:      isFirstSubscribe,
		}

		articlesToCreate = append(articlesToCreate, article)
	}

	if len(articlesToCreate) > 0 {
		if err := s.repo.CreateBatch(ctx, articlesToCreate); err != nil {
			return err
		}
	}

	// 更新最后更新时间为最新文章的发布时间
	latestTime, err := s.repo.GetLatestPublishedTime(ctx, friend.ID)
	if err != nil {
		return err
	}
	if latestTime == nil {
		return nil
	}
	return s.repo.UpdateFriendRSSLatime(ctx, friend.ID, *latestTime)
}

// parseRSSDate 尝试解析多种RSS日期格式
func parseRSSDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 -0700",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unable to parse date")
}

// CleanOrphanedArticles 清理孤立文章
func (s *RssFeedService) CleanOrphanedArticles() error {
	ctx := context.Background()
	_, err := s.repo.DeleteOrphaned(ctx)
	return err
}

// SendDailyPush 发送每日RSS订阅推送
func (s *RssFeedService) SendDailyPush() error {
	ctx := context.Background()

	unreadCount, err := s.repo.CountUnread(ctx)
	if err != nil {
		return err
	}

	if unreadCount == 0 {
		return nil
	}

	articles, err := s.repo.ListUnread(ctx, 20)
	if err != nil {
		return err
	}

	var articleItems []feishu.RssArticleItem
	for _, article := range articles {
		friendName := ""
		if article.Friend != nil {
			friendName = article.Friend.Name
		}
		articleItems = append(articleItems, feishu.RssArticleItem{
			Title:      article.Title,
			Link:       article.Link,
			FriendName: friendName,
		})
	}

	if s.notificationSvc != nil {
		go func() {
			if err := s.notificationSvc.NotifyRssFeedDaily(int(unreadCount), articleItems); err != nil {
				logger.Error("发送每日通知失败: %v", err)
			}
		}()
	}

	return nil
}
