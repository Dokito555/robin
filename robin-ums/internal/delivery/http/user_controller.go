package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin-ums/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-ums/internal/model"
	"github.com/Dokito555/robin-ums/internal/services"
	"github.com/Dokito555/robin-ums/utils/constants"
	"github.com/Dokito555/robin-ums/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	Service *services.UserService
}

func NewUserController(log *logrus.Logger, service *services.UserService) *UserController {
	return &UserController{
		Log:     log,
		Service: service,
	}
}

func (c *UserController) RegisterUser(ctx *gin.Context) {
	req := new(model.RegisterUserRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	if req.Role == constants.ROLE_ADMIN {
		c.Log.Warnf("failed to register as admin")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_UNAUTHORIZED))
		return
	}

	rsp, err := c.Service.Register(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to register user: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: rsp})
	return
}

func (c *UserController) RegisterAdmin(ctx *gin.Context) {
	req := new(model.RegisterUserRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	if req.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("failed to register")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_UNAUTHORIZED))
		return
	}

	rsp, err := c.Service.Register(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to register user: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: rsp})
	return
}

func (c *UserController) Login(ctx *gin.Context) {
	req := new(model.LoginUserRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
	}

	rsp, err := c.Service.Login(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to login user: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: rsp})
	return
}

func (c *UserController) Logout(ctx *gin.Context) {
	req := new(model.LogoutUserRequest)
	token := ctx.GetHeader("Authorization")
	if token == "" {
		c.Log.Warnf("token is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_UNAUTHORIZED))
		return
	}

	req.Token = token

	err := c.Service.Logout(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to logout user: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[interface{}]{Message: http.StatusOK, Data: nil})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	req := new(model.GetUserRequest)
	idStr := ctx.Param("id")
	if idStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = id

	rsp, err := c.Service.GetUser(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get user: %v", err)
		c.Log.Warnf("failed to logout user: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	req := new(model.DeleteUserRequest)
	auth := middleware.GetProfile(ctx)
	idStr := ctx.Param("id")

	if auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("unauthorized access")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_UNAUTHORIZED))
		return
	}

	if idStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = id

	err = c.Service.DeleteUser(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to delete user: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: nil})
}

func (c *UserController) VerifyUser(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)

	if auth == nil {
		c.Log.Warn("failed to get profile from ctx")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req := &model.GetUserRequest{ID: auth.UserID}
	rsp, err := c.Service.GetUser(ctx, req)
	if err != nil {
		c.Log.Warnf("Failed to verify user: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	req := new(model.UpdateUserRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
	}

	rsp, err := c.Service.UpdateUser(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update user: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: rsp})
	return
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	var (
		pageStr = ctx.DefaultQuery("page", "1")
		limitStr = ctx.DefaultQuery("limit", "10")
	)

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Log.Warnf("failed to convert page string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.Log.Warnf("failed to convert limit string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	rsps, err := c.Service.GetUsers(ctx.Request.Context(), page, limit)
	if err != nil {
		c.Log.Warnf("failed to login user: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[[]model.UserResponse]{Message: http.StatusOK, Data: rsps})
}