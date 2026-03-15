package services

import (
	"contract-manage/models"
	"errors"

	"gorm.io/gorm"
)

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) GetCustomerByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := models.DB.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) GetCustomerByCode(code string) (*models.Customer, error) {
	var customer models.Customer
	if err := models.DB.Where("code = ?", code).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) GetCustomers(skip, limit int, customerType string) ([]models.Customer, error) {
	var customers []models.Customer
	query := models.DB
	if customerType != "" {
		query = query.Where("type = ?", customerType)
	}
	if err := query.Offset(skip).Limit(limit).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

type CustomerCreateInput struct {
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type"`
	Code          string `json:"code" binding:"required"`
	ContactPerson string `json:"contact_person"`
	ContactPhone  string `json:"contact_phone"`
	ContactEmail  string `json:"contact_email"`
	Address       string `json:"address"`
	CreditRating  string `json:"credit_rating"`
}

func (s *CustomerService) CreateCustomer(input CustomerCreateInput) (*models.Customer, error) {
	if _, err := s.GetCustomerByCode(input.Code); err == nil {
		return nil, errors.New("customer code already exists")
	}

	customer := models.Customer{
		Name:          input.Name,
		Type:          input.Type,
		Code:          input.Code,
		ContactPerson: input.ContactPerson,
		ContactPhone:  input.ContactPhone,
		ContactEmail:  input.ContactEmail,
		Address:       input.Address,
		CreditRating:  input.CreditRating,
		IsActive:      true,
	}

	if err := models.DB.Create(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

type CustomerUpdateInput struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	ContactPerson string `json:"contact_person"`
	ContactPhone  string `json:"contact_phone"`
	ContactEmail  string `json:"contact_email"`
	Address       string `json:"address"`
	CreditRating  string `json:"credit_rating"`
	IsActive      *bool  `json:"is_active"`
}

func (s *CustomerService) UpdateCustomer(id uint, input CustomerUpdateInput) (*models.Customer, error) {
	customer, err := s.GetCustomerByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Type != "" {
		updates["type"] = input.Type
	}
	if input.ContactPerson != "" {
		updates["contact_person"] = input.ContactPerson
	}
	if input.ContactPhone != "" {
		updates["contact_phone"] = input.ContactPhone
	}
	if input.ContactEmail != "" {
		updates["contact_email"] = input.ContactEmail
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.CreditRating != "" {
		updates["credit_rating"] = input.CreditRating
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := models.DB.Model(customer).Updates(updates).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) DeleteCustomer(id uint) error {
	result := models.DB.Delete(&models.Customer{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *CustomerService) GetContractTypeByID(id uint) (*models.ContractType, error) {
	var contractType models.ContractType
	if err := models.DB.First(&contractType, id).Error; err != nil {
		return nil, err
	}
	return &contractType, nil
}

func (s *CustomerService) GetContractTypes(skip, limit int) ([]models.ContractType, error) {
	var contractTypes []models.ContractType
	if err := models.DB.Offset(skip).Limit(limit).Find(&contractTypes).Error; err != nil {
		return nil, err
	}
	return contractTypes, nil
}

type ContractTypeCreateInput struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

func (s *CustomerService) CreateContractType(input ContractTypeCreateInput) (*models.ContractType, error) {
	contractType := models.ContractType{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
	}

	if err := models.DB.Create(&contractType).Error; err != nil {
		return nil, err
	}
	return &contractType, nil
}

func (s *CustomerService) UpdateContractType(id uint, input ContractTypeCreateInput) (*models.ContractType, error) {
	var contractType models.ContractType
	if err := models.DB.First(&contractType, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":        input.Name,
		"code":        input.Code,
		"description": input.Description,
	}

	if err := models.DB.Model(&contractType).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &contractType, nil
}

func (s *CustomerService) DeleteContractType(id uint) error {
	return models.DB.Delete(&models.ContractType{}, id).Error
}
