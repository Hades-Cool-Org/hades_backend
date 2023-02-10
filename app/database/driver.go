package database

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hades_backend/app/config"
	"hades_backend/app/logging"
	"os"
)

var (
	// DB is the database connection
	DB *gorm.DB
)

func init() {
	DB = NewMySqlDB()
}

// NewMySqlDB creates a new MySQL database connection
func NewMySqlDB() *gorm.DB {
	cfg := config.Cfg

	l := logging.Initialize()

	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.Username
	pass := cfg.Database.Password
	dbName := cfg.Database.DbName

	l.Info("Connecting to database",
		zap.Field{Key: "host", Type: zapcore.StringType, String: host},
		zap.Field{Key: "port", Type: zapcore.StringType, String: port},
		zap.Field{Key: "user", Type: zapcore.StringType, String: user},
		zap.Field{Key: "db.name", Type: zapcore.StringType, String: dbName},
	)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		processError(err)
	}

	return db
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
