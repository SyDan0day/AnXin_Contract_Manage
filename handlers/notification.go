package handlers

import (
	"contract-manage/models"
	"contract-manage/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewNotificationService(models.DB),
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	notifications, err := h.notificationService.GetNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	notifications, err := h.notificationService.GetUnreadNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	var input struct {
		All bool `json:"all"`
	}
	c.ShouldBindJSON(&input)

	userID := c.GetUint("user_id")

	if input.All {
		err := h.notificationService.MarkAllAsRead(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(notificationID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err := h.notificationService.MarkAsRead(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	notifications, err := h.notificationService.GetUnreadNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": len(notifications)})
}
