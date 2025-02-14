package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin-ums/constants"
	"github.com/Dokito555/robin-ums/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-ums/internal/model"
	"github.com/Dokito555/robin-ums/internal/services"
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if req.Role == constants.ROLE_ADMIN {
		c.Log.Warnf("failed to register as admin")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	rsp, err := c.Service.Register(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to register user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if req.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("failed to register")
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	rsp, err := c.Service.Register(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to register user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusBadRequest, err)
	}

	rsp, err := c.Service.Login(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to login user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "unauthorized"})
		return
	}

	req.Token = token

	err := c.Service.Logout(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to logout user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[interface{}]{Message: http.StatusOK, Data: nil})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	req := new(model.GetUserRequest)
	idStr := ctx.Param("id")
	if idStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "id is empty"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	req.ID = id

	rsp, err := c.Service.GetUser(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"Message": ""})
		return
	}

	if idStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "id is empty"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	req.ID = id

	err = c.Service.DeleteUser(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.UserResponse]{Message: http.StatusOK, Data: nil})
}


func (c *UserController) VerifyUser(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)

	if auth == nil {
		c.Log.Warn("failed to get profile from ctx")
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "couldn't get profile from ctx"})
		return
	}

	req := &model.GetUserRequest{ID: auth.UserID}
	rsp, err := c.Service.GetUser(ctx, req)
	if err != nil {
		c.Log.Warnf("Failed to verify user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Message": "couldn't verify user"})
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
