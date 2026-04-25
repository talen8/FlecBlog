package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/utils"

	"gorm.io/gorm"
)

// 页面标题映射
var pageTitle = map[string]string{
	"friend":  "友链",
	"moment":  "动态",
	"message": "留言",
}

// CommentService 评论服务
type CommentService struct {
	repo                *repository.CommentRepository
	articleRepo         *repository.ArticleRepository
	userRepo            *repository.UserRepository
	notificationService *NotificationService
	fileService         *FileService
}

// NewCommentService 创建评论服务实例
func NewCommentService(repo *repository.CommentRepository, articleRepo *repository.ArticleRepository, userRepo *repository.UserRepository, notificationService *NotificationService, fileService *FileService) *CommentService {
	return &CommentService{
		repo:                repo,
		articleRepo:         articleRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		fileService:         fileService,
	}
}

// ============ 通用服务 ============

// GetForWeb 获取评论详情
func (s *CommentService) GetForWeb(ctx context.Context, id uint) (*dto.CommentResponse, error) {
	comment, err := s.repo.GetForWeb(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment), nil
}

// ============ 前台服务 ============

// ListForWeb 获取前台评论列表
func (s *CommentService) ListForWeb(ctx context.Context, req *dto.CommentQueryForWebRequest) ([]dto.CommentResponse, int64, error) {
	return s.GetByTarget(ctx, req.TargetType, req.TargetKey, req.Page, req.PageSize)
}

// GetByTarget 获取目标的评论列表（扁平化结构）
func (s *CommentService) GetByTarget(ctx context.Context, targetType, targetKey string, page, pageSize int) ([]dto.CommentResponse, int64, error) {
	// 获取顶级评论
	topComments, total, err := s.repo.GetByTarget(ctx, targetType, targetKey, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if len(topComments) == 0 {
		return []dto.CommentResponse{}, 0, nil
	}

	// 收集顶级评论ID
	rootIDs := make([]uint, len(topComments))
	for i, comment := range topComments {
		rootIDs[i] = comment.ID
	}

	// 批量获取所有回复
	replies, err := s.repo.GetRepliesByRootIDs(ctx, rootIDs)
	if err != nil {
		return nil, 0, err
	}

	// 构建回复映射表
	repliesMap := make(map[uint][]dto.CommentResponse)
	for _, reply := range replies {
		if reply.RootID != nil {
			replyDTO := s.toCommentResponse(&reply)
			replyDTO.Replies = []dto.CommentResponse{}
			repliesMap[*reply.RootID] = append(repliesMap[*reply.RootID], *replyDTO)
		}
	}

	// 构建扁平化结构，将正常评论和已删除/隐藏评论分开
	var normalComments []dto.CommentResponse
	var hiddenComments []dto.CommentResponse

	for _, comment := range topComments {
		commentResp := s.toCommentResponse(&comment)

		// 添加回复列表
		if replies, ok := repliesMap[comment.ID]; ok {
			commentResp.Replies = replies
		} else {
			commentResp.Replies = []dto.CommentResponse{}
		}

		// 区分正常评论和已删除/隐藏评论
		if comment.DeletedAt.Valid || comment.Status == 0 {
			hiddenComments = append(hiddenComments, *commentResp)
		} else {
			normalComments = append(normalComments, *commentResp)
		}
	}

	// 正常评论在前，已删除/隐藏评论在后（排序到末尾）
	//nolint:gocritic // 有意创建新切片，不修改原切片
	result := append(normalComments, hiddenComments...)

	return result, total, nil
}

// Create 创建评论
func (s *CommentService) Create(ctx context.Context, req *dto.CreateCommentRequest, userID uint) (*dto.CommentResponse, error) {
	// 如果 userID 为 0，说明是游客评论
	if userID == 0 {
		// 查找或创建游客用户
		guestUser, err := s.userRepo.GetGuestByEmail(req.Email)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			// 检查邮箱是否已被其他角色的用户使用（注册用户）
			if s.userRepo.ExistsByEmail(req.Email) {
				return nil, errors.New("邮箱已被注册，请登录")
			}

			// 创建新游客用户
			guestUser = &model.User{
				Email:    req.Email,
				Nickname: req.Nickname,
				Website:  req.Website,
				Role:     model.RoleGuest,
				Password: "",
			}

			if err := s.userRepo.Create(guestUser); err != nil {
				return nil, err
			}
		} else {
			// 游客已存在，更新昵称和网站
			needUpdate := false
			if guestUser.Nickname != req.Nickname {
				guestUser.Nickname = req.Nickname
				needUpdate = true
			}
			if guestUser.Website != req.Website {
				guestUser.Website = req.Website
				needUpdate = true
			}
			if needUpdate {
				if err := s.userRepo.Update(guestUser); err != nil {
					return nil, err
				}
			}
		}
		userID = guestUser.ID
	}

	return s.createComment(ctx, req, userID)
}

