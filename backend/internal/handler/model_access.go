package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ysmmc/backend/internal/model"
	"github.com/ysmmc/backend/internal/service"
	"github.com/ysmmc/backend/pkg/response"
)

func canManageModel(c *gin.Context, m *model.Model) bool {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	return userID != nil && (userID == m.UserID || role == "admin" || role == "super_admin")
}

func visibleModel(c *gin.Context, id uuid.UUID) (*model.Model, bool) {
	m, err := service.NewModelService().GetByID(id)
	if err != nil || !(m.Status == "approved" && m.IsPublic || canManageModel(c, m)) {
		response.NotFound(c, "model not found")
		return nil, false
	}
	return m, true
}

// 不向公开访问者返回待审修改、审核意见或投稿者的私有资料。
func publicModel(m model.Model) model.Model {
	m.PendingChanges = nil
	m.RejectionReason = nil
	if m.User != nil {
		user := *m.User
		user.Email = ""
		user.PendingChanges = nil
		user.BannedReason = nil
		user.BannedAt = nil
		m.User = &user
	}
	return m
}
