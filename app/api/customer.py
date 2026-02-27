from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from app.core.database import get_db
from app.schemas.customer import CustomerCreate, CustomerUpdate, CustomerResponse, ContractTypeCreate, ContractTypeResponse
from app.services.customer_service import (
    get_customer, get_customers, create_customer, 
    update_customer, delete_customer,
    get_contract_type, get_contract_types, create_contract_type
)

router = APIRouter()

@router.get("/customers", response_model=List[CustomerResponse])
def get_all_customers(skip: int = 0, limit: int = 100, type: Optional[str] = None, db: Session = Depends(get_db)):
    customers = get_customers(db, skip=skip, limit=limit, type=type)
    return customers

@router.get("/customers/{customer_id}", response_model=CustomerResponse)
def get_customer_by_id(customer_id: int, db: Session = Depends(get_db)):
    customer = get_customer(db, customer_id)
    if not customer:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Customer not found")
    return customer

@router.post("/customers", response_model=CustomerResponse, status_code=status.HTTP_201_CREATED)
def create_customer_endpoint(customer: CustomerCreate, db: Session = Depends(get_db)):
    db_customer = create_customer(db, customer)
    return db_customer

@router.put("/customers/{customer_id}", response_model=CustomerResponse)
def update_customer_endpoint(customer_id: int, customer_update: CustomerUpdate, db: Session = Depends(get_db)):
    customer = update_customer(db, customer_id, customer_update)
    if not customer:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Customer not found")
    return customer

@router.delete("/customers/{customer_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_customer_endpoint(customer_id: int, db: Session = Depends(get_db)):
    success = delete_customer(db, customer_id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Customer not found")
    return None

@router.get("/contract-types", response_model=List[ContractTypeResponse])
def get_all_contract_types(skip: int = 0, limit: int = 100, db: Session = Depends(get_db)):
    contract_types = get_contract_types(db, skip=skip, limit=limit)
    return contract_types

@router.post("/contract-types", response_model=ContractTypeResponse, status_code=status.HTTP_201_CREATED)
def create_contract_type_endpoint(contract_type: ContractTypeCreate, db: Session = Depends(get_db)):
    db_contract_type = create_contract_type(db, contract_type)
    return db_contract_type