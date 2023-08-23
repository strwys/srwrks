package middleware

import (
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Auth() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JwtCustomClaims)
		},
		SigningKey: []byte(viper.GetString("APP_JWT_SECRET")),
	}
	return echojwt.WithConfig(config)
}