// Update 更新评论
func (s *CommentService) Update(ctx context.Context, id uint, req *dto.UpdateCommentRequest, userID uint) (*dto.CommentResponse, error) {
	comment, err := s.repo.GetForWeb(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}

	// 权限检查
	if comment.UserID != userID {
		return nil, errors.New("无权修改此评论")
	}

	oldContent := comment.Content
	comment.Content = req.Content
	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, err
	}

	s.updateImageStatus(oldContent, req.Content)

	return s.toCommentResponse(comment), nil
}

// DeleteForWeb 软删除评论
func (s *CommentService) DeleteForWeb(ctx context.Context, id uint, userID uint) error {
	comment, err := s.repo.GetForWeb(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 权限检查
	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}

	// 只删除评论本身，子评论保留
	if err := s.repo.DeleteForWeb(ctx, id); err != nil {
		return err
	}

	return nil
}

// ============ 后台管理服务 ============

// List 获取评论列表
func (s *CommentService) List(ctx context.Context, req *dto.CommentQueryRequest) ([]dto.CommentListResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	comments, total, err := s.repo.List(
		ctx,
		offset, req.PageSize,
		req.Keyword,
		req.Status, req.IsDeleted, req.IsSub,
		req.StartTime, req.EndTime,
	)
	if err != nil {
		return nil, 0, err
	}

	// 转换为后台响应格式
	result := make([]dto.CommentListResponse, len(comments))
	for i, comment := range comments {
		result[i] = *s.toDTO(&comment)
	}

	return result, total, nil
}

// Get 获取评论详情
func (s *CommentService) Get(ctx context.Context, id uint) (*dto.CommentListResponse, error) {
	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDTO(comment), nil
}

// ToggleStatus 切换评论状态
func (s *CommentService) ToggleStatus(ctx context.Context, id uint) error {
	comment, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 切换状态
	if comment.Status == 0 {
		comment.Status = 1
	} else {
		comment.Status = 0
	}

	return s.repo.UpdateStatus(ctx, id, comment.Status)
}

