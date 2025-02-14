package postgres_test

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	client "timble/internal/connection/postgres"
)

type testStruct struct {
	Name string
}

func TestPostgres_NewClient(t *testing.T) {
	_, _, gormDb, _ := openMockDB(t)

	tests := []struct {
		name                 string
		expectedError        error
		mockeGormOpenFunc    func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)
		mockPostgresOpenFunc func(dsn string) gorm.Dialector
		mockGetGormDBFunc    func(db *gorm.DB) (*sql.DB, error)
	}{
		{
			name: "success",
			mockeGormOpenFunc: func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
				return gormDb, nil
			},
			mockPostgresOpenFunc: func(dsn string) gorm.Dialector {
				return postgres.Open(dsn)
			},
			mockGetGormDBFunc: func(db *gorm.DB) (sqlDb *sql.DB, err error) {
				return db.DB()
			},
		},
		{
			name:          "error gorm open",
			expectedError: errors.New("error gorm open"),
			mockeGormOpenFunc: func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
				return nil, errors.New("error gorm open")
			},
			mockPostgresOpenFunc: func(dsn string) gorm.Dialector {
				return postgres.Open(dsn)
			},
			mockGetGormDBFunc: func(db *gorm.DB) (sqlDb *sql.DB, err error) {
				return db.DB()
			},
		},
		{
			name:          "error gorm DB",
			expectedError: errors.New("error gorm DB"),
			mockeGormOpenFunc: func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
				return gormDb, nil
			},
			mockPostgresOpenFunc: func(dsn string) gorm.Dialector {
				return postgres.Open(dsn)
			},
			mockGetGormDBFunc: func(db *gorm.DB) (sqlDb *sql.DB, err error) {
				return nil, errors.New("error get db")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &client.PostgresClient{
				Name:             "postgres",
				GormOpenFunc:     tt.mockeGormOpenFunc,
				PostgresOpenFunc: tt.mockPostgresOpenFunc,
				GormGetDBFunc:    tt.mockGetGormDBFunc,
			}

			client, err := client.NewClient(
				db,
				"host",
				1234,
				"dbname",
				"username",
				"password",
				10,
				15,
			)

			if tt.expectedError != nil {
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, client.Client)
			}
		})
	}
}

func TestPostgres_GetSqlDb(t *testing.T) {
	_, _, gormDb, _ := openMockDB(t)

	db, _ := client.GetSQLDB(gormDb)

	assert.NotNil(t, db)
}

func TestPostgres_OpenGorm(t *testing.T) {
	_, _, dialector := initMockDB(t)

	db, err := client.OpenGorm(dialector)

	assert.NotNil(t, db)
	assert.Nil(t, err)
}

func TestPostgres_OpenPosgres(t *testing.T) {
	expectedDSN := "dsn"
	expectedDialector := postgres.Open(expectedDSN)

	dialector := client.OpenPostgres(expectedDSN)

	assert.Equal(t, expectedDialector, dialector, "Dialector mismatch")
}

func TestPostgres_GetFirst(t *testing.T) {
	name := "foo"

	tests := []struct {
		name           string
		rows           *sqlmock.Rows
		expectedResult string
		expectedError  error
	}{
		{
			name:           "successfully get first row",
			rows:           sqlmock.NewRows([]string{"name"}).AddRow(name),
			expectedResult: name,
		},
		{
			name: "successfully ran the query, but the row is not found",
			rows: sqlmock.NewRows([]string{"name"}),
		},
		{
			name:          "unexpected error from db",
			rows:          sqlmock.NewRows([]string{"namez"}).AddRow(name),
			expectedError: errors.New("model accessible fields required"),
		},
	}

	db, mock, gormDb, _ := openMockDB(t)
	defer db.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			query := `SELECT * FROM "test_structs" WHERE name = $1 ORDER BY "test_structs"."name" LIMIT $2`
			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(name, 1).
				WillReturnRows(tc.rows).
				RowsWillBeClosed()

			client := client.PostgresClient{
				Client: gormDb,
			}

			record := &testStruct{}
			err := client.GetFirst(record, "name = ?", name)
			if tc.expectedError != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				if tc.expectedResult != "" {
					assert.Equal(t, tc.expectedResult, record.Name)
				}
			}
		})
	}

}

func initMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, gorm.Dialector) {
	db := &sql.DB{}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	return db, mock, dialector
}

func openMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB, time.Time) {
	db, mock, dialector := initMockDB(t)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	now := time.Now()
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger:  newLogger,
		NowFunc: func() time.Time { return now },
	})

	assert.NoError(t, err)

	return db, mock, gormDB, now
}
