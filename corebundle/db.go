package corebundle

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var OrmDb *gorm.DB

func InitDb() *gorm.DB {
	log.Println("initializing database...")

	databaseLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", Config.Database.User, Config.Database.Password, Config.Database.Host, Config.Database.Port, Config.Database.Name)
	ormDb, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		Logger:                                   databaseLogger,
	})

	// retry connection if failed
	if err != nil {
		log.Println(err)
		time.Sleep(60 * time.Second)
		return InitDb()
	}

	log.Println("initializing database complete")

	return ormDb
}
