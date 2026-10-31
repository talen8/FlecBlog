package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/subosito/gotenv"

	"flec_blog/api/router"
	"flec_blog/config"
	"flec_blog/internal/service"
	"flec_blog/pkg/database"
	"flec_blog/pkg/logsetup"

	_ "flec_blog/docs" // swagger docs
)

// @title           FlecBlog
// @version         v1
// @description     一个基于 Go 语言的现代化博客后端服务

// @contact.name   Talen
// @contact.email  talen2004@163.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在请求头中添加 Bearer Token，格式：Bearer {token}

func main() {
	// 加载 .env 文件
	_ = gotenv.Load()

	// 初始化日志
	if err := logsetup.Setup(); err != nil {
		fatal("初始化日志失败", slog.Any("err", err))
	}
	defer logsetup.Close()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fatal("加载配置失败", slog.Any("err", err))
	}

	// 初始化数据库连接
	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		fatal("连接数据库失败", slog.Any("err", err))
	}
	defer func() {
		_ = db.Close()
	}()

	// 执行数据库迁移
	if err := database.RunMigrations(db.DB); err != nil {
		_ = db.Close()
		fatal("执行数据库迁移失败", slog.Any("err", err))
	}

	// 检查 API_URL 环境变量
	if os.Getenv("API_URL") == "" {
		fatal("API_URL 环境变量未设置")
	}

	// 从数据库加载运行时配置（邮箱、上传等）
	settingService := service.NewSettingService(db.DB)
	settingService.SetConfig(cfg)               // 设置全局配置对象，用于热重载
	_ = settingService.ApplyDatabaseConfig(cfg) // 应用数据库配置

	// 初始化路由与调度器
	r, sched := router.InitRouter(db, cfg)

	defer sched.Stop()

	// 启动服务器
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		defer logsetup.Recover("退出信号处理")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.ErrorContext(ctx, "HTTP 服务关闭失败", slog.Any("err", err))
		}
	}()

	slog.Info("服务已启动", slog.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fatal("启动服务器失败", slog.Any("err", err))
	}
}

// fatal 记录致命错误并退出进程
func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
