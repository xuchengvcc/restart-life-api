package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/xuchengvcc/restart-life-api/internal/dao"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID uint) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID uint) error
	Delete(ctx context.Context, userID uint) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// userRepository MySQL用户仓储实现
type userRepository struct {
	userDAO dao.UserDAO
}

// NewUserRepository 创建用户仓储
func NewUserRepository(userDAO dao.UserDAO) UserRepository {
	return &userRepository{userDAO: userDAO}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// 设置时间戳
	user.SetTimestamps()

	// 调用 DAO 层进行数据插入
	err := r.userDAO.Insert(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := r.userDAO.SelectByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := r.userDAO.SelectByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.userDAO.SelectByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	// 设置更新时间戳
	user.SetTimestamps()

	// 调用 DAO 层进行数据更新
	err := r.userDAO.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now().UnixMilli()

	// 调用 DAO 层进行最后登录时间更新
	err := r.userDAO.UpdateLastLogin(ctx, userID, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, userID uint) error {
	// 调用 DAO 层进行用户删除
	err := r.userDAO.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.userDAO.CountByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check username exists: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.userDAO.CountByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email exists: %w", err)
	}

	return count > 0, nil
}
