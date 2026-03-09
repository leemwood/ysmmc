package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/middleware"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type ModelImageHandler struct {
	imageService *service.ModelImageService
}

func NewModelImageHandler() *ModelImageHandler {
	return &ModelImageHandler{
		imageService: service.NewModelImageService(),
	}
}

func (h *ModelImageHandler) AddImage(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	var req service.AddModelImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	image, err := h.imageService.AddImage(modelID, userID, &req, isAdmin)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, image)
}

func (h *ModelImageHandler) ListImages(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	images, err := h.imageService.ListImages(modelID)
	if err != nil {
		response.InternalError(c, "failed to fetch images")
		return
	}

	response.Success(c, images)
}

func (h *ModelImageHandler) DeleteImage(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	fileID, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	if err := h.imageService.DeleteImage(modelID, fileID, userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "image deleted", nil)
}

func (h *ModelImageHandler) UpdateOrder(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	isAdmin := role == "admin" || role == "super_admin"

	var req service.UpdateImageOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.imageService.UpdateOrder(modelID, userID, &req, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "order updated", nil)
}