// Delete 硬删除评论
func (s *CommentService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ============ 辅助方法 ============

// createComment 创建评论的通用逻辑
func (s *CommentService) createComment(ctx context.Context, req *dto.CreateCommentRequest, userID uint) (*dto.CommentResponse, error) {
	// 验证目标类型
	switch req.TargetType {
	case "article":
		// 文章评论：验证 slug 是否存在
		if _, err := s.articleRepo.GetBySlug(req.TargetKey); err != nil {
			return nil, errors.New("文章不存在")
		}
	case "page":
		// 页面评论：验证页面key是否在配置中
		if _, ok := pageTitle[req.TargetKey]; !ok {
			return nil, errors.New("无效的页面标识")
		}
	}

	// 解析 User-Agent 获取浏览器和操作系统信息
	browser, os := utils.ParseUserAgent(req.UserAgent)

	// 获取 IP 地理位置
	location := utils.GetIPLocation(req.IP)

	comment := &model.Comment{
		Content:    req.Content,
		TargetType: req.TargetType,
		TargetKey:  req.TargetKey,
		UserID:     userID,
		ParentID:   req.ParentID,
		Status:     1, // 默认显示
		IP:         req.IP,
		Location:   location,
		Browser:    browser,
		OS:         os,
	}

	// 处理回复评论的关系
	if req.ParentID != nil {
		parentComment, err := s.repo.GetForWeb(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("父评论不存在")
		}

		// 验证父评论属于同一目标
		if parentComment.TargetType != req.TargetType || parentComment.TargetKey != req.TargetKey {
			return nil, errors.New("不能跨目标回复评论")
		}

		// 设置根评论ID
		if parentComment.RootID != nil {
			comment.RootID = parentComment.RootID
		} else {
			comment.RootID = &parentComment.ID
		}

		// 设置回复目标用户
		comment.ReplyTo = &parentComment.UserID
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// 标记评论中的图片为已使用
	s.markImagesAsUsed(comment.Content)

	// 异步发送通知
	go s.sendNotifications(ctx, comment, userID)

	return s.GetForWeb(ctx, comment.ID)
}

// toCommentResponse 转换为前台响应格式
func (s *CommentService) toCommentResponse(comment *model.Comment) *dto.CommentResponse {
	resp := &dto.CommentResponse{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		CreatedAt: utils.NewJSONTime(comment.CreatedAt),
		Replies:   []dto.CommentResponse{},
	}

	// 处理已删除评论
	if comment.DeletedAt.Valid {
		resp.IsDeleted = true
		resp.Content = "该评论已被删除"
		resp.User.ID = comment.User.ID
		resp.User.Nickname = "匿名用户"
		resp.User.Avatar = ""

		if comment.ReplyUser != nil {
			resp.ReplyUser = &struct {
				ID        uint   `json:"id"`
				Nickname  string `json:"nickname"`
				Avatar    string `json:"avatar"`
				Badge     string `json:"badge"`
				EmailHash string `json:"email_hash"`
				Website   string `json:"website"`
				Role      string `json:"role"`
			}{
				ID:        comment.ReplyUser.ID,
				Nickname:  "匿名用户",
				Avatar:    "",
				Badge:     "",
				EmailHash: "",
				Website:   "",
				Role:      "guest",
			}
		}
		return resp
	}

	// 处理隐藏评论
	if comment.Status == 0 {
		resp.Content = "该评论已被隐藏"
		resp.User.ID = comment.User.ID
		resp.User.Nickname = "匿名用户"
		resp.User.Avatar = ""
		resp.User.EmailHash = ""
		resp.User.Website = ""

		if comment.ReplyUser != nil {
			resp.ReplyUser = &struct {
				ID        uint   `json:"id"`
				Nickname  string `json:"nickname"`
				Avatar    string `json:"avatar"`
				Badge     string `json:"badge"`
				EmailHash string `json:"email_hash"`
				Website   string `json:"website"`
				Role      string `json:"role"`
			}{
				ID:        comment.ReplyUser.ID,
				Nickname:  "匿名用户",
				Avatar:    "",
				Badge:     "",
				EmailHash: "",
				Website:   "",
				Role:      "guest",
			}
		}
		return resp
	}

	// 正常显示评论
	resp.Content = comment.Content
	resp.Location = comment.Location
	resp.Browser = comment.Browser
	resp.OS = comment.OS

	// 处理用户信息
	resp.User.ID = comment.User.ID
	if comment.User.DeletedAt.Valid {
		resp.User.Nickname = "已删除用户"
		resp.User.Avatar = ""
		resp.User.Badge = ""
		resp.User.EmailHash = ""
		resp.User.Website = ""
		resp.User.Role = "guest"
	} else {
		resp.User.Nickname = comment.User.Nickname
		resp.User.Avatar = comment.User.Avatar
		resp.User.Badge = comment.User.Badge
		resp.User.EmailHash = utils.GetEmailHash(comment.User.Email)
		resp.User.Website = comment.User.Website
		resp.User.Role = string(comment.User.Role)
	}

	// 处理回复用户信息
	if comment.ReplyUser != nil {
		resp.ReplyUser = &struct {
			ID        uint   `json:"id"`
			Nickname  string `json:"nickname"`
			Avatar    string `json:"avatar"`
			Badge     string `json:"badge"`
			EmailHash string `json:"email_hash"`
			Website   string `json:"website"`
			Role      string `json:"role"`
		}{
			ID: comment.ReplyUser.ID,
		}

		if comment.ReplyUser.DeletedAt.Valid {
			resp.ReplyUser.Nickname = "已删除用户"
			resp.ReplyUser.Avatar = ""
			resp.ReplyUser.Badge = ""
			resp.ReplyUser.EmailHash = ""
			resp.ReplyUser.Website = ""
			resp.ReplyUser.Role = "guest"
		} else {
			resp.ReplyUser.Nickname = comment.ReplyUser.Nickname
			resp.ReplyUser.Avatar = comment.ReplyUser.Avatar
			resp.ReplyUser.Badge = comment.ReplyUser.Badge
			resp.ReplyUser.EmailHash = utils.GetEmailHash(comment.ReplyUser.Email)
			resp.ReplyUser.Website = comment.ReplyUser.Website
			resp.ReplyUser.Role = string(comment.ReplyUser.Role)
		}
	}

	return resp
}

// toDTO 转换为后台管理响应格式
func (s *CommentService) toDTO(comment *model.Comment) *dto.CommentListResponse {
	resp := &dto.CommentListResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		Status:    comment.Status,
		ParentID:  comment.ParentID,
		CreatedAt: utils.NewJSONTime(comment.CreatedAt),
	}

	// 处理软删除时间
	if comment.DeletedAt.Valid {
		deletedTime := utils.NewJSONTime(comment.DeletedAt.Time)
		resp.DeletedAt = &deletedTime
	}

	// 填充 Target 信息
	resp.Target.Type = comment.TargetType
	resp.Target.Key = comment.TargetKey
	resp.Target.Title = s.getTargetTitle(comment.TargetType, comment.TargetKey)

	// 填充用户信息
	resp.User.ID = comment.User.ID
	resp.User.Email = comment.User.Email // 后台显示邮箱
	resp.User.Nickname = comment.User.Nickname
	resp.User.Avatar = comment.User.Avatar
	resp.User.Badge = comment.User.Badge

	return resp
}

// getTargetTitle 获取目标标题（后台使用）
func (s *CommentService) getTargetTitle(targetType, targetKey string) string {
	switch targetType {
	case "article":
		// 文章：通过 slug 查询数据库获取标题
		article, err := s.articleRepo.GetBySlug(targetKey)
		if err != nil {
			return "文章已删除"
		}
		return article.Title

	case "page":
		// 页面：从配置获取固定标题
		if title, ok := pageTitle[targetKey]; ok {
			return title
		}
		return "未知页面"

	default:
		return "未知类型"
	}
}

// sendNotifications 发送评论通知
func (s *CommentService) sendNotifications(_ context.Context, comment *model.Comment, senderID uint) {
	if s.notificationService == nil {
		return
	}

	// 使用独立的 background context，避免原请求上下文取消影响通知发送
	notifyCtx := context.Background()

	// 获取目标标题
	targetTitle := s.getTargetTitle(comment.TargetType, comment.TargetKey)

	// 1. 如果是回复评论，通知被回复者
	if comment.ReplyTo != nil {
		_ = s.notificationService.NotifyCommentReply(notifyCtx, senderID, *comment.ReplyTo, comment, targetTitle)
	}

	// 2. 通知所有管理员（有新评论），排除发送者自己避免自通知
	_ = s.notificationService.NotifyCommentToAdmins(notifyCtx, senderID, comment, targetTitle, &senderID)
}

// markImagesAsUsed 标记评论内容中的图片为已使用
func (s *CommentService) markImagesAsUsed(content string) {
	if s.fileService == nil || content == "" {
		return
	}

	for _, imageURL := range extractCommentImageURLs(content) {
		_ = s.fileService.MarkAsUsed(imageURL)
	}
}

// updateImageStatus 对比评论内容变化，更新图片状态
func (s *CommentService) updateImageStatus(oldContent, newContent string) {
	if s.fileService == nil || oldContent == newContent {
		return
	}

	oldImages := extractCommentImageURLs(oldContent)
	newImages := extractCommentImageURLs(newContent)

	oldMap := make(map[string]bool, len(oldImages))
	for _, url := range oldImages {
		oldMap[url] = true
	}

	newMap := make(map[string]bool, len(newImages))
	for _, url := range newImages {
		newMap[url] = true
		if !oldMap[url] {
			_ = s.fileService.MarkAsUsed(url)
		}
	}

	for _, url := range oldImages {
		if !newMap[url] {
			_ = s.fileService.MarkAsUnused(url)
		}
	}
}

func extractCommentImageURLs(content string) []string {
	if content == "" {
		return nil
	}

	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	urls := make([]string, 0, len(matches))
	seen := make(map[string]bool, len(matches))

	for _, match := range matches {
		if len(match) <= 1 {
			continue
		}
		url := strings.TrimSpace(match[1])
		if url == "" || seen[url] {
			continue
		}
		seen[url] = true
		urls = append(urls, url)
	}

	return urls
}

// ============ 数据导入导出方法 ============

// ImportFromArtalk 从Artalk格式导入评论
func (s *CommentService) ImportFromArtalk(ctx context.Context, jsonData []byte) (*dto.ImportCommentsResult, error) {
	var artalkComments []dto.ArtalkCommentData
	if err := json.Unmarshal(jsonData, &artalkComments); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if len(artalkComments) == 0 {
		return nil, errors.New("没有找到有效的评论数据")
	}

	result := &dto.ImportCommentsResult{
		Total: len(artalkComments),
	}

	// 创建映射表和缓存
	commentIDMap := make(map[string]uint)
	userCache := make(map[string]*model.User)

	// 分两轮处理：顶级评论和回复
	topLevel, replies := s.separateCommentsByLevel(artalkComments)

	// 处理顶级评论
	s.processArtalkComments(ctx, topLevel, result, commentIDMap, userCache)

	// 处理回复
	s.processArtalkComments(ctx, replies, result, commentIDMap, userCache)

	return result, nil
}

// ============ 辅助方法 ============

func (s *CommentService) separateCommentsByLevel(comments []dto.ArtalkCommentData) (topLevel, replies []dto.ArtalkCommentData) {
	for _, c := range comments {
		if c.RID == "" || c.RID == "0" {
			topLevel = append(topLevel, c)
		} else {
			replies = append(replies, c)
		}
	}
	return
}

func (s *CommentService) processArtalkComments(
	ctx context.Context,
	comments []dto.ArtalkCommentData,
	result *dto.ImportCommentsResult,
	commentIDMap map[string]uint,
	userCache map[string]*model.User,
) {
	for i, artalk := range comments {
		// 解析评论数据
		parsed, err := s.parseArtalkComment(&artalk)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.ImportCommentError{
				Index:   i,
				Content: truncate(artalk.Content, 50),
				Error:   err.Error(),
			})
			continue
		}

		// 处理用户
		user, isNew, err := s.getOrCreateUser(parsed.Nick, parsed.Email, parsed.Link, userCache)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.ImportCommentError{
				Index:   i,
				Content: truncate(artalk.Content, 50),
				Error:   fmt.Sprintf("用户处理失败: %v", err),
			})
			continue
		}
		if isNew {
			result.UserCreated++
		}

		// 处理父评论
		var parentID *uint
		if artalk.RID != "" && artalk.RID != "0" {
			if pid, exists := commentIDMap[artalk.RID]; exists {
				parentID = &pid
			} else {
				result.Failed++
				result.Errors = append(result.Errors, dto.ImportCommentError{
					Index:   i,
					Content: truncate(artalk.Content, 50),
					Error:   fmt.Sprintf("父评论 %s 不存在", artalk.RID),
				})
				continue
			}
		}

		// 创建评论
		comment := &model.Comment{
			Content:    parsed.Content,
			TargetType: parsed.TargetType,
			TargetKey:  parsed.TargetKey,
			UserID:     user.ID,
			ParentID:   parentID,
			Status:     1,
			IP:         parsed.IP,
			Location:   utils.GetIPLocation(parsed.IP),
			Browser:    parsed.Browser,
			OS:         parsed.OS,
		}

		// 设置时间（gorm.Model的字段）
		comment.CreatedAt = parsed.CreatedAt
		comment.UpdatedAt = parsed.UpdatedAt

		// 处理回复关系
		if parentID != nil {
			parent, _ := s.repo.GetForWeb(ctx, *parentID)
			if parent != nil {
				if parent.RootID != nil {
					comment.RootID = parent.RootID
				} else {
					comment.RootID = &parent.ID
				}
				comment.ReplyTo = &parent.UserID
			}
		}

		// 保存
		if err := s.repo.Create(ctx, comment); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.ImportCommentError{
				Index:   i,
				Content: truncate(artalk.Content, 50),
				Error:   fmt.Sprintf("保存失败: %v", err),
			})
			continue
		}

		commentIDMap[artalk.ID] = comment.ID
		result.Success++
	}
}

