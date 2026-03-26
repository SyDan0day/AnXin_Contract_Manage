package main

import (
	"contract-manage/config"
	"contract-manage/models"
	"fmt"
	"log"
	"time"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	fmt.Println("开始创建测试数据...")
	fmt.Println()

	// 查询用户ID
	var sales1, sales2, admin models.User
	if err := models.DB.Where("username = ?", "sales1").First(&sales1).Error; err != nil {
		log.Fatalf("找不到用户 sales1: %v", err)
	}
	if err := models.DB.Where("username = ?", "sales2").First(&sales2).Error; err != nil {
		log.Fatalf("找不到用户 sales2: %v", err)
	}
	if err := models.DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		log.Fatalf("找不到用户 admin: %v", err)
	}

	// 确保有客户数据
	var customerCount int64
	models.DB.Model(&models.Customer{}).Count(&customerCount)
	if customerCount == 0 {
		// 创建测试客户
		customers := []models.Customer{
			{Name: "阿里巴巴（中国）有限公司", Code: "C001", Type: "customer", ContactPerson: "张三", ContactPhone: "13800138001"},
			{Name: "腾讯科技（深圳）有限公司", Code: "C002", Type: "customer", ContactPerson: "李四", ContactPhone: "13800138002"},
			{Name: "华为技术有限公司", Code: "C003", Type: "customer", ContactPerson: "王五", ContactPhone: "13800138003"},
			{Name: "字节跳动科技有限公司", Code: "C004", Type: "customer", ContactPerson: "赵六", ContactPhone: "13800138004"},
			{Name: "京东集团", Code: "C005", Type: "supplier", ContactPerson: "钱七", ContactPhone: "13800138005"},
		}
		models.DB.Create(&customers)
		fmt.Printf("✅ 创建了 %d 个测试客户\n", len(customers))
	}

	// 确保有合同类型数据
	var typeCount int64
	models.DB.Model(&models.ContractType{}).Count(&typeCount)
	if typeCount == 0 {
		types := []models.ContractType{
			{Name: "采购合同", Code: "PO"},
			{Name: "服务合同", Code: "SV"},
			{Name: "租赁合同", Code: "LC"},
			{Name: "软件开发合同", Code: "SD"},
			{Name: "咨询合同", Code: "CS"},
		}
		models.DB.Create(&types)
		fmt.Printf("✅ 创建了 %d 个合同类型\n", len(types))
	}

	// 获取客户和类型
	var customers []models.Customer
	models.DB.Find(&customers)
	var contractTypes []models.ContractType
	models.DB.Find(&contractTypes)

	// 创建测试合同 - sales1的合同
	sales1Contracts := []struct {
		no     string
		title  string
		amount float64
		status models.ContractStatus
		days   int // 距离到期的天数，-1表示已过期
	}{
		{"CT202401001", "阿里巴巴云服务采购合同", 150000, models.StatusActive, 15},
		{"CT202401002", "腾讯企业邮箱服务合同", 50000, models.StatusActive, 25},
		{"CT202401003", "华为设备采购合同", 300000, models.StatusPending, 30},
		{"CT202401004", "软件开发服务合同", 200000, models.StatusDraft, 60},
		{"CT202401005", "咨询服务合同", 80000, models.StatusCompleted, -10},
	}

	fmt.Println("========== 创建 sales1 的合同 ==========")
	for i, c := range sales1Contracts {
		contract := models.Contract{
			ContractNo:     c.no,
			Title:          c.title,
			CustomerID:     customers[i%len(customers)].ID,
			ContractTypeID: contractTypes[i%len(contractTypes)].ID,
			Amount:         c.amount,
			Currency:       "CNY",
			Status:         c.status,
			CreatorID:      sales1.ID,
		}

		// 设置日期
		now := time.Now()
		contract.SignDate = func() *time.Time { t := now.AddDate(0, -1, 0); return &t }()
		contract.StartDate = func() *time.Time { t := now.AddDate(0, -1, 0); return &t }()

		if c.days > 0 {
			contract.EndDate = func() *time.Time { t := now.AddDate(0, 0, c.days); return &t }()
		} else {
			contract.EndDate = func() *time.Time { t := now.AddDate(0, 0, c.days); return &t }()
		}

		if err := models.DB.Create(&contract).Error; err != nil {
			fmt.Printf("❌ 创建合同失败 [%s]: %v\n", c.no, err)
		} else {
			fmt.Printf("✅ 创建合同: %s - %s (金额: %.2f, 状态: %s)\n", c.no, c.title, c.amount, c.status)
		}
	}

	// 创建测试合同 - sales2的合同
	sales2Contracts := []struct {
		no     string
		title  string
		amount float64
		status models.ContractStatus
		days   int
	}{
		{"CT202401006", "京东物流服务合同", 120000, models.StatusActive, 10},
		{"CT202401007", "字节跳动广告投放合同", 250000, models.StatusActive, 20},
		{"CT202401008", "办公设备采购合同", 75000, models.StatusPending, 45},
		{"CT202401009", "IT咨询服务合同", 60000, models.StatusDraft, 90},
		{"CT202401010", "培训服务合同", 45000, models.StatusCompleted, -5},
	}

	fmt.Println()
	fmt.Println("========== 创建 sales2 的合同 ==========")
	for i, c := range sales2Contracts {
		contract := models.Contract{
			ContractNo:     c.no,
			Title:          c.title,
			CustomerID:     customers[(i+2)%len(customers)].ID,
			ContractTypeID: contractTypes[(i+1)%len(contractTypes)].ID,
			Amount:         c.amount,
			Currency:       "CNY",
			Status:         c.status,
			CreatorID:      sales2.ID,
		}

		now := time.Now()
		contract.SignDate = func() *time.Time { t := now.AddDate(0, -1, 0); return &t }()
		contract.StartDate = func() *time.Time { t := now.AddDate(0, -1, 0); return &t }()
		contract.EndDate = func() *time.Time { t := now.AddDate(0, 0, c.days); return &t }()

		if err := models.DB.Create(&contract).Error; err != nil {
			fmt.Printf("❌ 创建合同失败 [%s]: %v\n", c.no, err)
		} else {
			fmt.Printf("✅ 创建合同: %s - %s (金额: %.2f, 状态: %s)\n", c.no, c.title, c.amount, c.status)
		}
	}

	// 创建审批流程
	fmt.Println()
	fmt.Println("========== 创建审批流程 ==========")

	// 为待审批的合同创建审批记录
	var pendingContracts []models.Contract
	models.DB.Where("status = ?", models.StatusPending).Find(&pendingContracts)

	for _, contract := range pendingContracts {
		// 检查是否已有审批记录
		var existingCount int64
		models.DB.Model(&models.ApprovalRecord{}).Where("contract_id = ?", contract.ID).Count(&existingCount)
		if existingCount > 0 {
			fmt.Printf("⚠️  合同 %s 已有审批记录，跳过\n", contract.ContractNo)
			continue
		}

		// 创建审批记录 - 一级审批（销售负责人）
		approval1 := models.ApprovalRecord{
			ContractID:   contract.ID,
			ApproverID:   0, // 待审批
			Level:        1,
			ApproverRole: string(models.RoleSalesManager),
			Status:       models.ApprovalPending,
			Comment:      "",
		}
		models.DB.Create(&approval1)
		fmt.Printf("✅ 创建审批记录: %s - 一级审批(销售负责人)\n", contract.ContractNo)

		// 创建审批记录 - 二级审批（技术负责人）
		approval2 := models.ApprovalRecord{
			ContractID:   contract.ID,
			ApproverID:   0,
			Level:        2,
			ApproverRole: string(models.RoleTechLeader),
			Status:       models.ApprovalPending,
			Comment:      "",
		}
		models.DB.Create(&approval2)
		fmt.Printf("✅ 创建审批记录: %s - 二级审批(技术负责人)\n", contract.ContractNo)

		// 创建审批记录 - 三级审批（财务负责人）
		approval3 := models.ApprovalRecord{
			ContractID:   contract.ID,
			ApproverID:   0,
			Level:        3,
			ApproverRole: string(models.RoleFinanceLeader),
			Status:       models.ApprovalPending,
			Comment:      "",
		}
		models.DB.Create(&approval3)
		fmt.Printf("✅ 创建审批记录: %s - 三级审批(财务负责人)\n", contract.ContractNo)
	}

	// 统计
	fmt.Println()
	fmt.Println("========== 测试数据统计 ==========")
	fmt.Println()

	// 统计各角色创建的合同数量
	var sales1Count, sales2Count, adminCount int64
	models.DB.Model(&models.Contract{}).Where("creator_id = ?", sales1.ID).Count(&sales1Count)
	models.DB.Model(&models.Contract{}).Where("creator_id = ?", sales2.ID).Count(&sales2Count)
	models.DB.Model(&models.Contract{}).Where("creator_id = ?", admin.ID).Count(&adminCount)

	fmt.Printf("合同数量统计:\n")
	fmt.Printf("  - sales1 创建: %d 个\n", sales1Count)
	fmt.Printf("  - sales2 创建: %d 个\n", sales2Count)
	fmt.Printf("  - admin 创建: %d 个\n", adminCount)

	// 统计审批记录
	var pendingCount, approvedCount, rejectedCount int64
	models.DB.Model(&models.ApprovalRecord{}).Where("status = ?", models.ApprovalPending).Count(&pendingCount)
	models.DB.Model(&models.ApprovalRecord{}).Where("status = ?", models.ApprovalApproved).Count(&approvedCount)
	models.DB.Model(&models.ApprovalRecord{}).Where("status = ?", models.ApprovalRejected).Count(&rejectedCount)

	fmt.Printf("\n审批记录统计:\n")
	fmt.Printf("  - 待审批: %d 条\n", pendingCount)
	fmt.Printf("  - 已批准: %d 条\n", approvedCount)
	fmt.Printf("  - 已拒绝: %d 条\n", rejectedCount)

	fmt.Println()
	fmt.Println("========== 测试账号说明 ==========")
	fmt.Println()
	fmt.Println("| 用户名      | 密码         | 角色       | 合同数量 |")
	fmt.Println("|------------|-------------|-----------|---------|")
	fmt.Printf("| sales1     | Sales1@123456 | 销售      | %d 个   |\n", sales1Count)
	fmt.Printf("| sales2     | Sales2@123456 | 销售      | %d 个   |\n", sales2Count)
	fmt.Printf("| admin      | Admin@123456  | 管理员     | %d 个   |\n", adminCount)
	fmt.Println()
	fmt.Println("测试数据创建完成！")
}
