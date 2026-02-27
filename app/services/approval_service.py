from typing import List, Optional
from sqlalchemy.orm import Session
from fastapi import HTTPException, status
from datetime import datetime, date, timedelta

from app.models.models import ApprovalRecord, Reminder, ApprovalStatus
from app.schemas.approval import ApprovalRecordCreate, ApprovalRecordUpdate, ReminderCreate

def get_approval_record(db: Session, record_id: int) -> Optional[ApprovalRecord]:
    return db.query(ApprovalRecord).filter(ApprovalRecord.id == record_id).first()

def get_approval_records(db: Session, contract_id: int) -> List[ApprovalRecord]:
    return db.query(ApprovalRecord).filter(
        ApprovalRecord.contract_id == contract_id
    ).order_by(ApprovalRecord.created_at.desc()).all()

def create_approval_record(db: Session, approval: ApprovalRecordCreate, approver_id: int) -> ApprovalRecord:
    db_approval = ApprovalRecord(
        approver_id=approver_id,
        **approval.model_dump()
    )
    db.add(db_approval)
    db.commit()
    db.refresh(db_approval)
    return db_approval

def update_approval_record(db: Session, record_id: int, approval_update: ApprovalRecordUpdate) -> Optional[ApprovalRecord]:
    db_approval = get_approval_record(db, record_id)
    if not db_approval:
        return None
    
    if db_approval.status != ApprovalStatus.PENDING:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="This approval has already been processed"
        )
    
    db_approval.status = approval_update.status
    db_approval.comment = approval_update.comment
    db_approval.approved_at = datetime.utcnow()
    
    db.commit()
    db.refresh(db_approval)
    return db_approval

def get_reminder(db: Session, reminder_id: int) -> Optional[Reminder]:
    return db.query(Reminder).filter(Reminder.id == reminder_id).first()

def get_reminders(db: Session, contract_id: int) -> List[Reminder]:
    return db.query(Reminder).filter(
        Reminder.contract_id == contract_id
    ).order_by(Reminder.reminder_date.desc()).all()

def create_reminder(db: Session, reminder: ReminderCreate) -> Reminder:
    db_reminder = Reminder(**reminder.model_dump())
    db.add(db_reminder)
    db.commit()
    db.refresh(db_reminder)
    return db_reminder

def get_expiring_contracts(db: Session, days: int = 30) -> List:
    from app.models.models import Contract, ContractStatus
    expiry_date = date.today() + timedelta(days=days)
    return db.query(Contract).filter(
        Contract.end_date <= expiry_date,
        Contract.end_date >= date.today(),
        Contract.status == ContractStatus.ACTIVE
    ).all()

def update_reminder_sent(db: Session, reminder_id: int) -> bool:
    db_reminder = get_reminder(db, reminder_id)
    if not db_reminder:
        return False
    
    db_reminder.is_sent = True
    db_reminder.sent_at = datetime.utcnow()
    db.commit()
    return True

def get_statistics(db: Session) -> dict:
    from app.models.models import Contract, ContractStatus
    from sqlalchemy import func
    
    today = date.today()
    this_month_start = date(today.year, today.month, 1)
    
    total_contracts = db.query(Contract).count()
    active_contracts = db.query(Contract).filter(Contract.status == ContractStatus.ACTIVE).count()
    pending_contracts = db.query(Contract).filter(Contract.status == ContractStatus.PENDING).count()
    completed_contracts = db.query(Contract).filter(Contract.status == ContractStatus.COMPLETED).count()
    
    total_amount = db.query(func.sum(Contract.amount)).filter(
        Contract.amount.isnot(None)
    ).scalar() or 0
    
    this_month_contracts = db.query(Contract).filter(
        Contract.created_at >= this_month_start
    ).count()
    
    this_month_amount = db.query(func.sum(Contract.amount)).filter(
        Contract.created_at >= this_month_start,
        Contract.amount.isnot(None)
    ).scalar() or 0
    
    expiring_soon = len(get_expiring_contracts(db, days=30))
    
    return {
        "total_contracts": total_contracts,
        "active_contracts": active_contracts,
        "pending_contracts": pending_contracts,
        "completed_contracts": completed_contracts,
        "total_amount": float(total_amount),
        "this_month_contracts": this_month_contracts,
        "this_month_amount": float(this_month_amount),
        "expiring_soon": expiring_soon
    }