func (s *CommentService) parseArtalkComment(artalk *dto.ArtalkCommentData) (*parsedComment, error) {
	// 转换HTML内容为Markdown
	markdownContent := convertHTMLToMarkdown(strings.TrimSpace(artalk.Content))

	parsed := &parsedComment{
		Content: markdownContent,
		Nick:    strings.TrimSpace(artalk.Nick),
		Email:   strings.TrimSpace(artalk.Email),
		Link:    strings.TrimSpace(artalk.Link),
		IP:      artalk.IP,
	}

	if parsed.Content == "" {
		return nil, errors.New("评论内容为空")
	}
	if parsed.Email == "" {
		return nil, errors.New("邮箱为空")
	}
	if parsed.Nick == "" {
		parsed.Nick = "匿名用户"
	}

	// 解析时间
	if t, err := parseTime(artalk.CreatedAt); err == nil {
		parsed.CreatedAt = t
	} else {
		parsed.CreatedAt = time.Now()
	}
	if t, err := parseTime(artalk.UpdatedAt); err == nil {
		parsed.UpdatedAt = t
	} else {
		parsed.UpdatedAt = parsed.CreatedAt
	}

	// 解析User-Agent
	parsed.Browser, parsed.OS = utils.ParseUserAgent(artalk.UA)

	// 解析页面信息
	pageKey := artalk.PageKey

	// 判断类型
	if isArticlePath(pageKey) {
		parsed.TargetType = "article"
		parsed.TargetKey = extractArticleSlug(pageKey)
	} else {
		parsed.TargetType = "page"
		parsed.TargetKey = normalizePageKey(pageKey)
	}

	return parsed, nil
}

