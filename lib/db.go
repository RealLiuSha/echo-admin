package lib

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	ORM *gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(config Config, logger Logger) Database {
	mc := mysql.Config{
		DSN:                       config.Database.DSN(),
		DefaultStringSize:         191,   // default length of string type field
		SkipInitializeWithVersion: false, // Automatic configuration based on version
		DisableDatetimePrecision:  true,  // Disable datetime precision. Databases before MySQL 5.6 do not support it.
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
	}

	db, err := gorm.Open(mysql.New(mc), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		SkipDefaultTransaction: true,
		// disable foreign keys
		// specifying foreign keys does not create real foreign key constraints in mysql
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.Database.TablePrefix + "_",
		},
		// query all fields, and in some cases "*" does not walk the index
		QueryFields: true,
	})

	if err != nil {
		logger.Zap.Fatalf("Error to open database[%s] connection: %v", mc.DSN, err)
	}

	if config.Log.Level == "debug" {
		db = db.Debug()
	}

	logger.Zap.Info("Database connection established")
	return Database{
		ORM: db,
	}
}
