package user

import (
	"app/internal/interfaces/repositories/models"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
	suite.db.Exec("DROP TABLE IF EXISTS users CASCADE;")
	suite.rd.FlushAll(suite.ctx).Err()
}

func (suite *UserRepositoryTestSuite) BeforeTest(_ string, testName string) {
	if testName == "TestSaveUpdateUser" || testName == "TestDeleteUser" {
		user := models.User{
			Fullname:    "Test Updated Fullname",
			Email:       "test_updated@gmail.com",
			Password:    "test password",
			PhoneNumber: "0937409682",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC),
			Model: gorm.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		result := suite.db.Save(&user)
		suite.NoError(result.Error)
		suite.Greater(user.ID, uint(0))

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
		err = suite.rd.HSet(suite.ctx, idKey, "fullname", user.Fullname).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "phone_number", user.PhoneNumber).Err()
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
		err = suite.rd.HSet(suite.ctx, emailKey, "fullname", user.Fullname).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "phone_number", user.PhoneNumber).Err()
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
    Password: "my_password",
		PhoneNumber: "0937409682",
		Birthday:    time.Date(1996, 07, 19, 0, 0, 0, 0, time.UTC),
		LatestLogin: time.Date(2024, 06, 16, 0, 0, 0, 0, time.UTC),
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
	suite.Equal(newUser.Birthday.Unix(), users[0].Birthday.Unix())
	suite.Equal(newUser.LatestLogin.Unix(), users[0].LatestLogin.Unix())
}

