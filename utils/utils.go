package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	cs "github.com/cecepsprd/starworks-test/constans"
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mssola/user_agent"
	"golang.org/x/crypto/bcrypt"
)

func MappingInterface(request interface{}, model interface{}) error {
	// convert interface to json
	jsonRecords, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error encode records: %s", err.Error())
	}

	// bind json to struct
	if err := json.Unmarshal(jsonRecords, model); err != nil {
		return fmt.Errorf("error decode json to struct: %s", err.Error())
	}

	return nil
}

func HashPassword(password string) (hashedString string, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}

	return string(hashed), nil
}

func GenerateEncryptedAddress(username, email string) string {
	data := fmt.Sprintf("%s:%s", username, email)
	hash := sha256.Sum256([]byte(data))
	encryptedAddress := hex.EncodeToString(hash[:])
	return encryptedAddress
}

func GetUserIDByContext(ctx echo.Context) int64 {
	u := ctx.Get("user")
	claims := u.(*jwt.Token).Claims.(*model.JwtCustomClaims)
	return claims.UserID
}

func GetUserByContext(ctx echo.Context) model.User {
	u := ctx.Get("user")
	claims := u.(*jwt.Token).Claims.(*model.JwtCustomClaims)
	return model.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}
}

func GetBrowserName(ctx context.Context) string {
	userAgent := GetContextValue(ctx, cs.CtxUserAgent)
	browserName, _ := user_agent.New(userAgent).Browser()
	return browserName
}

func GetContextValue(ctx context.Context, key interface{}) string {
	val := ctx.Value(key)
	if val != nil {
		str, ok := val.(string)
		if ok {
			return str
		}
	}
	return ""
}

func SetHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err.Error() {
	case cs.ErrInternalServerError.Error():
		return http.StatusInternalServerError
	case cs.ErrNotFound.Error():
		return http.StatusNotFound
	// case cs.ErrConflict.Error():
	// return http.StatusConflict
	case cs.ErrWrongEmailOrPassword.Error():
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
