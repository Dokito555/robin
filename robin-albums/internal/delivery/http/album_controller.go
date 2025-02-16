package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin-albums/internal/delivery/http/middleware"
	model "github.com/Dokito555/robin-albums/internal/models"
	"github.com/Dokito555/robin-albums/internal/services"
	"github.com/Dokito555/robin-albums/utils/constants"
	"github.com/Dokito555/robin-albums/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AlbumController struct {
	Log          *logrus.Logger
	AlbumService *services.AlbumService
}

func NewAlbumService(log *logrus.Logger, service *services.AlbumService) *AlbumController {
	return &AlbumController{
		Log:          log,
		AlbumService: service,
	}
}

func (c *AlbumController) CreateAlbum(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.CreateAlbumRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	if auth.Role != constants.ROLE_ARTIST {
		c.Log.Warnf("auth role is not artist")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	req.ArtistID = auth.UserID
	rsp, err := c.AlbumService.CreateAlbum(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to create album: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusCreated, model.BaseResponse[*model.AlbumResponse]{Message: http.StatusCreated, Data: rsp})
}

func (c *AlbumController) GetAlbum(ctx *gin.Context) {
	req := new(model.GetAlbumRequest)
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

	rsp, err := c.AlbumService.GetAlbum(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get album: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, model.BaseResponse[*model.AlbumResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *AlbumController) UpdateAlbum(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.UpdateAlbumRequest)
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

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	if auth.Role != constants.ROLE_ARTIST && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not artist or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	rsp, err := c.AlbumService.UpdateAlbum(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update album: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, model.BaseResponse[*model.AlbumResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *AlbumController) DeleteAlbum(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.DeleteAlbumRequest)
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

	if auth.Role != constants.ROLE_ARTIST && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not artist or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	err = c.AlbumService.DeleteAlbum(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update album: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, model.BaseResponse[*model.AlbumResponse]{Message: http.StatusOK, Data: nil})
}
