package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/middleware"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: service.NewFavoriteService(),
	}
}

func (h *FavoriteHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))

	favorites, total, err := h.favoriteService.ListByUserID(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch favorites")
		return
	}

	response.Paginated(c, favorites, total, page, pageSize)
}

type AnnouncementHandler struct {
	announcementService *service.AnnouncementService
}

func NewAnnouncementHandler() *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: service.NewAnnouncementService(),
	}
}

func (h *AnnouncementHandler) List(c *gin.Context) {
	announcements, err := h.announcementService.ListActive()
	if err != nil {
		response.InternalError(c, "failed to fetch announcements")
		return
	}

	response.Success(c, announcements)
}

func (h *AnnouncementHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	announcements, total, err := h.announcementService.List(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to fetch announcements")
		return
	}

	response.Paginated(c, announcements, total, page, pageSize)
}

func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	announcement, err := h.announcementService.GetByID(id)
	if err != nil {
		response.NotFound(c, "announcement not found")
		return
	}

	response.Success(c, announcement)
}
