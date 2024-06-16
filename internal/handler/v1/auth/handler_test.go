package auth

import (
	"app/internal/interfaces/repositories/models"
	user_repo "app/internal/interfaces/repositories/user"
	"app/internal/interfaces/usecases"
	"app/internal/interfaces/usecases/user"
	mocks "app/internal/mocks/usecases/user"
	"app/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	rd "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
)

func setupDatabase(t *testing.T, ctx context.Context) (*gorm.DB, func()) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		// postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	assert.NoError(t, err, "failed to startup database")

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	db, err := gorm.Open(pg.Open(connStr))
	assert.NoError(t, err, "failed to connect to database")
	db.AutoMigrate(&models.User{})

	close_db := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		postgresContainer.Terminate(ctx)
	}

	return db, close_db
}

func setupRedis(t *testing.T, ctx context.Context) (*rd.Client, func()) {
	rdContainer, err := redis.RunContainer(ctx, testcontainers.WithImage("redis:6"))
	assert.NoError(t, err)

	rdConnectionString, err := rdContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	rdConnOpts, err := rd.ParseURL(rdConnectionString)
	assert.NoError(t, err)

	rd := rd.NewClient(rdConnOpts)

	err = rd.Ping(ctx).Err()
	assert.NoError(t, err)
	return rd, func() { rdContainer.Terminate(ctx) }
}

func setupRouter(uc user.UserUsecase) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/register", NewHandler(uc).HandleRegister)
	return r
}

func TestSuccessfulRegisterRoute(t *testing.T) {
	tests := []struct {
		name     string
		args     *models.User
		wantErr  error
		wantResp utils.ResponseObject
		wantCode int
	}{
		{
			"success: register with email",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "0123456789",
				Email:       "datshiro@gmail.com",
				Password:    "sTr0ngP@ssword",
			},
			nil,
			utils.ResponseObject{
				Message: "",
				Data: &models.User{
					Fullname:    "Nguyen Quoc Dat",
					PhoneNumber: "0123456789",
					Email:       "datshiro@gmail.com",
					Password:    "sTr0ngP@ssword",
					Birthday:    time.Time{},
					LatestLogin: time.Time{},
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						DeletedAt: gorm.DeletedAt{},
					},
				},
			},
			http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			db, close_db := setupDatabase(t, ctx)
			rd, close_rd := setupRedis(t, ctx)
			repo := user_repo.NewRepo(db, rd)
			uc := usecases.NewPostgresUsecase(repo)
			defer func() {
				close_db()
				close_rd()
			}()

			router := setupRouter(uc.UserUC)
			w := httptest.NewRecorder()
			b, err := json.Marshal(tt.args)
			assert.NoError(t, err, "Marshal request body failed")

			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
			assert.NoError(t, err, "Create test request failed")
			router.ServeHTTP(w, req)

			if got := w.Code; !reflect.DeepEqual(got, tt.wantCode) {
				t.Errorf("HandleRegister() = %d, want %d", got, tt.wantCode)
			}

			wantResp, _ := json.Marshal(tt.wantResp)

			if got := w.Body.String(); !reflect.DeepEqual(got, wantResp) {
				t.Errorf("HandleRegister() = %s, want %s", got, wantResp)
			}
		})
	}
}

func TestFailureRegisterRoute(t *testing.T) {
	tests := []struct {
		name     string
		args     *models.User
		expected error
	}{
		{
			"failed: register with empty email",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "datshiro@gmail.com",
				Password:    "sTr0ngP@ssword",
			},
			errors.New("email must not be empty"),
		},
		{
			"failed: register with empty password",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "",
				Password:    "sTr0ngP@ssword",
			},
			nil,
		},
		{
			"failed: existed user",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Password:    "",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mocks.UserUsecase)
			mock.
				On("RegisterUser", context.Background(), tt.args).
				Return(nil, tt.expected)

			router := setupRouter(mock)
			w := httptest.NewRecorder()
			b, _ := json.Marshal(tt.args)

			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
			router.ServeHTTP(w, req)

			assert.NotEqual(t, 201, w.Code, tt.name)
			exp := `{"success": false, "message": "Invalid request; ", "data": null}`

			if !reflect.DeepEqual(w.Body.String(), exp) {
				t.Errorf("RegisterUser() = %s , want=%s", w.Body.String(), exp)
			}
		})
	}
}
