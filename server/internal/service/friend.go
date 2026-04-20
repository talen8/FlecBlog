package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/utils"
)

// FriendService 友链服务
type FriendService struct {
	repo                *repository.FriendRepository
	fileService         *FileService
	notificationService *NotificationService
}

// NewFriendService 创建友链服务实例
func NewFriendService(repo *repository.FriendRepository, fileService *FileService, notificationService *NotificationService) *FriendService {
	return &FriendService{
		repo:                repo,
		fileService:         fileService,
		notificationService: notificationService,
	}
}

// ============ 友链类型 ============

// ListTypes 获取所有友链类型（后台管理用）
func (s *FriendService) ListTypes(ctx context.Context, page, pageSize int) ([]dto.FriendTypeListResponse, int64, error) {
	types, total, err := s.repo.ListTypes(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.FriendTypeListResponse, 0, len(types))
	for _, t := range types {
		// 统计该类型下的友链数量
		count, err := s.repo.CountFriendsByTypeID(ctx, t.ID)
		if err != nil {
			count = 0
		}

		result = append(result, dto.FriendTypeListResponse{
			ID:        t.ID,
			Name:      t.Name,
			Sort:      t.Sort,
			IsVisible: t.IsVisible,
			Count:     int(count),
		})
	}
	return result, total, nil
}

// GetType 获取友链类型详情
func (s *FriendService) GetType(ctx context.Context, id uint) (*model.FriendType, error) {
	return s.repo.GetTypeByID(ctx, id)
}

// CreateType 创建友链类型
func (s *FriendService) CreateType(ctx context.Context, req *dto.CreateFriendTypeRequest) error {
	friendType := &model.FriendType{
		Name:      req.Name,
		Sort:      req.Sort,
		IsVisible: true, // 默认展示
	}

	if req.IsVisible != nil {
		friendType.IsVisible = *req.IsVisible
	}

	return s.repo.CreateType(ctx, friendType)
}

// UpdateType 更新友链类型
func (s *FriendService) UpdateType(ctx context.Context, id uint, req *dto.UpdateFriendTypeRequest) error {
	// 检查类型是否存在
	existingType, err := s.repo.GetTypeByID(ctx, id)
	if err != nil {
		return errors.New("友链类型不存在")
	}

	// 更新字段
	existingType.Name = req.Name
	existingType.Sort = req.Sort

	if req.IsVisible != nil {
		existingType.IsVisible = *req.IsVisible
	}

	return s.repo.UpdateType(ctx, existingType)
}

// DeleteType 删除友链类型
func (s *FriendService) DeleteType(ctx context.Context, id uint) error {
	// 检查类型是否存在
	_, err := s.repo.GetTypeByID(ctx, id)
	if err != nil {
		return errors.New("友链类型不存在")
	}

	// 注意：由于外键设置了 ON DELETE SET NULL，删除类型时，关联的友链的 type_id 会被设置为 NULL
	return s.repo.DeleteType(ctx, id)
}

// ============ 前台服务 ============

