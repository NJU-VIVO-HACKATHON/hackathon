package repository

import (
	"github.com/NJU-VIVO-HACKATHON/hackathon/model"
	"gorm.io/gorm"
	"testing"
)

func TestGetDataBase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test successful database connection",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDataBase()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	db, _ := GetDataBase()

	type args struct {
		email string
		sns   string

		db *gorm.DB
	}
	tests := []struct {
		name             string
		args             args
		wantRowsAffected int64
		wantErr          bool
	}{
		{
			name: "normal case",
			args: args{
				email: "test@example.com",
				sns:   "example",
				db:    db,
			},

			wantRowsAffected: 1,
			wantErr:          false,
		},
		{
			name: "empty email",
			args: args{
				email: "",
				sns:   "example",
				db:    db,
			},

			wantRowsAffected: 1,
			wantErr:          false,
		},
		{
			name: "empty sns",
			args: args{
				email: "test@example.com",
				sns:   "",
				db:    db,
			},

			wantRowsAffected: 1,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotRowsAffected, err := CreateUser(&tt.args.email, &tt.args.sns, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotRowsAffected != tt.wantRowsAffected {
				t.Errorf("CreateUser() gotRowsAffected = %v, want %v", gotRowsAffected, tt.wantRowsAffected)
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	db, _ := GetDataBase()
	type args struct {
		id int64
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test01",
			args: args{
				id: 19,
				db: db,
			},
			wantErr: true,
		},
		{
			name: "test02",
			args: args{
				id: 1,
				db: db,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetUserInfo(tt.args.id, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Skip()
		})
	}
}

func TestUpdateUserInfo(t *testing.T) {
	db, _ := GetDataBase()
	Emailstr := "ssss346346"
	Smsstr := "ssss346346"
	type args struct {
		id   int64
		user model.User
		db   *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{

			name: "testOK",
			args: args{
				id: 2,
				db: db,
				user: model.User{
					Email: &Emailstr,
					Sms:   &Smsstr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateUserInfo(tt.args.id, tt.args.user, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
