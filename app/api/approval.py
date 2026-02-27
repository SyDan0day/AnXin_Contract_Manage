from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from app.core.database import get_db
from app.schemas.approval import ApprovalRecordCreate, ApprovalRecordUpdate, ApprovalRecordResponse, ReminderCreate, ReminderResponse, StatisticsResponse
from app.services.approval_service import (
    get_approval_record, get_approval_records, create_approval_record, update_approval_record,
    get_reminder, get_reminders, create_reminder, update_reminder_sent, get_expiring_contracts, get_statistics
)

router = APIRouter()

@router.get("/contracts/{contract_id}/approvals", response_model=List[ApprovalRecordResponse])
def get_contract_approvals(contract_id: int, db: Session = Depends(get_db)):
    approvals = get_approval_records(db, contract_id)
    return approvals

@router.post("/contracts/{contract_id}/approvals", response_model=ApprovalRecordResponse, status_code=status.HTTP_201_CREATED)
def create_approval_endpoint(contract_id: int, approval: ApprovalRecordCreate, db: Session = Depends(get_db)):
    db_approval = create_approval_record(db, approval, approver_id=1)
    return db_approval

@router.put("/approvals/{approval_id}", response_model=ApprovalRecordResponse)
def update_approval_endpoint(approval_id: int, approval_update: ApprovalRecordUpdate, db: Session = Depends(get_db)):
    approval = update_approval_record(db, approval_id, approval_update)
    if not approval:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Approval record not found")
    return approval

@router.get("/contracts/{contract_id}/reminders", response_model=List[ReminderResponse])
def get_contract_reminders(contract_id: int, db: Session = Depends(get_db)):
    reminders = get_reminders(db, contract_id)
    return reminders

@router.post("/contracts/{contract_id}/reminders", response_model=ReminderResponse, status_code=status.HTTP_201_CREATED)
def create_reminder_endpoint(contract_id: int, reminder: ReminderCreate, db: Session = Depends(get_db)):
    db_reminder = create_reminder(db, reminder)
    return db_reminder

@router.post("/reminders/{reminder_id}/send", status_code=status.HTTP_200_OK)
def send_reminder_endpoint(reminder_id: int, db: Session = Depends(get_db)):
    success = update_reminder_sent(db, reminder_id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Reminder not found")
    return {"message": "Reminder sent successfully"}

@router.get("/expiring-contracts")
def get_expiring_contracts_endpoint(days: int = 30, db: Session = Depends(get_db)):
    contracts = get_expiring_contracts(db, days)
    return {"contracts": contracts, "days": days}

@router.get("/statistics", response_model=StatisticsResponse)
def get_statistics_endpoint(db: Session = Depends(get_db)):
    stats = get_statistics(db)
    return stats