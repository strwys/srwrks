package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/cecepsprd/starworks-test/internal/model"
)

type UserRepository interface {
	BeginTx(ctx context.Context) *sql.Tx
	Create(ctx context.Context, tx *sql.Tx, user model.User) (userID int64, err error)
	ReadByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
	IsUserRegistered(ctx context.Context, username, email string) (bool, error)
	ReadLoginHistory(ctx context.Context, req model.LoginHistory) (history model.LoginHistory, err error)
	WriteLoginHistory(ctx context.Context, req model.LoginHistory) error
	UpdateLoginHistory(ctx context.Context, req model.LoginHistory) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{
		db: db,
	}
}

func (m *mysqlUserRepository) BeginTx(ctx context.Context) *sql.Tx {
	tx, _ := m.db.BeginTx(ctx, nil)
	return tx
}

func (m *mysqlUserRepository) Create(ctx context.Context, tx *sql.Tx, user model.User) (userID int64, err error) {
	query := `INSERT INTO user (first_name, last_name,birth_date,street_address,city,province,phone,email,username,password,token) VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.BirthDate, user.StreetAddress, user.City, user.Province, user.Phone, user.Email, user.Username, user.Password, user.Token)
	if err != nil {
		return 0, err
	}
	userID, err = res.LastInsertId()

	return userID, err
}

func (m *mysqlUserRepository) IsUserRegistered(ctx context.Context, username, email string) (bool, error) {
	total := 0
	query := `SELECT COUNT(id) FROM user WHERE username=? OR email=?`
	err := m.db.QueryRowContext(ctx, query, username, email).Scan(&total)
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func (m *mysqlUserRepository) ReadByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, first_name, last_name, birth_date, street_address, city, province, phone, email, username, password FROM user WHERE username=? OR email=?`
	err := m.db.QueryRowContext(ctx, query, username, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.StreetAddress,
		&user.City,
		&user.Province,
		&user.Phone,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

func (m *mysqlUserRepository) ReadLoginHistory(ctx context.Context, req model.LoginHistory) (history model.LoginHistory, err error) {
	query := `SELECT browser_name, login_succeed, login_failed, user_id FROM login_history WHERE browser_name=? AND user_id=?`

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, req.BrowserName, req.UserID).Scan(
		&history.BrowserName,
		&history.LoginSucceed,
		&history.LoginFailed,
		&history.UserID,
	)

	if err != nil {
		return
	}

	return history, nil
}

func (m *mysqlUserRepository) WriteLoginHistory(ctx context.Context, req model.LoginHistory) error {
	insertQuery := `INSERT INTO login_history (browser_name, login_succeed, login_failed, user_id) VALUES (?,?,?,?)`

	stmt, err := m.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req.BrowserName, req.LoginSucceed, req.LoginFailed, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlUserRepository) UpdateLoginHistory(ctx context.Context, req model.LoginHistory) error {
	query := `UPDATE login_history SET login_succeed=?, login_failed=?, updated_at=? WHERE browser_name=? AND user_id=?`

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, req.LoginSucceed, req.LoginFailed, time.Now(), req.BrowserName, req.UserID)
	if err != nil {
		return err
	}

	return nil
}
