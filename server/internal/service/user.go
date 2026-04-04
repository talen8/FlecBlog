package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/pkg/random"
	"flec_blog/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	repo        *repository.UserRepository
	fileService *FileService
	config      *config.Config
}

// NewUserService 创建用户服务实例
func NewUserService(repo *repository.UserRepository, fileService *FileService, cfg *config.Config) *UserService {
	return &UserService{
		repo:        repo,
		fileService: fileService,
		config:      cfg,
	}
}

// ============ 通用服务 ============

// Get 获取用户信息
func (s *UserService) Get(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(user), nil
}

// ValidateToken 验证token并返回用户信息
func (s *UserService) ValidateToken(token string) (*model.User, error) {
	claims, err := utils.ParseToken(token, &s.config.JWT)
	if err != nil {
		return nil, err
	}

	// 检查 token 是否在黑名单中
	tokenHash := hashToken(token)
	if s.repo.IsTokenBlacklisted(tokenHash) {
		return nil, errors.New("token已失效，请重新登录")
	}

	return s.repo.Get(claims.UserID)
}

// buildLoginResponse 构建登录响应（含token和用户信息）
func (s *UserService) buildLoginResponse(user *model.User) (*dto.LoginResponse, error) {
	// 生成access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role, &s.config.JWT)
	if err != nil {
		return nil, err
	}

	// 生成refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role, &s.config.JWT)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.NewUserResponse(user),
	}, nil
}

// ============ 前台服务 ============

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterRequest, host string) (*dto.LoginResponse, error) {
	// 检查邮箱是否存在
	existingUser, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		// 邮箱已存在，检查是否为游客用户
		if existingUser.Role == model.RoleGuest {
			// 游客账户升级为正式用户
			return s.upgradeGuest(existingUser, req, host)
		}
		// 已是正式用户，不能重复注册
		return nil, errors.New("邮箱已被注册")
	}

	// 邮箱不存在（或查询出错），继续检查
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	now := utils.Now().Time
	user := &model.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		HasPassword: true, // 账密注册的用户有密码
		Nickname:    req.Nickname,
		Website:     req.Website,
		Role:        model.RoleUser, // 普通用户角色
		LastLogin:   &now,           // 注册即登录，设置最后登录时间
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 异步下载Cravatar头像
	go func() {
		avatarURL, err := s.downloadAndSaveCravatarAvatar(req.Email, user.ID, host)
		if err == nil && avatarURL != "" {
			if err := s.repo.UpdateAvatar(user.ID, avatarURL); err == nil && s.fileService != nil {
				_ = s.fileService.MarkAsUsed(avatarURL)
			}
		}
	}()

	return s.buildLoginResponse(user)
}

// upgradeGuest 将游客账户升级为正式用户
func (s *UserService) upgradeGuest(guestUser *model.User, req *dto.RegisterRequest, host string) (*dto.LoginResponse, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 更新用户信息：从游客升级为正式用户
	now := utils.Now().Time
	guestUser.Password = string(hashedPassword) // 设置密码
	guestUser.HasPassword = true                // 标记已设置密码
	guestUser.Nickname = req.Nickname

	if req.Website != "" {
		guestUser.Website = req.Website
	}

	guestUser.Role = model.RoleUser // 升级为普通用户角色
	guestUser.LastLogin = &now

	if err := s.repo.Update(guestUser); err != nil {
		return nil, err
	}

	// 如果游客没有头像，异步下载Cravatar头像
	if guestUser.Avatar == "" {
		go func() {
			avatarURL, err := s.downloadAndSaveCravatarAvatar(req.Email, guestUser.ID, host)
			if err == nil && avatarURL != "" {
				if err := s.repo.UpdateAvatar(guestUser.ID, avatarURL); err == nil && s.fileService != nil {
					_ = s.fileService.MarkAsUsed(avatarURL)
				}
			}
		}()
	}

	return s.buildLoginResponse(guestUser)
}

