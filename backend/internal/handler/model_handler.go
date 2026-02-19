package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/middleware"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type ModelHandler struct {
	modelService     *service.ModelService
	favoriteService  *service.FavoriteService
}

func NewModelHandler() *ModelHandler {
	return &ModelHandler{
		modelService:    service.NewModelService(),
		favoriteService: service.NewFavoriteService(),
	}
}

func (h *ModelHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	models, total, err := h.modelService.ListPublic(page, pageSize, search)
	if err != nil {
		response.InternalError(c, "failed to fetch models")
		return
	}

	response.Paginated(c, models, total, page, pageSize)
}

func (h *ModelHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	model, err := h.modelService.GetByID(id)
	if err != nil {
		response.NotFound(c, "model not found")
		return
	}

	result := gin.H{
		"model": model,
	}

	userID, exists := c.Get("user_id")
	if exists {
		result["is_favorited"] = h.favoriteService.IsFavorited(userID.(uuid.UUID), id)
	}

	favoriteCount, _ := h.favoriteService.CountByModel(id)
	result["favorite_count"] = favoriteCount

	response.Success(c, result)
}

func (h *ModelHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	model, err := h.modelService.Create(userID, &req)
	if err != nil {
		c.Error(err)
		response.InternalError(c, "failed to create model: "+err.Error())
		return
	}

	response.Success(c, model)
}

func (h *ModelHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	var req service.UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	model, err := h.modelService.Update(id, userID, &req, role == "admin")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, model)
}

func (h *ModelHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	if err := h.modelService.Delete(id, userID, role == "admin"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "model deleted successfully", nil)
}

func (h *ModelHandler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	model, err := h.modelService.GetByID(id)
	if err != nil {
		response.NotFound(c, "model not found")
		return
	}

	h.modelService.IncrementDownloads(id)

	response.Success(c, gin.H{
		"file_path": model.FilePath,
		"file_name": model.Title,
	})
}

func (h *ModelHandler) AddFavorite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)

	if err := h.favoriteService.Add(userID, id); err != nil {
		response.InternalError(c, "failed to add favorite")
		return
	}

	response.SuccessWithMessage(c, "added to favorites", nil)
}

func (h *ModelHandler) RemoveFavorite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)

	if err := h.favoriteService.Remove(userID, id); err != nil {
		response.InternalError(c, "failed to remove favorite")
		return
	}

	response.SuccessWithMessage(c, "removed from favorites", nil)
}

func (h *ModelHandler) CheckFavorite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	userID := middleware.GetUserID(c)

	response.Success(c, gin.H{
		"is_favorited": h.favoriteService.IsFavorited(userID, id),
	})
}
