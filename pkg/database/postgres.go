package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// PostgresConfig PostgreSQL配置
type PostgresConfig struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// PostgresDB PostgreSQL数据库连接
type PostgresDB struct {
	DB     *sql.DB
	config *PostgresConfig
}

// NewPostgresDB 创建新的PostgreSQL连接
func NewPostgresDB(config *PostgresConfig) (*PostgresDB, error) {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)

	// 打开数据库连接
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
	}).Info("Successfully connected to PostgreSQL")

	return &PostgresDB{
		DB:     db,
		config: config,
	}, nil
}

// HealthCheck 检查数据库连接健康状态
func (p *PostgresDB) HealthCheck() error {
	if err := p.DB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Close 关闭数据库连接
func (p *PostgresDB) Close() error {
	if p.DB != nil {
		logrus.Info("Closing PostgreSQL connection")
		return p.DB.Close()
	}
	return nil
}

// GetStats 获取连接池统计信息
func (p *PostgresDB) GetStats() sql.DBStats {
	return p.DB.Stats()
}

// Begin 开始事务
func (p *PostgresDB) Begin() (*sql.Tx, error) {
	return p.DB.Begin()
}

// Exec 执行SQL语句
func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

// Query 查询数据
func (p *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(query, args...)
}

// QueryRow 查询单行数据
func (p *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRow(query, args...)
}