// LoginBySocial 第三方登录逻辑
func (s *UserService) LoginBySocial(provider, providerID, email, nickname, avatarURL, host string) (*dto.LoginResponse, error) {
	// 1. 先通过 OAuth ID 查找用户
	user, err := s.repo.GetByOAuthID(provider, providerID)
	if err == nil {
		// OAuth ID 已绑定，直接登录
		now := utils.Now().Time
		user.LastLogin = &now
		_ = s.repo.Update(user)
		return s.buildLoginResponse(user)
	}

	// 2. OAuth ID 未绑定，通过邮箱查找
	user, err = s.repo.GetByEmail(email)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// 3. 邮箱也不存在 -> 自动注册新用户
		now := utils.Now().Time
		user = &model.User{
			Email:       email,
			Nickname:    nickname,
			HasPassword: false, // OAuth 注册的用户没有密码
			Role:        model.RoleUser,
			LastLogin:   &now,
			IsEnabled:   true,
		}

		// 设置对应的 OAuth ID
		switch provider {
		case "github":
			user.GithubID = providerID
		case "google":
			user.GoogleID = providerID
		case "qq":
			user.QQID = providerID
		case "microsoft":
			user.MicrosoftID = providerID
		}

		if err := s.repo.Create(user); err != nil {
			return nil, err
		}

		// 异步下载头像：优先使用第三方头像，否则使用 Cravatar
		go s.downloadSocialAvatar(user.ID, email, avatarURL, host)

	} else {
		// 4. 邮箱存在 -> 绑定 OAuth 并登录
		now := utils.Now().Time
		user.LastLogin = &now

		// 绑定 OAuth ID
		switch provider {
		case "github":
			user.GithubID = providerID
		case "google":
			user.GoogleID = providerID
		case "qq":
			user.QQID = providerID
		case "microsoft":
			user.MicrosoftID = providerID
		}

		// 如果用户没有头像，异步下载
		if user.Avatar == "" {
			go s.downloadSocialAvatar(user.ID, email, avatarURL, host)
		}

		_ = s.repo.Update(user)
	}

	// 5. 签发 Token
	return s.buildLoginResponse(user)
}

// downloadSocialAvatar 下载第三方登录头像（优先用第三方头像，否则用 Cravatar）
func (s *UserService) downloadSocialAvatar(userID uint, email, avatarURL, host string) {
	if s.fileService == nil {
		return
	}

	var savedURL string
	var err error

	if avatarURL != "" {
		// 下载第三方头像
		savedURL, err = s.downloadAndSaveRemoteAvatar(avatarURL, userID, host)
	}

	// 如果第三方头像下载失败或没有提供，使用 Cravatar
	if savedURL == "" || err != nil {
		savedURL, err = s.downloadAndSaveCravatarAvatar(email, userID, host)
	}

	if err == nil && savedURL != "" {
		if err := s.repo.UpdateAvatar(userID, savedURL); err == nil {
			_ = s.fileService.MarkAsUsed(savedURL)
		}
	}
}

// downloadAndSaveRemoteAvatar 下载并保存远程头像
func (s *UserService) downloadAndSaveRemoteAvatar(avatarURL string, userID uint, host string) (string, error) {
	if s.fileService == nil {
		return "", nil
	}

	reader, err := utils.DownloadRemoteImage(avatarURL)
	if err != nil {
		return "", err
	}

	filename := "social_avatar_" + random.Code(8) + ".webp"
	return s.fileService.UploadFromReader(
		reader,
		filename,
		"image/webp",
		"用户头像",
		userID,
		host,
	)
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱未注册")
		}
		return nil, err
	}

	// 禁止游客用户登录（游客没有密码）
	if user.Role == model.RoleGuest {
		return nil, errors.New("游客账户无法登录，请先注册成为正式用户")
	}

	// 检查用户状态
	if !user.IsEnabled {
		return nil, errors.New("该账号已被禁用，请联系管理员")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 更新最后登录时间
	now := utils.Now().Time
	user.LastLogin = &now
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.buildLoginResponse(user)
}

