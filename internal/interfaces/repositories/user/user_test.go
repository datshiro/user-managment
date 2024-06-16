package user

import (
	"app/internal/interfaces/repositories/models"
	"context"
	"testing"
	"time"

	rd "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	ctx                context.Context
	db                 *gorm.DB
	pgContainer        *postgres.PostgresContainer
	pgConnectionString string
	rd                 *rd.Client
	rdContainer        *redis.RedisContainer
	rdConnectionString string
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := postgres.RunContainer(suite.ctx, testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("userdb"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("my_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	suite.NoError(err)

	pgConnectionString, err := pgContainer.ConnectionString(suite.ctx, "sslmode=disable")
	suite.NoError(err)

	db, err := gorm.Open(pg.Open(pgConnectionString), &gorm.Config{})
	suite.NoError(err)

	suite.db = db
	suite.pgContainer = pgContainer
	suite.pgConnectionString = pgConnectionString

	sqlDB, err := suite.db.DB()
	suite.NoError(err)

	err = sqlDB.Ping()
	suite.NoError(err)

	rdContainer, err := redis.RunContainer(suite.ctx, testcontainers.WithImage("redis:6"))
	suite.NoError(err)

	rdConnectionString, err := rdContainer.ConnectionString(suite.ctx)
	suite.NoError(err)

	rdOptions, err := rd.ParseURL(rdConnectionString)
	suite.NoError(err)

	rd := rd.NewClient(rdOptions)

	suite.rd = rd
	suite.rdConnectionString = rdConnectionString
	suite.rdContainer = rdContainer

	err = suite.rd.Ping(suite.ctx).Err()
	suite.NoError(err)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)
	err = suite.rdContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	err := suite.db.AutoMigrate(&models.User{})
	suite.NoError(err)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	err := suite.db.Exec("DELETE FROM users;").Error
	suite.NoError(err)

	err = suite.rd.FlushAll(suite.ctx).Err()
	suite.NoError(err)
}

func (suite *UserRepositoryTestSuite) BeforeTest(_ string, testName string) {
	if testName == "TestSaveUpdateUser" || testName == "TestDeleteUser" {
		user := models.User{
			Fullname:    "Test Updated Fullname",
			Email:       "test_updated@gmail.com",
			PhoneNumber: "0937409682",
			Birthday:    time.Now().AddDate(1996, int(time.July), 19),
			LatestLogin: time.Now().AddDate(2024, int(time.June), 16),
		}

		result := suite.db.Save(&user)
		suite.NoError(result.Error)

		idKey := getIdKey(user.ID)
		emailKey := getEmailKey(user.Email)

		err := suite.rd.HSet(suite.ctx, idKey, "id", user.ID).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "email", user.Email).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "password", user.Password).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "birthday", user.Birthday).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "latest_login", user.LatestLogin).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "created_at", user.CreatedAt).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "updated_at", user.UpdatedAt).Err()
		suite.NoError(err)

		err = suite.rd.HSet(suite.ctx, emailKey, "id", user.ID).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "email", user.Email).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "password", user.Password).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "birthday", user.Birthday).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "latest_login", user.LatestLogin).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "created_at", user.CreatedAt).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "updated_at", user.UpdatedAt).Err()
		suite.NoError(err)
	}
}

func (suite *UserRepositoryTestSuite) TestSaveNewUser() {
	// ensure that cache is empty
	keys, err := suite.rd.Keys(suite.ctx, "*").Result()
	suite.NoError(err)
	suite.Empty(keys)

	// ensure that database is empty
	var users []models.User
	result := suite.db.Find(&users)
	suite.NoError(result.Error)
	suite.Empty(users)

	// Create new user and save to database
	newUser := models.User{
		Fullname:    "Test Fullname",
		Email:       "test@gmail.com",
		PhoneNumber: "0937409682",
		Birthday:    time.Now().AddDate(1996, int(time.July), 19),
		LatestLogin: time.Now().AddDate(2024, int(time.June), 16),
	}
	result = suite.db.Save(&newUser)
	suite.NoError(result.Error)

	// ensure that cache is still empty
	keys, err = suite.rd.Keys(suite.ctx, "*").Result()
	suite.NoError(err)
	suite.Empty(keys)

	// ensure that we have saved the new user in database
	result = suite.db.Find(&users)
	suite.NoError(result.Error)
	suite.Equal(1, len(users))
  suite.Equal(newUser.ID, users[0].ID)
  suite.Equal(newUser.Email, users[0].Email)
  suite.Equal(newUser.Password, users[0].Password)
  suite.Equal(newUser.PhoneNumber, users[0].PhoneNumber)
  suite.Equal(newUser.Birthday, users[0].Birthday)
  suite.Equal(newUser.LatestLogin, users[0].LatestLogin)
}

func (repo *UserRepositoryTestSuite) TestSaveUpdatedUser() {

}

func (repo *UserRepositoryTestSuite) TestDeleteUser() {

}

func (repo *UserRepositoryTestSuite) TestGetUser() {

}
func TestUserRepository(t *testing.T)  {
  suite.Run(t, &UserRepositoryTestSuite{})
}
