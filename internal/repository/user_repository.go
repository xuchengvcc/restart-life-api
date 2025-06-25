package repository

import (
	"context"
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
	return r.userDAO.Insert(ctx, user) // 直接返回DAO层错误
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	return r.userDAO.SelectByID(ctx, userID) // 直接返回DAO层结果
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.userDAO.SelectByUsername(ctx, username) // 直接返回DAO层结果
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.userDAO.SelectByEmail(ctx, email) // 直接返回DAO层结果
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	// 设置更新时间戳
	user.SetTimestamps()

	// 调用 DAO 层进行数据更新
	return r.userDAO.Update(ctx, user)
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now().UnixMilli()

	// 调用 DAO 层进行最后登录时间更新
	return r.userDAO.UpdateLastLogin(ctx, userID, now)
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, userID uint) error {
	// 调用 DAO 层进行用户删除
	return r.userDAO.Delete(ctx, userID)
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.userDAO.CountByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.userDAO.CountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
