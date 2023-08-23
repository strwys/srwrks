package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cecepsprd/starworks-test/internal/model"
)

func Test_mysqlUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx         = context.Background()
		repo        = NewUserRepository(db)
		insertQuery = "INSERT INTO user (first_name, last_name,birth_date,street_address,city,province,phone,email,username,password,token) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	)

	tests := []struct {
		name       string
		repo       UserRepository
		wantUserID int64
		wantErr    bool
	}{
		{
			name:       "Succees scenario",
			repo:       repo,
			wantUserID: 1,
			wantErr:    false,
		},
		{
			name:       "Failed scenario",
			repo:       repo,
			wantUserID: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertQuery).ExpectExec().WillReturnError(fmt.Errorf("err"))
			} else {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertQuery).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			}

			gotUserID, err := tt.repo.Create(ctx, tt.repo.BeginTx(ctx), model.User{})
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlUserRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("mysqlUserRepository.Create() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func Test_mysqlUserRepository_IsUserRegistered(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		repo  = NewUserRepository(db)
		query = "SELECT COUNT\\(id\\) FROM user WHERE username=\\? OR email=\\?"
	)

	tests := []struct {
		name    string
		fields  UserRepository
		args    []string
		want    bool
		wantErr bool
	}{
		{
			name:    "User registered",
			fields:  repo,
			args:    []string{"starworks", "starworks@gmail.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "User registered",
			fields:  repo,
			args:    []string{"starworks", "starworks@gmail.com"},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.want {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(query).WithArgs(tt.args[0], tt.args[1]).WillReturnRows(rows)
			} else {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(0)
				mock.ExpectQuery(query).WithArgs(tt.args[0], tt.args[1]).WillReturnRows(rows)
			}

			got, err := tt.fields.IsUserRegistered(ctx, tt.args[0], tt.args[1])
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlUserRepository.IsUserRegistered() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mysqlUserRepository.IsUserRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysqlUserRepository_ReadByUsernameOrEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		repo  = NewUserRepository(db)
		query = "SELECT id, first_name, last_name, birth_date, street_address, city, province, phone, email, username, password FROM user WHERE username=\\? OR email=\\?"
	)

	user := model.User{}

	tests := []struct {
		name    string
		fields  UserRepository
		args    []string
		want    *model.User
		wantErr bool
	}{
		{
			name:    "success",
			fields:  repo,
			args:    []string{"starworks", "starworks@gmail.com"},
			want:    &user,
			wantErr: false,
		},
		{
			name:    "user not found",
			fields:  repo,
			args:    []string{"starworks", "starworks@gmail.com"},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "failed get data",
			fields:  repo,
			args:    []string{"starworks", "starworks@gmail.com"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectQuery(query).WithArgs(tt.args[0], tt.args[1]).WillReturnError(fmt.Errorf("some error"))
			} else if !tt.wantErr && tt.want == nil {
				mock.ExpectQuery(query).WithArgs(tt.args[0], tt.args[1]).WillReturnError(sql.ErrNoRows)
			} else {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birth_date", "street_address", "city", "province", "phone", "email", "username", "password"})
				rows.AddRow(user.ID, user.FirstName, user.LastName, user.BirthDate, user.StreetAddress, user.City, user.Province, user.Phone, user.Email, user.Username, user.Password)
				mock.ExpectQuery(query).WithArgs(tt.args[0], tt.args[1]).WillReturnRows(rows)
			}

			got, err := tt.fields.ReadByUsernameOrEmail(ctx, tt.args[0], tt.args[1])
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlUserRepository.ReadByUsernameOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mysqlUserRepository.ReadByUsernameOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
