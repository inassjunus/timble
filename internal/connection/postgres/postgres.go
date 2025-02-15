package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"timble/internal/utils"
)

type PostgresInterface interface {
	GetFirst(record interface{}, condition string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
}

var (
	ignoredErrors = map[string]bool{
		gorm.ErrRecordNotFound.Error(): true, // this error is ignorable since it is expected that some search terms don't have any recommendation
	}
)

type PostgresClient struct {
	Name   string
	Client *gorm.DB

	// library functions for mocking library behavior
	GormOpenFunc     func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)
	PostgresOpenFunc func(dsn string) gorm.Dialector
	GormGetDBFunc    func(db *gorm.DB) (*sql.DB, error)
}

func NewClient(client *PostgresClient, host string, port int, database string, username string, password string, maxIdleConns int, maxOpenConns int) (*PostgresClient, error) {
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v",
		host, username, password, database, port,
	)

	gormDb, err := client.GormOpenFunc(client.PostgresOpenFunc(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := client.GormGetDBFunc(gormDb)
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	client.Client = gormDb

	return client, nil
}

func GetSQLDB(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}

func OpenGorm(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(dialector, opts...)
}

func OpenPostgres(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func (c *PostgresClient) GetFirst(record interface{}, condition string, args ...interface{}) error {
	metricInfo := utils.NewClientMetric(c.Name, "get-first")
	result := c.Client.Where(condition, args...).First(record)
	err := c.wrapError(result.Error)
	metricInfo.TrackClientWithError(err)
	return err
}

func (c *PostgresClient) Exec(query string, args ...interface{}) error {
	metricInfo := utils.NewClientMetric(c.Name, "exec")
	result := c.Client.Exec(query, args...)
	err := c.wrapError(result.Error)
	metricInfo.TrackClientWithError(err)
	return err
}

func (c *PostgresClient) wrapError(err error) error {
	if err != nil && !ignoredErrors[err.Error()] {
		return err
	}

	return nil
}
