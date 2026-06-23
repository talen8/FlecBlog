package v1

import (
	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// ThemeHandler 主题处理器
type ThemeHandler struct {
	themeService *service.ThemeService
}

// NewThemeHandler 创建主题处理器
func NewThemeHandler(themeService *service.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

// GetActive 获取前台激活主题
//
//	@Summary		获取激活主题
//	@Description	获取当前激活主题详情，菜单仅返回启用项
//	@Tags			主题
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.ThemePublicResponse}
//	@Router			/themes [get]
func (h *ThemeHandler) GetActive(ctx *gin.Context) {
	result, err := h.themeService.GetActiveTheme()
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// SyncMeta 同步主题元数据
//
//	@Summary		同步主题元数据
//	@Description	主题镜像启动时上报元数据并激活当前主题
//	@Tags			主题
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.ThemeMetaSyncRequest	true	"主题元数据"
//	@Success		200		{object}	response.Response
//	@Router			/themes/_sync [post]
func (h *ThemeHandler) SyncMeta(ctx *gin.Context) {
	var req dto.ThemeMetaSyncRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := h.themeService.SyncThemeMeta(&req); err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// List 获取后台主题列表
//
//	@Summary		主题列表
//	@Description	获取全部主题实例
//	@Tags			主题管理
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]dto.ThemeResponse}
//	@Router			/admin/themes [get]
func (h *ThemeHandler) List(ctx *gin.Context) {
	result, err := h.themeService.ListThemes()
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// Get 获取后台主题详情
//
//	@Summary		主题详情
//	@Description	获取指定主题详情
//	@Tags			主题管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			slug	path	string	true	"主题 slug"
//	@Success		200		{object}	response.Response{data=dto.ThemeResponse}
//	@Router			/admin/themes/{slug} [get]
func (h *ThemeHandler) Get(ctx *gin.Context) {
	result, err := h.themeService.GetTheme(ctx.Param("slug"))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// UpdateConfig 更新主题配置
//
//	@Summary		更新主题配置
//	@Description	替换指定主题的 config JSON
//	@Tags			主题管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			slug	path	string					true	"主题 slug"
//	@Param			request	body	dto.ConfigUpdateRequest	true	"主题配置"
//	@Success		200		{object}	response.Response
//	@Router			/admin/themes/{slug}/config [put]
func (h *ThemeHandler) UpdateConfig(ctx *gin.Context) {
	var req dto.ConfigUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := h.themeService.UpdateConfig(ctx.Param("slug"), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// CheckUpdate 检查主题版本更新
//
//	@Summary		检查主题版本更新
//	@Description	检查指定主题是否有新版本
//	@Tags			主题管理
//	@Produce		json
//	@Security		BearerAuth
//	@Param			slug	path	string	true	"主题 slug"
//	@Success		200		{object}	response.Response{data=dto.ThemeUpdateCheckResponse}
//	@Router			/admin/themes/{slug}/check [post]
func (h *ThemeHandler) CheckUpdate(ctx *gin.Context) {
	result, err := h.themeService.CheckThemeUpdate(ctx.Request.Context(), ctx.Param("slug"))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}

// UpdateMenus 更新主题菜单
//
//	@Summary		更新主题菜单
//	@Description	整体替换指定主题菜单
//	@Tags			主题管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			slug	path	string					true	"主题 slug"
//	@Param			request	body	dto.MenuUpdateRequest	true	"主题菜单"
//	@Success		200		{object}	response.Response
//	@Router			/admin/themes/{slug}/menus [put]
func (h *ThemeHandler) UpdateMenus(ctx *gin.Context) {
	var req dto.MenuUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	result, err := h.themeService.UpdateMenus(ctx.Param("slug"), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}
	response.Success(ctx, result)
}
