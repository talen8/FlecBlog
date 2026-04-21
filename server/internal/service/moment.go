package service

import (
	"context"
	"encoding/json"
	"errors"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/utils"
	"flec_blog/pkg/videoparser"

	"gorm.io/gorm"
)

// MomentService 动态服务
type MomentService struct {
	repo        *repository.MomentRepository
	fileService *FileService
}

// NewMomentService 创建动态服务实例
func NewMomentService(repo *repository.MomentRepository, fileService *FileService) *MomentService {
	return &MomentService{
		repo:        repo,
		fileService: fileService,
	}
}

// ============ 前台服务 ============

// ListForWeb 获取前台动态列表
func (s *MomentService) ListForWeb(ctx context.Context, page, pageSize int) ([]dto.MomentForWebResponse, int64, error) {
	// 前台只显示已发布的动态
	isPublish := true
	moments, total, err := s.repo.List(ctx, page, pageSize, "", "", "", &isPublish, nil, nil, nil, nil, "", "")
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.MomentForWebResponse, 0, len(moments))
	for _, moment := range moments {
		var content dto.MomentContent
		if moment.Content != "" {
			_ = json.Unmarshal([]byte(moment.Content), &content)
		}

		momentResp := dto.MomentForWebResponse{
			ID:          moment.ID,
			Content:     content,
			IsPublish:   moment.IsPublish,
			PublishTime: utils.ToJSONTime(moment.PublishTime),
		}

		result = append(result, momentResp)
	}

	return result, total, nil
}

// ============ 后台管理服务 ============

// List 获取动态列表（管理）
func (s *MomentService) List(
	ctx context.Context,
	req dto.ListMomentsRequest,
) ([]dto.MomentListResponse, int64, error) {
	moments, total, err := s.repo.List(
		ctx,
		req.Page,
		req.PageSize,
		req.Keyword,
		req.Tags,
		req.Location,
		req.IsPublish,
		req.HasImages,
		req.HasVideo,
		req.HasMusic,
		req.HasLink,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.MomentListResponse, 0, len(moments))
	for _, moment := range moments {
		var content dto.MomentContent
		if moment.Content != "" {
			_ = json.Unmarshal([]byte(moment.Content), &content)
		}

		momentResp := dto.MomentListResponse{
			ID:          moment.ID,
			Content:     content,
			IsPublish:   moment.IsPublish,
			PublishTime: utils.ToJSONTime(moment.PublishTime),
		}

		result = append(result, momentResp)
	}

	return result, total, nil
}

// Get 获取动态详情（管理）
func (s *MomentService) Get(ctx context.Context, id uint) (*dto.MomentListResponse, error) {
	moment, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在")
		}
		return nil, err
	}

	var content dto.MomentContent
	if moment.Content != "" {
		_ = json.Unmarshal([]byte(moment.Content), &content)
	}

	result := &dto.MomentListResponse{
		ID:          moment.ID,
		Content:     content,
		IsPublish:   moment.IsPublish,
		PublishTime: utils.ToJSONTime(moment.PublishTime),
	}

	return result, nil
}

// Create 创建动态
func (s *MomentService) Create(ctx context.Context, req *dto.CreateMomentRequest) (*dto.MomentListResponse, error) {
	// 如果有视频且未提供platform和video_id，自动识别
	if req.Content.Video != nil && req.Content.Video.URL != "" {
		// 只在前端没有提供解析结果时才进行解析（避免重复处理）
		if req.Content.Video.Platform == "" || req.Content.Video.VideoID == "" {
			videoInfo := videoparser.ParseVideoURL(req.Content.Video.URL)
			if videoInfo != nil {
				req.Content.Video.Platform = videoInfo.Platform
				req.Content.Video.VideoID = videoInfo.VideoID
			}
		}
	}

	// 将内容转换为JSON字符串
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return nil, err
	}

	moment := &model.Moment{
		Content:     string(contentBytes),
		IsPublish:   req.IsPublish,
		PublishTime: utils.FromJSONTime(req.PublishTime),
	}

	if err := s.repo.Create(ctx, moment); err != nil {
		return nil, err
	}

	// 标记文件为使用中
	s.markFilesAsUsed(&req.Content)

	return s.Get(ctx, moment.ID)
}

