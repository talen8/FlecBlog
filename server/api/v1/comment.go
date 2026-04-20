package v1

import (
	"io"
	"strconv"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *service.CommentService
}

// NewCommentController 创建评论控制器
func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// ============ 前台接口 ============

// ListForWeb 获取评论列表
//
//	@Summary		评论列表
//	@Description	获取目标评论，扁平化显示所有评论和回复。需指定 target_type 和 target_key
//	@Tags			评论
//	@Accept			json
//	@Produce		json
//	@Param			target_type	query		string	true	"目标类型(article/page)"
//	@Param			target_key	query		string	true	"目标标识(文章slug或页面key)"
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页数量（不传则返回全部）"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Failure		400			{object}	response.Response
//	@Router			/comments [get]
func (c *CommentController) ListForWeb(ctx *gin.Context) {
	var req dto.CommentQueryForWebRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	comments, total, err := c.commentService.ListForWeb(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, comments, total, req.Page, req.PageSize)
}

// Create 创建评论
//
//	@Summary		发表评论
//	@Description	发表评论或回复。已登录用户直接评论，未登录用户需提供昵称、邮箱等游客信息
//	@Tags			评论
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateCommentRequest	true	"评论信息"
//	@Success		201		{object}	response.Response{data=dto.CommentResponse}
//	@Failure		400		{object}	response.Response
//	@Router			/comments [post]
func (c *CommentController) Create(ctx *gin.Context) {
	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	// 获取客户端信息
	req.IP = ctx.ClientIP()
	req.UserAgent = ctx.GetHeader("User-Agent")

	// 获取用户ID（登录用户有值，游客为0）
	userID, _ := ctx.Get("user_id")
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	// 游客评论需要验证游客信息
	if uid == 0 && (req.Email == "" || req.Nickname == "") {
		response.ValidateFailed(ctx, "游客评论需要提供昵称和邮箱")
		return
	}

	comment, err := c.commentService.Create(ctx.Request.Context(), &req, uid)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, comment)
}

// Update 更新评论
//
//	@Summary		修改评论
//	@Description	只能修改自己的评论内容
//	@Tags			评论
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"评论 ID"
//	@Param			request	body		dto.UpdateCommentRequest	true	"评论信息"
//	@Success		200		{object}	response.Response{data=dto.CommentResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/comments/{id} [put]
func (c *CommentController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	var req dto.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	comment, err := c.commentService.Update(ctx.Request.Context(), uint(id), &req, userID)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, comment)
}

// DeleteForWeb 删除评论（前台）
//
//	@Summary		删除评论
//	@Description	只能删除自己的评论，子评论会保留
//	@Tags			评论
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/comments/{id} [delete]
func (c *CommentController) DeleteForWeb(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	if err := c.commentService.DeleteForWeb(ctx.Request.Context(), uint(id), userID); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// ============ 后台管理接口 ============

// List 获取评论列表
//
//	@Summary		评论列表
//	@Description	获取所有评论，支持按状态筛选
//	@Tags			评论管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页数量（不传则返回全部）"
//	@Param			status		query		int	false	"状态筛选 0:隐藏 1:显示"
//	@Success		200			{object}	response.Response{data=response.PageResult}
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/comments [get]
func (c *CommentController) List(ctx *gin.Context) {
	var req dto.CommentQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	comments, total, err := c.commentService.List(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.PageSuccess(ctx, comments, total, req.Page, req.PageSize)
}

// Get 获取评论详情
//
//	@Summary		评论详情
//	@Description	查看评论详细信息
//	@Tags			评论管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	response.Response{data=dto.CommentListResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/comments/{id} [get]
func (c *CommentController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	comment, err := c.commentService.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, comment)
}

// ToggleStatus 切换评论状态
//
//	@Summary		切换评论状态
//	@Description	切换评论的显示状态，隐藏后博客端不可见或被替换
//	@Tags			评论管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/comments/{id}/toggle-status [put]
func (c *CommentController) ToggleStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.commentService.ToggleStatus(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// Delete 删除评论
//
//	@Summary		删除评论
//	@Description	硬删除评论，将从数据库中彻底删除，不可恢复
//	@Tags			评论管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"评论 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/comments/{id} [delete]
func (c *CommentController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	if err := c.commentService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// CreateForAdmin 创建评论
//
//	@Summary		创建评论
//	@Description	管理员创建评论，用于回复用户，但无法设置评论IP等数据
//	@Tags			评论管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateCommentRequest	true	"评论信息"
//	@Success		201	{object}	response.Response{data=dto.CommentResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/admin/comments [post]
func (c *CommentController) CreateForAdmin(ctx *gin.Context) {
	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	// 管理员创建评论，使用管理员的用户ID
	userID := ctx.GetUint("user_id")

	comment, err := c.commentService.Create(ctx.Request.Context(), &req, userID)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Created(ctx, comment)
}

// ============ 数据导入导出接口 ============

// ImportComments 导入评论数据
//
//	@Summary		导入评论数据
//	@Description	从Artalk评论系统导入评论数据，暂不支持更多评论系统
//	@Tags			评论管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			file			formData	file	true	"导入文件（JSON或Artrans格式）"
//	@Param			source_type		formData	string	true	"来源类型，目前支持：artalk"	Enums(artalk)
//	@Success		200				{object}	response.Response{data=dto.ImportCommentsResult}
//	@Failure		400				{object}	response.Response
//	@Failure		401				{object}	response.Response
//	@Failure		403				{object}	response.Response
//	@Router			/admin/comments/import [post]
func (c *CommentController) ImportComments(ctx *gin.Context) {
	// 解析表单数据
	var req dto.ImportCommentsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		response.ValidateFailed(ctx, "请选择要导入的文件")
		return
	}
	defer func() {
		_ = file.Close()
	}()

	// 检查文件大小（10MB限制）
	const maxFileSize = 10 << 20
	if header.Size > maxFileSize {
		response.ValidateFailed(ctx, "文件大小不能超过10MB")
		return
	}

	// 检查文件类型（支持 .json 和 .artrans 格式）
	filename := header.Filename
	if len(filename) <= 5 || (filename[len(filename)-5:] != ".json" && filename[len(filename)-8:] != ".artrans") {
		response.ValidateFailed(ctx, "只支持JSON格式的文件")
		return
	}

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.Failed(ctx, "文件读取失败")
		return
	}

	// 根据来源类型调用对应的导入方法
	var result *dto.ImportCommentsResult
	switch req.SourceType {
	case "artalk":
		result, err = c.commentService.ImportFromArtalk(ctx.Request.Context(), fileBytes)
	default:
		response.ValidateFailed(ctx, "暂不支持该导入类型")
		return
	}

	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, result)
}
