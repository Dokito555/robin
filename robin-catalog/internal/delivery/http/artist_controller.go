package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http/middleware"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
	"github.com/Dokito555/robin/robin-catalog/internal/services"
	"github.com/Dokito555/robin/robin-catalog/utils/constants"
	"github.com/Dokito555/robin/robin-catalog/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArtistController struct {
	Log     *logrus.Logger
	Service *services.ArtistService
}

func NewArtistController(log *logrus.Logger, service *services.ArtistService) *ArtistController {
	return &ArtistController{
		Log:     log,
		Service: service,
	}
}

func (c *ArtistController) RegisterArtist(ctx *gin.Context) {
	req := new(model.RegisterArtistRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	if req.Role == constants.ROLE_ADMIN {
		c.Log.Warnf("failed to register as admin")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	rsp, err := c.Service.RegisterArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to register artist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *ArtistController) LoginArtist(ctx *gin.Context) {
	req := new(model.LoginArtistRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest,  errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
	}

	rsp, err := c.Service.LoginArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to login artist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *ArtistController) GetArtist(ctx *gin.Context) {
	req := new(model.GetArtistRequest)
	idStr := ctx.Param("id")
	if idStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusInternalServerError, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = id

	rsp, err := c.Service.GetArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get artist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *ArtistController) UpdateArtist(ctx *gin.Context) {
	req := new(model.UpdateArtistRequest)
	auth := middleware.GetProfile(ctx)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
	}

	req.ID = auth.UserID

	rsp, err := c.Service.UpdateArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get artist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *ArtistController) LogoutArtist(ctx *gin.Context) {
	req := new(model.LogoutArtistRequest)
	token := ctx.GetHeader("Authorization")
	if token == "" {
		c.Log.Warnf("token is empty")
		ctx.JSON(http.StatusBadRequest,errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	req.Token = token

	err := c.Service.LogoutArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to logout artist: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: nil})
}

func (c *ArtistController) DeleteArtist(ctx *gin.Context) {
	req := new(model.DeleteArtistRequest)
	auth := middleware.GetProfile(ctx)
	idStr := ctx.Param("id")

	if auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("unauthorized access")
		ctx.JSON(http.StatusUnauthorized, errs.NewErrorResponse(errs.ERROR_UNAUTHORIZED))
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
		ctx.JSON(http.StatusInternalServerError, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = id

	err = c.Service.DeleteArtist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get user: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.ArtistResponse]{Message: http.StatusOK, Data: nil})
}

func (c *ArtistController) GetArtistList(ctx *gin.Context) {
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

	rsps, err := c.Service.GetArtistList(ctx.Request.Context(), page, limit)
	if err != nil {
		c.Log.Warnf("failed to get artist list: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, model.BaseResponse[[]model.ArtistResponse]{Message: http.StatusOK, Data: rsps})
}