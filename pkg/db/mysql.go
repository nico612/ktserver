package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// MySQLOptions defines optsions for mysql database.
type MySQLOptions struct {
	Host                  string           // 主机地址
	Username              string           // 用户名
	Password              string           // 密码
	Database              string           // 数据库名
	MaxIdleConnections    int              // 最大空闲连接数
	MaxOpenConnections    int              // 最大打开连接数
	MaxConnectionLifeTime time.Duration    // 连接最大生命周期
	Logger                logger.Interface // +optional 实现gorm.Logger接口
}

// DSN return DSN from MySQLOptions.
func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")
}

// NewMySQL create a new gorm db instance with the given options.
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	// Set default values to ensure all fields in opts are available.
	setDefaults(opts)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       opts.DSN(),
		DefaultStringSize:         200,  // string 类型字段的默认长度
		SkipInitializeWithVersion: true, // 根据版本自动配置
	}),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   opts.Logger,
			// PrepareStmt executes the given query in cached statement.
			// This can improve performance.
			PrepareStmt: true,
		},
	)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	return db, nil
}

func setDefaults(opts *MySQLOptions) {
	if opts.Host == "" {
		opts.Host = "127.0.0.1:3306"
	}
	if opts.MaxIdleConnections == 0 {
		opts.MaxIdleConnections = 100
	}
	if opts.MaxOpenConnections == 0 {
		opts.MaxOpenConnections = 100
	}
	if opts.MaxConnectionLifeTime == 0 {
		opts.MaxConnectionLifeTime = time.Duration(10) * time.Second
	}
	if opts.Logger == nil {
		opts.Logger = logger.Default
	}
}