func (s *CommentService) getOrCreateUser(nickname, email, website string, cache map[string]*model.User) (*model.User, bool, error) {
	if user, exists := cache[email]; exists {
		return user, false, nil
	}

	user, err := s.userRepo.GetGuestByEmail(email)
	if err == nil {
		cache[email] = user
		return user, false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 检查是否被注册用户使用，如果是则复用已有用户
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil {
		// 找到已注册用户，直接复用
		cache[email] = existingUser
		return existingUser, false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 邮箱不存在，创建新游客用户
	newUser := &model.User{
		Email:    email,
		Nickname: nickname,
		Website:  website,
		Role:     model.RoleGuest,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, false, err
	}

	cache[email] = newUser
	return newUser, true, nil
}

// parsedComment 解析后的评论数据
type parsedComment struct {
	Content    string
	TargetType string
	TargetKey  string
	Nick       string
	Email      string
	Link       string
	IP         string
	Browser    string
	OS         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// 辅助函数
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// convertHTMLToMarkdown 将HTML内容转换为Markdown格式
func convertHTMLToMarkdown(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// 解码HTML实体
	content := htmlContent
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&quot;", "\"")
	content = strings.ReplaceAll(content, "&#39;", "'")

	// Unicode转义解码
	content = strings.ReplaceAll(content, "\\u003c", "<")
	content = strings.ReplaceAll(content, "\\u003e", ">")

	// 处理代码块 <pre><code class="language-xxx">...</code></pre>
	content = regexp.MustCompile(`<pre><code[^>]*class="language-([^"]*)"[^>]*>([\s\S]*?)</code></pre>`).ReplaceAllStringFunc(content, func(match string) string {
		re := regexp.MustCompile(`<pre><code[^>]*class="language-([^"]*)"[^>]*>([\s\S]*?)</code></pre>`)
		matches := re.FindStringSubmatch(match)
		if len(matches) >= 3 {
			language := matches[1]
			codeContent := matches[2]
			return fmt.Sprintf("```%s\n%s\n```", language, strings.TrimSpace(codeContent))
		}
		return match
	})

	// 处理无语言的代码块 <pre><code>...</code></pre>
	content = regexp.MustCompile(`<pre><code[^>]*>([\s\S]*?)</code></pre>`).ReplaceAllStringFunc(content, func(match string) string {
		re := regexp.MustCompile(`<pre><code[^>]*>([\s\S]*?)</code></pre>`)
		matches := re.FindStringSubmatch(match)
		if len(matches) >= 2 {
			codeContent := matches[1]
			return fmt.Sprintf("```\n%s\n```", strings.TrimSpace(codeContent))
		}
		return match
	})

	// HTML标签转换为Markdown
	replacements := []struct{ from, to string }{
		// 段落标签
		{"<p>", ""},
		{"</p>", "\n\n"},

		// 换行标签
		{"<br>", "\n"},
		{"<br/>", "\n"},
		{"<br />", "\n"},

		// 强调标签
		{"<strong>", "**"},
		{"</strong>", "**"},
		{"<b>", "**"},
		{"</b>", "**"},

		{"<em>", "*"},
		{"</em>", "*"},
		{"<i>", "*"},
		{"</i>", "*"},

		// 行内代码标签（注意：代码块已在上面处理）
		{"<code>", "`"},
		{"</code>", "`"},

		// 删除线
		{"<del>", "~~"},
		{"</del>", "~~"},

		// 引用
		{"<blockquote>", "> "},
		{"</blockquote>", "\n\n"},
	}

	for _, r := range replacements {
		content = strings.ReplaceAll(content, r.from, r.to)
	}

	// 处理链接 <a href="url">text</a> -> [text](url)
	content = regexp.MustCompile(`<a[^>]+href="([^"]*)"[^>]*>([^<]*)</a>`).ReplaceAllString(content, "[$2]($1)")

	// 处理图片 <img src="url" alt="alt"> -> ![alt](url) 或 ![](url)
	content = regexp.MustCompile(`<img[^>]+src="([^"]*)"[^>]*alt="([^"]*)"[^>]*/?>`).ReplaceAllString(content, "![$2]($1)")
	content = regexp.MustCompile(`<img[^>]+src="([^"]*)"[^>]*/?>`).ReplaceAllString(content, "![]($1)")

	// 清理多余的空行
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	return strings.TrimSpace(content)
}

func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("无法解析时间")
}

func isArticlePath(path string) bool {
	patterns := []string{`^/posts?/\d+`, `^/articles?/\d+`, `^/posts?/[^/]+/?$`}
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, path); matched {
			return true
		}
	}
	return false
}

