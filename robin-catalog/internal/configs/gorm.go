package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Dokito555/robin/robin-catalog/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetInt("DB_PORT")
	database := viper.GetString("DB_NAME")
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		username,
		password,
		database,
		strconv.Itoa(port),
	)

	log.Infof("Connecting to DB with DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		os.Exit(1)
	}

	log.Info("Starting database migration...")

	tables := []interface{}{
		&entity.Artist{},
		&entity.Album{},
		&entity.Song{},
	}

	for _, table := range tables {
		log.Info("starting to migrate: %T", table)
		if err := db.AutoMigrate(table); err != nil {
			log.Fatalf("failed to migrate %T: %v", table, err)
			os.Exit(1)
		}
		log.Infof("Successfully migrated: %T", table)
	}

	log.Info("Database migration completed!")

	log.Infof("artists exists? %v", db.Migrator().HasTable(&entity.Artist{}))
	log.Infof("albums exists?  %v", db.Migrator().HasTable(&entity.Album{}))
	log.Infof("songs exists?   %v", db.Migrator().HasTable(&entity.Song{}))
	
	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}