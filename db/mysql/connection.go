package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLService struct {
	sqlDB    *sql.DB
	gormDB   *gorm.DB
	initOnce sync.Once
	connStr  string
}

var instance *MySQLService
var once sync.Once

const DBTimeout = 8 * time.Second

// GetMySQLService returns a singleton instance of MySQLService
func GetMySQLService() *MySQLService {
	once.Do(func() {
		instance = &MySQLService{}
	})
	return instance
}

// Initialize initializes the MySQL connection with the provided connection string
func (m *MySQLService) Initialize(ctx context.Context, connectionString string) error {
	var err error
	m.initOnce.Do(func() {
		m.connStr = connectionString

		// SQL DB connection
		m.sqlDB, err = sql.Open("mysql", connectionString)
		if err != nil {
			err = fmt.Errorf("failed to connect to MySQL: %w", err)
			return
		}

		// GORM DB connection
		m.gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn: m.sqlDB,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			err = fmt.Errorf("failed to connect to GORM MySQL: %w", err)
			return
		}
		log.Println("Connected to MySQL!")
	})
	return err
}

// GetSQLDB returns the raw SQL DB connection
func (m *MySQLService) GetSQLDB() (*sql.DB, error) {
	if m.sqlDB == nil {
		return nil, fmt.Errorf("SQL DB is not initialized. Call Initialize first")
	}
	return m.sqlDB, nil
}

// GetGORMDB returns the GORM DB connection
func (m *MySQLService) GetGORMDB() (*gorm.DB, error) {
	if m.gormDB == nil {
		return nil, fmt.Errorf("GORM DB is not initialized. Call Initialize first")
	}
	return m.gormDB, nil
}
