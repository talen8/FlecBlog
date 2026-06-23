package v1

import (
	"flec_blog/internal/model"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"
	"flec_blog/pkg/upload"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingController 配置控制器
type SettingController struct {
	settingService *service.SettingService
	db             *gorm.DB
	uploadManager  *upload.Manager // 上传管理器（用于热重载上传配置）
}

// NewSettingController 创建配置控制器
func NewSettingController(settingService *service.SettingService, db *gorm.DB, uploadManager *upload.Manager) *SettingController {
	return &SettingController{
		settingService: settingService,
		db:             db,
		uploadManager:  uploadManager,
	}
}

// ============ 后台管理接口 ============

// GetGroup 获取某个分组的所有配置
//
//	@Summary		获取配置分组
//	@Description	获取指定分组的所有配置项
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			group	path		string	true	"配置分组"	Enums(basic, notification, upload, ai, oauth)
//	@Success		200		{object}	response.Response{data=dto.SettingGroupResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/admin/settings/{group} [get]
func (c *SettingController) GetGroup(ctx *gin.Context) {
	group := ctx.Param("group")
	if group == "" {
		response.ValidateFailed(ctx, "配置分组不能为空")
		return
	}

	// 使用服务层方法获取所有配置项
	settingsMap, err := c.settingService.GetByGroup(group)
	if err != nil {
		response.Failed(ctx, "获取配置失败: "+err.Error())
		return
	}

	response.Success(ctx, settingsMap)
}

// UpdateGroup 更新某个分组的配置
//
//	@Summary		更新配置分组
//	@Description	批量更新指定分组的配置项，限超级管理员
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			group	path		string							true	"配置分组"	Enums(basic, notification, upload, ai, oauth)
//	@Param			request	body		dto.UpdateSettingGroupRequest	true	"配置更新内容"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/admin/settings/{group} [patch]
func (c *SettingController) UpdateGroup(ctx *gin.Context) {
	group := ctx.Param("group")
	if group == "" {
		response.ValidateFailed(ctx, "配置分组不能为空")
		return
	}

	var settings map[string]string
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if len(settings) == 0 {
		response.ValidateFailed(ctx, "配置内容不能为空")
		return
	}

	// 调用 service 更新配置
	if err := c.settingService.UpdateGroup(group, settings); err != nil {
		response.Failed(ctx, "更新配置失败: "+err.Error())
		return
	}

	// 上传配置更新后重新加载存储实例
	if group == model.SettingGroupUpload && c.uploadManager != nil {
		_ = c.uploadManager.ReloadStorage()
	}

	response.Success(ctx, nil, "配置更新成功")
}

// ResetMCPSecret 重置 MCP Secret
//
//	@Summary		重置 MCP Secret
//	@Description	重新生成并返回 MCP Secret，限超级管理员
//	@Tags			配置管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=map[string]string}
//	@Failure		401	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/admin/settings/ai/mcp-secret/reset [put]
func (c *SettingController) ResetMCPSecret(ctx *gin.Context) {
	secret, err := c.settingService.ResetMCPSecret()
	if err != nil {
		response.Failed(ctx, "重置 MCP Secret 失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{"secret": secret}, "MCP Secret 重置成功")
}

// ============ 前台公开接口 ============

// GetPublicSettingGroup 获取公开的配置分组
//
//	@Summary		获取配置分组
//	@Description	获取指定分组的公开配置项，根据配置项可见性规则筛选
//	@Tags			配置
//	@Accept			json
//	@Produce		json
//	@Param			group	path		string	true	"配置分组"	Enums(basic, oauth, upload)
//	@Success		200	{object}	response.Response{data=map[string]string}
//	@Failure		400	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/settings/{group} [get]
func (c *SettingController) GetPublicSettingGroup(ctx *gin.Context) {
	group := ctx.Param("group")
	if group == "" {
		response.ValidateFailed(ctx, "配置分组不能为空")
		return
	}

	// 使用通用方法获取该分组的公开配置
	settings, err := c.settingService.GetByGroup(group, true)
	if err != nil {
		response.Failed(ctx, "获取配置失败: "+err.Error())
		return
	}

	response.Success(ctx, settings)
}
