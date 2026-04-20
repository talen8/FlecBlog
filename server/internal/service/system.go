package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/pkg/email"
	feishupkg "flec_blog/pkg/feishu"
	"flec_blog/pkg/upload"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"gorm.io/gorm"
)

const githubLatestReleaseAPI = "https://api.github.com/repos/talen8/FlecBlog/releases/latest"

// AppVersion 由构建参数注入，默认 dev
var AppVersion = "dev"

// VersionStatus 版本检测状态
type VersionStatus struct {
	LatestVersion  string `json:"latest_version"`
	LastCheckError string `json:"last_check_error"`
}

// SystemService 系统服务
type SystemService struct {
	db                  *gorm.DB
	uploadManager       *upload.Manager
	emailClient         *email.Client
	feishuClient        *feishupkg.Client
	notificationService *NotificationService
	httpClient          *http.Client
	mu                  sync.RWMutex
	versionStatus       VersionStatus
}

type versionManifest struct {
	LatestVersion string
	ReleaseURL    string
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

type parsedVersion struct {
	numbers    [3]int
	prerelease []string
}

// NewSystemService 创建系统服务
func NewSystemService(db *gorm.DB, uploadManager *upload.Manager, emailClient *email.Client, feishuClient *feishupkg.Client, notificationService *NotificationService) *SystemService {
	return &SystemService{
		db:                  db,
		uploadManager:       uploadManager,
		emailClient:         emailClient,
		feishuClient:        feishuClient,
		notificationService: notificationService,
		httpClient:          &http.Client{Timeout: 10 * time.Second},
	}
}

// GetStaticInfo 获取系统静态信息
func (s *SystemService) GetStaticInfo() *dto.SystemStaticInfo {
	info := &dto.SystemStaticInfo{
		CPUCore: runtime.NumCPU(),
		CPUArch: runtime.GOARCH,
		OS:      runtime.GOOS,
		DbType:  s.getDBType(),
	}

	if ci, err := cpu.Info(); err == nil && len(ci) > 0 {
		info.CPUModel = ci[0].ModelName
	}

	info.Hostname, _ = os.Hostname()
	info.Timezone = time.Now().Location().String()
	info.ServerIP = getServerIP()

	if m, err := mem.VirtualMemory(); err == nil {
		info.MemoryTotal = m.Total
	}
	if swap, err := mem.SwapMemory(); err == nil {
		info.SwapTotal = swap.Total
	}
	if usage, err := disk.Usage(systemDiskPath()); err == nil {
		info.DiskTotal = usage.Total
	}

	info.DbTables = s.getTableCount()
	info.StorageStatus = s.checkStorage()
	info.EmailStatus = s.checkEmail()
	info.FeishuStatus = s.checkFeishu()
	info.AppVersion = currentVersion()

	return info
}

// GetDynamicInfo 获取系统动态信息
func (s *SystemService) GetDynamicInfo() *dto.SystemDynamicInfo {
	info := &dto.SystemDynamicInfo{}

	s.setDynamicCPU(info)
	s.setDynamicMemory(info)
	s.setDynamicHost(info)
	s.setDynamicDisk(info)
	s.setDynamicDB(info)

	status := s.GetVersionStatus()
	info.VersionLatestVersion = status.LatestVersion
	info.VersionLastCheckError = status.LastCheckError

	return info
}

// GetSystemStatus 获取飞书系统状态
func (s *SystemService) GetSystemStatus(_ context.Context) (*feishupkg.SystemStatus, error) {
	status := &feishupkg.SystemStatus{
		DBStatus:      s.checkDB(),
		StorageStatus: s.checkStorage(),
		EmailStatus:   s.checkEmail(),
		FeishuStatus:  s.checkFeishu(),
	}

	if p, err := cpu.Percent(time.Second, false); err == nil && len(p) > 0 {
		status.CPUUsage = p[0]
	}
	if m, err := mem.VirtualMemory(); err == nil {
		status.MemoryTotal = m.Total
		status.MemoryUsed = m.Used
	}
	if usage, err := disk.Usage(systemDiskPath()); err == nil {
		status.DiskTotal = usage.Total
		status.DiskUsed = usage.Used
	}

	return status, nil
}

// CheckForUpdates 检查是否有新版本
func (s *SystemService) CheckForUpdates() error {
	manifest, err := s.fetchManifest(context.Background())
	if err != nil {
		s.setVersionStatus(VersionStatus{LastCheckError: err.Error()})
		return err
	}

	latestVersion := strings.TrimSpace(manifest.LatestVersion)
	if latestVersion == "" {
		err = fmt.Errorf("版本清单缺少 latest_version")
		s.setVersionStatus(VersionStatus{LastCheckError: err.Error()})
		return err
	}

	status := VersionStatus{
		LatestVersion:  latestVersion,
		LastCheckError: "",
	}

	currentVersion := currentVersion()
	if strings.EqualFold(currentVersion, "dev") {
		s.setVersionStatus(status)
		return nil
	}

	compareResult, err := compareVersion(latestVersion, currentVersion)
	if err != nil {
		err = fmt.Errorf("比较版本失败: %w", err)
		status.LastCheckError = err.Error()
		s.setVersionStatus(status)
		return err
	}

	if compareResult > 0 {
		exists, err := s.notificationService.HasVersionUpdateNotification(context.Background(), latestVersion)
		if err != nil {
			status.LastCheckError = err.Error()
			s.setVersionStatus(status)
			return err
		}

		if !exists {
			if err := s.notificationService.NotifyVersionUpdateToSuperAdmins(context.Background(), currentVersion, latestVersion, strings.TrimSpace(manifest.ReleaseURL)); err != nil {
				status.LastCheckError = err.Error()
				s.setVersionStatus(status)
				return err
			}
		}
	}

	s.setVersionStatus(status)
	return nil
}

// GetVersionStatus 获取最近一次版本检测状态
func (s *SystemService) GetVersionStatus() VersionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.versionStatus
}

