package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// PreviewToken 预览令牌
type PreviewToken struct {
	Token      string
	DocumentID uint
	UserID     uint
	ExpiresAt  time.Time
}

// PreviewService 预览令牌服务
// 使用内存存储，生产环境应使用Redis
type PreviewService struct {
	tokens map[string]*PreviewToken
	mu     sync.RWMutex
}

func NewPreviewService() *PreviewService {
	service := &PreviewService{
		tokens: make(map[string]*PreviewToken),
	}

	// 清理过期令牌
	go service.cleanupExpired()

	return service
}

// GenerateToken 生成预览令牌
func (s *PreviewService) GenerateToken(documentID, userID uint) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)

	s.mu.Lock()
	s.tokens[token] = &PreviewToken{
		Token:      token,
		DocumentID: documentID,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(5 * time.Minute), // 5分钟有效
	}
	s.mu.Unlock()

	return token, nil
}

// ValidateToken 验证预览令牌
func (s *PreviewService) ValidateToken(token string) (*PreviewToken, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	previewToken, exists := s.tokens[token]
	if !exists {
		return nil, false
	}

	if time.Now().After(previewToken.ExpiresAt) {
		return nil, false
	}

	return previewToken, true
}

// RevokeToken 撤销令牌
func (s *PreviewService) RevokeToken(token string) {
	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()
}

// cleanupExpired 清理过期令牌
func (s *PreviewService) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for token, previewToken := range s.tokens {
			if time.Now().After(previewToken.ExpiresAt) {
				delete(s.tokens, token)
			}
		}
		s.mu.Unlock()
	}
}
