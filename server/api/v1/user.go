package v1

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"
	"flec_blog/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// UserController 用户控制器
type UserController struct {
	userService         *service.UserService
	verificationService *service.VerificationService
	config              *config.Config
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService, verificationService *service.VerificationService, cfg *config.Config) *UserController {
	return &UserController{
		userService:         userService,
		verificationService: verificationService,
		config:              cfg,
	}
}

// ============ 认证接口 ============

// BeginAuth 第三方认证
//
//	@Summary		第三方认证
//	@Description	跳转到第三方认证页面，支持登录和绑定
//	@Tags			认证
//	@Param			provider	path	string	true	"提供商 (github/google/qq)"
//	@Param			action		query	string	false	"操作类型：login(默认) 或 bind"
//	@Param			token		query	string	false	"绑定时的认证token"
//	@Param			redirect	query	string	false	"成功后跳转的页面路径"
//	@Router			/auth/{provider} [get]
func (c *UserController) BeginAuth(ctx *gin.Context) {
	provider := ctx.Param("provider")
	action := ctx.DefaultQuery("action", "login")
	redirect := ctx.Query("redirect")

	// 构建state参数
	var state string
	if action == "bind" {
		token := ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "绑定操作缺少认证token"})
			return
		}
		user, err := c.userService.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的认证token: " + err.Error()})
			return
		}
		state = "bind:" + strconv.FormatUint(uint64(user.ID), 10)
	} else if redirect != "" {
		state = "redirect:" + redirect
	}

	// 添加redirect到绑定state
	if action == "bind" && redirect != "" {
		state += "|redirect:" + redirect
	}

	// 设置OAuth参数并跳转
	q := ctx.Request.URL.Query()
	q.Set("provider", provider)
	if state != "" {
		q.Set("state", state)
	}
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// AuthCallback 第三方认证回调
//
//	@Summary		第三方认证回调
//	@Description	处理第三方登录回调，返回 Token 并跳转
//	@Tags			认证
//	@Param			provider	path	string	true	"提供商 (github/google/qq)"
//	@Router			/auth/{provider}/callback [get]
func (c *UserController) AuthCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	state := ctx.Query("state")

	// 解析state
	var action string
	var bindUserID uint
	var redirect string
	switch {
	case strings.HasPrefix(state, "bind:"):
		action = "bind"
		parts := strings.Split(state, "|")
		if len(parts) > 0 {
			idStr := strings.TrimPrefix(parts[0], "bind:")
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				bindUserID = uint(id)
			}
		}
		if len(parts) > 1 && strings.HasPrefix(parts[1], "redirect:") {
			redirect = strings.TrimPrefix(parts[1], "redirect:")
		}
	case strings.HasPrefix(state, "redirect:"):
		action = "login"
		redirect = strings.TrimPrefix(state, "redirect:")
	default:
		action = "login"
	}

	// 完成OAuth认证
	q := ctx.Request.URL.Query()
	q.Set("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "认证失败: " + err.Error()})
		return
	}

	// 获取昵称和邮箱（公共逻辑）
	nickname := oauthUser.NickName
	if nickname == "" {
		nickname = oauthUser.Name
	}
	email := oauthUser.Email
	if email == "" {
		email = provider + "_" + oauthUser.UserID + "@virtual.local"
	}

	host := upload.ExtractHostFromContext(ctx)

	// 获取前端基础URL（公共逻辑）
	frontendBaseURL := c.config.Basic.BlogURL
	if frontendBaseURL == "" {
		frontendBaseURL = "http://localhost:3000"
	}

	// 处理绑定或登录
	if action == "bind" {
		_, err := c.userService.BindOAuth(bindUserID, provider, oauthUser.UserID, email, oauthUser.AvatarURL, host)
		if err != nil {
			// 直接内联错误重定向逻辑
			targetPath := "/profile"
			if redirect != "" {
				targetPath = redirect
			}
			ctx.Redirect(http.StatusFound, frontendBaseURL+targetPath+"?bind=error&message="+url.QueryEscape("绑定失败: "+err.Error()))
			return
		}
		// 跳转到个人资料页
		targetPath := "/profile"
		if redirect != "" {
			targetPath = redirect
		}
		ctx.Redirect(http.StatusFound, frontendBaseURL+targetPath+"?bind=success&provider="+provider)
	} else {
		// 处理登录
		loginResp, err := c.userService.LoginBySocial(provider, oauthUser.UserID, email, nickname, oauthUser.AvatarURL, host)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败: " + err.Error()})
			return
		}
		// 跳转到前端回调页
		frontendURL := frontendBaseURL + "/oauth/callback?token=" + loginResp.AccessToken + "&refresh_token=" + loginResp.RefreshToken
		if redirect != "" {
			frontendURL += "&redirect=" + url.QueryEscape(redirect)
		}
		ctx.Redirect(http.StatusFound, frontendURL)
	}
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	邮箱密码注册
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"注册信息"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Router			/auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	host := upload.ExtractHostFromContext(ctx)
	loginResp, err := c.userService.Register(&req, host)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, loginResp)
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	邮箱密码登录，返回 JWT token 和用户基本信息
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"登录信息"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	loginResp, err := c.userService.Login(&req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, loginResp)
}

