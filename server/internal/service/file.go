package service

import (
	"fmt"
	"io"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/logger"
	"flec_blog/pkg/upload"
	"flec_blog/pkg/utils"
)

var reconciledSettingImageKeys = []string{
	KeyBasicAuthorAvatar,
	KeyBasicAuthorPhoto,
	KeyBlogFavicon,
	KeyBlogBackgroundImage,
	KeyBlogAboutExhibition,
	KeyBlogScreenshot,
}

// FileUsageChecker 文件引用检查器
type FileUsageChecker struct {
	articleRepo  *repository.ArticleRepository
	friendRepo   *repository.FriendRepository
	momentRepo   *repository.MomentRepository
	settingRepo  *repository.SettingRepository
	userRepo     *repository.UserRepository
	menuRepo     *repository.MenuRepository
	feedbackRepo *repository.FeedbackRepository
	commentRepo  *repository.CommentRepository
}

// NewFileUsageChecker 创建文件引用检查器
func NewFileUsageChecker(
	articleRepo *repository.ArticleRepository,
	friendRepo *repository.FriendRepository,
	momentRepo *repository.MomentRepository,
	settingRepo *repository.SettingRepository,
	userRepo *repository.UserRepository,
	menuRepo *repository.MenuRepository,
	feedbackRepo *repository.FeedbackRepository,
	commentRepo *repository.CommentRepository,
) *FileUsageChecker {
	return &FileUsageChecker{
		articleRepo:  articleRepo,
		friendRepo:   friendRepo,
		momentRepo:   momentRepo,
		settingRepo:  settingRepo,
		userRepo:     userRepo,
		menuRepo:     menuRepo,
		feedbackRepo: feedbackRepo,
		commentRepo:  commentRepo,
	}
}

// IsActuallyUsed 检查文件是否仍被业务引用
func (c *FileUsageChecker) IsActuallyUsed(fileURL string) (bool, string, error) {
	checks := []struct {
		name string
		fn   func(string) (bool, error)
	}{
		{name: "文章封面", fn: c.articleRepo.ExistsByCover},
		{name: "文章正文", fn: c.articleRepo.ExistsByContentURL},
		{name: "友链图片", fn: c.friendRepo.ExistsByAvatarOrScreenshot},
		{name: "动态内容", fn: c.momentRepo.ExistsByContentURL},
		{name: "设置图片", fn: func(url string) (bool, error) {
			return c.settingRepo.ExistsByValueAndKeys(url, reconciledSettingImageKeys)
		}},
		{name: "用户头像", fn: c.userRepo.ExistsByAvatar},
		{name: "菜单图标", fn: c.menuRepo.ExistsByIcon},
		{name: "反馈附件", fn: c.feedbackRepo.ExistsByAttachmentURL},
		{name: "评论内容", fn: c.commentRepo.ExistsByContentURL},
	}

	for _, check := range checks {
		used, err := check.fn(fileURL)
		if err != nil {
			return false, "", err
		}
		if used {
			return true, check.name, nil
		}
	}

	return false, "", nil
}

// FileService 文件服务
type FileService struct {
	fileRepo      *repository.FileRepository
	uploadManager *upload.Manager
	usageChecker  *FileUsageChecker
}

// NewFileService 创建文件服务
func NewFileService(fileRepo *repository.FileRepository, uploadManager *upload.Manager) *FileService {
	return &FileService{
		fileRepo:      fileRepo,
		uploadManager: uploadManager,
	}
}

// SetUsageChecker 设置文件引用检查器
func (s *FileService) SetUsageChecker(checker *FileUsageChecker) {
	s.usageChecker = checker
}

// ============ 通用服务 ============