func (suite *UserRepositoryTestSuite) TestSaveUpdateUser() {
	// Ensure that we have user in database
	var user models.User
	result := suite.db.Find(&user)
	suite.NoError(result.Error)
	suite.NotZero(user)

	// Ensure that we have user in cache
	idKey := getIdKey(user.ID)
	emailKey := getEmailKey(user.Email)

	res, err := suite.rd.Exists(suite.ctx, idKey).Result()
	suite.NoError(err)
	suite.Greater(res, int64(0))
	res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
	suite.NoError(err)
	suite.Greater(res, int64(0))

	// Update and save user to the database
	user.Email = "updated@gmail.com"
	user.Fullname = "Updated fullname"
	user.Birthday = time.Now()
	user.Password = "updated_password"
	user.PhoneNumber = "12345678909"
	user.LatestLogin = time.Now()
	repo := NewRepo(suite.db, suite.rd)
	err = repo.SaveUser(suite.ctx, &user)
	suite.NoError(err)

	// Ensure that cache is invalidated
	idKey = getIdKey(user.ID)
	emailKey = getEmailKey(user.Email)
	res, err = suite.rd.Exists(suite.ctx, idKey).Result()
	suite.NoError(err)
	suite.Equal(int64(0), res)
	res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
	suite.NoError(err)
	suite.Equal(int64(0), res)

	// Ensure that updated user save to database
	var users []models.User
	result = suite.db.Find(&users)
	suite.NoError(result.Error)
	suite.Equal(1, len(users))
	suite.Equal(user.ID, users[0].ID)
	suite.Equal(user.Email, users[0].Email)
	suite.Equal(user.Fullname, users[0].Fullname)
	suite.Equal(user.Password, users[0].Password)
	suite.Equal(user.Birthday.Unix(), users[0].Birthday.Unix())
	suite.Equal(user.LatestLogin.Unix(), users[0].LatestLogin.Unix())
	suite.Equal(user.PhoneNumber, users[0].PhoneNumber)
	suite.Equal(user.CreatedAt, users[0].CreatedAt)
	suite.Equal(user.UpdatedAt, users[0].UpdatedAt)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser() {
	// ensure that user have in database
	var user models.User
	result := suite.db.Find(&user)
	suite.NoError(result.Error)
	suite.NotZero(user)

	// ensure that user exists in cache
	idKey := getIdKey(user.ID)
	emailKey := getEmailKey(user.Email)
	res, err := suite.rd.Exists(suite.ctx, idKey).Result()
	suite.NoError(err)
	suite.Greater(res, int64(0))
	res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
	suite.NoError(err)
	suite.Greater(res, int64(0))

	// delete user
	repo := NewRepo(suite.db, suite.rd)
	err = repo.DeleteUser(suite.ctx, int(user.ID))
	suite.NoError(err)

	// ensure that user not exists in cache
	res, err = suite.rd.Exists(suite.ctx, idKey).Result()
	suite.NoError(err)
	suite.Equal(int64(0), res)
	res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
	suite.NoError(err)
	suite.Equal(int64(0), res)

	// ensure that user not exists in database
	users := []models.User{}
	result = suite.db.Find(&users)
	suite.NoError(result.Error)
	suite.Empty(users)
}

func (suite *UserRepositoryTestSuite) TestGetUser() {
	suite.Run("Get User By ID when user exist in database", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// Insert dbUser to database
		dbUser := models.User{
			Fullname:    "test fullname",
			PhoneNumber: "0987654321",
			Email:       "test@gmail.com",
			Password:    "my_password",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Now(),
		}
		result := suite.db.Save(&dbUser)
		suite.NoError(result.Error)
		suite.Greater(dbUser.ID, uint(0))

		// Ensure that user not exists in cache
		idKey := getIdKey(dbUser.ID)
		emailKey := getEmailKey(dbUser.Email)
		res, err := suite.rd.Exists(suite.ctx, idKey).Result()
		suite.NoError(err)
		suite.Equal(int64(0), res)
		res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
		suite.NoError(err)
		suite.Equal(int64(0), res)

		// Get user from the database
		repo := NewRepo(suite.db, suite.rd)
		user, err := repo.GetUser(suite.ctx, int(dbUser.ID))
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(user.ID, dbUser.ID)
		suite.Equal(user.Email, dbUser.Email)
		suite.Equal(user.Fullname, dbUser.Fullname)
		suite.Equal(user.Password, dbUser.Password)
		suite.Equal(user.Birthday.Unix(), dbUser.Birthday.Unix())
		suite.Equal(user.LatestLogin.Unix(), dbUser.LatestLogin.Unix())
		suite.Equal(user.PhoneNumber, dbUser.PhoneNumber)
		suite.Equal(user.CreatedAt, dbUser.CreatedAt)
		suite.Equal(user.UpdatedAt, dbUser.UpdatedAt)

		// ensure that user exists in cache
		res, err = suite.rd.Exists(suite.ctx, idKey).Result()
		suite.NoError(err)
		suite.Greater(res, int64(0))
		res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
		suite.NoError(err)
		suite.Greater(res, int64(0))

		mapUser, err := suite.rd.HGetAll(suite.ctx, idKey).Result()
		suite.NoError(err)
		suite.Equal(strconv.Itoa(int(dbUser.ID)), mapUser["id"])
		suite.Equal(dbUser.Email, mapUser["email"])
		suite.Equal(dbUser.Fullname, mapUser["fullname"])
		suite.Equal(dbUser.Password, mapUser["password"])
		birthday, err := time.Parse(time.RFC3339Nano, mapUser["birthday"])
		suite.NoError(err)
		suite.Equal(dbUser.Birthday.Unix(), birthday.Unix())
		suite.Equal(dbUser.PhoneNumber, mapUser["phone_number"])
		latestLogin, err := time.Parse(time.RFC3339Nano, mapUser["latest_login"])
		suite.NoError(err)
		suite.Equal(dbUser.LatestLogin.Unix(), latestLogin.Unix())
	})

	suite.Run("Get User By ID when user not exists in database", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// ensure user not exists in cache
		res, err := suite.rd.Keys(suite.ctx, "*").Result()
		suite.NoError(err)
		suite.Empty(res)

		// ensure user not exists in database
		users := []models.User{}
		result := suite.db.Find(&users)
		suite.Equal(int64(0), result.RowsAffected)
		suite.Empty(users)

		// Get user and return nil
		repo := NewRepo(suite.db, suite.rd)
		user, err := repo.GetUser(suite.ctx, 1)
		suite.Zero(user)
		suite.Equal(nil, err)
	})

	suite.Run("Get User By ID when user exists in cache", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// Insert dbUser to database
		dbUser := models.User{
			Fullname:    "test fullname",
			PhoneNumber: "0987654321",
			Email:       "test@gmail.com",
			Password:    "my_password",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Now(),
		}
		result := suite.db.Save(&dbUser)
		suite.NoError(result.Error)
		suite.Greater(dbUser.ID, uint(0))

		// cache the user
		idKey := getIdKey(dbUser.ID)

		err := suite.rd.HSet(suite.ctx, idKey, "id", dbUser.ID).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "email", dbUser.Email).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "password", dbUser.Password).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "birthday", dbUser.Birthday).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "fullname", dbUser.Fullname).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "phone_number", dbUser.PhoneNumber).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "latest_login", dbUser.LatestLogin).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "created_at", dbUser.CreatedAt).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, idKey, "updated_at", dbUser.UpdatedAt).Err()
		suite.NoError(err)

		// test if we not hit database
		mockDB, mock, err := sqlmock.New()
		suite.NoError(err)
		suite.T().Cleanup(func() {
			mockDB.Close()
		})
		dialector := pg.New(pg.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		suite.NoError(err)

		// get the user by id and ensure not hit database
		repo := NewRepo(db, suite.rd)
		user, err := repo.GetUser(suite.ctx, int(dbUser.ID))
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(user.ID, dbUser.ID)
		suite.Equal(user.Email, dbUser.Email)
		suite.Equal(user.Fullname, dbUser.Fullname)
		suite.Equal(user.Password, dbUser.Password)
		suite.Equal(user.Birthday.Unix(), dbUser.Birthday.Unix())
		suite.Equal(user.LatestLogin.Unix(), dbUser.LatestLogin.Unix())
		suite.Equal(user.PhoneNumber, dbUser.PhoneNumber)
		suite.Equal(user.CreatedAt, dbUser.CreatedAt)
		suite.Equal(user.UpdatedAt, dbUser.UpdatedAt)

		// ensure that we all met expectation and we didn't expect any
		err = mock.ExpectationsWereMet()
		suite.NoError(err)
	})

	suite.Run("Get User By email when user exist in database", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// Insert dbUser to database
		dbUser := models.User{
			Fullname:    "test fullname",
			PhoneNumber: "0987654321",
			Email:       "test@gmail.com",
			Password:    "my_password",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Now(),
		}
		result := suite.db.Save(&dbUser)
		suite.NoError(result.Error)
		suite.Greater(dbUser.ID, uint(0))
		suite.Equal(int64(1), result.RowsAffected)

		// Ensure that user not exists in cache
		emailKey := getEmailKey(dbUser.Email)
		res, err := suite.rd.Exists(suite.ctx, emailKey).Result()
		suite.NoError(err)
		suite.Equal(int64(0), res)

		// Get user from the database
		repo := NewRepo(suite.db, suite.rd)
		user, err := repo.GetUserByEmail(suite.ctx, dbUser.Email)
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(user.ID, dbUser.ID)
		suite.Equal(user.Email, dbUser.Email)
		suite.Equal(user.Fullname, dbUser.Fullname)
		suite.Equal(user.Password, dbUser.Password)
		suite.Equal(user.Birthday.Unix(), dbUser.Birthday.Unix())
		suite.Equal(user.LatestLogin.Unix(), dbUser.LatestLogin.Unix())
		suite.Equal(user.PhoneNumber, dbUser.PhoneNumber)
		suite.Equal(user.CreatedAt, dbUser.CreatedAt)
		suite.Equal(user.UpdatedAt, dbUser.UpdatedAt)

		// ensure that user exists in cache
		res, err = suite.rd.Exists(suite.ctx, emailKey).Result()
		suite.NoError(err)
		suite.Greater(res, int64(0))

		mapUser, err := suite.rd.HGetAll(suite.ctx, emailKey).Result()
		suite.NoError(err)
		suite.Equal(strconv.Itoa(int(dbUser.ID)), mapUser["id"])
		suite.Equal(dbUser.Email, mapUser["email"])
		suite.Equal(dbUser.Fullname, mapUser["fullname"])
		suite.Equal(dbUser.Password, mapUser["password"])
		birthday, err := time.Parse(time.RFC3339Nano, mapUser["birthday"])
		suite.NoError(err)
		suite.Equal(dbUser.Birthday.Unix(), birthday.Unix())
		suite.Equal(dbUser.PhoneNumber, mapUser["phone_number"])
		latestLogin, err := time.Parse(time.RFC3339Nano, mapUser["latest_login"])
		suite.NoError(err)
		suite.Equal(dbUser.LatestLogin.Unix(), latestLogin.Unix())
	})

	suite.Run("Get User By email when user not exists in database", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// ensure user not exists in cache
		res, err := suite.rd.Keys(suite.ctx, "*").Result()
		suite.NoError(err)
		suite.Empty(res)

		// ensure user not exists in database
		users := []models.User{}
		result := suite.db.Find(&users)
		suite.Equal(int64(0), result.RowsAffected)
		suite.Empty(users)

		// Get user and return nil
		repo := NewRepo(suite.db, suite.rd)
		user, err := repo.GetUserByEmail(suite.ctx, "not_exist@gmail.com")
		suite.Zero(user)
		suite.Equal(gorm.ErrRecordNotFound, err)
		suite.Equal(int64(0), result.RowsAffected)
	})

	suite.Run("Get User By email when user exists in cache", func() {
		suite.T().Cleanup(func() {
			suite.db.Exec("DELETE FROM users;")
			suite.rd.FlushAll(suite.ctx)
		})
		// Insert dbUser to database
		dbUser := models.User{
			Fullname:    "test fullname",
			PhoneNumber: "0987654321",
			Email:       "test@gmail.com",
			Password:    "my_password",
			Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
			LatestLogin: time.Now(),
		}
		result := suite.db.Save(&dbUser)
		suite.NoError(result.Error)
		suite.Greater(dbUser.ID, uint(0))
		suite.Equal(int64(1), result.RowsAffected)

		// cache the user
		emailKey := getEmailKey(dbUser.Email)

		err := suite.rd.HSet(suite.ctx, emailKey, "id", dbUser.ID).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "email", dbUser.Email).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "password", dbUser.Password).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "fullname", dbUser.Fullname).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "phone_number", dbUser.PhoneNumber).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "birthday", dbUser.Birthday).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "latest_login", dbUser.LatestLogin).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "created_at", dbUser.CreatedAt).Err()
		suite.NoError(err)
		err = suite.rd.HSet(suite.ctx, emailKey, "updated_at", dbUser.UpdatedAt).Err()
		suite.NoError(err)

		// test if we not hit database
		mockDB, mock, err := sqlmock.New()
		suite.NoError(err)
		suite.T().Cleanup(func() {
			mockDB.Close()
		})
		dialector := pg.New(pg.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		db, err := gorm.Open(dialector, &gorm.Config{})
		suite.NoError(err)

		// get the user by email and ensure not hit database
		repo := NewRepo(db, suite.rd)
		user, err := repo.GetUserByEmail(suite.ctx, dbUser.Email)
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(user.ID, dbUser.ID)
		suite.Equal(user.Email, dbUser.Email)
		suite.Equal(user.Fullname, dbUser.Fullname)
		suite.Equal(user.Password, dbUser.Password)
		suite.Equal(user.Birthday.Unix(), dbUser.Birthday.Unix())
		suite.Equal(user.LatestLogin.Unix(), dbUser.LatestLogin.Unix())
		suite.Equal(user.PhoneNumber, dbUser.PhoneNumber)
		suite.Equal(user.CreatedAt, dbUser.CreatedAt)
		suite.Equal(user.UpdatedAt, dbUser.UpdatedAt)

		// ensure that we all met expectation and we didn't expect any
		err = mock.ExpectationsWereMet()
		suite.NoError(err)
	})
}
func TestUserRepository(t *testing.T) {
	suite.Run(t, &UserRepositoryTestSuite{})
}