// ListForWeb 获取友链分组列表（包括失效友链，前端可根据 is_invalid 字段处理展示）
func (s *FriendService) ListForWeb(ctx context.Context) (*dto.GroupedFriendsResponse, error) {
	// 1. 获取所有数据
	types, allFriends, err := s.repo.GetFriendsForWeb(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 在内存中按类型分组（使用 0 表示未分类）
	const uncategorizedKey uint = 0
	typeMap := make(map[uint][]model.Friend)

	for _, friend := range allFriends {
		key := uncategorizedKey
		if friend.TypeID != nil {
			key = *friend.TypeID
		}
		typeMap[key] = append(typeMap[key], friend)
	}

	// 3. 组装响应数据（空分组不返回）
	groups := make([]dto.FriendGroupResponse, 0, len(types)+1)
	totalFriends := 0

	for _, t := range types {
		if friends := typeMap[t.ID]; len(friends) > 0 {
			groups = append(groups, buildGroupResponse(&t.ID, t.Name, t.Sort, friends))
			totalFriends += len(friends)
		}
	}

	// 添加未分类友链
	if friends := typeMap[uncategorizedKey]; len(friends) > 0 {
		groups = append(groups, buildGroupResponse(nil, "友情链接", 0, friends))
		totalFriends += len(friends)
	}

	return &dto.GroupedFriendsResponse{
		Groups:       groups,
		TotalGroups:  len(groups),
		TotalFriends: totalFriends,
	}, nil
}

// buildGroupResponse 构建分组响应（辅助函数）
func buildGroupResponse(typeID *uint, name string, sort int, friends []model.Friend) dto.FriendGroupResponse {
	return dto.FriendGroupResponse{
		TypeID:   typeID,
		TypeName: name,
		TypeSort: sort,
		Friends:  convertToFriendInGroupResponse(friends),
	}
}

// convertToFriendInGroupResponse 转换为分组友链响应格式（辅助函数）
func convertToFriendInGroupResponse(friends []model.Friend) []dto.FriendInGroupResponse {
	result := make([]dto.FriendInGroupResponse, len(friends))
	for i, f := range friends {
		result[i] = dto.FriendInGroupResponse{
			ID:          f.ID,
			Name:        f.Name,
			URL:         f.URL,
			Description: f.Description,
			Avatar:      f.Avatar,
			Screenshot:  f.Screenshot,
			Sort:        f.Sort,
			IsInvalid:   f.IsInvalid,
		}
	}
	return result
}

// ============ 后台管理服务 ============

// List 获取友链列表
func (s *FriendService) List(ctx context.Context, req *dto.ListFriendRequest) ([]dto.FriendListResponse, int64, error) {
	friends, total, err := s.repo.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.FriendListResponse, 0, len(friends))
	for _, friend := range friends {
		response := dto.FriendListResponse{
			ID:          friend.ID,
			Name:        friend.Name,
			URL:         friend.URL,
			Description: friend.Description,
			Avatar:      friend.Avatar,
			Screenshot:  friend.Screenshot,
			Sort:        friend.Sort,
			IsInvalid:   friend.IsInvalid,
			IsPending:   friend.IsPending,
			TypeID:      friend.TypeID,
			RSSUrl:      friend.RSSUrl,
			Accessible:  friend.Accessible,
		}

		// 如果有类型，返回类型名称（后台展示用）
		if friend.Type != nil {
			response.TypeName = friend.Type.Name
		}

		// 如果有RSS最后更新时间，格式化返回
		response.RSSLatime = utils.ToJSONTime(friend.RSSLatime)

		result = append(result, response)
	}

	return result, total, nil
}