// UploadFromReader 从Reader上传文件
func (s *FileService) UploadFromReader(reader io.Reader, originalName, fileType string, uploadType upload.Type, userID uint, host string) (string, error) {
	// 读取文件数据并计算hash
	data, fileHash, err := s.uploadManager.HandleUploadFromReader(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 保存文件数据（使用相对路径）
	fileInfo, err := s.uploadManager.SaveFileData(data, fileHash, originalName, fileType, uploadType, userID, host)
	if err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建数据库记录
	file := s.createFileFromUploadInfo(fileInfo)
	file.Status = 0 // 默认未使用

	if err := s.fileRepo.Create(file); err != nil {
		_ = s.uploadManager.DeleteFile(fileInfo.FilePath)
		return "", fmt.Errorf("保存记录失败: %w", err)
	}

	return file.FileURL, nil
}

// MarkAsUsed 标记文件为使用中
func (s *FileService) MarkAsUsed(fileUrl string) error {
	if fileUrl == "" {
		return nil
	}
	return s.fileRepo.UpdateStatus(fileUrl, 1)
}

// MarkAsUnused 标记文件为未使用
func (s *FileService) MarkAsUnused(fileUrl string) error {
	if fileUrl == "" {
		return nil
	}
	return s.fileRepo.UpdateStatus(fileUrl, 0)
}

// ============ 前台服务 ============

// UploadForWeb 前台文件上传
func (s *FileService) UploadForWeb(req *upload.Request, host string) (*dto.FileUploadForWebResponse, error) {
	// 验证上传类型
	if string(req.UploadType) == "" {
		return nil, fmt.Errorf("上传类型不能为空")
	}

	// 文件大小限制（从配置获取，单位MB）
	maxFileSizeMB := s.uploadManager.GetMaxFileSize()
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5 // 默认 5MB
	}
	maxWebFileSize := maxFileSizeMB * 1024 * 1024
	if req.File.Size > maxWebFileSize {
		return nil, fmt.Errorf("文件大小超出限制，前台上传最大允许 %dMB", maxFileSizeMB)
	}

	// 文件类型白名单验证（具体场景限制由前端控制）
	contentType := req.File.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg":         true,
		"image/jpg":          true,
		"image/png":          true,
		"image/gif":          true,
		"image/webp":         true,
		"image/avif":         true,
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	// 调用通用上传方法（传递 host）
	file, err := s.handleUpload(req, host)
	if err != nil {
		return nil, err
	}

	// 返回简化响应
	return &dto.FileUploadForWebResponse{
		OriginalName: file.OriginalName,
		FileURL:      file.FileURL,
	}, nil
}

// ============ 后台管理服务 ============

// Upload 文件上传
func (s *FileService) Upload(req *upload.Request, host string) (*dto.FileResponse, error) {
	// 验证上传类型
	if string(req.UploadType) == "" {
		return nil, fmt.Errorf("上传类型不能为空")
	}

	// 调用通用上传方法（传递 host）
	file, err := s.handleUpload(req, host)
	if err != nil {
		return nil, err
	}

	return &dto.FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileURL:      file.FileURL,
		UploadType:   upload.Type(file.UploadType),
		UserID:       file.UserID,
		Status:       file.Status,
		UploadTime:   utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// List 获取文件列表
func (s *FileService) List(req *dto.ListFilesRequest) ([]dto.FileResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	var files []model.File
	var total int64
	var err error

	// 根据类型过滤
	if req.Type != "" {
		files, total, err = s.fileRepo.GetByUploadType(req.Type, offset, req.PageSize)
	} else {
		files, total, err = s.fileRepo.List(offset, req.PageSize)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 转换为响应格式
	fileResponses := make([]dto.FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = dto.FileResponse{
			ID:           file.ID,
			OriginalName: file.OriginalName,
			FileName:     file.FileName,
			FileSize:     file.FileSize,
			FileType:     file.FileType,
			FileURL:      file.FileURL,
			UploadType:   upload.Type(file.UploadType),
			UserID:       file.UserID,
			Status:       file.Status,
			UploadTime:   utils.NewJSONTime(file.CreatedAt),
		}
	}

	return fileResponses, total, nil
}

// Get 获取文件详情
func (s *FileService) Get(id uint) (*dto.FileResponse, error) {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	return &dto.FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileURL:      file.FileURL,
		UploadType:   upload.Type(file.UploadType),
		UserID:       file.UserID,
		Status:       file.Status,
		UploadTime:   utils.NewJSONTime(file.CreatedAt),
	}, nil
}

// Delete 删除文件
func (s *FileService) Delete(id uint) error {
	file, err := s.fileRepo.Get(id)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	// 删除物理文件
	if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
		fmt.Printf("删除文件失败，文件路径: %s, 存储类型: %s, 错误: %v\n",
			file.FilePath, file.StorageType, err)
	}

	return nil
}

