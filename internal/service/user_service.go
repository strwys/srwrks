package service

import (
	"context"
	"errors"
	"time"

	cs "github.com/cecepsprd/starworks-test/constans"
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/internal/repository"
	"github.com/cecepsprd/starworks-test/utils"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(context.Context, model.LoginRequest) (user *model.User, token string, err error)
	Create(ctx context.Context, request model.User) error
	WriteLoginHistory(ctx context.Context, userID int64, isLoginFailed bool) error
}

type userService struct {
	repo           repository.UserRepository
	walletRepo     repository.WalletRepository
	JWTSecret      string
	contextTimeout time.Duration
}

func NewUserService(urepo repository.UserRepository, walletRepo repository.WalletRepository, JWTSecret string, timeout time.Duration) UserService {
	return &userService{
		repo:           urepo,
		walletRepo:     walletRepo,
		JWTSecret:      JWTSecret,
		contextTimeout: timeout,
	}
}

func (s *userService) Login(ctx context.Context, request model.LoginRequest) (user *model.User, token string, err error) {
	user, err = s.repo.ReadByUsernameOrEmail(ctx, request.Username, request.Email)
	if err != nil {
		logger.Log.Error(err.Error())
		return user, "", err
	}

	if user == nil {
		return nil, "", errors.New("user with this email/username does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		logger.Log.Error(err.Error())
		return user, "", errors.New("login failed, incorrect password")
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := jwtToken.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return user, "", err
	}

	return user, signedToken, nil
}

func (s *userService) Create(ctx context.Context, request model.User) error {
	isUserRegistered, err := s.repo.IsUserRegistered(ctx, request.Username, request.Email)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	if isUserRegistered {
		return cs.ErrUserAlreadyExist
	}

	request.Password, err = utils.HashPassword(request.Password)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	tx := s.repo.BeginTx(ctx)

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	userID, err := s.repo.Create(ctx, tx, request)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	walletAddress := utils.GenerateEncryptedAddress(request.Username, request.Email)

	wallet := model.Wallet{
		Address: walletAddress,
		UserID:  userID,
	}

	if err = s.walletRepo.AddWallet(ctx, tx, wallet); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	tx.Commit()

	return nil
}

func (s *userService) WriteLoginHistory(ctx context.Context, userID int64, isLoginFailed bool) error {
	var req = model.LoginHistory{
		BrowserName: utils.GetBrowserName(ctx),
		UserID:      userID,
	}

	history, err := s.repo.ReadLoginHistory(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	isHistoryExist := true

	if history.UserID == 0 {
		history.BrowserName = req.BrowserName
		history.UserID = req.UserID
		isHistoryExist = false
	}

	if isLoginFailed {
		history.LoginFailed++
	} else {
		history.LoginSucceed++
	}

	if !isHistoryExist {
		if err := s.repo.WriteLoginHistory(ctx, history); err != nil {
			logger.Log.Error(err.Error())
		}
	} else {
		if err := s.repo.UpdateLoginHistory(ctx, history); err != nil {
			logger.Log.Error(err.Error())
		}
	}

	return err
}
