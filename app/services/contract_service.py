from typing import List, Optional
from sqlalchemy.orm import Session
from fastapi import HTTPException, status
from datetime import datetime, date

from app.models.models import Contract, ContractExecution, Document, ContractStatus
from app.schemas.contract import ContractCreate, ContractUpdate, ContractExecutionCreate, DocumentCreate

def generate_contract_no(db: Session) -> str:
    from sqlalchemy import func
    today = datetime.now()
    prefix = f"CT{today.strftime('%Y%m')}"
    
    last_contract = db.query(Contract).filter(
        Contract.contract_no.like(f"{prefix}%")
    ).order_by(Contract.contract_no.desc()).first()
    
    if last_contract:
        last_no = int(last_contract.contract_no[-4:])
        new_no = str(last_no + 1).zfill(4)
    else:
        new_no = "0001"
    
    return f"{prefix}{new_no}"

def get_contract(db: Session, contract_id: int) -> Optional[Contract]:
    return db.query(Contract).filter(Contract.id == contract_id).first()

def get_contract_by_no(db: Session, contract_no: str) -> Optional[Contract]:
    return db.query(Contract).filter(Contract.contract_no == contract_no).first()

def get_contracts(
    db: Session, 
    skip: int = 0, 
    limit: int = 100,
    customer_id: Optional[int] = None,
    contract_type_id: Optional[int] = None,
    status: Optional[str] = None
) -> List[Contract]:
    query = db.query(Contract)
    if customer_id:
        query = query.filter(Contract.customer_id == customer_id)
    if contract_type_id:
        query = query.filter(Contract.contract_type_id == contract_type_id)
    if status:
        query = query.filter(Contract.status == status)
    return query.order_by(Contract.created_at.desc()).offset(skip).limit(limit).all()

def create_contract(db: Session, contract: ContractCreate, creator_id: int) -> Contract:
    contract_no = generate_contract_no(db)
    db_contract = Contract(
        contract_no=contract_no,
        creator_id=creator_id,
        status=ContractStatus.DRAFT,
        **contract.model_dump()
    )
    db.add(db_contract)
    db.commit()
    db.refresh(db_contract)
    return db_contract

def update_contract(db: Session, contract_id: int, contract_update: ContractUpdate) -> Optional[Contract]:
    db_contract = get_contract(db, contract_id)
    if not db_contract:
        return None
    
    update_data = contract_update.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        setattr(db_contract, field, value)
    
    db.commit()
    db.refresh(db_contract)
    return db_contract

def delete_contract(db: Session, contract_id: int) -> bool:
    db_contract = get_contract(db, contract_id)
    if not db_contract:
        return False
    
    db.delete(db_contract)
    db.commit()
    return True

def get_contract_execution(db: Session, execution_id: int) -> Optional[ContractExecution]:
    return db.query(ContractExecution).filter(ContractExecution.id == execution_id).first()

def get_contract_executions(db: Session, contract_id: int) -> List[ContractExecution]:
    return db.query(ContractExecution).filter(
        ContractExecution.contract_id == contract_id
    ).order_by(ContractExecution.created_at.desc()).all()

def create_contract_execution(db: Session, execution: ContractExecutionCreate, operator_id: int) -> ContractExecution:
    db_execution = ContractExecution(
        operator_id=operator_id,
        **execution.model_dump()
    )
    db.add(db_execution)
    db.commit()
    db.refresh(db_execution)
    
    contract = get_contract(db, execution.contract_id)
    if contract and execution.progress:
        contract.content = str(execution.progress)
        db.commit()
    
    return db_execution

def get_document(db: Session, document_id: int) -> Optional[Document]:
    return db.query(Document).filter(Document.id == document_id).first()

def get_documents(db: Session, contract_id: int) -> List[Document]:
    return db.query(Document).filter(
        Document.contract_id == contract_id
    ).order_by(Document.created_at.desc()).all()

def create_document(db: Session, document: DocumentCreate, uploader_id: int) -> Document:
    db_document = Document(
        uploader_id=uploader_id,
        **document.model_dump()
    )
    db.add(db_document)
    db.commit()
    db.refresh(db_document)
    return db_document

def delete_document(db: Session, document_id: int) -> bool:
    db_document = get_document(db, document_id)
    if not db_document:
        return False
    
    db.delete(db_document)
    db.commit()
    return True