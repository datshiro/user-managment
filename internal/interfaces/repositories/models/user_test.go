package models

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserModelTestSuite struct {
	suite.Suite
	ctx                context.Context
	pgContainer        *postgres.PostgresContainer
	pgConnectionString string
	db                 *gorm.DB
}

func (suite *UserModelTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.ctx = ctx
	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithUsername("myadmin"),
		postgres.WithPassword("my_password"),
		postgres.WithDatabase("userdb"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	suite.NoError(err)
	pgConnectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	suite.NoError(err)
	db, err := gorm.Open(pg.Open(pgConnectionString), &gorm.Config{})
	suite.NoError(err)

	suite.db = db
	suite.pgConnectionString = pgConnectionString
	suite.pgContainer = pgContainer

	sqlDB, err := db.DB()
	suite.NoError(err)
	err = sqlDB.Ping()
	suite.NoError(err)
}

func (suite *UserModelTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *UserModelTestSuite) SetupTest() {
	suite.db.AutoMigrate(&User{})
}

func (suite *UserModelTestSuite) TearDownTest() {
	suite.db.Exec("DROP TABLE IF EXISTS users CASCADE;")
}

func (suite *UserModelTestSuite) TestSaveUser() {
	suite.Run("Save with accepted password", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
		})
		user := User{
			Fullname:    "Test Updated Fullname",
			Email:       "test_updated@gmail.com",
      Password:    "normal password",
			PhoneNumber: "0937409682",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC),
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		result := suite.db.Save(&user)
		suite.NoError(result.Error )
	})
	suite.Run("Save with short password", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
		})
		user := User{
			Fullname:    "Test Updated Fullname",
			Email:       "test_updated@gmail.com",
			Password:    "short", // Short password
			PhoneNumber: "0937409682",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC),
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		err := suite.db.Save(&user).Error
		suite.Equal(ErrPasswordTooShort, err)
	})

	suite.Run("Save with long password", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
		})
		user := User{
			Fullname:    "Test Updated Fullname",
			Email:       "test_updated@gmail.com",
      Password:    strings.Repeat("ASDKLJWKL", 100)  , // Long password
			PhoneNumber: "0937409682",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC),
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		err := suite.db.Save(&user).Error
		suite.Equal(ErrPasswordTooLong, err)
	})
}

func TestUserModel(t *testing.T)  {
  suite.Run(t, &UserModelTestSuite{})
}
