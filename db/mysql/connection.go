package mysql
/*
	Mysql 연결 및 초기화 
*/
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JokerTrickster/common/aws"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLService struct {
	sqlDB    *sql.DB
	gormDB   *gorm.DB
	initOnce sync.Once
	isLocal  bool
	ssmKeys  []string
	region   string
	connStr  string
}

var instance *MySQLService
var once sync.Once

const DBTimeout = 8 * time.Second

// GetMySQLService returns a singleton instance of MySQLService
func GetMySQLService(isLocal bool, region string, ssmKeys []string) *MySQLService {
	once.Do(func() {
		instance = &MySQLService{
			isLocal: isLocal,
			region:  region,
			ssmKeys: ssmKeys,
		}
	})
	return instance
}

// Initialize initializes the MySQL connection
func (m *MySQLService) Initialize(ctx context.Context) error {
	var err error
	m.initOnce.Do(func() {
		connectionString, e := m.getConnectionString(ctx)
		if e != nil {
			err = e
			return
		}
		m.connStr = connectionString

		// SQL DB connection
		m.sqlDB, e = sql.Open("mysql", connectionString)
		if e != nil {
			err = fmt.Errorf("failed to connect to MySQL: %w", e)
			return
		}

		// GORM DB connection
		m.gormDB, e = gorm.Open(mysql.New(mysql.Config{
			Conn: m.sqlDB,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if e != nil {
			err = fmt.Errorf("failed to connect to GORM MySQL: %w", e)
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

// getConnectionString generates the MySQL connection string
func (m *MySQLService) getConnectionString(ctx context.Context) (string, error) {
	if m.isLocal {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			getEnvOrFallback("MYSQL_USER", "root"),
			getEnvOrFallback("MYSQL_PASSWORD", ""),
			getEnvOrFallback("MYSQL_HOST", "localhost"),
			getEnvOrFallback("MYSQL_PORT", "3306"),
			getEnvOrFallback("MYSQL_DATABASE", "test_db"),
		), nil
	}

	// Fetch credentials from SSM
	ssmService := aws.SSMService{}
	dbInfos, err := ssmService.AwsSsmGetParams(ctx, m.ssmKeys)
	if err != nil {
		return "", fmt.Errorf("failed to fetch MySQL SSM parameters: %w", err)
	}

	// SSM Keys order: [host, port, db, user, password]
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbInfos[3], // user
		dbInfos[4], // password
		dbInfos[0], // host
		dbInfos[1], // port
		dbInfos[2], // db
	), nil
}
