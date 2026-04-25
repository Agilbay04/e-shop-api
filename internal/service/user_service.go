package service

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"e-shop-api/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService interface {
	Profile(user dto.CurrentUser) (dto.UserResponse, error)
	UploadPicture(req dto.UploadPictureRequest, user dto.CurrentUser) (dto.UserResponse, error)
}

type userService struct {
	db            *gorm.DB
	userRepo      repository.UserRepository
	userQueryRepo repository.UserQueryRepository
	rdb           *redis.Client
}

func NewUserService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	userQueryRepo repository.UserQueryRepository,
	rdb *redis.Client,
) UserService {
	return &userService{
		db,
		userRepo,
		userQueryRepo,
		rdb,
	}
}

func (s *userService) Profile(user dto.CurrentUser) (dto.UserResponse, error) {
	cacheKey := "profile:user:" + user.ID

	cached, err := util.GetCache[dto.UserResponse](s.rdb, cacheKey)
	if err == nil {
		return cached, nil
	}

	userData, err := s.userQueryRepo.FindByID(user.ID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	res := dto.UserResponse{
		ID:       userData.ID.String(),
		Username: userData.Username,
		Email:    userData.Email,
		Role:     userData.Role,
		Picture:  userData.Picture,
	}

	ttl := util.GetEnvTime("24h", "24h")
	_ = util.SetCache(s.rdb, cacheKey, res, ttl)

	return res, nil
}

func (s *userService) UploadPicture(
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

	userData, err := s.userQueryRepo.FindByID(user.ID)
	if err != nil {
		tx.Rollback()
		return dto.UserResponse{}, err
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
		uploader.DeleteFile(newPath)
		return dto.UserResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		uploader.DeleteFile(newPath)
		return dto.UserResponse{}, err
	}

	if oldPath != "" {
		uploader.DeleteFile(oldPath)
	}

	cacheKey := "profile:user:" + user.ID
	_ = util.DeleteCache(s.rdb, cacheKey)

	return dto.UserResponse{
		ID:       userData.ID.String(),
		Username: userData.Username,
		Email:    userData.Email,
		Role:     userData.Role,
		Picture:  userData.Picture,
	}, nil
}