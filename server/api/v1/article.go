package v1

import (
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"
	"flec_blog/pkg/upload"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService *service.ArticleService
}

// NewArticleController 创建文章控制器
func NewArticleController(articleService *service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

// ============ 前台接口 ============

// ListForWeb 文章列表
//
//	@Summary		获取文章列表
//	@Description	获取已发布文章，置顶文章在前，支持按年/月/分类/标签筛选
//	@Tags			文章
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页数量（不传则返回全部）"
//	@Param			year		query		string	false	"年份，如 2025"
//	@Param			month		query		string	false	"月份 1-12，需配合 year"
//	@Param			category	query		string	false	"分类 slug"
//	@Param			tag			query		string	false	"标签 slug"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Failure		400			{object}	response.Response
//	@Router			/articles [get]
func (c *ArticleController) ListForWeb(ctx *gin.Context) {
	var req dto.ListArticlesForWebRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	articles, total, err := c.articleService.ListForWeb(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, articles, total, req.Page, req.PageSize)
}

// Search 搜索文章
//
//	@Summary		搜索文章
//	@Description	全文搜索标题和正文
//	@Tags			文章
//	@Accept			json
//	@Produce		json
//	@Param			keyword		query		string	true	"搜索词"
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Failure		400			{object}	response.Response
//	@Router			/articles/search [get]
func (c *ArticleController) Search(ctx *gin.Context) {
	var req dto.SearchArticlesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	articles, total, err := c.articleService.Search(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, articles, total, req.Page, req.PageSize)
}

// GetBySlug 文章详情
//
//	@Summary		文章详情
//	@Description	通过 slug 读取文章完整内容
//	@Tags			文章
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"文章 slug"
//	@Success		200		{object}	response.Response{data=dto.ArticleDetailResponse}
//	@Failure		404		{object}	response.Response
//	@Router			/articles/{slug} [get]
func (c *ArticleController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.ValidateFailed(ctx, "slug不能为空")
		return
	}

	article, err := c.articleService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, article)
}

// ============ 后台管理接口 ============

// List 获取文章列表
//
//	@Summary		文章列表
//	@Description	获取包含草稿的所有文章，支持多种筛选条件
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page			query		int		false	"页码"
//	@Param			page_size		query		int		false	"每页数量（不传则返回全部）"
//	@Param			keyword			query		string	false	"搜索关键词（标题/内容）"
//	@Param			category_id		query		uint	false	"分类ID"
//	@Param			tag_ids			query		[]uint	false	"标签ID列表"
//	@Param			is_publish		query		bool	false	"是否发布"
//	@Param			is_top			query		bool	false	"是否置顶"
//	@Param			is_essence		query		bool	false	"是否精选"
//	@Param			is_outdated		query		bool	false	"是否过时"
//	@Param			start_time		query		string	false	"发布开始时间（格式：2006-01-02）"
//	@Param			end_time		query		string	false	"发布结束时间（格式：2006-01-02）"
//	@Success		200				{object}	response.Response{data=response.PageResult}
//	@Failure		401				{object}	response.Response
//	@Failure		403				{object}	response.Response
//	@Router			/admin/articles [get]
func (c *ArticleController) List(ctx *gin.Context) {
	var req dto.ListArticlesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	articles, total, err := c.articleService.List(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, articles, total, req.Page, req.PageSize)
}

// Get 获取文章详情
//
//	@Summary		文章详情
//	@Description	文章详细信息，用于编辑器回显
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章 ID"
//	@Success		200	{object}	response.Response{data=dto.ArticleAdminDetailResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/articles/{id} [get]
func (c *ArticleController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	article, err := c.articleService.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, article)
}

// Create 创建文章
//
//	@Summary		新建文章
//	@Description	创建草稿或发布文章，支持设置文章各种信息，发布时自动设置发布时间
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateArticleRequest	true	"文章信息"
//	@Success		201		{object}	response.Response{data=dto.ArticleAdminDetailResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/articles [post]
func (c *ArticleController) Create(ctx *gin.Context) {
	var req dto.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	article, err := c.articleService.Create(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, article)
}

