from datetime import datetime
from typing import Optional
from pydantic import BaseModel

class CustomerBase(BaseModel):
    name: str
    type: str = "customer"
    contact_person: Optional[str] = None
    contact_phone: Optional[str] = None
    contact_email: Optional[str] = None
    address: Optional[str] = None
    credit_rating: Optional[str] = None

class CustomerCreate(CustomerBase):
    code: str

class CustomerUpdate(BaseModel):
    name: Optional[str] = None
    type: Optional[str] = None
    contact_person: Optional[str] = None
    contact_phone: Optional[str] = None
    contact_email: Optional[str] = None
    address: Optional[str] = None
    credit_rating: Optional[str] = None
    is_active: Optional[bool] = None

class CustomerResponse(CustomerBase):
    id: int
    code: str
    is_active: bool
    created_at: datetime
    
    class Config:
        from_attributes = True

class ContractTypeBase(BaseModel):
    name: str
    description: Optional[str] = None

class ContractTypeCreate(ContractTypeBase):
    code: str

class ContractTypeResponse(ContractTypeBase):
    id: int
    code: str
    created_at: datetime
    
    class Config:
        from_attributes = True