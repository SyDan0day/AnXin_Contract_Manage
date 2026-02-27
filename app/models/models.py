from sqlalchemy import Column, Integer, String, DateTime, Text, Boolean, ForeignKey, Float, Date, Enum
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
import enum

from app.core.database import Base

class UserRole(str, enum.Enum):
    ADMIN = "admin"
    MANAGER = "manager"
    USER = "user"

class User(Base):
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(50), unique=True, index=True, nullable=False)
    email = Column(String(100), unique=True, index=True)
    hashed_password = Column(String(200), nullable=False)
    full_name = Column(String(100))
    role = Column(Enum(UserRole), default=UserRole.USER)
    department = Column(String(100))
    phone = Column(String(20))
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
    
    contracts = relationship("Contract", back_populates="creator")
    approval_records = relationship("ApprovalRecord", back_populates="approver")

class Role(Base):
    __tablename__ = "roles"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(50), unique=True, nullable=False)
    description = Column(Text)
    permissions = Column(Text)
    created_at = Column(DateTime, server_default=func.now())

class Customer(Base):
    __tablename__ = "customers"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(200), nullable=False, index=True)
    type = Column(String(20), default="customer")
    code = Column(String(50), unique=True, index=True)
    contact_person = Column(String(100))
    contact_phone = Column(String(20))
    contact_email = Column(String(100))
    address = Column(Text)
    credit_rating = Column(String(20))
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
    
    contracts = relationship("Contract", back_populates="customer")

class ContractType(Base):
    __tablename__ = "contract_types"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(100), unique=True, nullable=False)
    code = Column(String(50), unique=True)
    description = Column(Text)
    created_at = Column(DateTime, server_default=func.now())

class ContractStatus(str, enum.Enum):
    DRAFT = "draft"
    PENDING = "pending"
    APPROVED = "approved"
    ACTIVE = "active"
    COMPLETED = "completed"
    TERMINATED = "terminated"

class Contract(Base):
    __tablename__ = "contracts"
    
    id = Column(Integer, primary_key=True, index=True)
    contract_no = Column(String(50), unique=True, index=True, nullable=False)
    title = Column(String(200), nullable=False, index=True)
    customer_id = Column(Integer, ForeignKey("customers.id"))
    contract_type_id = Column(Integer, ForeignKey("contract_types.id"))
    amount = Column(Float)
    currency = Column(String(10), default="CNY")
    status = Column(Enum(ContractStatus), default=ContractStatus.DRAFT)
    sign_date = Column(Date)
    start_date = Column(Date)
    end_date = Column(Date)
    payment_terms = Column(Text)
    content = Column(Text)
    notes = Column(Text)
    creator_id = Column(Integer, ForeignKey("users.id"))
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
    
    customer = relationship("Customer", back_populates="contracts")
    creator = relationship("User", back_populates="contracts")
    contract_type = relationship("ContractType")
    executions = relationship("ContractExecution", back_populates="contract")
    documents = relationship("Document", back_populates="contract")
    approval_records = relationship("ApprovalRecord", back_populates="contract")

class ContractExecution(Base):
    __tablename__ = "contract_executions"
    
    id = Column(Integer, primary_key=True, index=True)
    contract_id = Column(Integer, ForeignKey("contracts.id"))
    stage = Column(String(100))
    stage_date = Column(Date)
    progress = Column(Float, default=0)
    payment_amount = Column(Float)
    payment_date = Column(Date)
    description = Column(Text)
    operator_id = Column(Integer, ForeignKey("users.id"))
    created_at = Column(DateTime, server_default=func.now())
    
    contract = relationship("Contract", back_populates="executions")

class ApprovalStatus(str, enum.Enum):
    PENDING = "pending"
    APPROVED = "approved"
    REJECTED = "rejected"

class ApprovalRecord(Base):
    __tablename__ = "approval_records"
    
    id = Column(Integer, primary_key=True, index=True)
    contract_id = Column(Integer, ForeignKey("contracts.id"))
    approver_id = Column(Integer, ForeignKey("users.id"))
    status = Column(Enum(ApprovalStatus), default=ApprovalStatus.PENDING)
    comment = Column(Text)
    approved_at = Column(DateTime)
    created_at = Column(DateTime, server_default=func.now())
    
    contract = relationship("Contract", back_populates="approval_records")
    approver = relationship("User", back_populates="approval_records")

class Document(Base):
    __tablename__ = "documents"
    
    id = Column(Integer, primary_key=True, index=True)
    contract_id = Column(Integer, ForeignKey("contracts.id"))
    name = Column(String(200))
    file_path = Column(String(500))
    file_size = Column(Integer)
    file_type = Column(String(50))
    version = Column(String(20), default="1.0")
    uploader_id = Column(Integer, ForeignKey("users.id"))
    created_at = Column(DateTime, server_default=func.now())
    
    contract = relationship("Contract", back_populates="documents")

class Reminder(Base):
    __tablename__ = "reminders"
    
    id = Column(Integer, primary_key=True, index=True)
    contract_id = Column(Integer, ForeignKey("contracts.id"))
    type = Column(String(50))
    reminder_date = Column(Date)
    days_before = Column(Integer)
    is_sent = Column(Boolean, default=False)
    sent_at = Column(DateTime)
    created_at = Column(DateTime, server_default=func.now())