// setVersionStatus 更新最近一次版本检测状态
func (s *SystemService) setVersionStatus(status VersionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.versionStatus = status
}

// fetchManifest 获取最新版本信息
func (s *SystemService) fetchManifest(ctx context.Context) (*versionManifest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubLatestReleaseAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "FlecBlog-VersionChecker")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("请求版本信息失败: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析版本信息失败: %w", err)
	}

	return &versionManifest{
		LatestVersion: strings.TrimSpace(release.TagName),
		ReleaseURL:    strings.TrimSpace(release.HTMLURL),
	}, nil
}

// setDynamicCPU 填充 CPU 动态信息
func (s *SystemService) setDynamicCPU(info *dto.SystemDynamicInfo) {
	if p, err := cpu.Percent(time.Second, false); err == nil && len(p) > 0 {
		info.CPUUsage = p[0]
	}
	if l, err := load.Avg(); err == nil {
		info.Load1 = l.Load1
		info.Load5 = l.Load5
		info.Load15 = l.Load15
	}
}

// setDynamicMemory 填充内存动态信息
func (s *SystemService) setDynamicMemory(info *dto.SystemDynamicInfo) {
	if m, err := mem.VirtualMemory(); err == nil {
		info.MemoryUsed = m.Used
		info.MemoryAvailable = m.Available
	}
	if swap, err := mem.SwapMemory(); err == nil {
		info.SwapUsed = swap.Used
	}
}

// setDynamicHost 填充主机动态信息
func (s *SystemService) setDynamicHost(info *dto.SystemDynamicInfo) {
	if hi, err := host.Info(); err == nil {
		// #nosec G115 - Uptime 是从系统获取的正常运行时间，值在合理范围内
		info.HostUptime = int64(hi.Uptime)
	}
}

// setDynamicDisk 填充磁盘动态信息
func (s *SystemService) setDynamicDisk(info *dto.SystemDynamicInfo) {
	if usage, err := disk.Usage(systemDiskPath()); err == nil {
		info.DiskUsed = usage.Used
		info.DiskFree = usage.Free
	}
}

// setDynamicDB 填充数据库动态信息
func (s *SystemService) setDynamicDB(info *dto.SystemDynamicInfo) {
	info.DbStatus = s.checkDB()
	info.DbSize = s.getDBSize()
	info.DbConnCount = s.getConnCount()
}

// checkDB 检查数据库状态
func (s *SystemService) checkDB() string {
	db, err := s.db.DB()
	if err != nil || db.Ping() != nil {
		return "连接失败"
	}
	return "正常"
}

// getDBType 获取数据库类型
func (s *SystemService) getDBType() string {
	return s.db.Name()
}

// checkStorage 检查存储状态
func (s *SystemService) checkStorage() string {
	if s.uploadManager == nil {
		return "未配置"
	}
	if err := s.uploadManager.HealthCheck(); err != nil {
		return "异常"
	}
	return "正常"
}

// checkEmail 检查邮件状态
func (s *SystemService) checkEmail() string {
	if s.emailClient == nil {
		return "未配置"
	}
	if err := s.emailClient.HealthCheck(); err != nil {
		return "异常"
	}
	return "正常"
}

// checkFeishu 检查飞书状态
func (s *SystemService) checkFeishu() string {
	if s.feishuClient == nil {
		return "未配置"
	}
	if err := s.feishuClient.HealthCheck(); err != nil {
		return "异常"
	}
	return "正常"
}

