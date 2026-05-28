package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/account/internal/dal"
)

// AccountService 账号服务
type AccountService struct {
	accountRepo *dal.AccountRepository
}

// NewAccountService 创建账号服务
func NewAccountService(accountRepo *dal.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

// CreateAccount 创建账号
func (s *AccountService) CreateAccount(ctx context.Context, userID int64, platform, accountName, accountID string) (*dal.Account, error) {
	account := &dal.Account{
		UserID:       userID,
		Platform:     platform,
		AccountName:  accountName,
		AccountIDStr: accountID,
		Status:       1,
		HealthScore:  100.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccount 获取账号
func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (*dal.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// UpdateAccount 更新账号
func (s *AccountService) UpdateAccount(ctx context.Context, accountID int64, accountName string, status int32) (*dal.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if accountName != "" {
		account.AccountName = accountName
	}
	if status > 0 {
		account.Status = status
	}
	account.UpdatedAt = time.Now()

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

// DeleteAccount 删除账号
func (s *AccountService) DeleteAccount(ctx context.Context, accountID int64) error {
	if err := s.accountRepo.Delete(ctx, accountID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// ListAccounts 列出账号
func (s *AccountService) ListAccounts(ctx context.Context, userID int64, platform string, page, pageSize int32) ([]*dal.Account, int32, error) {
	accounts, total, err := s.accountRepo.List(ctx, userID, platform, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, total, nil
}

// GetAccountHealth 获取账号健康状态
func (s *AccountService) GetAccountHealth(ctx context.Context, accountID int64) (*dal.AccountHealth, error) {
	health, err := s.accountRepo.GetHealth(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account health: %w", err)
	}

	return health, nil
}