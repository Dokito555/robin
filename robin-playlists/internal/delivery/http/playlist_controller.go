package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin-playlists/internal/delivery/http/middleware"
	model "github.com/Dokito555/robin-playlists/internal/models"
	"github.com/Dokito555/robin-playlists/internal/services"
	"github.com/Dokito555/robin-playlists/internal/utils/constants"
	"github.com/Dokito555/robin-playlists/internal/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PlaylistController struct {
	Log     *logrus.Logger
	PlaylistService *services.PlaylistService
}

func NewPlaylistController(log *logrus.Logger, service *services.PlaylistService) *PlaylistController {
	return &PlaylistController{
		Log:     log,
		PlaylistService: service,
	}
}

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.CreatePlaylistRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	if auth.Role != constants.ROLE_USER {
		c.Log.Warnf("auth role is not user")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	req.UserID = auth.UserID
	rsp, err := c.PlaylistService.CreatePlaylist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to create playlist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, model.BaseResponse[*model.PlaylistResponse]{Message: http.StatusCreated, Data: rsp})
}

func (c *PlaylistController) GetPlaylist(ctx *gin.Context) {
	req := new(model.GetPlaylistRequest)
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

	rsp, err := c.PlaylistService.GetPlaylist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get playlist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.PlaylistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.UpdatePlaylistRequest)
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

	
	if auth.Role != constants.ROLE_USER && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not user or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	rsp, err := c.PlaylistService.UpdatePlaylist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update playlist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.PlaylistResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.DeletePlaylistRequest)
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

	if auth.Role != constants.ROLE_USER && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not artist or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	err = c.PlaylistService.DeletePlaylist(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to delete playlist: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[bool]{Message: http.StatusOK, Data: true})
}