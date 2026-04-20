package router

import (
	"flec_blog/api/middleware"
	v1 "flec_blog/api/v1"
	"flec_blog/api/v1/feeds"
	"flec_blog/config"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/pkg/database"
	"flec_blog/pkg/email"
	"flec_blog/pkg/feishu"
	mcpserver "flec_blog/pkg/mcp"
	"flec_blog/pkg/notification"
	"flec_blog/pkg/scheduler"
	"flec_blog/pkg/upload"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter(db *database.Database, conf *config.Config) *gin.Engine {
	// 设置为 release 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// 配置代理信任
	_ = r.SetTrustedProxies([]string{"127.0.0.1/32", "::1/128", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})

	// 初始化仓库
	articleRepo := repository.NewArticleRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	tagRepo := repository.NewTagRepository(db.DB)
	commentRepo := repository.NewCommentRepository(db.DB)
	fileRepo := repository.NewFileRepository(db.DB)
	verificationRepo := repository.NewVerificationRepository(db.DB)
	statsRepo := repository.NewStatsRepository(db.DB)
	friendRepo := repository.NewFriendRepository(db.DB)
	momentRepo := repository.NewMomentRepository(db.DB)
	menuRepo := repository.NewMenuRepository(db.DB)
	notificationRepo := repository.NewNotificationRepository(db.DB)
	feedbackRepo := repository.NewFeedbackRepository(db.DB)
	subscriberRepo := repository.NewSubscriberRepository(db.DB)
	rssFeedRepo := repository.NewRssFeedRepository(db.DB)
	settingRepo := repository.NewSettingRepository(db.DB)

	// 初始化上传系统
	uploadManager := upload.InitializeUploadSystem(conf)

	// 使用中间件
	r.Use(middleware.CORS(conf))              // CORS跨域
	r.Use(middleware.Logger())                // 日志记录
	r.Use(middleware.RateLimit(500, 1, "ip")) // 全局IP限流: 500次/分钟
	r.Use(middleware.Recovery())              // 错误恢复

	// 注册本地静态文件服务
	r.Static("/uploads", "/app/data/uploads")

	// 根路径欢迎页面
	r.GET("/", func(c *gin.Context) { c.String(200, "Flec-Server 运行成功") })

	// Swagger API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fileService := service.NewFileService(fileRepo, uploadManager)
	fileUsageChecker := service.NewFileUsageChecker(articleRepo, friendRepo, momentRepo, settingRepo, userRepo, menuRepo, feedbackRepo, commentRepo)
	fileService.SetUsageChecker(fileUsageChecker)

	// 初始化服务
	emailClient := email.Initialize(conf)
	feishuClient := feishu.Initialize(conf)
	notificationSvc := notification.NewService(emailClient, feishuClient, conf)
	userService := service.NewUserService(userRepo, fileService, conf)
	verificationService := service.NewVerificationService(verificationRepo, userRepo, emailClient, conf)
	articleService := service.NewArticleService(articleRepo, tagRepo, categoryRepo, commentRepo, fileService, db.DB)
	tagService := service.NewTagService(tagRepo, articleRepo)
	categoryService := service.NewCategoryService(categoryRepo, articleRepo)
	notificationService := service.NewNotificationService(notificationRepo, notificationSvc)
	commentService := service.NewCommentService(commentRepo, articleRepo, userRepo, notificationService, fileService)
	statsService := service.NewStatsService(statsRepo, conf)
	friendService := service.NewFriendService(friendRepo, fileService, notificationService)
	momentService := service.NewMomentService(momentRepo, fileService)
	menuService := service.NewMenuService(menuRepo, fileService)
	feedbackService := service.NewFeedbackService(feedbackRepo, notificationService, fileService)
	subscriberService := service.NewSubscriberService(subscriberRepo, emailClient, conf)
	rssFeedService := service.NewRssFeedService(rssFeedRepo, notificationService)
	systemService := service.NewSystemService(db.DB, uploadManager, emailClient, feishuClient, notificationService)
	feishu.InitCardHandlers(friendService, commentService, userService, statsService, systemService, rssFeedService)
	settingService := service.NewSettingService(db.DB)
	settingService.SetConfig(conf)                         // 设置全局配置对象，用于热重载
	settingService.SetFileService(fileService)             // 设置文件服务，用于文件状态管理
	articleService.SetSubscriberService(subscriberService) // 设置订阅服务，用于文章推送

	// 初始化并启动定时任务调度器
	initScheduler(fileService, userService, verificationService, rssFeedService, friendService, systemService)

	// 初始化控制器
	userController := v1.NewUserController(userService, verificationService, conf)
	articleController := v1.NewArticleController(articleService)
	categoryController := v1.NewCategoryController(categoryService)
	tagController := v1.NewTagController(tagService)
	commentController := v1.NewCommentController(commentService)
	fileController := v1.NewFileController(fileService)
	statsHandler := v1.NewStatsHandler(statsService)
	friendController := v1.NewFriendController(friendService)
	momentController := v1.NewMomentController(momentService)
	menuHandler := v1.NewMenuHandler(menuService)
	notificationController := v1.NewNotificationController(notificationService)
	feedbackHandler := v1.NewFeedbackHandler(feedbackService)
	subscriberHandler := v1.NewSubscriberHandler(subscriberService)
	settingController := v1.NewSettingController(settingService, db.DB, uploadManager)
	systemController := v1.NewSystemController(systemService)
	atomController := feeds.NewAtomController(articleService, conf)
	rssController := feeds.NewRSSController(articleService, conf)
	toolsHandler := v1.NewToolsController()
	aiController := v1.NewAIController(settingService)
	rssFeedController := v1.NewRssFeedController(rssFeedService)

	// MCP 接口
	mcpHandler := gin.WrapH(mcpserver.NewPublicHandler(
		articleService, categoryService, tagService, commentService, friendService, rssFeedService, momentService,
		userService, statsService,
	))
	r.Any("/mcp", middleware.MCPAuth(conf), mcpHandler)

	// Atom 订阅
	r.GET("/atom.xml", atomController.GetAtomFeed)
	// RSS 2.0 订阅
	r.GET("/rss.xml", rssController.GetRSSFeed)

	// ============================
	// 前台 API 路由组 - 面向普通用户的公开接口
	// ============================
	frontendAPI := r.Group("/api/v1")
	{
		// ==================== 站点配置（公开接口）====================
		frontendAPI.GET("/settings/:group", settingController.GetPublicSettingGroup)

		// ==================== 统计数据收集（公开接口）====================
		frontendAPI.POST("/collect", statsHandler.Collect)

		// ==================== 统计信息（公开接口）====================
		statsGroup := frontendAPI.Group("/stats")
		{
			statsGroup.GET("/site", statsHandler.GetSiteStats)    // 网站统计信息
			statsGroup.GET("/archives", statsHandler.GetArchives) // 文章归档数据
		}

		// ==================== 用户认证相关 ====================
		authGroup := frontendAPI.Group("/auth")
		{
			// OAuth 第三方登录
			authGroup.GET("/:provider", userController.BeginAuth)
			authGroup.GET("/:provider/callback", userController.AuthCallback)

			// 公开接口
			authGroup.POST("/register", userController.Register)
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/refresh", userController.RefreshToken)

			// 密码找回接口
			authGroup.POST("/forgot-password", userController.ForgotPassword)
			authGroup.POST("/reset-password", userController.ResetPassword)

			// 需要登录的接口
			authGroup.POST("/logout", middleware.Auth(userService), userController.Logout)
		}

		// ==================== 用户信息相关 ====================
		userGroup := frontendAPI.Group("/user")
		userGroup.Use(middleware.Auth(userService), middleware.IsUserOrAbove())
		{
			userGroup.GET("/profile", userController.GetProfile)              // 获取当前用户信息
			userGroup.PATCH("/profile", userController.UpdateForWeb)          // 更新当前用户信息
			userGroup.PUT("/password", userController.ChangePassword)         // 修改当前用户密码
			userGroup.POST("/password", userController.SetPassword)           // 设置密码（OAuth用户首次设置）
			userGroup.DELETE("/deactivate", userController.DeactivateAccount) // 注销账号
			userGroup.DELETE("/oauth/:provider", userController.UnbindOAuth)  // 解绑第三方账号
		}

		// ==================== 文章相关 ====================
		articleGroup := frontendAPI.Group("/articles")
		{
			// 公开接口
			articleGroup.GET("", articleController.ListForWeb)      // 获取前台文章列表
			articleGroup.GET("/search", articleController.Search)   // 搜索文章
			articleGroup.GET("/:slug", articleController.GetBySlug) // 通过slug获取文章详情
		}

		// ==================== 标签相关 ====================
		tagGroup := frontendAPI.Group("/tags")
		{
			// 公开接口
			tagGroup.GET("", tagController.ListForWeb)      // 获取标签列表
			tagGroup.GET("/:slug", tagController.GetBySlug) // 通过slug获取标签详情
		}

		// ==================== 分类相关 ====================
		categoryGroup := frontendAPI.Group("/categories")
		{
			// 公开接口
			categoryGroup.GET("", categoryController.ListForWeb)      // 获取分类列表
			categoryGroup.GET("/:slug", categoryController.GetBySlug) // 通过slug获取分类详情
		}

		// ==================== 友链相关 ====================
		friendGroup := frontendAPI.Group("/friends")
		{
			// 公开接口
			friendGroup.GET("", friendController.ListForWeb) // 获取友链分组列表

			// 需要登录的接口
			authFriend := friendGroup.Group("")
			authFriend.Use(middleware.Auth(userService), middleware.IsUserOrAbove())
			{
				authFriend.POST("/apply", friendController.ApplyFriend) // 申请友链
			}
		}

		// ==================== 动态相关 ====================
		momentGroup := frontendAPI.Group("/moments")
		{
			// 公开接口
			momentGroup.GET("", momentController.ListForWeb) // 获取动态列表
		}

		// ==================== 菜单相关 ====================
		// 公开接口：获取菜单树（支持按类型筛选）
		frontendAPI.GET("/menus", menuHandler.ListForWeb)

		// ==================== 评论相关 ====================
		commentGroup := frontendAPI.Group("/comments")
		{
			// 公开接口
			commentGroup.GET("", commentController.ListForWeb)                                    // 获取评论列表
			commentGroup.POST("", middleware.OptionalAuth(userService), commentController.Create) // 发表评论

			// 需要登录的接口
			authComment := commentGroup.Group("")
			authComment.Use(middleware.Auth(userService), middleware.IsUserOrAbove())
			{
				authComment.PUT("/:id", commentController.Update)
				authComment.DELETE("/:id", commentController.DeleteForWeb)
			}
		}

		// ==================== 通知相关（前台用户）====================
		notificationGroup := frontendAPI.Group("/notifications")
		notificationGroup.Use(middleware.Auth(userService), middleware.IsUserOrAbove())
		{
			notificationGroup.GET("", notificationController.ListForWeb)             // 获取通知列表（含未读数量）
			notificationGroup.PUT("/:id/read", notificationController.MarkAsRead)    // 标记已读
			notificationGroup.PUT("/read-all", notificationController.MarkAllAsRead) // 全部标记已读
		}

		// ==================== 文件上传相关 ====================
		// 使用可选认证：用户头像需要登录，评论贴图/反馈投诉允许匿名
		uploadGroup := frontendAPI.Group("/upload")
		uploadGroup.Use(middleware.OptionalAuth(userService))
		{
			uploadGroup.POST("", fileController.UploadForWeb)
		}

		// ==================== 反馈投诉相关 ====================
		// 公开接口：反馈接口组
		feedbackGroup := frontendAPI.Group("/feedback")
		{
			feedbackGroup.POST("", feedbackHandler.Submit)                         // 提交反馈
			feedbackGroup.GET("/ticket/:ticket_no", feedbackHandler.GetByTicketNo) // 查询工单
		}

		// ==================== 邮件订阅相关 ====================
		// 公开接口：订阅接口组
		subscribeGroup := frontendAPI.Group("/subscribe")
		{
			subscribeGroup.POST("", subscriberHandler.Subscribe)              // 订阅
			subscribeGroup.GET("/unsubscribe", subscriberHandler.Unsubscribe) // 退订
		}
	}

	// ============================
	// 后台管理 API 路由组 - 面向管理员的管理接口
	// ============================
	adminAPI := r.Group("/api/v1/admin")
	adminAPI.Use(middleware.Auth(userService), middleware.IsAdminOrAbove())
	{
		// ==================== 用户管理 ====================
		userManagement := adminAPI.Group("/users")
		{
			userManagement.GET("", userController.List)          // 获取用户列表
			userManagement.GET("/:id", userController.Get)       // 获取用户详情
			userManagement.POST("", userController.Create)       // 创建用户
			userManagement.PUT("/:id", userController.Update)    // 更新用户信息
			userManagement.DELETE("/:id", userController.Delete) // 删除用户
		}

		// ==================== 文章管理 ====================
		articleManagement := adminAPI.Group("/articles")
		{
			articleManagement.GET("", articleController.List)          // 获取文章列表
			articleManagement.GET("/:id", articleController.Get)       // 获取文章详情
			articleManagement.POST("", articleController.Create)       // 创建文章
			articleManagement.PUT("/:id", articleController.Update)    // 更新文章
			articleManagement.DELETE("/:id", articleController.Delete) // 删除文章

			// 数据导入
			articleManagement.POST("/import", articleController.ImportArticles) // 导入文章数据（Hexo等）

			// 微信公众号格式生成
			articleManagement.POST("/:id/wechat/export", articleController.ExportToWeChat) // 生成微信公众号 HTML

			// 文章下载
			articleManagement.GET("/:id/download/zip", articleController.DownloadZip) // 下载为 Markdown
		}

		// ==================== 标签管理 ====================
		tagManagement := adminAPI.Group("/tags")
		{
			tagManagement.GET("", tagController.List)          // 获取标签列表
			tagManagement.GET("/:id", tagController.Get)       // 获取标签详情
			tagManagement.POST("", tagController.Create)       // 创建标签
			tagManagement.PUT("/:id", tagController.Update)    // 更新标签
			tagManagement.DELETE("/:id", tagController.Delete) // 删除标签
		}

		// ==================== 分类管理 ====================
		categoryManagement := adminAPI.Group("/categories")
		{
			categoryManagement.GET("", categoryController.List)          // 获取分类列表
			categoryManagement.GET("/:id", categoryController.Get)       // 获取分类详情
			categoryManagement.POST("", categoryController.Create)       // 创建分类
			categoryManagement.PUT("/:id", categoryController.Update)    // 更新分类
			categoryManagement.DELETE("/:id", categoryController.Delete) // 删除分类
		}

		// ==================== 友链管理 ====================
		friendManagement := adminAPI.Group("/friends")
		{
			// 友链类型管理
			friendManagement.GET("/types", friendController.ListTypes)         // 获取类型列表
			friendManagement.GET("/types/:id", friendController.GetType)       // 获取类型详情
			friendManagement.POST("/types", friendController.CreateType)       // 创建类型
			friendManagement.PUT("/types/:id", friendController.UpdateType)    // 更新类型
			friendManagement.DELETE("/types/:id", friendController.DeleteType) // 删除类型

			// 友链管理
			friendManagement.GET("", friendController.List)          // 获取友链列表
			friendManagement.GET("/:id", friendController.Get)       // 获取友链详情
			friendManagement.POST("", friendController.Create)       // 创建友链
			friendManagement.PUT("/:id", friendController.Update)    // 更新友链
			friendManagement.DELETE("/:id", friendController.Delete) // 删除友链
		}

		// ==================== 动态管理 ====================
		momentManagement := adminAPI.Group("/moments")
		{
			momentManagement.GET("", momentController.List)          // 获取动态列表
			momentManagement.GET("/:id", momentController.Get)       // 获取动态详情
			momentManagement.POST("", momentController.Create)       // 创建动态
			momentManagement.PUT("/:id", momentController.Update)    // 更新动态
			momentManagement.DELETE("/:id", momentController.Delete) // 删除动态
		}

		// ==================== 评论管理 ====================
		commentManagement := adminAPI.Group("/comments")
		{
			commentManagement.POST("", commentController.CreateForAdmin)                // 创建评论（管理员回复）
			commentManagement.GET("", commentController.List)                           // 获取评论列表
			commentManagement.GET("/:id", commentController.Get)                        // 获取评论详情
			commentManagement.PUT("/:id/toggle-status", commentController.ToggleStatus) // 切换评论状态
			commentManagement.DELETE("/:id", commentController.Delete)                  // 删除评论（管理员硬删除）

			// 数据导入
			commentManagement.POST("/import", commentController.ImportComments) // 导入评论数据（Artalk等）
		}

		// ==================== 文件管理 ====================
		fileManagement := adminAPI.Group("/files")
		{
			fileManagement.POST("", fileController.Upload)       // 上传文件
			fileManagement.GET("", fileController.List)          // 获取文件列表
			fileManagement.GET("/:id", fileController.Get)       // 获取文件详情
			fileManagement.DELETE("/:id", fileController.Delete) // 删除文件
		}

		// ==================== 统计管理 ====================
		statsManagement := adminAPI.Group("/stats")
		{
			statsManagement.GET("/dashboard", statsHandler.GetDashboard)              // 获取仪表盘统计
			statsManagement.GET("/trend", statsHandler.GetTrend)                      // 获取访问趋势
			statsManagement.GET("/category", statsHandler.GetCategoryStats)           // 获取分类统计
			statsManagement.GET("/tag", statsHandler.GetTagStats)                     // 获取标签统计
			statsManagement.GET("/contribution", statsHandler.GetArticleContribution) // 获取文章贡献数据
			statsManagement.GET("/visits", statsHandler.GetVisitLogs)                 // 获取访问日志
		}

		// ==================== 菜单管理 ====================
		menuManagement := adminAPI.Group("/menus")
		{
			menuManagement.GET("", menuHandler.List)          // 获取菜单树
			menuManagement.GET("/:id", menuHandler.Get)       // 获取菜单详情
			menuManagement.POST("", menuHandler.Create)       // 创建菜单
			menuManagement.PUT("/:id", menuHandler.Update)    // 更新菜单
			menuManagement.DELETE("/:id", menuHandler.Delete) // 删除菜单
		}

		// ==================== 反馈管理 ====================
		feedbackManagement := adminAPI.Group("/feedback")
		{
			feedbackManagement.GET("", feedbackHandler.List)          // 获取反馈列表
			feedbackManagement.GET("/:id", feedbackHandler.Get)       // 获取反馈详情
			feedbackManagement.PUT("/:id", feedbackHandler.Update)    // 更新反馈
			feedbackManagement.DELETE("/:id", feedbackHandler.Delete) // 删除反馈
		}

		// ==================== 通知管理 ====================
		notificationManagement := adminAPI.Group("/notifications")
		{
			notificationManagement.GET("", notificationController.List)                   // 获取通知列表（含未读数量）
			notificationManagement.PUT("/:id/read", notificationController.MarkAsRead)    // 标记已读
			notificationManagement.PUT("/read-all", notificationController.MarkAllAsRead) // 全部标记已读
		}

		// ==================== 配置管理 ====================
		settingManagement := adminAPI.Group("/settings")
		{
			settingManagement.GET("/:group", settingController.GetGroup)                                               // 获取指定分组的配置
			settingManagement.PATCH("/:group", middleware.IsSuperAdmin(), settingController.UpdateGroup)               // 更新指定分组的配置（仅超级管理员）
			settingManagement.PUT("/ai/mcp-secret/reset", middleware.IsSuperAdmin(), settingController.ResetMCPSecret) // 重置 MCP Secret（仅超级管理员）
		}

		// ==================== 系统信息 ====================
		systemManagement := adminAPI.Group("/system")
		{
			systemManagement.GET("/static", systemController.GetSystemStatic)   // 获取系统静态信息
			systemManagement.GET("/dynamic", systemController.GetSystemDynamic) // 获取系统动态信息
		}

		// ==================== 管理工具相关 ====================
		toolsManagement := adminAPI.Group("/tools")
		toolsManagement.POST("/parse-video", toolsHandler.ParseVideo)
		toolsManagement.POST("/fetch-linkmeta", toolsHandler.FetchLinkMetadata)
		toolsManagement.POST("/download-image", toolsHandler.DownloadImage)

		// ==================== AI功能相关 ====================
		aiManagement := adminAPI.Group("/ai")
		aiManagement.POST("/test", aiController.TestConfig)
		aiManagement.POST("/summary", aiController.Summary)
		aiManagement.POST("/ai-summary", aiController.AISummary)
		aiManagement.POST("/title", aiController.Title)

		// ==================== RSS订阅管理 ====================
		rssFeedManagement := adminAPI.Group("/rssfeed")
		{
			rssFeedManagement.GET("", rssFeedController.List)                                            // 获取RSS文章列表
			rssFeedManagement.PUT("/:id/read", middleware.IsSuperAdmin(), rssFeedController.MarkRead)    // 标记文章已读（仅超级管理员）
			rssFeedManagement.PUT("/read-all", middleware.IsSuperAdmin(), rssFeedController.MarkAllRead) // 全部标记已读（仅超级管理员）
		}

		// ==================== 邮件订阅者管理 ====================
		subscriberManagement := adminAPI.Group("/subscribers")
		{
			subscriberManagement.GET("", subscriberHandler.List)          // 获取订阅者列表
			subscriberManagement.DELETE("/:id", subscriberHandler.Delete) // 删除订阅者
		}
	}

	return r
}

// initScheduler 初始化并启动定时任务调度器
func initScheduler(fileService *service.FileService, userService *service.UserService, verificationService *service.VerificationService, rssFeedService *service.RssFeedService, friendService *service.FriendService, systemService *service.SystemService) {
	s := scheduler.NewScheduler()

	// 注册清理任务
	_ = s.AddJob(scheduler.NewJob("清理过期验证码", "0 0 2 * * *", verificationService.CleanExpiredVerifications))
	_ = s.AddJob(scheduler.NewJob("清理未使用文件", "0 0 3 * * *", fileService.DeleteUnusedFiles))
	_ = s.AddJob(scheduler.NewJob("清理过期Token黑名单", "0 0 4 * * *", userService.CleanupExpiredTokens))

	// RSS订阅相关任务
	_ = s.AddJob(scheduler.NewJob("刷新RSS订阅源", "0 0 * * * *", rssFeedService.RefreshAllFeeds))
	_ = s.AddJob(scheduler.NewJob("清理孤立RSS文章", "0 0 5 * * *", rssFeedService.CleanOrphanedArticles))
	_ = s.AddJob(scheduler.NewJob("每日RSS订阅推送", "0 0 9 * * *", rssFeedService.SendDailyPush))

	// 友链检测任务
	_ = s.AddJob(scheduler.NewJob("友链状态检测", "0 0 2 * * 3", friendService.CheckAllFriends))

	_ = s.AddJob(scheduler.NewJob("版本更新检测", "0 0 8 * * *", systemService.CheckForUpdates))

	s.Start()
}
