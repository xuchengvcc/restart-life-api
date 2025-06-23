package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// RedisDB Redis数据库连接
type RedisDB struct {
	Client *redis.Client
	config *RedisConfig
}

// NewRedisDB 创建新的Redis连接
func NewRedisDB(config *RedisConfig) (*RedisDB, error) {
	// 设置默认值
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
	}).Info("Successfully connected to Redis")

	return &RedisDB{
		Client: client,
		config: config,
	}, nil
}

// InitRedisFromConfig 根据全局配置初始化 Redis 连接
func InitRedisFromConfig(cfg *config.Config) (*RedisDB, error) {
	return NewRedisDB(&RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		Database:     cfg.Redis.Database,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
}

// HealthCheck 检查Redis连接健康状态
func (r *RedisDB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}

// Close 关闭Redis连接
func (r *RedisDB) Close() error {
	if r.Client != nil {
		logrus.Info("Closing Redis connection")
		return r.Client.Close()
	}
	return nil
}

// GetStats 获取Redis连接统计信息
func (r *RedisDB) GetStats() *redis.PoolStats {
	return r.Client.PoolStats()
}

// Set 设置键值对
func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisDB) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Del(ctx, keys...).Result()
}

// Exists 检查键是否存在
func (r *RedisDB) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Expire 设置键的过期时间
func (r *RedisDB) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(ctx, key, expiration).Result()
}

// HSet 设置哈希字段
func (r *RedisDB) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.Client.HSet(ctx, key, values...).Result()
}

// HGet 获取哈希字段值
func (r *RedisDB) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希的所有字段和值
func (r *RedisDB) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisDB) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return r.Client.HDel(ctx, key, fields...).Result()
}
