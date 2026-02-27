from typing import List, Optional
from sqlalchemy.orm import Session
from fastapi import HTTPException, status

from app.models.models import Customer, ContractType
from app.schemas.customer import CustomerCreate, CustomerUpdate, ContractTypeCreate

def get_customer(db: Session, customer_id: int) -> Optional[Customer]:
    return db.query(Customer).filter(Customer.id == customer_id).first()

def get_customer_by_code(db: Session, code: str) -> Optional[Customer]:
    return db.query(Customer).filter(Customer.code == code).first()

def get_customers(db: Session, skip: int = 0, limit: int = 100, type: Optional[str] = None) -> List[Customer]:
    query = db.query(Customer)
    if type:
        query = query.filter(Customer.type == type)
    return query.offset(skip).limit(limit).all()

def create_customer(db: Session, customer: CustomerCreate) -> Customer:
    if get_customer_by_code(db, customer.code):
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Customer code already exists"
        )
    
    db_customer = Customer(**customer.model_dump())
    db.add(db_customer)
    db.commit()
    db.refresh(db_customer)
    return db_customer

def update_customer(db: Session, customer_id: int, customer_update: CustomerUpdate) -> Optional[Customer]:
    db_customer = get_customer(db, customer_id)
    if not db_customer:
        return None
    
    update_data = customer_update.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        setattr(db_customer, field, value)
    
    db.commit()
    db.refresh(db_customer)
    return db_customer

def delete_customer(db: Session, customer_id: int) -> bool:
    db_customer = get_customer(db, customer_id)
    if not db_customer:
        return False
    
    db.delete(db_customer)
    db.commit()
    return True

def get_contract_type(db: Session, type_id: int) -> Optional[ContractType]:
    return db.query(ContractType).filter(ContractType.id == type_id).first()

def get_contract_types(db: Session, skip: int = 0, limit: int = 100) -> List[ContractType]:
    return db.query(ContractType).offset(skip).limit(limit).all()

def create_contract_type(db: Session, contract_type: ContractTypeCreate) -> ContractType:
    db_contract_type = ContractType(**contract_type.model_dump())
    db.add(db_contract_type)
    db.commit()
    db.refresh(db_contract_type)
    return db_contract_type