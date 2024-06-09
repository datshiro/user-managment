package repositories

import (
	"app/internal/models"
	"context"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewRepo(t *testing.T) {
	type args struct {
		dbc *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want UserRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepo(tt.args.dbc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_CreateUser(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx  context.Context
		data *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.CreateUser(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_GetUserByID(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_UpdateLastLogin(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.UpdateLastLogin(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.UpdateLastLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userRepo.UpdateLastLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_GetUserByEmail(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.GetUserByEmail(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_GetUserByPhone(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx      context.Context
		phone    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.GetUserByPhone(tt.args.ctx, tt.args.phone, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.GetUserByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.GetUserByPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_GetUserByUsername(t *testing.T) {
	type fields struct {
		dbc *gorm.DB
	}
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepo{
				dbc: tt.fields.dbc,
			}
			got, err := repo.GetUserByUsername(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
