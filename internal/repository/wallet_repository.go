package repository

import (
	"context"
	"database/sql"

	"github.com/cecepsprd/starworks-test/internal/model"
)

type WalletRepository interface {
	ReadBalance(ctx context.Context, req model.CheckBalanceRequest) (*model.Wallet, error)
	UpdateBalance(context.Context, model.Wallet) error
	AddWallet(ctx context.Context, tx *sql.Tx, wallet model.Wallet) error
}

type mysqlWalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &mysqlWalletRepository{
		db: db,
	}
}

func (m *mysqlWalletRepository) AddWallet(ctx context.Context, tx *sql.Tx, wallet model.Wallet) error {
	query := `INSERT INTO wallet (address, balance, user_id) VALUES (?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, wallet.Address, wallet.Balance, wallet.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlWalletRepository) ReadBalance(ctx context.Context, req model.CheckBalanceRequest) (*model.Wallet, error) {
	query := `SELECT user_id, address, balance FROM wallet WHERE address = ? and user_id = ?`

	var wallet model.Wallet

	err := m.db.QueryRowContext(ctx, query, req.Address, req.UserID).Scan(
		&wallet.UserID,
		&wallet.Address,
		&wallet.Balance,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &wallet, nil
}

func (m *mysqlWalletRepository) UpdateBalance(ctx context.Context, wallet model.Wallet) error {
	query := `UPDATE wallet SET balance=? WHERE address = ? AND user_id = ?`

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, wallet.Balance, wallet.Address, wallet.UserID)
	if err != nil {
		return err
	}

	return nil
}
