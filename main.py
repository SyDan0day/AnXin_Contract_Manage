from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import settings
from app.core.database import Base, engine
from app.api import auth, customer, contract, approval

Base.metadata.create_all(bind=engine)

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(auth.router, prefix="/api/auth", tags=["认证"])
app.include_router(customer.router, prefix="/api", tags=["客户管理"])
app.include_router(contract.router, prefix="/api", tags=["合同管理"])
app.include_router(approval.router, prefix="/api", tags=["审批与提醒"])

@app.get("/")
def read_root():
    return {"message": "合同管理系统 API", "version": settings.APP_VERSION}

@app.get("/health")
def health_check():
    return {"status": "healthy"}