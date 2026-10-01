package v1

import (
	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// ThemeHandler 主题配置处理器
type ThemeHandler struct {
	themeService *service.ThemeService
}

// NewThemeHandler 创建主题配置处理器
func NewThemeHandler(themeService *service.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

// GetForWeb 获取前台主题配置
//
//	@Summary		获取主题配置
//	@Description	获取前台主题配置，菜单仅返回启用项
//	@Tags			主题
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.ThemeConfigResponse}
//	@Router			/themes [get]
func (h *ThemeHandler) GetForWeb(ctx *gin.Context) {
	result, err := h.themeService.GetThemeConfigForWeb()
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// Get 获取后台主题配置
//
//	@Summary		主题配置详情
//	@Description	获取后台主题配置，返回全部配置项与完整菜单
//	@Tags			主题管理
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=dto.ThemeConfigResponse}
//	@Router			/admin/themes [get]
func (h *ThemeHandler) Get(ctx *gin.Context) {
	result, err := h.themeService.GetThemeConfig()
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// UpdateConfig 更新主题配置
//
//	@Summary		更新主题配置
//	@Description	按提交的配置项逐个覆盖（patch 语义），菜单需通过菜单接口更新
//	@Tags			主题管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	dto.ConfigUpdateRequest	true	"主题配置"
//	@Success		200		{object}	response.Response
//	@Router			/admin/themes/config [put]
func (h *ThemeHandler) UpdateConfig(ctx *gin.Context) {
	var req dto.ConfigUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := h.themeService.UpdateConfig(&req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// UpdateMenus 更新主题菜单
//
//	@Summary		更新主题菜单
//	@Description	整体替换主题菜单
//	@Tags			主题管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	dto.MenuUpdateRequest	true	"主题菜单"
//	@Success		200		{object}	response.Response
//	@Router			/admin/themes/menus [put]
func (h *ThemeHandler) UpdateMenus(ctx *gin.Context) {
	var req dto.MenuUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := h.themeService.UpdateMenus(&req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}
