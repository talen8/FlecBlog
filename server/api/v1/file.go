package v1

import (
	"strconv"
	"strings"

	"flec_blog/internal/dto"
	"flec_blog/internal/service"
	"flec_blog/pkg/errcode"
	"flec_blog/pkg/response"
	"flec_blog/pkg/upload"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService *service.FileService
}

// NewFileController 创建文件控制器
func NewFileController(fileService *service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// ============ 前台接口 ============

// UploadForWeb 文件上传
//
//	@Summary		文件上传
//	@Description	用户上传文件，存在限制
//	@Tags			文件
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"上传的文件"
//	@Param			type	formData	string	true	"上传类型：用户头像/评论贴图/反馈投诉"
//	@Success		200		{object}	response.Response{data=dto.FileUploadForWebResponse}
//	@Failure		400		{object}	response.Response
//	@Router			/upload [post]
func (ctrl *FileController) UploadForWeb(c *gin.Context) {
	var req dto.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	// 获取用户ID（如果已登录）
	var userID uint
	if id, exists := c.Get("user_id"); exists {
		userID = id.(uint)
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.FileUploadError.WithDetails(err.Error()))
		return
	}

	uploadReq := &upload.Request{
		File:       file,
		UploadType: upload.Type(req.Type),
		UserID:     userID,
	}

	// 从请求中提取 host
	host := upload.ExtractHostFromContext(c)

	fileInfo, err := ctrl.fileService.UploadForWeb(uploadReq, host)
	if err != nil {
		response.Error(c, errcode.FileUploadError.WithDetails(err.Error()))
		return
	}

	response.Success(c, fileInfo)
}

// ============ 后台管理接口 ============

// Upload 文件上传
//
//	@Summary		文件上传
//	@Description	管理员上传文件，限制相对宽松
//	@Tags			文件管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			file	formData	file	true	"上传的文件"
//	@Param			type	formData	string	true	"上传类型"
//	@Success		200		{object}	response.Response{data=dto.FileResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/admin/files [post]
func (ctrl *FileController) Upload(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.Unauthorized)
		return
	}

	var req dto.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.FileUploadError.WithDetails(err.Error()))
		return
	}

	uploadReq := &upload.Request{
		File:               file,
		UploadType:         upload.Type(req.Type),
		UserID:             userID.(uint),
		SkipSizeValidation: true, // 后台管理不限制文件大小
	}

	// 从请求中提取 host
	host := upload.ExtractHostFromContext(c)

	fileInfo, err := ctrl.fileService.Upload(uploadReq, host)
	if err != nil {
		response.Error(c, errcode.FileUploadError.WithDetails(err.Error()))
		return
	}

	response.Success(c, fileInfo)
}

// List 获取文件列表
//
//	@Summary		文件列表
//	@Description	获取已上传的所有文件
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			keyword		query		string	false	"关键词搜索（文件名）"
//	@Param			file_type	query		string	false	"文件类型（image/video/audio/document等）"
//	@Param			status		query		int		false	"状态（0=未使用 1=使用中）"
//	@Param			upload_type	query		string	false	"上传类型/用途"
//	@Param			min_size	query		int		false	"最小文件大小（字节）"
//	@Param			max_size	query		int		false	"最大文件大小（字节）"
//	@Param			start_time	query		string	false	"开始时间"
//	@Param			end_time	query		string	false	"结束时间"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Failure		403			{object}	response.Response
//	@Router			/admin/files [get]
func (ctrl *FileController) List(c *gin.Context) {
	var req dto.ListFilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileList, total, err := ctrl.fileService.List(&req)
	if err != nil {
		response.Error(c, errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.PageSuccess(c, fileList, total, req.Page, req.PageSize)
}

// Get 获取文件详情
//
//	@Summary		文件详情
//	@Description	获取文件详细信息
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文件 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/files/{id} [get]
func (ctrl *FileController) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileInfo, err := ctrl.fileService.Get(uint(id))
	if err != nil {
		response.Error(c, errcode.FileNotFound)
		return
	}

	response.Success(c, fileInfo)
}

// Delete 删除文件
//
//	@Summary		删除文件
//	@Description	删除指定文件，不可恢复，若支持也将同步删除存储桶文件
//	@Tags			文件管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"文件 ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/admin/files/{id} [delete]
func (ctrl *FileController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if err := ctrl.fileService.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "文件不存在") {
			response.Error(c, errcode.FileNotFound)
			return
		}
		response.Error(c, errcode.FileProcessError.WithDetails(err.Error()))
		return
	}

	response.Success(c, nil)
}
