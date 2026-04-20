package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.UserResponse, error)
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
}

type authService struct {
	db            *gorm.DB
	userRepo      repository.UserRepository
	userQueryRepo repository.UserQueryRepository
}

func NewAuthService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	userQueryRepo repository.UserQueryRepository,
) AuthService {
	return &authService{
		db,
		userRepo,
		userQueryRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (dto.UserResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	u, err := s.userQueryRepo.FindByEmail(req.Email)
	if err == nil && u != nil {
		tx.Rollback()
		return dto.UserResponse{}, errors.New("user already exists")
	}

	newUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(tx, &newUser); err != nil {
		tx.Rollback()
		return dto.UserResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	u, err := s.userQueryRepo.FindByEmail(req.Email)
	if err != nil {
		return dto.LoginResponse{}, util.UnprocessableEntityException("User email " + req.Email + " is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, util.UnauthorizedException("Invalid email or password")
	}

	token, err := util.GenerateToken(u.ID, u.Username, u.Email, u.Role)
	if err != nil {
		return dto.LoginResponse{}, util.UnauthorizedException("Token is invalid or expired")
	}

	return dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		},
	}, nil
}