// getDBSize 获取数据库大小
func (s *SystemService) getDBSize() int64 {
	var name string
	if err := s.db.Raw("SELECT current_database()").Scan(&name).Error; err != nil || name == "" {
		return 0
	}
	var size int64
	s.db.Raw(fmt.Sprintf("SELECT pg_database_size('%s')", name)).Scan(&size)
	return size
}

// getTableCount 获取数据库表数量
func (s *SystemService) getTableCount() int64 {
	var name string
	var count int64
	if err := s.db.Raw("SELECT current_database()").Scan(&name).Error; err != nil || name == "" {
		return 0
	}
	s.db.Raw(fmt.Sprintf("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_catalog = '%s'", name)).Scan(&count)
	return count
}

// getConnCount 获取数据库连接数
func (s *SystemService) getConnCount() int {
	var name string
	var count int
	if err := s.db.Raw("SELECT current_database()").Scan(&name).Error; err != nil || name == "" {
		return 0
	}
	s.db.Raw(fmt.Sprintf("SELECT count(*) FROM pg_stat_activity WHERE datname = '%s'", name)).Scan(&count)
	return count
}

// currentVersion 获取当前运行版本
func currentVersion() string {
	return strings.TrimSpace(AppVersion)
}

// compareVersion 比较两个语义化版本
func compareVersion(a, b string) (int, error) {
	av, err := parseVersion(a)
	if err != nil {
		return 0, err
	}

	bv, err := parseVersion(b)
	if err != nil {
		return 0, err
	}

	for i := range av.numbers {
		switch {
		case av.numbers[i] > bv.numbers[i]:
			return 1, nil
		case av.numbers[i] < bv.numbers[i]:
			return -1, nil
		}
	}

	switch {
	case len(av.prerelease) == 0 && len(bv.prerelease) > 0:
		return 1, nil
	case len(av.prerelease) > 0 && len(bv.prerelease) == 0:
		return -1, nil
	default:
		return comparePrerelease(av.prerelease, bv.prerelease), nil
	}
}

// parseVersion 解析语义化版本号
func parseVersion(raw string) (parsedVersion, error) {
	value := strings.TrimSpace(strings.TrimPrefix(raw, "v"))
	if value == "" {
		return parsedVersion{}, fmt.Errorf("版本号不能为空")
	}

	core, _, _ := strings.Cut(value, "+")
	mainPart, prereleasePart, hasPrerelease := strings.Cut(core, "-")

	parts := strings.Split(mainPart, ".")
	if len(parts) != 3 {
		return parsedVersion{}, fmt.Errorf("版本号格式无效: %s", raw)
	}

	parsed := parsedVersion{}
	for i, part := range parts {
		number, err := strconv.Atoi(part)
		if err != nil || number < 0 {
			return parsedVersion{}, fmt.Errorf("版本号格式无效: %s", raw)
		}
		parsed.numbers[i] = number
	}

	if hasPrerelease {
		identifiers := strings.Split(prereleasePart, ".")
		for _, identifier := range identifiers {
			if strings.TrimSpace(identifier) == "" {
				return parsedVersion{}, fmt.Errorf("版本号格式无效: %s", raw)
			}
		}
		parsed.prerelease = identifiers
	}

	return parsed, nil
}

// comparePrerelease 比较预发布版本标识
func comparePrerelease(a, b []string) int {
	for i := 0; i < len(a) || i < len(b); i++ {
		switch {
		case i >= len(a):
			return -1
		case i >= len(b):
			return 1
		}

		if a[i] == b[i] {
			continue
		}

		return comparePrereleaseIdentifier(a[i], b[i])
	}

	return 0
}

// comparePrereleaseIdentifier 比较单个预发布标识
func comparePrereleaseIdentifier(a, b string) int {
	av, aErr := strconv.Atoi(a)
	bv, bErr := strconv.Atoi(b)

	switch {
	case aErr == nil && bErr == nil:
		switch {
		case av > bv:
			return 1
		case av < bv:
			return -1
		default:
			return 0
		}
	case aErr == nil:
		return -1
	case bErr == nil:
		return 1
	case a > b:
		return 1
	default:
		return -1
	}
}

// systemDiskPath 获取系统磁盘路径
func systemDiskPath() string {
	if runtime.GOOS == "windows" {
		return "C:"
	}
	return "/"
}

// getServerIP 获取服务器 IP
func getServerIP() string {
	if ifs, err := net.Interfaces(); err == nil {
		for _, iface := range ifs {
			for _, addr := range iface.Addrs {
				ip := addr.Addr
				if idx := strings.IndexByte(ip, '/'); idx > 0 {
					ip = ip[:idx]
				}
				if len(ip) > 0 && ip != "127.0.0.1" && ip[0] != ':' && ip[0] >= '1' && ip[0] <= '9' {
					return ip
				}
			}
		}
	}
	return "N/A"
}
