from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from app.core.database import get_db
from app.schemas.contract import ContractCreate, ContractUpdate, ContractResponse, ContractExecutionCreate, ContractExecutionResponse, DocumentCreate, DocumentResponse
from app.services.contract_service import (
    get_contract, get_contracts, create_contract, 
    update_contract, delete_contract,
    get_contract_execution, get_contract_executions, create_contract_execution,
    get_document, get_documents, create_document, delete_document
)

router = APIRouter()

@router.get("/contracts", response_model=List[ContractResponse])
def get_all_contracts(
    skip: int = 0, 
    limit: int = 100,
    customer_id: Optional[int] = None,
    contract_type_id: Optional[int] = None,
    status: Optional[str] = None,
    db: Session = Depends(get_db)
):
    contracts = get_contracts(db, skip=skip, limit=limit, customer_id=customer_id, contract_type_id=contract_type_id, status=status)
    return contracts

@router.get("/contracts/{contract_id}", response_model=ContractResponse)
def get_contract_by_id(contract_id: int, db: Session = Depends(get_db)):
    contract = get_contract(db, contract_id)
    if not contract:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Contract not found")
    return contract

@router.post("/contracts", response_model=ContractResponse, status_code=status.HTTP_201_CREATED)
def create_contract_endpoint(contract: ContractCreate, db: Session = Depends(get_db)):
    db_contract = create_contract(db, contract, creator_id=1)
    return db_contract

@router.put("/contracts/{contract_id}", response_model=ContractResponse)
def update_contract_endpoint(contract_id: int, contract_update: ContractUpdate, db: Session = Depends(get_db)):
    contract = update_contract(db, contract_id, contract_update)
    if not contract:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Contract not found")
    return contract

@router.delete("/contracts/{contract_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_contract_endpoint(contract_id: int, db: Session = Depends(get_db)):
    success = delete_contract(db, contract_id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Contract not found")
    return None

@router.get("/contracts/{contract_id}/executions", response_model=List[ContractExecutionResponse])
def get_contract_executions_endpoint(contract_id: int, db: Session = Depends(get_db)):
    executions = get_contract_executions(db, contract_id)
    return executions

@router.post("/contracts/{contract_id}/executions", response_model=ContractExecutionResponse, status_code=status.HTTP_201_CREATED)
def create_contract_execution_endpoint(execution: ContractExecutionCreate, db: Session = Depends(get_db)):
    db_execution = create_contract_execution(db, execution, operator_id=1)
    return db_execution

@router.get("/contracts/{contract_id}/documents", response_model=List[DocumentResponse])
def get_contract_documents_endpoint(contract_id: int, db: Session = Depends(get_db)):
    documents = get_documents(db, contract_id)
    return documents

@router.post("/contracts/{contract_id}/documents", response_model=DocumentResponse, status_code=status.HTTP_201_CREATED)
def create_contract_document_endpoint(document: DocumentCreate, db: Session = Depends(get_db)):
    db_document = create_document(db, document, uploader_id=1)
    return db_document

@router.delete("/documents/{document_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_document_endpoint(document_id: int, db: Session = Depends(get_db)):
    success = delete_document(db, document_id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Document not found")
    return None