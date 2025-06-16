package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	Charset      string
	ParseTime    bool
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// MySQLDB MySQL数据库连接
type MySQLDB struct {
	DB     *sql.DB
	config *MySQLConfig
}

// NewMySQLDB 创建新的MySQL连接
func NewMySQLDB(config *MySQLConfig) (*MySQLDB, error) {
	// 设置默认值
	if config.Charset == "" {
		config.Charset = "utf8mb4"
	}
	if !config.ParseTime {
		config.ParseTime = true
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		config.Username, config.Password, config.Host, config.Port,
		config.Database, config.Charset, config.ParseTime)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
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
	}).Info("Successfully connected to MySQL")

	return &MySQLDB{
		DB:     db,
		config: config,
	}, nil
}

// HealthCheck 检查数据库连接健康状态
func (m *MySQLDB) HealthCheck() error {
	if err := m.DB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Close 关闭数据库连接
func (m *MySQLDB) Close() error {
	if m.DB != nil {
		logrus.Info("Closing MySQL connection")
		return m.DB.Close()
	}
	return nil
}

// GetStats 获取连接池统计信息
func (m *MySQLDB) GetStats() sql.DBStats {
	return m.DB.Stats()
}

// Begin 开始事务
func (m *MySQLDB) Begin() (*sql.Tx, error) {
	return m.DB.Begin()
}

// Exec 执行SQL语句
func (m *MySQLDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.DB.Exec(query, args...)
}

// Query 查询数据
func (m *MySQLDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.DB.Query(query, args...)
}

// QueryRow 查询单行数据
func (m *MySQLDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.DB.QueryRow(query, args...)
}
