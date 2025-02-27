package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dokito555/robin-ums/internal/entity"
	"github.com/Dokito555/robin-ums/internal/model"
	"github.com/Dokito555/robin-ums/internal/model/converter"
	"github.com/Dokito555/robin-ums/internal/repository"
	"github.com/Dokito555/robin-ums/utils/constants"
	"github.com/Dokito555/robin-ums/utils/errs"
	"github.com/IBM/sarama"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	Producer       sarama.SyncProducer
	Config         *viper.Viper
	UserRepository *repository.UserRepository
	TokenService   *TokenService
	Messaging      *MessagingService
}

func NewUserService(db *gorm.DB, logger *logrus.Logger, config *viper.Viper, validate *validator.Validate, producer sarama.SyncProducer,
	userRepository *repository.UserRepository, tokenSvc *TokenService, Messaging *MessagingService) *UserService {
	return &UserService{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		Producer:       producer,
		Config:         config,
		UserRepository: userRepository,
		TokenService:   tokenSvc,
		Messaging:      Messaging,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	s.Log.Info("starting Register function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" || req.Email == "" {
		return nil, errs.NewError(http.StatusBadRequest, "password or email are requried")
	}

	user, err := s.UserRepository.FindByEmail(s.DB, &entity.User{}, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching user: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if user != nil {
		s.Log.Warnf("user with that email already exists")
		return nil, errs.ERROR_USER_EXIST
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	user = &entity.User{
		Password: string(password),
		Email:    req.Email,
		UserName: req.Username,
		Role:     req.Role,
	}

	if err := s.UserRepository.Create(s.DB, user); err != nil {
		s.Log.Warnf("failed to create user in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	kafkaPayload := model.RegisterPayload{
		UserName: req.Username,
		Role:     req.Role,
	}

	jsonPayload, err := json.Marshal(kafkaPayload)
	if err != nil {
		// TODO: honestly sending email is optional for registration
		// might seperate this later
		s.Log.Warn("failed to marshal kafka payload")
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	err = s.Messaging.ProduceKafkaMessage(s.Config.GetString("KAFKA_REGISTER_TOPIC"), jsonPayload)
	if err != nil {
		// TODO: if kafka failed the registration system also failed
		// error should be optional?
		s.Log.Warn("failed from kafka")
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.UserToResponse(user), nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	s.Log.Info("starting Login function")
	s.Log.Infof("request received: %+v", req)

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body: %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" || req.Email == "" {
		return nil, errs.NewError(http.StatusBadRequest, "password or email are required")
	}

	newUser := new(entity.User)
	user, err := s.UserRepository.FindByEmail(s.DB, newUser, req.Email)
	if err != nil {
		s.Log.Warnf("database error fetching user: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.Log.Warnf("failed to compare hashed password: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	token, err := s.TokenService.GenerateToken(ctx, user.ID, constants.TOKEN_TYPE_TOKEN, user.Email, user.Role, user.UserName)
	if err != nil {
		s.Log.Warnf("failed to genereate token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	refreshToken, err := s.TokenService.GenerateToken(ctx, user.ID, constants.TOKEN_TYPE_REFRESH, user.Email, user.Role, user.UserName)
	if err != nil {
		s.Log.Warnf("failed to genereate refresh token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	newUser.Token = token
	newUser.RefreshToken = refreshToken

	err = s.UserRepository.Update(s.DB, newUser)
	if err != nil {
		s.Log.Warnf("failed to update user token and refresh token: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.UserToResponse(newUser), nil
}

func (s *UserService) Logout(ctx context.Context, req *model.LogoutUserRequest) error {
	s.Log.Info("starting Logout function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	user := new(entity.User)
	_, err = s.UserRepository.FindByToken(s.DB, user, req.Token)
	if err != nil {
		s.Log.Warnf("failed to find user in database: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	user.Token = ""

	err = s.UserRepository.Update(tx, user)
	if err != nil {
		s.Log.Warnf("failed to update user: %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	s.Log.Info("starting Get User function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	user := new(entity.User)
	if err := s.UserRepository.FindById(s.DB, user, req.ID); err != nil {
		s.Log.Warnf("failed to find user in database: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err == gorm.ErrRecordNotFound {
		s.Log.Warnf("user not found")
		return nil, errs.ERROR_NOT_FOUND
	}

	user.Token = ""
	user.RefreshToken = ""

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.UserToResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *model.DeleteUserRequest) error {
	s.Log.Info("starting Delete User function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return errs.ERROR_BAD_REQUEST
	}

	user := new(entity.User)
	err = s.UserRepository.FindById(tx, user, req.ID)
	if err != nil {
		s.Log.Warnf("failed to find user by token : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err == gorm.ErrRecordNotFound {
		s.Log.Warnf("user not found: %+v", err)
		return errs.ERROR_NOT_FOUND
	}

	err = s.UserRepository.Delete(s.DB, user)
	if err != nil {
		s.Log.Warnf("failed to delete user by token : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	s.Log.Info("starting Update User function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	if req.Password == "" {
		return nil, errors.New("password is empty")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("failed to generate bcrypt hash: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	user := &entity.User{
		Password: string(password),
		UserName: req.Username,
	}

	err = s.UserRepository.Update(s.DB, user)
	if err != nil {
		s.Log.Warnf("failed to update: %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.UserToResponse(user), nil
}

func (s *UserService) Verify(ctx context.Context, req *model.VerifyRequest) (*model.UserResponse, error) {
	s.Log.Info("starting Verify function")
	s.Log.Infof("request received: %+v", req)
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Warnf("invalid request body : %+v", err)
		return nil, errs.ERROR_BAD_REQUEST
	}

	user := new(entity.User)
	rsp, err := s.UserRepository.FindByToken(tx, user, req.Token)
	if err != nil {
		s.Log.Warnf("failed find user by token : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return converter.UserToResponse(rsp), nil
}

func (s *UserService) GetUsers(ctx context.Context, page int, limit int) ([]model.UserResponse, error) {
	s.Log.Info("starting get users function")
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page == 0 || limit == 0 {
		s.Log.Warn("page or limit is invalid")
		return nil, errs.ERROR_BAD_REQUEST
	}

	users, err := s.UserRepository.GetUserList(s.DB, page, limit)
	if err != nil {
		s.Log.Warn("failed to get users in database")
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	rsps := make([]model.UserResponse, len(users))
	for i, user := range users {
		rsps[i] = *converter.UserToResponse(&user)
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("failed commit transaction : %+v", err)
		return nil, errs.ERROR_INTERNAL_SERVER_ERROR
	}

	return rsps, nil
}
