const mammoth = require('mammoth');
const fs = require('fs');

const docxPath = '/root/AnXin_Contract_Manage_v1/uploads/1/test_contract_formatted.docx';

async function test() {
  try {
    const result = await mammoth.extractRawText({ path: docxPath });
    console.log('=== 提取的文本内容 ===');
    console.log(result.value);
    console.log('\n=== 提取的合同信息 ===');
    
    const text = result.value;
    
    // 提取合同编号
    const contractNoMatch = text.match(/合同编号[：:]\s*([A-Z0-9\-]+)/);
    console.log('合同编号:', contractNoMatch ? contractNoMatch[1] : '未找到');
    
    // 提取合同名称
    const titleMatch = text.match(/合同名称[：:]\s*([^\n]{2,50})/);
    console.log('合同名称:', titleMatch ? titleMatch[1].trim() : '未找到');
    
    // 提取客户名称
    const customerMatch = text.match(/甲方[（(]?客户[）)]?[：:]\s*([^\n]{2,50})/);
    console.log('客户名称:', customerMatch ? customerMatch[1].trim() : '未找到');
    
    // 提取金额
    const amountMatch = text.match(/合同金额[：:]\s*([\d,]+\.?\d*)\s*(?:元|万)?/);
    console.log('合同金额:', amountMatch ? amountMatch[1] : '未找到');
    
    // 提取日期
    const signDateMatch = text.match(/签订日期[：:]\s*(\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日]?)/);
    console.log('签订日期:', signDateMatch ? signDateMatch[1] : '未找到');
    
    // 提取联系人
    const personMatch = text.match(/联系人[：:]\s*([^\n]{2,20})/);
    console.log('联系人:', personMatch ? personMatch[1].trim() : '未找到');
    
    // 提取电话
    const phoneMatch = text.match(/(?:联系电话|电话)[：:]\s*(\d{3}[-\s]?\d{4}[-\s]?\d{4}|\d{11})/);
    console.log('联系电话:', phoneMatch ? phoneMatch[1] : '未找到');
    
  } catch (error) {
    console.error('错误:', error.message);
  }
}

test();