// extractArticleSlug 从 URL路径中提取文章 slug
// 支持格式：/posts/slug-name、/articles/slug-name、/post/123（数字ID也当作 slug）
func extractArticleSlug(path string) string {
	// 优先匹配 /posts/xxx 格式（支持 slug 和数字 ID）
	if re := regexp.MustCompile(`/posts?/([^/]+)/?$`); re.MatchString(path) {
		matches := re.FindStringSubmatch(path)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	// 匹配 /articles/xxx 格式
	if re := regexp.MustCompile(`/articles?/([^/]+)/?$`); re.MatchString(path) {
		matches := re.FindStringSubmatch(path)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func normalizePageKey(path string) string {
	path = strings.TrimSuffix(path, "/")
	if path == "" || path == "/" {
		return "home"
	}
	switch path {
	case "/friend", "/friends", "/links":
		return "friend"
	case "/moment", "/moments", "/essay", "/essays":
		return "moment"
	default:
		return strings.TrimPrefix(path, "/")
	}
}

// ============ 飞书回复评论 ============

// ReplyCommentFromFeishu 从飞书回复评论
func (s *CommentService) ReplyCommentFromFeishu(ctx context.Context, commentID uint, content, openID string) error {
	// 通过 OpenID 获取用户
	user, err := s.userRepo.GetByOAuthID("feishu", openID)
	if err != nil {
		return fmt.Errorf("未绑定飞书账号")
	}

	// 获取原评论
	comment, err := s.repo.Get(ctx, commentID)
	if err != nil {
		return fmt.Errorf("评论不存在")
	}

	// 创建回复
	req := &dto.CreateCommentRequest{
		Content:    content,
		TargetType: comment.TargetType,
		TargetKey:  comment.TargetKey,
		ParentID:   &commentID,
	}

	_, err = s.Create(ctx, req, user.ID)
	return err
}
