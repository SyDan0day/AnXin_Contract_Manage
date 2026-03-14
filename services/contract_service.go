package services

import (
	"contract-manage/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ContractService struct{}

func NewContractService() *ContractService {
	return &ContractService{}
}

func (s *ContractService) generateContractNo() string {
	today := time.Now()
	prefix := fmt.Sprintf("CT%s", today.Format("200601"))

	var lastContract models.Contract
	models.DB.Where("contract_no LIKE ?", prefix+"%").Order("contract_no DESC").First(&lastContract)

	var newNo string
	if lastContract.ID != 0 {
		lastNo := lastContract.ContractNo[len(lastContract.ContractNo)-4:]
		var num int
		fmt.Sscanf(lastNo, "%d", &num)
		newNo = fmt.Sprintf("%04d", num+1)
	} else {
		newNo = "0001"
	}

	return prefix + newNo
}

func (s *ContractService) GetContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.Preload("Customer").Preload("Creator").Preload("ContractType").First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (s *ContractService) GetContractByNo(contractNo string) (*models.Contract, error) {
	var contract models.Contract
	if err := models.DB.Where("contract_no = ?", contractNo).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (s *ContractService) GetContracts(skip, limit int, customerID, contractTypeID uint, status string) ([]models.Contract, error) {
	var contracts []models.Contract
	query := models.DB.Preload("Customer").Preload("Creator").Preload("ContractType")

	if customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}
	if contractTypeID > 0 {
		query = query.Where("contract_type_id = ?", contractTypeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Offset(skip).Limit(limit).Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

type ContractCreateInput struct {
	Title          string  `json:"title" binding:"required"`
	CustomerID     uint    `json:"customer_id" binding:"required"`
	ContractTypeID uint    `json:"contract_type_id" binding:"required"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	SignDate       *time.Time `json:"sign_date"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	PaymentTerms   string  `json:"payment_terms"`
	Content        string  `json:"content"`
	Notes          string  `json:"notes"`
}

func (s *ContractService) CreateContract(input ContractCreateInput, creatorID uint) (*models.Contract, error) {
	contract := models.Contract{
		ContractNo:     s.generateContractNo(),
		Title:          input.Title,
		CustomerID:     input.CustomerID,
		ContractTypeID: input.ContractTypeID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		SignDate:       input.SignDate,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		PaymentTerms:   input.PaymentTerms,
		Content:        input.Content,
		Notes:          input.Notes,
		CreatorID:      creatorID,
		Status:         models.StatusDraft,
	}

	if contract.Currency == "" {
		contract.Currency = "CNY"
	}

	if err := models.DB.Create(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

type ContractUpdateInput struct {
	Title          string             `json:"title"`
	CustomerID     uint               `json:"customer_id"`
	ContractTypeID uint               `json:"contract_type_id"`
	Amount         float64            `json:"amount"`
	Currency       string             `json:"currency"`
	Status         models.ContractStatus `json:"status"`
	SignDate       *time.Time         `json:"sign_date"`
	StartDate      *time.Time         `json:"start_date"`
	EndDate        *time.Time         `json:"end_date"`
	PaymentTerms   string             `json:"payment_terms"`
	Content        string             `json:"content"`
	Notes          string             `json:"notes"`
}

func (s *ContractService) UpdateContract(id uint, input ContractUpdateInput) (*models.Contract, error) {
	contract, err := s.GetContractByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.CustomerID > 0 {
		updates["customer_id"] = input.CustomerID
	}
	if input.ContractTypeID > 0 {
		updates["contract_type_id"] = input.ContractTypeID
	}
	if input.Amount > 0 {
		updates["amount"] = input.Amount
	}
	if input.Currency != "" {
		updates["currency"] = input.Currency
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.SignDate != nil {
		updates["sign_date"] = input.SignDate
	}
	if input.StartDate != nil {
		updates["start_date"] = input.StartDate
	}
	if input.EndDate != nil {
		updates["end_date"] = input.EndDate
	}
	if input.PaymentTerms != "" {
		updates["payment_terms"] = input.PaymentTerms
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	if err := models.DB.Model(contract).Updates(updates).Error; err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *ContractService) DeleteContract(id uint) error {
	result := models.DB.Delete(&models.Contract{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *ContractService) GetContractExecutionByID(id uint) (*models.ContractExecution, error) {
	var execution models.ContractExecution
	if err := models.DB.First(&execution, id).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

func (s *ContractService) GetContractExecutions(contractID uint) ([]models.ContractExecution, error) {
	var executions []models.ContractExecution
	if err := models.DB.Where("contract_id = ?", contractID).Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

type ContractExecutionCreateInput struct {
	ContractID    uint       `json:"contract_id" binding:"required"`
	Stage         string     `json:"stage"`
	StageDate     *time.Time `json:"stage_date"`
	Progress      float64    `json:"progress"`
	PaymentAmount float64    `json:"payment_amount"`
	PaymentDate   *time.Time `json:"payment_date"`
	Description   string     `json:"description"`
}

func (s *ContractService) CreateContractExecution(input ContractExecutionCreateInput, operatorID uint) (*models.ContractExecution, error) {
	execution := models.ContractExecution{
		ContractID:    input.ContractID,
		Stage:         input.Stage,
		StageDate:     input.StageDate,
		Progress:      input.Progress,
		PaymentAmount: input.PaymentAmount,
		PaymentDate:   input.PaymentDate,
		Description:   input.Description,
		OperatorID:    operatorID,
	}

	if err := models.DB.Create(&execution).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

func (s *ContractService) GetDocumentByID(id uint) (*models.Document, error) {
	var document models.Document
	if err := models.DB.First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (s *ContractService) GetDocuments(contractID uint) ([]models.Document, error) {
	var documents []models.Document
	if err := models.DB.Where("contract_id = ?", contractID).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

type DocumentCreateInput struct {
	ContractID uint   `json:"contract_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	FilePath   string `json:"file_path" binding:"required"`
	FileSize   int    `json:"file_size" binding:"required"`
	FileType   string `json:"file_type"`
}

func (s *ContractService) CreateDocument(input DocumentCreateInput, uploaderID uint) (*models.Document, error) {
	document := models.Document{
		ContractID: input.ContractID,
		Name:       input.Name,
		FilePath:   input.FilePath,
		FileSize:   input.FileSize,
		FileType:   input.FileType,
		Version:    "1.0",
		UploaderID: uploaderID,
	}

	if err := models.DB.Create(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (s *ContractService) DeleteDocument(id uint) error {
	result := models.DB.Delete(&models.Document{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}