package database

import (
	"testing"

	"gorm.io/gorm"
)

func TestOpts_GetDSN(t *testing.T) {
	type fields struct {
		DB        string
		User      string
		Password  string
		Host      string
		Port      string
		SSLMode   string
		TimeZone  string
		Dialector gorm.Dialector
	}
	type args struct {
		db_type string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"test create successful",
      fields{
        DB: "mydb",
        User: "admin",
        Password: "my_password",
        Host: "localhost",
        Port: "5432",
        SSLMode: "disable",
        TimeZone: "Asia/Shanghai",
        Dialector: nil,
      },
      "host=localhost user=admin password=my_password dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Opts{
				DB:        tt.fields.DB,
				User:      tt.fields.User,
				Password:  tt.fields.Password,
				Host:      tt.fields.Host,
				Port:      tt.fields.Port,
				SSLMode:   tt.fields.SSLMode,
				TimeZone:  tt.fields.TimeZone,
				Dialector: tt.fields.Dialector,
			}
			if got := opt.GetDSN(); got != tt.want {
				t.Errorf("Opts.GetDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
