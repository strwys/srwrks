package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cecepsprd/starworks-test/internal/mocks"
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/utils"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Init(-1, "2006-01-02T15:04:05.999999999Z07:00")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func dbConn() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}

	return db, mock
}

func beginTx(db *sql.DB, mock sqlmock.Sqlmock) *sql.Tx {
	mock.ExpectBegin()
	tx, _ := db.BeginTx(context.Background(), nil)
	return tx
}

var user = model.User{
	Username: "cecep",
	Email:    "cecep@mail.com",
	Password: "123",
}

func Test_userService_Login(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		args      model.LoginRequest
		wantUser  *model.User
		wantToken string
		wantErr   bool
	}{
		{
			name: "positif",
			args: model.LoginRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: user.Password,
			},
			wantUser: &user,
			wantErr:  false,
		},
		{
			name: "negatif: user not found",
			args: model.LoginRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: user.Password,
			},
			wantUser: &model.User{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := mocks.UserRepository{}
			mockWalletRepo := mocks.WalletRepository{}

			if tt.wantErr {
				mockUserRepo.On("ReadByUsernameOrEmail", ctx, tt.args.Username, tt.args.Email).Return(nil, fmt.Errorf("some err"))
			} else {
				user.Password, _ = utils.HashPassword(user.Password)
				mockUserRepo.On("ReadByUsernameOrEmail", ctx, tt.args.Username, tt.args.Email).Return(&user, nil)
			}

			s := NewUserService(&mockUserRepo, &mockWalletRepo, "jwtSecret", 5*time.Second)

			gotUser, _, err := s.Login(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotUser == nil {
				return
			}

			tt.wantUser.Password, _ = utils.HashPassword(user.Password)

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authService.Login() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_userService_Create(t *testing.T) {
	db, mockDB := dbConn()

	ctx := context.Background()

	tests := []struct {
		name    string
		args    model.User
		wantErr bool
	}{
		{
			name: "positif",
			args: model.User{
				Username: "elon",
				Email:    "elon@spacex.com",
			},
			wantErr: false,
		},
		{
			name: "negatif: user already exist",
			args: model.User{
				Username: "elon",
				Email:    "elon@spacex.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := beginTx(db, mockDB)

			mockUserRepo := mocks.UserRepository{}
			mockWalletRepo := mocks.WalletRepository{}

			if tt.wantErr {
				mockUserRepo.On("IsUserRegistered", ctx, tt.args.Username, tt.args.Email).Return(true, nil)
			} else {
				mockUserRepo.On("IsUserRegistered", ctx, tt.args.Username, tt.args.Email).Return(false, nil)
				mockUserRepo.On("BeginTx", ctx).Return(tx)
				mockUserRepo.On("Create", ctx, tx, mock.Anything).Return(int64(1), nil)
				mockWalletRepo.On("AddWallet", ctx, tx, mock.Anything).Return(nil)
			}
			s := NewUserService(&mockUserRepo, &mockWalletRepo, "secret", 5*time.Second)

			if err := s.Create(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("userService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
