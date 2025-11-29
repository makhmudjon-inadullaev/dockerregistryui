package persistence

import (
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	// Using SQLite by default.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// RegistriesVolumePathEnvironmentVariableName Name of the environment variable for the registries volume path.
const RegistriesVolumePathEnvironmentVariableName string = "REGISTRIES_VOLUME_PATH"

// DefaultDBPathPrefix Default path prefix for the database file.
const DefaultDBPathPrefix string = "/data/"

// DBConfig Configuration settings for persistence access.
type DBConfig struct {
	initialized  bool
	DBPathPrefix string
	DBName       string
	DBType       string
}

// DBHandle Handle for accessing the persistence context.
type DBHandle struct {
	db     *gorm.DB
	Config DBConfig
}

// NewDBConfig Creates a new DBConfig with default values.
// The database path can be configured using the REGISTRIES_VOLUME_PATH environment variable.
func NewDBConfig() DBConfig {
	dbPathPrefix := DefaultDBPathPrefix
	if volumePath := os.Getenv(RegistriesVolumePathEnvironmentVariableName); volumePath != "" {
		if strings.HasSuffix(volumePath, "/") {
			dbPathPrefix = volumePath
		} else {
			dbPathPrefix = volumePath + "/"
		}
	}
	return DBConfig{initialized: true, DBPathPrefix: dbPathPrefix, DBName: "registryui.db", DBType: "sqlite3"}
}

// StartPersistenceContext Starts a persistence context and returns the handle to that context.
func StartPersistenceContext(config DBConfig) *DBHandle {
	config.initIfNecessary()
	db, err := gorm.Open(config.DBType, config.DBPathPrefix+config.DBName)
	if err != nil {
		log.Fatalln("Failed to connect database: " + config.DBType + " @ " + config.DBPathPrefix + config.DBName)
		return nil
	}
	db.AutoMigrate(&ImageCategory{})
	db.AutoMigrate(&ImageDescription{})
	db.AutoMigrate(&HelloMessage{})
	log.Println("Database connected: " + config.DBType + " @ " + config.DBPathPrefix + config.DBName)
	return &DBHandle{db: db, Config: config}
}

// StopPersistenceContext Stops the context managed using the handle.
func (handle *DBHandle) StopPersistenceContext() {
	handle.db.Close()
}

func (config *DBConfig) initIfNecessary() {
	if !config.initialized {
		*config = NewDBConfig()
	}
}
