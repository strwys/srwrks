package handler

import (
	"context"
	"net/http"

	cs "github.com/cecepsprd/starworks-test/constans"
	"github.com/cecepsprd/starworks-test/internal/model"
	"github.com/cecepsprd/starworks-test/internal/service"
	"github.com/cecepsprd/starworks-test/utils"
	"github.com/cecepsprd/starworks-test/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(e *echo.Echo, userService service.UserService) {
	handler := &UserHandler{
		userService: userService,
	}

	e.Use(middleware.Recover())
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/register", handler.Register)
}

// @Description  Login endpoint
// @Tags         user
// @Param        request   body    model.LoginRequest  true  "Login Request"
// @Success      200  {object}  model.LoginResponse
// @Failure      400  {object}  model.ResponseError
// @Router       /api/user/login [get]
func (h *UserHandler) Login(c echo.Context) error {
	var (
		req = model.LoginRequest{}
		res = model.LoginResponse{}
		ctx = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	ctx = context.WithValue(ctx, cs.CtxUserAgent, c.Request().UserAgent())

	user, token, err := h.userService.Login(ctx, req)

	if user != nil {
		h.userService.WriteLoginHistory(ctx, user.ID, err != nil)
	}

	if err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	res.Token = token

	if err := utils.MappingInterface(user, &res.Data); err != nil {
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// @Description  If successful, 'data' will contain an instance of model.User. If an error occurs, 'data' will not be shown.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request   body    model.RegisterRequest  true  "Register Request"
// @Success      200  {object}  model.APIResponse
// @Failure      400  {object}  model.ResponseError
// @Router       /api/user/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		req  = model.RegisterRequest{}
		res  = model.UserPresenter{}
		user = model.User{}
	)

	err := c.Bind(&req)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	if err = c.Validate(req); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, model.ResponseError{Message: err.Error()})
	}

	if err = utils.MappingInterface(req, &user); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: err.Error()})
	}

	err = h.userService.Create(ctx, user)
	if err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	if err = utils.MappingInterface(user, &res); err != nil {
		logger.Log.Error(err.Error())
		return c.JSON(utils.SetHTTPStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, model.APIResponse{
		Code:    http.StatusCreated,
		Message: cs.MessageSuccess,
		Data:    res,
	})
}
