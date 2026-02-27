from datetime import datetime, date
from typing import Optional
from pydantic import BaseModel

class ContractBase(BaseModel):
    title: str
    customer_id: int
    contract_type_id: int
    amount: Optional[float] = None
    currency: str = "CNY"
    sign_date: Optional[date] = None
    start_date: Optional[date] = None
    end_date: Optional[date] = None
    payment_terms: Optional[str] = None
    content: Optional[str] = None
    notes: Optional[str] = None

class ContractCreate(ContractBase):
    pass

class ContractUpdate(BaseModel):
    title: Optional[str] = None
    customer_id: Optional[int] = None
    contract_type_id: Optional[int] = None
    amount: Optional[float] = None
    currency: Optional[str] = None
    status: Optional[str] = None
    sign_date: Optional[date] = None
    start_date: Optional[date] = None
    end_date: Optional[date] = None
    payment_terms: Optional[str] = None
    content: Optional[str] = None
    notes: Optional[str] = None

class ContractResponse(ContractBase):
    id: int
    contract_no: str
    status: str
    creator_id: int
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class ContractExecutionBase(BaseModel):
    contract_id: int
    stage: Optional[str] = None
    stage_date: Optional[date] = None
    progress: Optional[float] = None
    payment_amount: Optional[float] = None
    payment_date: Optional[date] = None
    description: Optional[str] = None

class ContractExecutionCreate(ContractExecutionBase):
    pass

class ContractExecutionResponse(ContractExecutionBase):
    id: int
    operator_id: int
    created_at: datetime
    
    class Config:
        from_attributes = True

class DocumentBase(BaseModel):
    contract_id: int
    name: str
    file_type: Optional[str] = None

class DocumentCreate(DocumentBase):
    file_path: str
    file_size: int

class DocumentResponse(DocumentBase):
    id: int
    file_path: str
    file_size: int
    version: str
    uploader_id: int
    created_at: datetime
    
    class Config:
        from_attributes = True