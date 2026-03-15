package handlers

import (
	"contract-manage/middleware"
	"contract-manage/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	approvalService *services.ApprovalService
}

func NewApprovalHandler() *ApprovalHandler {
	return &ApprovalHandler{
		approvalService: services.NewApprovalService(),
	}
}

func (h *ApprovalHandler) GetContractApprovals(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	approvals, err := h.approvalService.GetApprovalRecords(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, approvals)
}

func (h *ApprovalHandler) CreateApproval(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input services.ApprovalRecordCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	input.ContractID = uint(contractID)
	input.ApproverRole = role

	approval, err := h.approvalService.CreateApprovalRecord(input, userID, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "user" || role == "manager" || role == "admin" {
		h.approvalService.SubmitForApproval(uint(contractID), userID)
	}

	c.JSON(http.StatusCreated, approval)
}

func (h *ApprovalHandler) UpdateApproval(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("approval_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approval ID"})
		return
	}

	var input services.ApprovalRecordUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, _ := h.approvalService.GetApprovalRecordByID(uint(id))
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval record not found"})
		return
	}

	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	contractStatus := ""
	if input.Status == "approved" {
		contractStatus = "active"
	} else if input.Status == "rejected" {
		contractStatus = "draft"
	}

	userID, _ := middleware.GetCurrentUserID(c)
	approval, err := h.approvalService.UpdateApprovalRecord(uint(id), input, contractStatus, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, approval)
}

func (h *ApprovalHandler) GetContractReminders(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	reminders, err := h.approvalService.GetReminders(uint(contractID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *ApprovalHandler) CreateReminder(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input services.ReminderCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ContractID = uint(contractID)
	reminder, err := h.approvalService.CreateReminder(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *ApprovalHandler) SendReminder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("reminder_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	if err := h.approvalService.UpdateReminderSent(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder sent successfully"})
}

func (h *ApprovalHandler) GetExpiringContracts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	contracts, err := h.approvalService.GetExpiringContracts(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contracts": contracts,
		"days":      days,
	})
}

func (h *ApprovalHandler) GetStatistics(c *gin.Context) {
	stats, err := h.approvalService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ApprovalHandler) GetPendingApprovals(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	approvals, err := h.approvalService.GetPendingApprovalsByRole(role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, approvals)
}

func (h *ApprovalHandler) GetNotificationCounts(c *gin.Context) {
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	counts := map[string]int{}

	pendingApprovals, _ := h.approvalService.GetPendingApprovalsByRole(role, 0)
	counts["pendingApprovals"] = len(pendingApprovals)

	if role == "manager" || role == "admin" {
		pendingStatusChanges, _ := h.approvalService.GetPendingStatusChangesCount()
		counts["pendingStatusChanges"] = pendingStatusChanges
	} else {
		counts["pendingStatusChanges"] = 0
	}

	expiringContracts, _ := h.approvalService.GetExpiringContracts(30)
	counts["expiringContracts"] = len(expiringContracts)

	counts["total"] = counts["pendingApprovals"] + counts["pendingStatusChanges"] + counts["expiringContracts"]

	c.JSON(http.StatusOK, counts)
}
