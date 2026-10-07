package service

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/feishu"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/utils"

	"github.com/mmcdole/gofeed"
)

// friendCircleCache 朋友圈数据缓存
type friendCircleCache struct {
	mu        sync.RWMutex
	items     []dto.FriendCircleItemResponse
	expiresAt time.Time
}

// RssFeedService RSS订阅服务
type RssFeedService struct {
	repo            *repository.RssFeedRepository
	parser          *gofeed.Parser
	notificationSvc *NotificationService
	circleCache     *friendCircleCache
}

// NewRssFeedService 创建RSS订阅服务实例
func NewRssFeedService(repo *repository.RssFeedRepository, notificationSvc *NotificationService) *RssFeedService {
	return &RssFeedService{
		repo:            repo,
		parser:          gofeed.NewParser(),
		notificationSvc: notificationSvc,
		circleCache:     &friendCircleCache{},
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

	articles, total, err := s.repo.List(
		ctx,
		req.Page, req.PageSize,
		req.Keyword, req.FriendID, req.IsRead,
		req.StartTime, req.EndTime,
	)
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
			BlockCircle: article.BlockCircle,
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

// GetFriendCircle 获取友圈文章
func (s *RssFeedService) GetFriendCircle(ctx context.Context, req *dto.FriendCircleRequest) (*dto.FriendCircleResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	items, err := s.friendCircleItems(ctx)
	if err != nil {
		return nil, err
	}

	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		lowerKeyword := strings.ToLower(keyword)
		matched := make([]dto.FriendCircleItemResponse, 0, len(items))
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Title), lowerKeyword) ||
				(item.Author != nil && strings.Contains(strings.ToLower(item.Author.Name), lowerKeyword)) {
				matched = append(matched, item)
			}
		}
		items = matched
	}

	total := int64(len(items))
	start := (page - 1) * pageSize
	if start >= len(items) {
		return &dto.FriendCircleResponse{
			List:     []dto.FriendCircleItemResponse{},
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, nil
	}

	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return &dto.FriendCircleResponse{
		List:     items[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetFriendCircleRandom 随机获取一篇友圈文章
func (s *RssFeedService) GetFriendCircleRandom(ctx context.Context) (*dto.FriendCircleItemResponse, error) {
	items, err := s.friendCircleItems(ctx)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	picked := items[rand.Intn(len(items))] //nolint:gosec
	return &picked, nil
}

// GetFriendCircleStats 统计友圈数据
func (s *RssFeedService) GetFriendCircleStats(ctx context.Context) (*dto.FriendCircleStatsResponse, error) {
	items, err := s.friendCircleItems(ctx)
	if err != nil {
		return nil, err
	}

	siteCount, err := s.repo.CountSubscribedFriends(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var todayCount int64
	for _, item := range items {
		if item.PublishedAt == nil {
			continue
		}

		if !item.PublishedAt.Before(todayStart) {
			todayCount++
		}
	}

	return &dto.FriendCircleStatsResponse{
		Total:      int64(len(items)),
		SiteCount:  siteCount,
		TodayCount: todayCount,
	}, nil
}

// friendCircleItems 获取友圈文章
func (s *RssFeedService) friendCircleItems(ctx context.Context) ([]dto.FriendCircleItemResponse, error) {
	s.circleCache.mu.RLock()
	if s.circleCache.items != nil && time.Now().Before(s.circleCache.expiresAt) {
		items := s.circleCache.items
		s.circleCache.mu.RUnlock()
		return items, nil
	}
	s.circleCache.mu.RUnlock()

	s.circleCache.mu.Lock()
	defer s.circleCache.mu.Unlock()

	if s.circleCache.items != nil && time.Now().Before(s.circleCache.expiresAt) {
		return s.circleCache.items, nil
	}

	articles, err := s.repo.ListFriendCircle(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.FriendCircleItemResponse, 0, len(articles))
	for _, article := range articles {
		item := dto.FriendCircleItemResponse{
			ID:          article.ID,
			Title:       article.Title,
			Link:        article.Link,
			PublishedAt: utils.ToJSONTime(article.PublishedAt),
		}

		if article.Friend != nil {
			item.Author = &dto.FriendCircleAuthorResponse{
				ID:     article.Friend.ID,
				Name:   article.Friend.Name,
				Avatar: article.Friend.Avatar,
				URL:    article.Friend.URL,
			}
		}

		items = append(items, item)
	}

	s.circleCache.items = items
	s.circleCache.expiresAt = time.Now().Add(10 * time.Minute)
	return items, nil
}

// SetArticleBlockCircle 设置单篇 RSS 文章的友圈屏蔽状态
func (s *RssFeedService) SetArticleBlockCircle(ctx context.Context, id uint, block bool) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errors.New("文章不存在")
	}
	if err := s.repo.SetBlockCircle(ctx, id, block); err != nil {
		return err
	}
	s.InvalidateCircleCache()
	return nil
}

// InvalidateCircleCache 使友圈缓存失效，下次请求重新查库
func (s *RssFeedService) InvalidateCircleCache() {
	s.circleCache.mu.Lock()
	defer s.circleCache.mu.Unlock()
	s.circleCache.items = nil
	s.circleCache.expiresAt = time.Time{}
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

	s.InvalidateCircleCache()

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
	seenTitles := make(map[string]struct{})

	allowedHosts := make([]string, 0, 2)
	for _, rawURL := range []string{friend.RSSUrl, friend.URL} {
		if host := hostOf(rawURL); host != "" {
			allowedHosts = append(allowedHosts, host)
		}
	}

	for _, item := range feed.Items {
		link := item.Link
		if link == "" {
			continue
		}

		u, err := url.Parse(link)
		if err != nil {
			continue
		}
		if u.Host == "" {
			base, err := url.Parse(friend.URL)
			if err != nil || base.Host == "" {
				continue
			}
			link = base.ResolveReference(u).String()
		}

		title := strings.TrimSpace(item.Title)
		if title == "" {
			continue
		}

		articleHost := hostOf(link)
		if len(allowedHosts) > 0 && !matchesAnyHost(articleHost, allowedHosts) {
			continue
		}

		if _, dup := seenTitles[title]; dup {
			continue
		}

		exists, err := s.repo.ExistsByLink(ctx, link)
		if err != nil || exists {
			continue
		}

		duplicated, err := s.repo.ExistsByFriendAndTitle(ctx, friend.ID, title)
		if err != nil || duplicated {
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

		articlesToCreate = append(articlesToCreate, model.RssArticle{
			FriendID:    friend.ID,
			Title:       title,
			Link:        link,
			PublishedAt: publishedAt,
			IsRead:      isFirstSubscribe,
		})
		seenTitles[title] = struct{}{}
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

// matchesAnyHost 判断域名是否命中任一基准：相同或互为子域
func matchesAnyHost(host string, bases []string) bool {
	if host == "" {
		return false
	}
	for _, base := range bases {
		if base == "" {
			continue
		}
		if host == base ||
			strings.HasSuffix(host, "."+base) ||
			strings.HasSuffix(base, "."+host) {
			return true
		}
	}
	return false
}

// hostOf 提取 URL 主机名：小写、去端口、去 www 前缀
func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := strings.ToLower(u.Hostname())
	if host == "" {
		return ""
	}
	return strings.TrimPrefix(host, "www.")
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
	affected, err := s.repo.DeleteOrphaned(ctx)
	if err != nil {
		return err
	}

	if affected > 0 {
		s.InvalidateCircleCache()
	}

	return nil
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
