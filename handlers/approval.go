package handlers

import (
	"contract-manage/middleware"
	"contract-manage/models"
	"contract-manage/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	approvalService *services.ApprovalService
	workflowService *services.WorkflowService
}

func NewApprovalHandler() *ApprovalHandler {
	return &ApprovalHandler{
		approvalService: services.NewApprovalService(),
		workflowService: services.NewWorkflowService(models.DB),
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

	// 获取合同信息，检查当前状态
	var contract models.Contract
	if err := models.DB.First(&contract, uint(contractID)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "合同不存在"})
		return
	}

	// 检查合同状态：只有在草稿或被拒绝状态下才能提交审批
	if contract.Status == models.StatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同正在审批中，不能重复提交"})
		return
	}
	if contract.Status == models.StatusArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同已归档，不能提交审批"})
		return
	}
	if contract.Status == models.StatusTerminated {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同已终止，不能提交审批"})
		return
	}

	// 检查是否有未取消的审批流程
	var existingWorkflow models.ApprovalWorkflow
	if err := models.DB.Where("contract_id = ? AND status = ?", uint(contractID), models.WorkflowStatusPending).First(&existingWorkflow).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该合同已有待处理的审批流程，请勿重复提交"})
		return
	}

	// 如果存在已拒绝/已完成/已取消的旧审批流程，删除它们以便重新提交
	models.DB.Where("contract_id = ?", uint(contractID)).Delete(&models.ApprovalWorkflow{})
	models.DB.Where("contract_id = ?", uint(contractID)).Delete(&models.WorkflowApproval{})

	approval, err := h.approvalService.CreateApprovalRecord(input, userID, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "user" || role == "manager" || role == "admin" {
		var contract models.Contract
		if err := models.DB.First(&contract, uint(contractID)).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "合同不存在"})
			return
		}

		oldStatus := string(contract.Status)
		contract.Status = models.StatusPending
		if err := models.DB.Save(&contract).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新合同状态失败"})
			return
		}

		workflow, err := h.workflowService.CreateWorkflow(uint64(contractID), role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "创建审批流程失败: " + err.Error()})
			return
		}

		services.NewContractService().AddLifecycleEvent(uint(contractID), services.LifecycleEventInput{
			EventType:   string(models.LifecycleSubmitted),
			FromStatus:  oldStatus,
			ToStatus:    string(models.StatusPending),
			Description: "合同提交审批",
		}, userID)

		c.JSON(http.StatusCreated, gin.H{
			"approval": approval,
			"workflow": workflow,
		})
		return
	}

	c.JSON(http.StatusCreated, approval)
}

// WithdrawApproval 撤回审批申请
// POST /api/contracts/:contract_id/approvals/withdraw
func (h *ApprovalHandler) WithdrawApproval(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contract_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var contract models.Contract
	if err := models.DB.First(&contract, uint(contractID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	if contract.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有合同创建者可以撤回审批申请"})
		return
	}

	if contract.Status != models.StatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前状态不允许撤回，只有待审批状态的合同才能撤回"})
		return
	}

	var workflow models.ApprovalWorkflow
	if err := models.DB.Where("contract_id = ?", uint(contractID)).First(&workflow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "审批流程不存在"})
		return
	}

	oldStatus := string(contract.Status)
	contract.Status = models.StatusDraft
	if err := models.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新合同状态失败"})
		return
	}

	if err := models.DB.Model(&workflow).Update("status", models.WorkflowStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新审批流程状态失败"})
		return
	}

	if err := models.DB.Model(&models.WorkflowApproval{}).Where("workflow_id = ?", workflow.ID).
		Update("status", models.WorkflowStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新审批记录状态失败"})
		return
	}

	services.NewContractService().AddLifecycleEvent(uint(contractID), services.LifecycleEventInput{
		EventType:   string(models.LifecycleRejected),
		FromStatus:  oldStatus,
		ToStatus:    string(models.StatusDraft),
		Description: "撤回审批申请",
	}, userID)

	c.JSON(http.StatusOK, gin.H{"message": "审批申请已撤回"})
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

	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var workflowApproval models.WorkflowApproval
	if err := models.DB.First(&workflowApproval, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approval record not found"})
		return
	}

	var workflow models.ApprovalWorkflow
	if err := models.DB.First(&workflow, workflowApproval.WorkflowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if workflowApproval.ApproverRole != role {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to approve this"})
		return
	}

	if input.Status == "approved" {
		h.workflowService.Approve(workflow.ID, 0, workflowApproval.Level, uint64(userID), input.Comment)
	} else if input.Status == "rejected" {
		h.workflowService.Reject(workflow.ID, workflowApproval.Level, uint64(userID), input.Comment)
	}

	contractStatus := ""
	if input.Status == "approved" {
		if workflowApproval.Level >= workflow.MaxLevel {
			contractStatus = "active"
		} else {
			contractStatus = "pending"
		}
	} else if input.Status == "rejected" {
		contractStatus = "draft"
	}

	if contractStatus != "" {
		models.DB.Model(&models.Contract{}).Where("id = ?", workflow.ContractID).Update("status", contractStatus)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval updated successfully"})
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

	// 获取当前用户信息用于权限过滤
	userID, _ := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	contracts, err := h.approvalService.GetExpiringContractsWithAuth(days, userID, role)
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
	// 获取当前用户信息
	userID, exists := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	fmt.Printf("[DEBUG] GetStatistics - userID: %d, role: %s, exists: %v\n", userID, role, exists)

	stats, err := h.approvalService.GetStatistics(services.GetStatisticsParams{
		UserID: userID,
		Role:   role,
	})
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
	userID, _ := middleware.GetCurrentUserID(c)
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "" {
		role = "user"
	}

	counts := map[string]int{}

	// 待审批数量（只有审批角色才能看到）
	if role == "admin" || role == "sales_manager" || role == "tech_leader" || role == "finance_leader" {
		pendingApprovals, _ := h.approvalService.GetPendingApprovalsByRole(role, 0)
		counts["pendingApprovals"] = len(pendingApprovals)
	} else {
		counts["pendingApprovals"] = 0
	}

	// 待审批状态变更数量（只有管理员和销售负责人能看到）
	if role == "admin" || role == "sales_manager" {
		pendingStatusChanges, _ := h.approvalService.GetPendingStatusChangesCount()
		counts["pendingStatusChanges"] = pendingStatusChanges
	} else {
		counts["pendingStatusChanges"] = 0
	}

	// 即将到期合同数量（所有角色都可以看到，但数据根据角色过滤）
	expiringContracts, _ := h.approvalService.GetExpiringContractsWithAuth(30, userID, role)
	counts["expiringContracts"] = len(expiringContracts)

	counts["total"] = counts["pendingApprovals"] + counts["pendingStatusChanges"] + counts["expiringContracts"]

	c.JSON(http.StatusOK, counts)
}