// RefreshToken 刷新token
//
//	@Summary		刷新token
//	@Description	使用refresh token获取新的access token
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh Token"
//	@Success		200		{object}	response.Response{data=dto.LoginResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/auth/refresh [post]
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	refreshResp, err := c.userService.RefreshToken(&req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, refreshResp)
}

// ForgotPassword 重置密码请求
//
//	@Summary		重置密码请求
//	@Description	通过邮箱找回密码，发送重置验证码
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ForgotPasswordRequest	true	"邮箱地址"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/auth/forgot-password [post]
func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.verificationService.SendPasswordReset(&req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ResetPassword 重置密码
//
//	@Summary		重置密码
//	@Description	凭验证码设置新密码，成功后需重新登录
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ResetPasswordRequest	true	"重置密码信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/auth/reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.verificationService.ResetPassword(&req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	将当前 token 加入黑名单使其失效
//	@Tags			认证
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Router			/auth/logout [post]
func (c *UserController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		response.Failed(ctx, "未提供令牌")
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := c.userService.Logout(token); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ============ 用户信息接口 ============

// GetProfile 获取用户信息
//
//	@Summary		个人资料
//	@Description	获取当前登录用户的信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.UserResponse}
//	@Failure		401	{object}	response.Response
//	@Router			/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未找到用户信息")
		return
	}

	userInfo, err := c.userService.Get(userID.(uint))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, userInfo)
}

// UpdateForWeb 更新用户信息
//
//	@Summary		更新资料
//	@Description	修改用户资料信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateUserRequest	true	"用户信息"
//	@Success		200		{object}	response.Response{data=dto.UserResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/user/profile [patch]
func (c *UserController) UpdateForWeb(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未找到用户信息")
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.userService.UpdateForWeb(userID.(uint), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	userInfo, err := c.userService.Get(userID.(uint))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, userInfo)
}

// ChangePassword 修改当前用户密码
//
//	@Summary		修改密码
//	@Description	修改密码，需提供旧密码验证
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ChangePasswordRequest	true	"密码信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/user/password [put]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未找到用户信息")
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// SetPassword OAuth用户设置密码
//
//	@Summary		设置密码
//	@Description	OAuth 注册用户首次设置密码，无需旧密码验证
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.SetPasswordRequest	true	"密码信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/user/password [post]
func (c *UserController) SetPassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未找到用户信息")
		return
	}

	var req dto.SetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.userService.SetPassword(userID.(uint), req.Password, req.ConfirmPassword); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// DeactivateAccount 用户注销账号
//
//	@Summary		注销账号
//	@Description	主动注销自己的账号，需提供密码验证，注销后账号将被软删除
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.DeactivateAccountRequest	true	"注销账号信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/user/deactivate [delete]
func (c *UserController) DeactivateAccount(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未找到用户信息")
		return
	}

	var req dto.DeactivateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.userService.DeactivateAccount(userID.(uint), req.Password); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// UnbindOAuth 解绑第三方账号
