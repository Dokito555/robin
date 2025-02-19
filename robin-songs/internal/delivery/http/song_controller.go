package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/robin-songs/internal/delivery/http/middleware"
	"github.com/Dokito555/robin-songs/internal/model"
	"github.com/Dokito555/robin-songs/internal/services"
	"github.com/Dokito555/robin-songs/utils/constants"
	"github.com/Dokito555/robin-songs/utils/errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongController struct {
	Log         *logrus.Logger
	SongService *services.SongService
}

func NewSongController(log *logrus.Logger, service *services.SongService) *SongController {
	return &SongController{
		Log:         log,
		SongService: service,
	}
}

func (c *SongController) CreateSong(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.CreateSongRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	albumIdStr := ctx.Param("id")
	if albumIdStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	albumId, err := strconv.Atoi(albumIdStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.AlbumID = albumId
	req.ArtistID = auth.UserID

	if auth.Role != constants.ROLE_ARTIST {
		c.Log.Warnf("auth role is not artist")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	rsp, err := c.SongService.CreateNewSong(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to create new song: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, model.BaseResponse[*model.SongResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *SongController) GetSong(ctx *gin.Context) {
	req := new(model.GetSongRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	IdStr := ctx.Param("id")
	if IdStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = Id

	rsp, err := c.SongService.GetSong(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to get song: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.SongResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *SongController) UpdateSong(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.UpdateSongRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	IdStr := ctx.Param("id")
	if IdStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = Id


	if auth.Role != constants.ROLE_ARTIST && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not artist or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	rsp, err := c.SongService.UpdateSong(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update song: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.SongResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *SongController) DeleteSong(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.DeleteSongRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		c.Log.Warnf("failed to bind request to JSON: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	IdStr := ctx.Param("id")
	if IdStr == "" {
		c.Log.Warnf("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
		return
	}

	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		c.Log.Warnf("failed to convert id string to int")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_INTERNAL_SERVER_ERROR))
		return
	}

	req.ID = Id


	if auth.Role != constants.ROLE_ARTIST && auth.Role != constants.ROLE_ADMIN {
		c.Log.Warnf("auth role is not artist or admin")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_UNAUTHORIZED)
		return
	}

	err = c.SongService.DeleteSong(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Warnf("failed to update song: %+v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[bool]{Message: http.StatusOK, Data: true})
}