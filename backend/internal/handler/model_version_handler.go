package handler

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/middleware"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type ModelVersionHandler struct {
	versionService *service.ModelVersionService
	modelService   *service.ModelService
}

func NewModelVersionHandler() *ModelVersionHandler {
	return &ModelVersionHandler{
		versionService: service.NewModelVersionService(),
		modelService:   service.NewModelService(),
	}
}

func (h *ModelVersionHandler) CreateVersion(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)

	var req service.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	version, err := h.versionService.CreateVersion(modelID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, version)
}

func (h *ModelVersionHandler) ListVersions(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	if _, ok := visibleModel(c, modelID); !ok {
		return
	}
	versions, err := h.versionService.ListVersions(modelID)
	if err != nil {
		response.InternalError(c, "failed to fetch versions")
		return
	}

	response.Success(c, versions)
}

func (h *ModelVersionHandler) GetVersion(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	version, err := h.versionService.GetVersion(versionID)
	if err != nil {
		response.NotFound(c, "version not found")
		return
	}
	if version.ModelID.String() != c.Param("id") {
		response.NotFound(c, "version not found")
		return
	}
	m, ok := visibleModel(c, version.ModelID)
	if !ok {
		return
	}
	if !canManageModel(c, m) {
		public := publicModel(*m)
		version.Model = &public
	}

	response.Success(c, version)
}

func (h *ModelVersionHandler) SetCurrentVersion(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	if err := h.versionService.SetCurrentVersion(modelID, versionID, userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "current version updated", nil)
}

func (h *ModelVersionHandler) UpdateVersion(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	var req service.UpdateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	version, err := h.versionService.UpdateVersion(versionID, userID, &req, isAdmin)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, version)
}

func (h *ModelVersionHandler) DeleteVersion(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	if err := h.versionService.DeleteVersion(modelID, versionID, userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "version deleted", nil)
}

func (h *ModelVersionHandler) DownloadVersion(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	version, err := h.versionService.GetVersion(versionID)
	if err != nil {
		response.NotFound(c, "version not found")
		return
	}

	if version.Model == nil {
		response.NotFound(c, "model not found")
		return
	}

	if version.ModelID.String() != c.Param("id") || version.Model.Status != "approved" || !version.Model.IsPublic {
		response.Forbidden(c, "model is not available for download")
		return
	}

	h.versionService.IncrementDownloads(versionID)

	response.Success(c, gin.H{
		"download_url": "/api/models/" + version.ModelID.String() + "/versions/" + versionID.String() + "/file",
		"file_name":    version.Model.Title + "-" + version.VersionNumber,
	})
}

func (h *ModelVersionHandler) ServeVersionFile(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version id")
		return
	}

	version, err := h.versionService.GetVersion(versionID)
	if err != nil {
		response.NotFound(c, "version not found")
		return
	}

	if version.Model == nil {
		response.NotFound(c, "model not found")
		return
	}

	if version.ModelID.String() != c.Param("id") || version.Model.Status != "approved" || !version.Model.IsPublic {
		response.Forbidden(c, "model is not available for download")
		return
	}

	fullPath, err := service.ResolveUploadPath(version.FilePath, "models")
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}

	fileName := version.Model.Title + "-" + version.VersionNumber
	ext := filepath.Ext(version.FilePath)
	if ext != "" {
		fileName = fileName + ext
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("X-Content-Type-Options", "nosniff")
	c.FileAttachment(fullPath, fileName)
}
