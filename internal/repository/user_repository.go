package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	db *sql.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.SetTimestamps()

	query := `
		INSERT INTO user_tab (username, email, password_hash, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.UserID = uint(lastID)
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE user_id = ?
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE username = ?
	`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE email = ?
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	user.SetTimestamps()

	query := `
		UPDATE user_tab
		SET username = ?, email = ?, password_hash = ?, updated_at = ?,
			last_login = ?, is_active = ?, avatar_url = ?, bio = ?,
			birth_date = ?, gender = ?, country = ?
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash,
		user.UpdatedAt, user.LastLogin, user.IsActive, user.AvatarURL,
		user.Bio, user.BirthDate, user.Gender, user.Country, user.UserID)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now().UnixMilli()
	query := `UPDATE user_tab SET last_login = ? WHERE user_id = ?`

	_, err := r.db.ExecContext(ctx, query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, userID uint) error {
	query := `DELETE FROM user_tab WHERE user_id = ?`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tab WHERE username = ?`

	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username exists: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tab WHERE email = ?`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email exists: %w", err)
	}

	return count > 0, nil
}
