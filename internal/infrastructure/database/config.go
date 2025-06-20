package database

import (
	"context"
	"fmt"
	"time"

	config "github.com/Azzurriii/slythr/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
}

type ConnectionManager struct {
	postgres *gorm.DB
	redis    *redis.Client
}

func NewDatabaseConfig(cfg *config.Config) *DatabaseConfig {
	return &DatabaseConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.Name,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: 25,
		MaxIdleConns: 5,
		MaxLifetime:  5 * time.Minute,
	}
}

func NewRedisConfig(cfg *config.Config) *RedisConfig {
	return &RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
	}
}

func NewConnectionManager(dbConfig *DatabaseConfig, redisConfig *RedisConfig) (*ConnectionManager, error) {
	postgres, err := newPostgresConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	redisClient, err := newRedisConnection(redisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	return &ConnectionManager{
		postgres: postgres,
		redis:    redisClient,
	}, nil
}

func (cm *ConnectionManager) GetPostgres() *gorm.DB {
	return cm.postgres
}

func (cm *ConnectionManager) GetRedis() *redis.Client {
	return cm.redis
}

func (cm *ConnectionManager) Close() error {
	var errors []error

	if cm.postgres != nil {
		if sqlDB, err := cm.postgres.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close postgres: %w", err))
			}
		}
	}

	if cm.redis != nil {
		if err := cm.redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close redis: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	return nil
}

func (cm *ConnectionManager) HealthCheck(ctx context.Context) error {
	if err := cm.checkPostgresHealth(ctx); err != nil {
		return fmt.Errorf("postgres health check failed: %w", err)
	}

	if err := cm.checkRedisHealth(ctx); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

func newPostgresConnection(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		CreateBatchSize: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

func newRedisConnection(config *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

func (cm *ConnectionManager) checkPostgresHealth(ctx context.Context) error {
	sqlDB, err := cm.postgres.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (cm *ConnectionManager) checkRedisHealth(ctx context.Context) error {
	return cm.redis.Ping(ctx).Err()
}