// ============ 辅助方法 ============

// handleUpload 处理文件上传
func (s *FileService) handleUpload(req *upload.Request, host string) (*model.File, error) {
	result, err := s.uploadManager.HandleUpload(req, host)
	if err != nil {
		return nil, fmt.Errorf("文件上传失败: %w", err)
	}

	if !result.Success || result.FileInfo == nil {
		return nil, fmt.Errorf("文件上传失败: %s", result.Message)
	}

	// 创建文件记录
	file := s.createFileFromUploadInfo(result.FileInfo)
	file.Status = 0 // 默认未使用

	if err := s.fileRepo.Create(file); err != nil {
		_ = s.uploadManager.DeleteFile(result.FileInfo.FilePath)
		return nil, fmt.Errorf("保存记录失败: %w", err)
	}

	return file, nil
}

// createFileFromUploadInfo 从上传信息创建文件模型
func (s *FileService) createFileFromUploadInfo(info *upload.FileInfo) *model.File {
	// 处理 UserID：0 表示匿名上传，转为 nil
	var userID *uint
	if info.UserID > 0 {
		userID = &info.UserID
	}

	return &model.File{
		FileName:     info.FileName,
		OriginalName: info.OriginalName,
		FilePath:     info.FilePath,
		FileSize:     info.FileSize,
		FileType:     info.FileType,
		UploadType:   string(info.UploadType),
		StorageType:  info.StorageType,
		UserID:       userID,
		FileURL:      info.FileURL,
	}
}

// ============ 定时任务方法 ============

// DeleteUnusedFiles 删除未使用文件，先纠正误标，再清理超过15天仍未使用的文件
func (s *FileService) DeleteUnusedFiles() error {
	if s.usageChecker == nil {
		return nil
	}

	files, err := s.fileRepo.GetByStatus(0)
	if err != nil {
		return fmt.Errorf("获取未使用文件失败: %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	deleteBefore := time.Now().AddDate(0, 0, -15)
	usedURLs := make([]string, 0)
	deletableFiles := make([]model.File, 0)

	for _, file := range files {
		used, source, err := s.usageChecker.IsActuallyUsed(file.FileURL)
		if err != nil {
			logger.Warn("检查文件引用失败 %s: %v", file.FileURL, err)
			continue
		}
		if used {
			usedURLs = append(usedURLs, file.FileURL)
			logger.Info("文件引用自检纠正成功 %s -> %s", file.FileURL, source)
			continue
		}

		if file.CreatedAt.Before(deleteBefore) {
			deletableFiles = append(deletableFiles, file)
		}
	}

	if err := s.fileRepo.UpdateFileStatusByUrls(usedURLs, 1); err != nil {
		return fmt.Errorf("批量纠正文件状态失败: %w", err)
	}

	deletedIDs := make([]uint, 0, len(deletableFiles))
	for _, file := range deletableFiles {
		if err := s.uploadManager.DeleteFileByStorageType(file.FilePath, file.StorageType); err != nil {
			logger.Warn("删除物理文件失败 %s: %v", file.FilePath, err)
			continue
		}
		deletedIDs = append(deletedIDs, file.ID)
	}

	if err := s.fileRepo.DeleteByIDs(deletedIDs); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	return nil
}
