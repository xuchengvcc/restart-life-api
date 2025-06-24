package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xuchengvcc/restart-life-api/internal/models"
)

// UserDAO 用户数据访问对象接口
type UserDAO interface {
	Insert(ctx context.Context, user *models.User) error
	SelectByID(ctx context.Context, userID uint) (*models.User, error)
	SelectByUsername(ctx context.Context, username string) (*models.User, error)
	SelectByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID uint, lastLogin int64) error
	Delete(ctx context.Context, userID uint) error
	CountByUsername(ctx context.Context, username string) (int, error)
	CountByEmail(ctx context.Context, email string) (int, error)
}

// userDAO MySQL用户数据访问对象实现
type userDAO struct {
	db *sql.DB
}

// NewUserDAO 创建用户数据访问对象
func NewUserDAO(db *sql.DB) UserDAO {
	return &userDAO{db: db}
}

// Insert 插入用户
func (d *userDAO) Insert(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO user_tab (username, email, password_hash, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.ExecContext(ctx, query,
		user.Username, user.Email, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.UserID = uint(lastID)
	return nil
}

// SelectByID 根据ID查询用户
func (d *userDAO) SelectByID(ctx context.Context, userID uint) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE user_id = ?
	`

	err := d.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to select user by id: %w", err)
	}

	return user, nil
}

// SelectByUsername 根据用户名查询用户
func (d *userDAO) SelectByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE username = ?
	`

	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to select user by username: %w", err)
	}

	return user, nil
}

// SelectByEmail 根据邮箱查询用户
func (d *userDAO) SelectByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT user_id, username, email, password_hash, created_at, updated_at,
			   last_login, is_active, avatar_url, bio, birth_date, gender, country
		FROM user_tab
		WHERE email = ?
	`

	err := d.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.AvatarURL, &user.Bio, &user.BirthDate, &user.Gender, &user.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to select user by email: %w", err)
	}

	return user, nil
}

// Update 更新用户
func (d *userDAO) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE user_tab
		SET username = ?, email = ?, password_hash = ?, updated_at = ?,
			last_login = ?, is_active = ?, avatar_url = ?, bio = ?,
			birth_date = ?, gender = ?, country = ?
		WHERE user_id = ?
	`

	_, err := d.db.ExecContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.UpdatedAt,
		user.LastLogin, user.IsActive, user.AvatarURL, user.Bio,
		user.BirthDate, user.Gender, user.Country, user.UserID)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (d *userDAO) UpdateLastLogin(ctx context.Context, userID uint, lastLogin int64) error {
	query := `UPDATE user_tab SET last_login = ? WHERE user_id = ?`

	_, err := d.db.ExecContext(ctx, query, lastLogin, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Delete 删除用户
func (d *userDAO) Delete(ctx context.Context, userID uint) error {
	query := `DELETE FROM user_tab WHERE user_id = ?`

	_, err := d.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// CountByUsername 根据用户名统计数量
func (d *userDAO) CountByUsername(ctx context.Context, username string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tab WHERE username = ?`

	err := d.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count by username: %w", err)
	}

	return count, nil
}

// CountByEmail 根据邮箱统计数量
func (d *userDAO) CountByEmail(ctx context.Context, email string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tab WHERE email = ?`

	err := d.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count by email: %w", err)
	}

	return count, nil
}
