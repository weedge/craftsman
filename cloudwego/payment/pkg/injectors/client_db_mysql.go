package client

import (
	"time"

	kitexGorm "github.com/weedge/craftsman/cloudwego/kitex-contrib/gorm"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/utils/logutils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlDBClientOptions struct {
	Name                 string         `mapstructure:"name"`
	DbEngine             string         `mapstructure:"dbEngine"`
	DSN                  string         `mapstructure:"dsn"`
	MaxIdleConns         int            `mapstructure:"maxidleconns"`
	MaxOpenConns         int            `mapstructure:"maxopenconns"`
	ConnMaxLifeTime      time.Duration  `mapstructure:"connMaxLifeTime"`
	SlowSqlTimeThreshold time.Duration  `mapstructure:"slowSqlTimeThreshold"`
	TraceLogLevel        logutils.Level `mapstructure:"traceLogLevel"`
}

func DefaultMysqlDBClientOptions() *MysqlDBClientOptions {
	return &MysqlDBClientOptions{
		Name:            "default",
		DbEngine:        "mysql-innodb",
		DSN:             "",
		MaxIdleConns:    10,
		MaxOpenConns:    1000,
		ConnMaxLifeTime: 3600 * time.Second,
	}
}

func InitMysqlDBClient(opts *MysqlDBClientOptions, kvLogger kitexGorm.IkvLogger) (dbClient *gorm.DB, err error) {

	conf := &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,

		//Logger:                 logger.Default.LogMode(logger.Info),
		Logger: kitexGorm.NewGormLogger(
			kitexGorm.WithKvLogger(kvLogger),
			kitexGorm.WithSlowThreshold(opts.SlowSqlTimeThreshold),
			kitexGorm.WithTraceLogLevel(opts.TraceLogLevel.KitexLogLevel()),
		),

		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           map[string]clause.ClauseBuilder{},
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  map[string]gorm.Plugin{},
	}

	dbClient, err = gorm.Open(mysql.Open(opts.DSN), conf)
	if err != nil {
		return
	}
	sqlDB, err := dbClient.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.ConnMaxLifeTime)

	return
}