//
//	@Summary		解绑第三方账号
//	@Description	解绑已绑定的第三方账号
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			provider	path		string	true	"提供商 (github/google)"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Router			/user/oauth/{provider} [delete]
func (c *UserController) UnbindOAuth(ctx *gin.Context) {
	// 验证用户登录状态
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Failed(ctx, "未登录")
		return
	}

	// 验证提供商类型
	provider := ctx.Param("provider")
	if provider != "github" && provider != "google" && provider != "qq" {
		response.ValidateFailed(ctx, "不支持的登录方式")
		return
	}

	// 执行解绑
	if err := c.userService.UnbindOAuth(userID.(uint), provider); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ============ 后台管理接口 ============

// List 获取用户列表
//
//	@Summary		用户列表
//	@Description	获取所有用户，支持多种筛选条件
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(10)
//	@Param			keyword		query		string	false	"搜索关键词（邮箱、昵称）"
//	@Param			role		query		string	false	"角色筛选（super_admin/admin/user/guest）"
//	@Param			is_enabled	query		bool	false	"状态筛选（true:启用 false:禁用）"
//	@Param			is_deleted	query		bool	false	"是否已删除（true:已删除 false:未删除）"
//	@Param			login_method	query		string	false	"登录方式筛选（password/github/google/qq/microsoft）"
//	@Param			last_login_start	query		string	false	"最后登录开始时间（格式：2006-01-02）"
//	@Param			last_login_end	query		string	false	"最后登录结束时间（格式：2006-01-02）"
//	@Param			start_time	query		string	false	"注册开始时间（格式：2006-01-02）"
//	@Param			end_time	query		string	false	"注册结束时间（格式：2006-01-02）"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/users [get]
func (c *UserController) List(ctx *gin.Context) {
	var req dto.ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	users, total, err := c.userService.List(&req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, users, total, req.Page, req.PageSize)
}

// Get 获取用户详情
//
//	@Summary		用户详情
//	@Description	查看用户详细信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	response.Response{data=dto.UserResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/admin/users/{id} [get]
func (c *UserController) Get(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, "无效的用户ID")
		return
	}

	userInfo, err := c.userService.Get(uint(userID))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, userInfo)
}

// Create 创建用户
//
//	@Summary		创建用户
//	@Description	创建新用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AdminCreateUserRequest	true	"用户信息"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Router			/admin/users [post]
func (c *UserController) Create(ctx *gin.Context) {
	var req dto.AdminCreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	operator, ok := ctx.Get("user")
	if !ok {
		response.Failed(ctx, "未找到当前用户信息")
		return
	}

	authUser, ok := operator.(*model.User)
	if !ok || authUser == nil {
		response.Failed(ctx, "当前用户信息无效")
		return
	}

	host := upload.ExtractHostFromContext(ctx)
	if err := c.userService.Create(authUser, &req, host); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, nil)
}

// Update 更新用户
//
//	@Summary		更新用户
//	@Description	修改用户信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"用户 ID"
//	@Param			request	body		dto.AdminUpdateUserRequest	true	"用户信息（包含可选的密码字段）"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/users/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, "无效的用户ID")
		return
	}

	var req dto.AdminUpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	operator, ok := ctx.Get("user")
	if !ok {
		response.Failed(ctx, "未找到当前用户信息")
		return
	}

	authUser, ok := operator.(*model.User)
	if !ok || authUser == nil {
		response.Failed(ctx, "当前用户信息无效")
		return
	}

	if err := c.userService.Update(authUser, uint(userID), &req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// Delete 删除用户
//
//	@Summary		删除用户
//	@Description	软删除用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/admin/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, "无效的用户ID")
		return
	}

	operator, ok := ctx.Get("user")
	if !ok {
		response.Failed(ctx, "未找到当前用户信息")
		return
	}

	authUser, ok := operator.(*model.User)
	if !ok || authUser == nil {
		response.Failed(ctx, "当前用户信息无效")
		return
	}

	if err := c.userService.Delete(authUser, uint(userID)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}
