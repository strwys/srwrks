package repository

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cecepsprd/starworks-test/internal/model"
)

func Test_mysqlWalletRepository_AddWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		repo  = NewWalletRepository(db)
		query = "INSERT INTO wallet \\(address, balance, user_id\\) VALUES \\(\\?, \\?, \\?\\)"
	)

	tests := []struct {
		name    string
		fields  WalletRepository
		args    model.Wallet
		wantErr bool
	}{
		{
			name:   "success",
			fields: repo,
			args: model.Wallet{
				Address: "abcde",
				Balance: 1000,
				UserID:  5,
			},
			wantErr: false,
		},
		{
			name:   "failed",
			fields: repo,
			args: model.Wallet{
				Address: "abcde",
				Balance: 1000,
				UserID:  5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(fmt.Errorf("some error"))
			} else {
				mock.ExpectBegin()
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			}

			tx, _ := db.BeginTx(ctx, nil)

			if err := tt.fields.AddWallet(ctx, tx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("mysqlUserRepository.AddWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysqlWalletRepository_ReadBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		repo  = NewWalletRepository(db)
		query = "SELECT user_id, address, balance FROM wallet WHERE address = \\? and user_id = \\?"
	)

	tests := []struct {
		name    string
		fields  WalletRepository
		args    model.CheckBalanceRequest
		want    *model.Wallet
		wantErr bool
	}{
		{
			name:   "success",
			fields: repo,
			args: model.CheckBalanceRequest{
				Address: "abcde",
				UserID:  1,
			},
			want: &model.Wallet{
				UserID:  1,
				Address: "abcde",
				Balance: 1000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectQuery(query).WithArgs(tt.args.Address, tt.args.UserID).WillReturnError(fmt.Errorf("some error"))
			} else {
				rows := sqlmock.NewRows([]string{"user_id", "address", "balance"}).AddRow(1, "abcde", 1000)
				mock.ExpectQuery(query).WithArgs(tt.args.Address, tt.args.UserID).WillReturnRows(rows)
			}
			got, err := tt.fields.ReadBalance(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlWalletRepository.ReadBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mysqlWalletRepository.ReadBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysqlWalletRepository_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("error creating mock database\n")
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		repo  = NewWalletRepository(db)
		query = "UPDATE wallet SET balance=\\? WHERE address = \\? AND user_id = \\?"
	)

	tests := []struct {
		name    string
		fields  WalletRepository
		args    model.Wallet
		wantErr bool
	}{
		{
			name:   "positif",
			fields: repo,
			args: model.Wallet{
				Address: "addressx",
				UserID:  1,
			},
			wantErr: false,
		},
		{
			name:   "negatif",
			fields: repo,
			args: model.Wallet{
				Address: "addressx",
				UserID:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(fmt.Errorf("somer error"))
			} else {
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if err := tt.fields.UpdateBalance(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("mysqlWalletRepository.UpdateBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
