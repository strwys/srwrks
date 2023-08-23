package handler

import (
	"net/http"

	"github.com/cecepsprd/starworks-test/constans"
	m "github.com/cecepsprd/starworks-test/internal/handler/middleware"
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/internal/service"
	"github.com/cecepsprd/starworks-test/utils"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(e *echo.Echo, walletService service.WalletService) {
	handler := &WalletHandler{
		walletService: walletService,
	}

	e.GET("/api/wallet/check-balance", handler.CheckBalance, m.Auth())
	e.POST("/api/wallet/top-up", handler.TopUp, m.Auth())
	e.POST("/api/wallet/pay", handler.Pay, m.Auth())

}

// @Summary      Check Balance
// @Description  Check Balance endpoint
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        checkBalanceRequest   body    model.CheckBalanceRequest  true  "Check Balance Request"
// @Success      200  {object}  model.APIResponse
// @Failure      400  {object}  model.ResponseError
// @Router       /api/wallet/check-balance [get]
func (h *WalletHandler) CheckBalance(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = model.CheckBalanceRequest{}
		res = model.CheckBalanceResponse{}
	)

	user := utils.GetUserByContext(c)
	req.UserID = user.ID
	req.Address = utils.GenerateEncryptedAddress(user.Username, user.Email)

	wallet, err := h.walletService.CheckBalance(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	if err := utils.MappingInterface(wallet, &res); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: constans.MessageSuccess,
		Data:    res,
	})
}

// @Summary      Top Up
// @Description  Top Up endpoint
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        topUpRequest   body    model.TopUpRequest  true  "Top Up Request"
// @Success      200  {object}  model.APIResponse
// @Failure      400  {object}  model.ResponseError
// @Router       /api/wallet/top-up [post]
func (h *WalletHandler) TopUp(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = model.TopUpRequest{}
	)

	err := c.Bind(&req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	user := utils.GetUserByContext(c)
	req.UserID = user.ID
	req.Address = utils.GenerateEncryptedAddress(user.Username, user.Email)

	if err = c.Validate(req); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	err = h.walletService.TopUp(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: constans.MessageSuccess,
	})
}

// @Summary      Pay
// @Description  Pay endpoint
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        payRequest   body    model.PayRequest  true  "Pay Request"
// @Success      200  {object}  model.APIResponse
// @Failure      400  {object}  model.ResponseError
// @Router       /api/wallet/pay [post]
func (h *WalletHandler) Pay(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = model.PayRequest{}
	)

	err := c.Bind(&req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	user := utils.GetUserByContext(c)
	req.UserID = user.ID
	req.Address = utils.GenerateEncryptedAddress(user.Username, user.Email)

	if err = c.Validate(req); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	err = h.walletService.Pay(ctx, req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusOK,
		Message: constans.MessageSuccess,
	})
}
