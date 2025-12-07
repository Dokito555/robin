package http

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Dokito555/robin/robin-catalog/internal/delivery/http/middleware"
	"github.com/Dokito555/robin/robin-catalog/internal/model"
	"github.com/Dokito555/robin/robin-catalog/internal/services"
	"github.com/Dokito555/robin/robin-catalog/utils/constants"
	"github.com/Dokito555/robin/robin-catalog/utils/errs"
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
	if auth.Role != constants.ROLE_ARTIST {
		c.Log.Warnf("unauthorized: user role is %s", auth.Role)
		ctx.JSON(http.StatusUnauthorized, errs.ERROR_UNAUTHORIZED)
		return
	}

	name := ctx.PostForm("name")
	if name == "" {
		c.Log.Warn("song name is required")
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	albumId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Log.Warnf("invalid album ID: %v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.Warnf("file not found: %+v", err)
		ctx.JSON(http.StatusBadRequest, errs.ERROR_BAD_REQUEST)
		return
	}

	if filepath.Ext(file.Filename) != ".mp3" {
		c.Log.Warnf("invalid file type: %s", file.Filename)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "only MP3 files are allowed"})
		return
	}

	// max size 50mb
	const maxFileSize = 50 * 1024 * 1024
	if file.Size > maxFileSize {
		c.Log.Warnf("file too large: %d bytes", file.Size)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 50MB limit"})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		c.Log.Warnf("failed to open file: %+v", err)
		ctx.JSON(http.StatusInternalServerError, errs.ERROR_INTERNAL_SERVER_ERROR)
		return
	}
	defer uploadedFile.Close()

	duration, err := c.SongService.CalculateMP3Duration(uploadedFile)
	if err != nil {
		c.Log.Warnf("failed to calculate MP3 duration: %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid MP3 file"})
		return
	}

	if duration == 0 {
		c.Log.Warn("MP3 duration is zero")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid MP3 file or empty audio"})
		return
	}

	uploadedFile.Close()
	uploadedFile, err = file.Open()
	if err != nil {
		c.Log.Warnf("failed to reopen file: %+v", err)
		ctx.JSON(http.StatusInternalServerError, errs.ERROR_INTERNAL_SERVER_ERROR)
		return
	}
	defer uploadedFile.Close()

	fileReq := &model.UploadFileRequest{
		File:       uploadedFile,
		FileHeader: file,
	}

	req := &model.CreateSongRequest{
		Name:     name,
		AlbumID:  albumId,
		ArtistID: auth.UserID,
		Duration: int(duration.Milliseconds()),
	}

	rsp, err := c.SongService.CreateNewSong(ctx.Request.Context(), req, fileReq)
	if err != nil {
		c.Log.Warnf("failed to create song: %+v", err)
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

func (c *SongController) GetSongStream(ctx *gin.Context) {
	req := new(model.GetSongRequest)
	idStr := ctx.Param("id")

	if idStr == "" {
		c.Log.Warn("id is empty")
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

	rsp, err := c.SongService.GetSongStream(ctx, req)
	if err != nil {
		c.Log.Warnf("failed to get song stream: %v", err)
		appErr, ok := err.(*errs.AppError)
		if !ok {
			appErr = errs.ERROR_INTERNAL_SERVER_ERROR
		}
		ctx.JSON(appErr.Code, errs.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, model.BaseResponse[*model.SongStreamResponse]{Message: http.StatusOK, Data: rsp})
}

func (c *SongController) DeleteSong(ctx *gin.Context) {
	auth := middleware.GetProfile(ctx)
	req := new(model.DeleteSongRequest)
	idStr := ctx.Param("id")

	if idStr == "" {
		c.Log.Warn("id is empty")
		ctx.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.ERROR_BAD_REQUEST))
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