// Get 获取友链详情
func (s *FriendService) Get(ctx context.Context, id uint) (*model.Friend, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建友链
func (s *FriendService) Create(ctx context.Context, req *dto.CreateFriendRequest) (*model.Friend, error) {
	if req.Sort == 0 {
		req.Sort = 5
	}

	// 如果指定了类型，验证类型是否存在
	if req.TypeID != nil {
		_, err := s.repo.GetTypeByID(ctx, *req.TypeID)
		if err != nil {
			return nil, errors.New("指定的友链类型不存在")
		}
	}

	friend := &model.Friend{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Avatar:      req.Avatar,
		Screenshot:  req.Screenshot,
		Sort:        req.Sort,
		IsInvalid:   false, // 新创建的友链默认不失效
		TypeID:      req.TypeID,
		RSSUrl:      req.RSSUrl,
	}

	if err := s.repo.Create(ctx, friend); err != nil {
		return nil, err
	}

	// 标记文件为使用中
	if s.fileService != nil {
		if req.Avatar != "" {
			_ = s.fileService.MarkAsUsed(req.Avatar)
		}
		if req.Screenshot != "" {
			_ = s.fileService.MarkAsUsed(req.Screenshot)
		}
	}

	return friend, nil
}

// Update 更新友链
func (s *FriendService) Update(ctx context.Context, id uint, req *dto.UpdateFriendRequest) error {
	// 检查友链是否存在
	existingFriend, err := s.repo.Get(ctx, id)
	if err != nil {
		return errors.New("友链不存在")
	}

	// 如果指定了类型，验证类型是否存在
	if req.TypeID != nil {
		_, err := s.repo.GetTypeByID(ctx, *req.TypeID)
		if err != nil {
			return errors.New("指定的友链类型不存在")
		}
	}

	// 保存旧文件
	oldAvatar := existingFriend.Avatar
	oldScreenshot := existingFriend.Screenshot

	// 更新字段
	existingFriend.Name = req.Name
	existingFriend.URL = req.URL
	existingFriend.Description = req.Description
	existingFriend.Avatar = req.Avatar
	existingFriend.Screenshot = req.Screenshot
	existingFriend.Sort = req.Sort
	existingFriend.TypeID = req.TypeID
	existingFriend.RSSUrl = req.RSSUrl
	existingFriend.UpdatedAt = time.Now()

	// 如果请求中指定了失效状态
	if req.IsInvalid != nil {
		existingFriend.IsInvalid = *req.IsInvalid
	}

	// 如果请求中指定了待审核状态
	if req.IsPending != nil {
		existingFriend.IsPending = *req.IsPending
	}

	// 如果请求中指定了可访问性状态
	if req.Accessible != nil {
		existingFriend.Accessible = *req.Accessible
	}

	if err := s.repo.Update(ctx, existingFriend); err != nil {
		return err
	}

	// 处理文件变化
	if s.fileService != nil {
		if oldAvatar != req.Avatar {
			if oldAvatar != "" {
				_ = s.fileService.MarkAsUnused(oldAvatar)
			}
			if req.Avatar != "" {
				_ = s.fileService.MarkAsUsed(req.Avatar)
			}
		}
		if oldScreenshot != req.Screenshot {
			if oldScreenshot != "" {
				_ = s.fileService.MarkAsUnused(oldScreenshot)
			}
			if req.Screenshot != "" {
				_ = s.fileService.MarkAsUsed(req.Screenshot)
			}
		}
	}

	return nil
}

// Delete 删除友链
func (s *FriendService) Delete(ctx context.Context, id uint) error {
	// 检查友链是否存在
	friend, err := s.repo.Get(ctx, id)
	if err != nil {
		return errors.New("友链不存在")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 标记文件为未使用
	if s.fileService != nil {
		if friend.Avatar != "" {
			_ = s.fileService.MarkAsUnused(friend.Avatar)
		}
		if friend.Screenshot != "" {
			_ = s.fileService.MarkAsUnused(friend.Screenshot)
		}
	}

	return nil
}

// ============ 前台友链申请 ============

// ApplyFriend 申请友链（前台用户提交申请）
func (s *FriendService) ApplyFriend(ctx context.Context, req *dto.ApplyFriendRequest, userID uint) error {
	// 创建待审核的友链记录
	friend := &model.Friend{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Avatar:      req.Avatar,
		Screenshot:  req.Screenshot,
		IsPending:   true, // 标记为待审核
	}

	if err := s.repo.Create(ctx, friend); err != nil {
		return err
	}

	if s.fileService != nil {
		if req.Avatar != "" {
			_ = s.fileService.MarkAsUsed(req.Avatar)
		}
		if req.Screenshot != "" {
			_ = s.fileService.MarkAsUsed(req.Screenshot)
		}
	}

	// 发送友链申请通知给管理员
	if s.notificationService != nil {
		_ = s.notificationService.NotifyFriendApply(ctx, friend.ID, req.Name, req.URL, req.Description, req.Avatar, req.Screenshot, &userID)
	}

	return nil
}

// ApproveFriend 通过友链申请（将待审核状态改为已通过）
func (s *FriendService) ApproveFriend(ctx context.Context, id uint) error {
	// 检查友链是否存在
	friend, err := s.repo.Get(ctx, id)
	if err != nil {
		return errors.New("友链不存在")
	}

	// 将待审核状态改为 false
	friend.IsPending = false
	friend.UpdatedAt = time.Now()

	return s.repo.Update(ctx, friend)
}

// ============ 友链检测 ============

// CheckAllFriends 检测所有友链（定时任务调用）
func (s *FriendService) CheckAllFriends() error {
	ctx := context.Background()

	friends, err := s.repo.GetAllForCheck(ctx)
	if err != nil {
		return err
	}

	for _, friend := range friends {
		newStatus := 0
		if !s.checkAccessibility(ctx, friend.URL) {
			newStatus = friend.Accessible + 1
			if s.notificationService != nil {
				_ = s.notificationService.NotifyFriendAbnormal(ctx, friend.ID, friend.Name, newStatus)
			}
		}
		_ = s.repo.UpdateCheckStatus(ctx, friend.ID, newStatus)
	}

	return nil
}

// checkAccessibility 检测网站可访问性
func (s *FriendService) checkAccessibility(ctx context.Context, targetURL string) bool {
	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