// hashToken 对token进行SHA256哈希
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// RefreshToken 刷新token
func (s *UserService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// 解析并验证refresh token
	claims, err := utils.ParseRefreshToken(req.RefreshToken, &s.config.JWT)
	if err != nil {
		return nil, errors.New("无效的refresh token")
	}

	// 检查token是否在黑名单中
	tokenHash := hashToken(req.RefreshToken)
	if s.repo.IsTokenBlacklisted(tokenHash) {
		return nil, errors.New("token已失效，请重新登录")
	}

	// 获取用户信息（验证用户是否存在且未被禁用）
	user, err := s.repo.Get(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.IsEnabled {
		return nil, errors.New("该账号已被禁用")
	}

	// 生成新的access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role, &s.config.JWT)
	if err != nil {
		return nil, err
	}

	// 生成新的refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Role, &s.config.JWT)
	if err != nil {
		return nil, err
	}

	// 将旧的refresh token加入黑名单（Refresh Token轮换）
	expiresAt := claims.ExpiresAt.Time
	if err = s.repo.AddTokenToBlacklist(tokenHash, user.ID, expiresAt); err != nil {
		// 记录错误但不中断流程（降级处理）
		_ = err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(token string) error {
	// 解析并验证 access token
	claims, err := utils.ParseToken(token, &s.config.JWT)
	if err != nil {
		return errors.New("无效的token")
	}

	// 将 token 加入黑名单
	tokenHash := hashToken(token)
	expiresAt := claims.ExpiresAt.Time
	err = s.repo.AddTokenToBlacklist(tokenHash, claims.UserID, expiresAt)
	if err != nil {
		return errors.New("登出失败")
	}

	return nil
}

// CleanupExpiredTokens 清理过期的黑名单token（定时任务调用）
func (s *UserService) CleanupExpiredTokens() error {
	return s.repo.CleanupExpiredTokens()
}

// RevokeAllUserTokens 撤销某用户的所有token（用于强制下线，如账号被盗、密码修改等场景）
func (s *UserService) RevokeAllUserTokens(userID uint) error {
	return s.repo.RevokeAllUserTokens(userID)
}

// UpdateForWeb 更新用户信息
func (s *UserService) UpdateForWeb(id uint, req *dto.UpdateUserRequest) error {
	user, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	// 检查邮箱冲突
	if req.Email != "" && req.Email != user.Email {
		if s.repo.ExistsByEmail(req.Email) {
			return errors.New("邮箱已存在")
		}
	}

	oldAvatar := user.Avatar

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Badge != "" {
		if err := s.validateBadge(req.Badge); err != nil {
			return err
		}
		user.Badge = req.Badge
	}
	if req.Website != "" {
		user.Website = req.Website
	}

	// 处理头像变化
	if s.fileService != nil && oldAvatar != user.Avatar {
		if oldAvatar != "" {
			_ = s.fileService.MarkAsUnused(oldAvatar)
		}
		if user.Avatar != "" {
			_ = s.fileService.MarkAsUsed(user.Avatar)
		}
	}

	return s.repo.Update(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.repo.Get(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedNewPassword)
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// 修改密码后撤销所有现有token，强制用户重新登录
	_ = s.repo.RevokeAllUserTokens(userID)

	return nil
}

// SetPassword 设置密码（针对 OAuth 注册用户首次设置密码）
func (s *UserService) SetPassword(userID uint, password, confirmPassword string) error {
	user, err := s.repo.Get(userID)
	if err != nil {
		return err
	}

	// 检查用户是否已有密码
	if user.HasPassword {
		return errors.New("已设置密码，请使用修改密码功能")
	}

	// 验证两次密码是否一致
	if password != confirmPassword {
		return errors.New("两次输入的密码不一致")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.HasPassword = true

	return s.repo.Update(user)
}

// DeactivateAccount 用户注销账号
func (s *UserService) DeactivateAccount(userID uint, password string) error {
	user, err := s.repo.Get(userID)
	if err != nil {
		return err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("密码错误")
	}

	// 标记头像为未使用
	if s.fileService != nil && user.Avatar != "" {
		_ = s.fileService.MarkAsUnused(user.Avatar)
	}

	// 删除用户账号
	return s.repo.Delete(userID)
}

// ============ 后台管理服务 ============

// List 获取用户列表
func (s *UserService) List(req *dto.ListUsersRequest) ([]dto.UserListResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	users, total, err := s.repo.List(offset, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	userResponses := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		resp := dto.UserListResponse{
			ID:           user.ID,
			Email:        user.Email,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Badge:        user.Badge,
			Website:      user.Website,
			Role:         user.Role,
			IsEnabled:    user.IsEnabled,
			LastLogin:    utils.ToJSONTime(user.LastLogin),
			CreatedAt:    utils.NewJSONTime(user.CreatedAt),
			HasPassword:  user.HasPassword,
			GithubID:     user.GithubID,
			GoogleID:     user.GoogleID,
			QQID:         user.QQID,
			MicrosoftID:  user.MicrosoftID,
			FeishuOpenID: user.FeishuOpenID,
		}

		if user.DeletedAt.Valid {
			deletedTime := utils.NewJSONTime(user.DeletedAt.Time)
			resp.DeletedAt = &deletedTime
		}

		userResponses[i] = resp
	}

	return userResponses, total, nil
}

// Create 管理员创建用户
func (s *UserService) Create(operator *model.User, req *dto.AdminCreateUserRequest, host string) error {
	if err := s.ensureCanCreateUser(operator, req.Role); err != nil {
		return err
	}

	// 检查邮箱是否存在
	if s.repo.ExistsByEmail(req.Email) {
		return errors.New("邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Badge:     req.Badge,
		Website:   req.Website,
		Role:      req.Role,
		IsEnabled: true,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	// 如果没有提供头像，使用Cravatar
	if user.Avatar == "" {
		go func() {
			avatarURL, err := s.downloadAndSaveCravatarAvatar(req.Email, user.ID, host)
			if err == nil && avatarURL != "" {
				_ = s.repo.UpdateAvatar(user.ID, avatarURL)
				if s.fileService != nil {
					_ = s.fileService.MarkAsUsed(avatarURL)
				}
			}
		}()
	} else if user.Avatar != "" {
		// 标记手动设置的头像为使用中
		if s.fileService != nil {
			_ = s.fileService.MarkAsUsed(user.Avatar)
		}
	}

	return nil
}

// Update 管理员更新用户
func (s *UserService) Update(operator *model.User, id uint, req *dto.AdminUpdateUserRequest) error {
	user, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	if err := s.ensureCanUpdateUser(operator, user, req); err != nil {
		return err
	}

	// 检查邮箱冲突
	if req.Email != "" && req.Email != user.Email {
		if s.repo.ExistsByEmail(req.Email) {
			return errors.New("邮箱已存在")
		}
		user.Email = req.Email
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	oldAvatar := user.Avatar

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.Badge != "" {
		user.Badge = req.Badge
	}

	if req.Website != "" {
		user.Website = req.Website
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.IsEnabled != nil {
		user.IsEnabled = *req.IsEnabled
	}

	// 更新密码（如果提供）
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// 处理头像变化
	if s.fileService != nil && oldAvatar != user.Avatar {
		if oldAvatar != "" {
			_ = s.fileService.MarkAsUnused(oldAvatar)
		}
		if user.Avatar != "" {
			_ = s.fileService.MarkAsUsed(user.Avatar)
		}
	}

	return s.repo.Update(user)
}

// Delete 软删除用户
func (s *UserService) Delete(operator *model.User, id uint) error {
	user, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	if err := s.ensureCanDeleteUser(operator, user); err != nil {
		return err
	}

	// 标记头像为未使用
	if s.fileService != nil && user.Avatar != "" {
		_ = s.fileService.MarkAsUnused(user.Avatar)
	}

	return s.repo.Delete(id)
}

// ============ 辅助方法 ============

// ensureCanCreateUser 校验当前操作人是否可以创建指定角色的用户
func (s *UserService) ensureCanCreateUser(operator *model.User, role model.UserRole) error {
	if operator == nil {
		return errors.New("未找到当前用户信息")
	}
	if operator.Role == model.RoleSuperAdmin {
		return nil
	}
	if operator.Role != model.RoleAdmin {
		return errors.New("无权限执行该操作")
	}
	if isAdminManagedRole(role) {
		return errors.New("管理员不能创建管理员或超级管理员")
	}
	return nil
}

// ensureCanUpdateUser 校验当前操作人是否可以更新目标用户
func (s *UserService) ensureCanUpdateUser(operator, target *model.User, req *dto.AdminUpdateUserRequest) error {
	if operator == nil {
		return errors.New("未找到当前用户信息")
	}
	if target == nil {
		return errors.New("用户不存在")
	}
	if operator.Role == model.RoleSuperAdmin {
		return nil
	}
	if operator.Role != model.RoleAdmin {
		return errors.New("无权限执行该操作")
	}
	if isAdminManagedRole(target.Role) {
		return errors.New("管理员不能操作管理员或超级管理员账号")
	}
	if isAdminManagedRole(req.Role) {
		return errors.New("管理员不能将用户提升为管理员或超级管理员")
	}
	return nil
}

// ensureCanDeleteUser 校验当前操作人是否可以删除目标用户
func (s *UserService) ensureCanDeleteUser(operator, target *model.User) error {
	if operator == nil {
		return errors.New("未找到当前用户信息")
	}
	if target == nil {
		return errors.New("用户不存在")
	}
	if operator.Role == model.RoleSuperAdmin {
		return nil
	}
	if operator.Role != model.RoleAdmin {
		return errors.New("无权限执行该操作")
	}
	if isAdminManagedRole(target.Role) {
		return errors.New("管理员不能操作管理员或超级管理员账号")
	}
	return nil
}

// isAdminManagedRole 判断是否为管理员可管理受限角色
func isAdminManagedRole(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}

// downloadAndSaveCravatarAvatar 下载并保存Cravatar头像
func (s *UserService) downloadAndSaveCravatarAvatar(email string, userID uint, host string) (string, error) {
	if s.fileService == nil {
		return "", nil
	}

	// 下载头像
	reader, err := utils.DownloadCravatarAvatar(email)
	if err != nil {
		return "", nil
	}

	// 生成文件名
	emailHash := utils.GetEmailHash(email)
	filename := "cravatar_" + emailHash + ".webp"

	// 保存文件
	fileURL, err := s.fileService.UploadFromReader(
		reader,
		filename,
		"image/webp",
		"用户头像",
		userID,
		host,
	)

	if err != nil {
		return "", nil
	}

	return fileURL, nil
}

// validateBadge 校验铭牌是否合法
func (s *UserService) validateBadge(badge string) error {
	if badge == "" {
		return nil
	}
	// 特殊铭牌黑名单
	forbiddenBadges := []string{"站长", "博主", "管理员", "admin", "root", "super_admin"}
	for _, fb := range forbiddenBadges {
		if badge == fb {
			return errors.New("禁止使用该铭牌: " + fb)
		}
	}
	return nil
}

// ============ OAuth 绑定服务 ============

// BindOAuth 绑定第三方账号（已登录用户）
func (s *UserService) BindOAuth(userID uint, provider, providerID, email, avatarURL, host string) (*dto.UserResponse, error) {
	// 获取当前用户
	user, err := s.repo.Get(userID)
	if err != nil {
		return nil, err
	}

	// 检查该 OAuth ID 是否已被其他账户绑定
	existingUser, err := s.repo.GetByOAuthID(provider, providerID)
	if err == nil && existingUser.ID != userID {
		// 已绑定其他账户，执行合并
		if err := s.mergeAccounts(user, existingUser); err != nil {
			return nil, err
		}
		// 重新获取用户信息
		user, err = s.repo.Get(userID)
		if err != nil {
			return nil, err
		}
		return dto.NewUserResponse(user), nil
	}

	// 绑定到当前账户
	if err := s.repo.UpdateOAuthBinding(userID, provider, providerID); err != nil {
		return nil, err
	}

	// 如果用户没有头像，下载第三方头像
	if user.Avatar == "" && avatarURL != "" {
		go s.downloadSocialAvatar(userID, email, avatarURL, host)
	}

	// 返回更新后的用户信息
	return s.Get(userID)
}

// UnbindOAuth 解绑第三方账号
func (s *UserService) UnbindOAuth(userID uint, provider string) error {
	// 获取用户信息
	user, err := s.repo.Get(userID)
	if err != nil {
		return err
	}

	// 检查是否已绑定该提供商
	var isBound bool
	switch provider {
	case "github":
		isBound = user.GithubID != ""
	case "google":
		isBound = user.GoogleID != ""
	case "qq":
		isBound = user.QQID != ""
	case "microsoft":
		isBound = user.MicrosoftID != ""
	default:
		return fmt.Errorf("不支持的登录方式: %s", provider)
	}

	if !isBound {
		return fmt.Errorf("未绑定 %s 登录方式", provider)
	}

	// 统计登录方式数量
	loginCount := 0
	if user.HasPassword {
		loginCount++
	}
	if user.GithubID != "" {
		loginCount++
	}
	if user.GoogleID != "" {
		loginCount++
	}
	if user.QQID != "" {
		loginCount++
	}
	if user.MicrosoftID != "" {
		loginCount++
	}

	// 至少保留一种登录方式
	if loginCount <= 1 {
		return fmt.Errorf("至少保留一种登录方式")
	}

	// 执行解绑
	return s.repo.ClearOAuthBinding(userID, provider)
}

// mergeAccounts 合并账户
func (s *UserService) mergeAccounts(primary, secondary *model.User) error {
	// 转移 OAuth 绑定
	if secondary.GithubID != "" && primary.GithubID == "" {
		primary.GithubID = secondary.GithubID
	}
	if secondary.GoogleID != "" && primary.GoogleID == "" {
		primary.GoogleID = secondary.GoogleID
	}
	if secondary.QQID != "" && primary.QQID == "" {
		primary.QQID = secondary.QQID
	}
	if secondary.MicrosoftID != "" && primary.MicrosoftID == "" {
		primary.MicrosoftID = secondary.MicrosoftID
	}
	if secondary.FeishuOpenID != "" && primary.FeishuOpenID == "" {
		primary.FeishuOpenID = secondary.FeishuOpenID
	}

	// 转移头像
	if primary.Avatar == "" && secondary.Avatar != "" {
		primary.Avatar = secondary.Avatar
		if s.fileService != nil {
			_ = s.fileService.MarkAsUsed(primary.Avatar)
		}
	}

	// 更新主账户
	if err := s.repo.Update(primary); err != nil {
		return err
	}

	// 删除副账户
	return s.repo.Delete(secondary.ID)
}

// ============ 飞书绑定服务 ============

// BindFeishuByEmail 通过邮箱绑定飞书
func (s *UserService) BindFeishuByEmail(email, openID string) error {
	// 查询用户
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证用户是管理员
	if user.Role != model.RoleAdmin && user.Role != model.RoleSuperAdmin {
		return errors.New("只有管理员可以绑定飞书")
	}

	// 检查 OpenID 是否已被其他用户绑定，如果是则清除旧绑定
	existingUser, err := s.repo.GetByOAuthID("feishu", openID)
	if err == nil && existingUser.ID != user.ID {
		_ = s.repo.ClearOAuthBinding(existingUser.ID, "feishu")
	}

	return s.repo.UpdateOAuthBinding(user.ID, "feishu", openID)
}
