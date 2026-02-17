package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type AdminHandler struct {
	modelService        *service.ModelService
	userService         *service.UserService
	announcementService *service.AnnouncementService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		modelService:        service.NewModelService(),
		userService:         service.NewUserService(),
		announcementService: service.NewAnnouncementService(),
	}
}

func (h *AdminHandler) ListPendingModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	models, total, err := h.modelService.ListPending(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch pending models")
		return
	}

	response.Paginated(c, models, total, page, pageSize)
}

func (h *AdminHandler) ListPendingUpdates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	models, total, err := h.modelService.ListPendingUpdates(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch pending updates")
		return
	}

	response.Paginated(c, models, total, page, pageSize)
}

func (h *AdminHandler) ApproveModel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	if err := h.modelService.Approve(id); err != nil {
		response.InternalError(c, "failed to approve model")
		return
	}

	response.SuccessWithMessage(c, "model approved successfully", nil)
}

func (h *AdminHandler) RejectModel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid model id")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.modelService.Reject(id, req.Reason); err != nil {
		response.InternalError(c, "failed to reject model")
		return
	}

	response.SuccessWithMessage(c, "model rejected", nil)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.userService.List(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch users")
		return
	}

	response.Paginated(c, users, total, page, pageSize)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateRole(id, req.Role); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "role updated successfully", nil)
}

func (h *AdminHandler) ListPendingProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.userService.ListPendingProfiles(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch pending profiles")
		return
	}

	response.Paginated(c, users, total, page, pageSize)
}

func (h *AdminHandler) ApproveProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.ApproveProfile(id); err != nil {
		response.InternalError(c, "failed to approve profile")
		return
	}

	response.SuccessWithMessage(c, "profile approved successfully", nil)
}

func (h *AdminHandler) RejectProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userService.RejectProfile(id); err != nil {
		response.InternalError(c, "failed to reject profile")
		return
	}

	response.SuccessWithMessage(c, "profile changes rejected", nil)
}

func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	var req service.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	announcement, err := h.announcementService.Create(&req)
	if err != nil {
		response.InternalError(c, "failed to create announcement")
		return
	}

	response.Success(c, announcement)
}

func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	var req service.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	announcement, err := h.announcementService.Update(id, &req)
	if err != nil {
		response.InternalError(c, "failed to update announcement")
		return
	}

	response.Success(c, announcement)
}

func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	if err := h.announcementService.Delete(id); err != nil {
		response.InternalError(c, "failed to delete announcement")
		return
	}

	response.SuccessWithMessage(c, "announcement deleted successfully", nil)
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	totalUsers, _ := h.userService.Count()
	totalModels, _ := h.modelService.Count()
	pendingModels, _ := h.modelService.CountByStatus("pending")
	totalDownloads, _ := h.modelService.SumDownloads()

	response.Success(c, gin.H{
		"total_users":     totalUsers,
		"total_models":    totalModels,
		"pending_models":  pendingModels,
		"total_downloads": totalDownloads,
	})
}