// Update 更新动态
func (s *MomentService) Update(ctx context.Context, id uint, req *dto.UpdateMomentRequest) (*dto.MomentListResponse, error) {
	moment, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在")
		}
		return nil, err
	}

	// 如果有视频且未提供platform和video_id，自动识别
	if req.Content.Video != nil && req.Content.Video.URL != "" {
		// 只在前端没有提供解析结果时才进行解析（避免重复处理）
		if req.Content.Video.Platform == "" || req.Content.Video.VideoID == "" {
			videoInfo := videoparser.ParseVideoURL(req.Content.Video.URL)
			if videoInfo != nil {
				req.Content.Video.Platform = videoInfo.Platform
				req.Content.Video.VideoID = videoInfo.VideoID
			}
		}
	}

	// 将内容转换为JSON字符串
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return nil, err
	}

	// 获取旧内容，用于对比文件变化
	var oldContent dto.MomentContent
	if moment.Content != "" {
		_ = json.Unmarshal([]byte(moment.Content), &oldContent)
	}

	moment.Content = string(contentBytes)
	moment.IsPublish = req.IsPublish
	moment.PublishTime = utils.FromJSONTime(req.PublishTime)

	if err := s.repo.Update(ctx, moment); err != nil {
		return nil, err
	}

	// 更新文件使用状态
	s.updateFileStatus(&oldContent, &req.Content)

	return s.Get(ctx, id)
}

// Delete 删除动态
func (s *MomentService) Delete(ctx context.Context, id uint) error {
	moment, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("动态不存在")
		}
		return err
	}

	// 获取动态内容，用于标记文件为未使用
	var content dto.MomentContent
	if moment.Content != "" {
		_ = json.Unmarshal([]byte(moment.Content), &content)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 标记文件为未使用
	s.markFilesAsUnused(&content)

	return nil
}

// ============ 文件状态管理 ============

// markFilesAsUsed 标记内容中的所有文件为使用中
func (s *MomentService) markFilesAsUsed(content *dto.MomentContent) {
	if s.fileService == nil {
		return
	}

	// 标记图片
	for _, imgURL := range content.Images {
		_ = s.fileService.MarkAsUsed(imgURL)
	}

	// 标记视频（如果是本地视频）
	if content.Video != nil && content.Video.URL != "" && content.Video.Platform == "" {
		_ = s.fileService.MarkAsUsed(content.Video.URL)
	}
}

// markFilesAsUnused 标记内容中的所有文件为未使用
func (s *MomentService) markFilesAsUnused(content *dto.MomentContent) {
	if s.fileService == nil {
		return
	}

	// 标记图片
	for _, imgURL := range content.Images {
		_ = s.fileService.MarkAsUnused(imgURL)
	}

	// 标记视频（如果是本地视频）
	if content.Video != nil && content.Video.URL != "" && content.Video.Platform == "" {
		_ = s.fileService.MarkAsUnused(content.Video.URL)
	}
}

// updateFileStatus 更新文件使用状态（对比新旧内容）
func (s *MomentService) updateFileStatus(oldContent, newContent *dto.MomentContent) {
	if s.fileService == nil {
		return
	}

	// 对比图片变化
	oldImages := make(map[string]bool)
	for _, img := range oldContent.Images {
		oldImages[img] = true
	}

	newImages := make(map[string]bool)
	for _, img := range newContent.Images {
		newImages[img] = true
		// 新增的图片标记为使用中
		if !oldImages[img] {
			_ = s.fileService.MarkAsUsed(img)
		}
	}

	// 移除的图片标记为未使用
	for _, img := range oldContent.Images {
		if !newImages[img] {
			_ = s.fileService.MarkAsUnused(img)
		}
	}

	// 对比视频变化（仅本地视频）
	oldVideoURL := ""
	if oldContent.Video != nil && oldContent.Video.Platform == "" {
		oldVideoURL = oldContent.Video.URL
	}

	newVideoURL := ""
	if newContent.Video != nil && newContent.Video.Platform == "" {
		newVideoURL = newContent.Video.URL
	}

	if oldVideoURL != newVideoURL {
		if oldVideoURL != "" {
			_ = s.fileService.MarkAsUnused(oldVideoURL)
		}
		if newVideoURL != "" {
			_ = s.fileService.MarkAsUsed(newVideoURL)
		}
	}
}
