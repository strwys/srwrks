package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/cecepsprd/starworks-test/internal/mocks"
	"github.com/cecepsprd/starworks-test/internal/model"
)

func Test_walletService_CheckBalance(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name    string
		args    model.CheckBalanceRequest
		want    *model.Wallet
		wantErr bool
	}{
		{
			name: "positif",
			args: model.CheckBalanceRequest{
				UserID:  1,
				Address: "d49b7e51ca34f55e1e40d922cc42a134d1c731d5f03f3ecd661b7c5501f8c954",
			},
			want: &model.Wallet{
				UserID:  1,
				Address: "d49b7e51ca34f55e1e40d922cc42a134d1c731d5f03f3ecd661b7c5501f8c954",
				Balance: 1000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.WalletRepository{}

			mockRepo.On("ReadBalance", ctx, tt.args).Return(tt.want, nil)

			s := NewWalletService(&mockRepo)
			got, err := s.CheckBalance(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.CheckBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.CheckBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_walletService_TopUp(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		args    model.TopUpRequest
		wantErr bool
	}{
		{
			name: "positif",
			args: model.TopUpRequest{
				Nominal: 1000,
				UserID:  1,
				Address: "d49b7e51ca34f55e1e40d922cc42a134d1c731d5f03f3ecd661b7c5501f8c954",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.WalletRepository{}

			checkBalReq := model.CheckBalanceRequest{
				UserID:  tt.args.UserID,
				Address: tt.args.Address,
			}

			wallet := model.Wallet{
				Balance: 1000,
			}

			mockRepo.On("ReadBalance", ctx, checkBalReq).Return(&wallet, nil)

			mockRepo.On("UpdateBalance", ctx, model.Wallet{
				Balance: tt.args.Nominal + wallet.Balance,
				Address: tt.args.Address,
				UserID:  tt.args.UserID,
			}).Return(nil)

			s := NewWalletService(&mockRepo)
			if err := s.TopUp(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("walletService.TopUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_walletService_Pay(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		args    model.PayRequest
		wantErr bool
	}{
		{
			name: "positif",
			args: model.PayRequest{
				NominalPayment: 1000,
				UserID:         1,
				Address:        "d49b7e51ca34f55e1e40d922cc42a134d1c731d5f03f3ecd661b7c5501f8c954",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.WalletRepository{}

			checkBalReq := model.CheckBalanceRequest{
				UserID:  tt.args.UserID,
				Address: tt.args.Address,
			}

			wallet := model.Wallet{
				Balance: 5000,
			}

			mockRepo.On("ReadBalance", ctx, checkBalReq).Return(&wallet, nil)

			mockRepo.On("UpdateBalance", ctx, model.Wallet{
				Balance: wallet.Balance - tt.args.NominalPayment,
				Address: tt.args.Address,
				UserID:  tt.args.UserID,
			}).Return(nil)

			s := NewWalletService(&mockRepo)
			if err := s.Pay(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("walletService.Pay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
