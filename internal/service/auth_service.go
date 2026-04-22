package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/model"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (dto.UserResponse, error)
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	Profile(user dto.CurrentUser) (dto.UserResponse, error)
	UploadPicture(req dto.UploadPictureRequest, user dto.CurrentUser) (dto.UserResponse, error)
}

type authService struct {
	db            *gorm.DB
	userRepo      repository.UserRepository
	userQueryRepo repository.UserQueryRepository
	rdb           *redis.Client
}

func NewAuthService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	userQueryRepo repository.UserQueryRepository,
	rdb *redis.Client,
) AuthService {
	return &authService{
		db,
		userRepo,
		userQueryRepo,
		rdb,
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
		return dto.UserResponse{}, util.BadRequestException("Email already use by another account", err)
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
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	// Define cache key
    cacheKey := "user:email:" + req.Email
    
    // Try to get data from Redis
    u, err := util.GetCache[*model.User](s.rdb, cacheKey)
    
    // If data not found in Redis
    if err != nil {
        // Get data from database
        u, err = s.userQueryRepo.FindByEmail(req.Email)
        if err != nil {
            return dto.LoginResponse{}, util.UnprocessableEntityException("User email " + req.Email + " is not registered")
        }

        // Set data to Redis
		ttl := util.GetEnvInt(os.Getenv("REDIS_CACHE_TTL"), 5)
        _ = util.SetCache(s.rdb, cacheKey, u, time.Duration(ttl)*time.Minute)
    }

    // Validate password
    err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
    if err != nil {
        return dto.LoginResponse{}, util.UnauthorizedException("Invalid email or password")
    }

    // Generate token
    token, err := util.GenerateToken(u.ID, u.Username, u.Email, u.Picture, u.Role)
    if err != nil {
        return dto.LoginResponse{}, util.UnauthorizedException("Token is invalid or expired")
    }

    return dto.LoginResponse{
        Token: token,
        User: dto.UserResponse{
            ID:       u.ID.String(),
            Username: u.Username,
            Email:    u.Email,
            Role:     u.Role,
        },
    }, nil
}

func (s *authService) Profile(user dto.CurrentUser) (dto.UserResponse, error) {
	return dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Picture:  user.Picture,
	}, nil
}

func (s *authService) UploadPicture(
	req dto.UploadPictureRequest, 
	user dto.CurrentUser,
) (dto.UserResponse, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	uploader := util.NewFileUploader(
		util.WithDirectory("uploads/avatars"),
		util.WithMaxSize(2),
		util.WithExtensions([]string{".jpg", ".jpeg", ".png", ".webp"}),
	)

	userData, err := s.userQueryRepo.FindByID(user.ID.String())
	if err != nil {
		tx.Rollback()
		return dto.UserResponse{}, err;
	}

	oldPath := userData.Picture

	newPath, err := uploader.UploadFile(req.Picture)
	if err != nil {
		tx.Rollback()
		return dto.UserResponse{}, err
	}

	userData.Picture = newPath
	if err := s.userRepo.Update(tx, userData); err != nil {
		tx.Rollback()
		os.Remove(newPath)
		return dto.UserResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		os.Remove(newPath)
		return dto.UserResponse{}, err
	}

	if oldPath != "" {
		os.Remove(oldPath)
	}

	return dto.UserResponse{
		ID:       userData.ID.String(),
		Username: userData.Username,
		Email:    userData.Email,
		Role:     userData.Role,
		Picture:  userData.Picture,
	}, nil
}

