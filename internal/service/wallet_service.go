package service

import (
	"context"
	"errors"

	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/internal/repository"
	"github.com/cecepsprd/starworks-test/utils/logger"
)

type walletService struct {
	repo repository.WalletRepository
}

type WalletService interface {
	CheckBalance(context.Context, model.CheckBalanceRequest) (*model.Wallet, error)
	TopUp(ctx context.Context, req model.TopUpRequest) error
	Pay(ctx context.Context, req model.PayRequest) error
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return &walletService{
		repo: walletRepo,
	}
}

func (s *walletService) CheckBalance(ctx context.Context, req model.CheckBalanceRequest) (*model.Wallet, error) {
	res, err := s.repo.ReadBalance(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return res, nil
}

func (s *walletService) TopUp(ctx context.Context, req model.TopUpRequest) error {
	wallet, err := s.repo.ReadBalance(ctx, model.CheckBalanceRequest{
		UserID:  req.UserID,
		Address: req.Address,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	err = s.repo.UpdateBalance(ctx, model.Wallet{
		Balance: wallet.Balance + req.Nominal,
		Address: req.Address,
		UserID:  req.UserID,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func (s *walletService) Pay(ctx context.Context, req model.PayRequest) error {

	wallet, err := s.repo.ReadBalance(ctx, model.CheckBalanceRequest{
		UserID:  req.UserID,
		Address: req.Address,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if wallet.Balance < req.NominalPayment {
		return errors.New("your current balance is insufficient")
	}

	err = s.repo.UpdateBalance(ctx, model.Wallet{
		Balance: wallet.Balance - req.NominalPayment,
		Address: req.Address,
		UserID:  req.UserID,
	})

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}
