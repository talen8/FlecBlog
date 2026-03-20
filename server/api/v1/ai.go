package v1

import (
	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/ai"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// AIController AI功能控制器
type AIController struct {
	settingService *service.SettingService
}

// NewAIController 创建AI控制器
func NewAIController(settingService *service.SettingService) *AIController {
	return &AIController{
		settingService: settingService,
	}
}

// getAIProvider 获取AI服务提供商
func (c *AIController) getAIProvider() (ai.Provider, error) {
	cfg, err := c.settingService.GetAIConfig()
	if err != nil {
		return nil, err
	}

	provider, err := ai.GetProvider(cfg)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// TestConfig 测试AI配置是否可用
//
//	@Summary		测试AI配置
//	@Description	使用提供的配置测试AI服务连通性
//	@Tags			AI功能
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AITestRequest	true	"AI配置"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/admin/ai/test [post]
func (c *AIController) TestConfig(ctx *gin.Context) {
	var req dto.AITestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	provider := ai.NewOpenAIClient(req.BaseURL, req.APIKey, req.Model)
	if err := provider.Test(); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// Summary 生成文章摘要
//
//	@Summary		生成文章摘要
//	@Description	基于文章内容自动生成摘要（50-100字，创作者角度）
//	@Tags			AI功能
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AISummaryRequest	true	"文章内容"
//	@Success		200		{object}	response.Response{data=dto.AISummaryResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/admin/ai/summary [post]
func (c *AIController) Summary(ctx *gin.Context) {
	var req dto.AISummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	provider, err := c.getAIProvider()
	if err != nil {
		response.Failed(ctx, "AI配置未设置或配置错误")
		return
	}

	summary, err := provider.GenerateSummary(req.Content)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, dto.AISummaryResponse{Summary: summary})
}

// AISummary 生成AI摘要
//
//	@Summary		生成AI摘要
//	@Description	基于文章内容生成AI摘要（150-200字，旁观者角度）
//	@Tags			AI功能
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AIAISummaryRequest	true	"文章内容"
//	@Success		200		{object}	response.Response{data=dto.AIAISummaryResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/admin/ai/ai-summary [post]
func (c *AIController) AISummary(ctx *gin.Context) {
	var req dto.AIAISummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	provider, err := c.getAIProvider()
	if err != nil {
		response.Failed(ctx, "AI配置未设置或配置错误")
		return
	}

	aiSummary, err := provider.GenerateAISummary(req.Content)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, dto.AIAISummaryResponse{Summary: aiSummary})
}

// Title 生成标题建议
//
//	@Summary		生成标题建议
//	@Description	根据内容生成多个标题建议
//	@Tags			AI功能
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.AITitleRequest	true	"文章内容"
//	@Success		200		{object}	response.Response{data=dto.AITitleResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/admin/ai/title [post]
func (c *AIController) Title(ctx *gin.Context) {
	var req dto.AITitleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	provider, err := c.getAIProvider()
	if err != nil {
		response.Failed(ctx, "AI配置未设置或配置错误")
		return
	}

	titles, err := provider.GenerateTitle(req.Content)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, dto.AITitleResponse{Title: titles[0]})
}