// Update 更新文章
//
//	@Summary		更新文章
//	@Description	修改文章各种信息，支持调整发布时间
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"文章 ID"
//	@Param			request	body		dto.UpdateArticleRequest	true	"文章信息"
//	@Success		200		{object}	response.Response{data=dto.ArticleAdminDetailResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/admin/articles/{id} [put]
func (c *ArticleController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	updatedArticle, err := c.articleService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, updatedArticle)
}

// Delete 删除文章
//
//	@Summary		删除文章
//	@Description	硬删除文章，不可恢复
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/articles/{id} [delete]
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.articleService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ============ 数据导入导出接口 ============

// ImportArticles 导入文章数据
//
//	@Summary		导入文章数据
//	@Description	从Hexo或Markdown格式导入文章数据，支持多文件
//	@Tags			文章管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			source_type	formData	string	true	"来源类型，支持：hexo, markdown"	Enums(hexo, markdown)
//	@Param			upload_images formData	bool	false	"是否上传文章中的图片（默认false）"
//	@Param			files		formData	[]file	true	"文章文件（.md或.markdown格式，支持多文件）"
//	@Success		200			{object}	response.Response{data=dto.ImportArticlesResult}
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/articles/import [post]
func (c *ArticleController) ImportArticles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		response.ValidateFailed(ctx, "文件上传失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.ValidateFailed(ctx, "请选择要导入的文件")
		return
	}

	if len(files) > 100 {
		response.ValidateFailed(ctx, "单次最多导入100个文件")
		return
	}

	sourceType := ctx.PostForm("source_type")
	if sourceType != "hexo" && sourceType != "markdown" {
		response.ValidateFailed(ctx, "不支持的来源类型，仅支持 hexo 和 markdown")
		return
	}

	uploadImages := ctx.PostForm("upload_images") == "true"

	// 提取 host 用于生成完整的图片 URL
	host := upload.ExtractHostFromContext(ctx)

	// Markdown/Hexo 文件导入
	const maxFileSize = 10 << 20 // 10MB
	fileContents := make(map[string]string)

	for _, fileHeader := range files {
		if !isMarkdownFile(fileHeader.Filename) {
			response.ValidateFailed(ctx, fmt.Sprintf("文件 %s 不是Markdown格式", fileHeader.Filename))
			return
		}

		fileBytes, err := readUploadFile(fileHeader, maxFileSize)
		if err != nil {
			response.Failed(ctx, err.Error())
			return
		}
		fileContents[fileHeader.Filename] = string(fileBytes)
	}

	var result *dto.ImportArticlesResult
	result, err = c.articleService.ImportArticles(ctx.Request.Context(), fileContents, sourceType, uploadImages, host)

	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}

// readUploadFile 读取上传文件内容
func readUploadFile(fileHeader *multipart.FileHeader, maxSize int64) ([]byte, error) {
	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("文件 %s 超过大小限制", fileHeader.Filename)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	return io.ReadAll(file)
}

// isMarkdownFile 检查文件是否为Markdown格式
func isMarkdownFile(filename string) bool {
	return len(filename) > 3 && (filename[len(filename)-3:] == ".md" || filename[len(filename)-9:] == ".markdown")
}

// ============ 微信公众号导出接口 ============

// ExportToWeChat 将文章渲染为微信公众号格式
//
//	@Summary		生成微信公众号 HTML
//	@Description	将文章 Markdown 转换为微信公众号 HTML 格式，供复制粘贴到微信公众平台
//	@Tags			文章管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int	true	"文章 ID"
//	@Success		200		{object}	response.Response{data=dto.WeChatExportResult}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Router			/admin/articles/{id}/wechat/export [post]
func (c *ArticleController) ExportToWeChat(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.ValidateFailed(ctx, "无效的文章 ID")
		return
	}

	result := c.articleService.ExportToWeChat(ctx.Request.Context(), uint(id))
	response.Success(ctx, result)
}

// DownloadZip 下载文章为压缩包
//
//	@Summary		下载为 Markdown
//	@Description	下载文章为压缩包，包含 Markdown 文件、配图、封面图等资源
//	@Tags			文章管理
//	@Accept			json
//	@Produce		application/zip
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文章 ID"
//	@Success		200	{file}		byte
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/articles/{id}/download/zip [get]
func (c *ArticleController) DownloadZip(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.ValidateFailed(ctx, "无效的文章 ID")
		return
	}

	data, filename, err := c.articleService.DownloadZip(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Data(200, "application/zip", data)
}
