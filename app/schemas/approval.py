from datetime import datetime
from typing import Optional
from pydantic import BaseModel

class ApprovalRecordBase(BaseModel):
    contract_id: int
    status: str = "pending"

class ApprovalRecordCreate(ApprovalRecordBase):
    comment: Optional[str] = None

class ApprovalRecordUpdate(BaseModel):
    status: str
    comment: Optional[str] = None

class ApprovalRecordResponse(ApprovalRecordBase):
    id: int
    approver_id: int
    comment: Optional[str] = None
    approved_at: Optional[datetime] = None
    created_at: datetime
    
    class Config:
        from_attributes = True

class ReminderBase(BaseModel):
    contract_id: int
    type: str
    reminder_date: datetime
    days_before: int

class ReminderCreate(ReminderBase):
    pass

class ReminderResponse(ReminderBase):
    id: int
    is_sent: bool
    sent_at: Optional[datetime] = None
    created_at: datetime
    
    class Config:
        from_attributes = True

class StatisticsResponse(BaseModel):
    total_contracts: int
    active_contracts: int
    pending_contracts: int
    completed_contracts: int
    total_amount: float
    this_month_contracts: int
    this_month_amount: float
    expiring_soon: int